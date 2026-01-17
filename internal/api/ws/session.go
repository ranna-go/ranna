package ws

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/zekrotja/rogu"
	"github.com/zekrotja/rogu/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
)

var sessionPool = sync.Pool{
	New: func() any {
		return &session{}
	},
}

type session struct {
	manager SandboxManager
	logger  rogu.Logger
	conn    *websocket.Conn
	rlm     *RateLimitManager
}

func newSession(rlm *RateLimitManager, manager SandboxManager) (t *session) {
	t = sessionPool.Get().(*session)
	t.conn = nil
	t.manager = manager
	t.logger = log.Tagged("WS")
	t.rlm = rlm
	return t
}

func (t *session) Close() {
	t.logger.Debug().
		Field("addr", getAddr(t.conn)).
		Msg("websocket connection closed")
	sessionPool.Put(t)
}

func (t *session) Handler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		t.logger.Debug().
			Field("addr", getAddr(c)).
			Msg("new websocket connection")

		t.conn = c
		var (
			typ   int
			msg   []byte
			err   error
			nonce int
		)
		for {
			if typ, msg, err = c.ReadMessage(); err != nil {
				t.Close()
				break
			}
			if typ != websocket.TextMessage {
				t.SendError(0, models.ErrInvalidMessageType)
			}
			go func() {
				if nonce, err = t.HandleOp(msg); err != nil {
					t.SendError(nonce, err)
				}
			}()
		}
	})
}

func (t *session) Send(v models.Event) (err error) {
	return t.conn.WriteJSON(v)
}

func (t *session) SendError(nonce int, err error) error {
	var data models.WsError
	if !errors.As(err, &data) {
		data = models.WsError{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	return t.Send(models.Event{
		Code:  models.EventError,
		Nonce: nonce,
		Data:  data,
	})
}

func (t *session) HandleOp(msg []byte) (nonce int, err error) {
	var op models.Operation
	if err = json.Unmarshal(msg, &op); err != nil {
		return 0, err
	}

	if !t.rlm.GetLimiter(t.conn, op.Op).Allow() {
		return op.Nonce, models.ErrRateLimited
	}

	var event models.Event
	event.Nonce = op.Nonce

	switch op.Op {
	case models.OpPing:
		event.Code = models.EventPong
		event.Data = "Pong!"
		err = t.Send(event)
	case models.OpExec:
		var eop models.OperationExec
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = t.handleExec(eop)
		}
	case models.OpKill:
		var eop models.OperationKill
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = t.handleKill(eop)
		}
	default:
		err = models.ErrInvalidOpCode
	}

	return op.Nonce, err
}

func (t *session) handleExec(op models.OperationExec) (err error) {
	if op.Args.Code == "" {
		return models.ErrEmptyCode
	}

	cSpn := make(chan string, 1)
	cStdOut := make(chan []byte)
	cStdErr := make(chan []byte)
	cStop := make(chan struct{}, 1)

	var runId string

	go func() {
		for {
			select {
			case <-cStop:
				return
			case runId = <-cSpn:
				err = t.Send(models.Event{
					Code:  models.EventSpawn,
					Nonce: op.Nonce,
					Data: models.DataSpawn{
						DataRunId: models.DataRunId{
							RunId: runId,
						},
					},
				})
			case p := <-cStdOut:
				err = t.Send(models.Event{
					Code:  models.EventLog,
					Nonce: op.Nonce,
					Data: models.DataLog{
						DataRunId: models.DataRunId{
							RunId: runId,
						},
						StdOut: string(p),
					},
				})
			case p := <-cStdErr:
				err = t.Send(models.Event{
					Code:  models.EventLog,
					Nonce: op.Nonce,
					Data: models.DataLog{
						DataRunId: models.DataRunId{
							RunId: runId,
						},
						StdErr: string(p),
					},
				})
			}
			if err != nil {
				t.logger.Error().Err(err).Msg("Failed sending event")
				if err = t.SendError(op.Nonce, err); err != nil {
					t.logger.Error().Err(err).Msg("Failed sending error event")
					return
				}
			}
		}
	}()
	defer func() {
		cStop <- struct{}{}
	}()

	execTime := util.MeasureTime(func() {
		err = t.manager.RunInSandbox(context.TODO(), &op.Args, cSpn, cStdOut, cStdErr)
	})

	if err != nil {
		if sandbox.IsSystemError(err) {
			return t.SendError(op.Nonce, err)
		}
		return t.SendError(op.Nonce, models.WsError{Code: http.StatusBadRequest, Message: err.Error()})
	}

	err = t.Send(models.Event{
		Code:  models.EventStop,
		Nonce: op.Nonce,
		Data: models.DataStop{
			DataRunId: models.DataRunId{
				RunId: runId,
			},
			ExecTimeMS: int(execTime.Milliseconds()),
		},
	})

	return err
}

func (t *session) handleKill(op models.OperationKill) (err error) {
	ok, err := t.manager.KillAndCleanUp(context.TODO(), op.Args.RunId)
	if err != nil {
		return
	}
	if !ok {
		err = models.ErrSandboxNotRunning
	}
	return
}
