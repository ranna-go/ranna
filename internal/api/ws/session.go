package ws

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
	"github.com/ranna-go/ranna/pkg/models"
	"github.com/sarulabs/di/v2"
	"github.com/sirupsen/logrus"
)

var sessionPool = sync.Pool{
	New: func() interface{} {
		return &session{}
	},
}

type session struct {
	conn    *websocket.Conn
	manager sandbox.Manager
	rlm     *RateLimitManager
}

func newSession(rlm *RateLimitManager, ctn di.Container) (s *session) {
	s = sessionPool.Get().(*session)
	s.conn = nil
	s.manager = ctn.Get(static.DiSandboxManager).(sandbox.Manager)
	s.rlm = rlm
	return
}

func (s *session) Close() {
	logrus.
		WithField("addr", getAddr(s.conn)).
		Debug("websocket connection closed")
	sessionPool.Put(s)
}

func (s *session) Handler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		logrus.
			WithField("addr", getAddr(c)).
			Debug("new websocket connection")

		s.conn = c
		var (
			typ   int
			msg   []byte
			err   error
			nonce int
		)
		for {
			if typ, msg, err = c.ReadMessage(); err != nil {
				s.Close()
				break
			}
			if typ != websocket.TextMessage {
				s.SendError(models.ErrInvalidMessageType, 0)
			}
			go func() {
				if err, nonce = s.HandleOp(msg); err != nil {
					s.SendError(err, nonce)
				}
			}()
		}
	})
}

func (s *session) Send(v models.Event) (err error) {
	err = s.conn.WriteJSON(v)
	return
}

func (s *session) SendError(err error, nonce int) error {
	var data models.WsError
	if wsErr, ok := err.(models.WsError); ok {
		data = wsErr
	} else {
		data = models.WsError{500, err.Error()}
	}
	return s.Send(models.Event{
		Code:  models.EventError,
		Nonce: nonce,
		Data:  data,
	})
}

func (s *session) HandleOp(msg []byte) (err error, nonce int) {
	var op models.Operation
	if err = json.Unmarshal(msg, &op); err != nil {
		return
	}
	nonce = op.Nonce

	if !s.rlm.GetLimiter(s.conn, op.Op).Allow() {
		err = models.ErrRateLimited
		return
	}

	var event models.Event
	event.Nonce = op.Nonce

	switch op.Op {
	case models.OpPing:
		event.Code = models.EventPong
		event.Data = "Pong!"
		err = s.Send(event)
	case models.OpExec:
		var eop models.OperationExec
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = s.handleExec(eop)
		}
	case models.OpKill:
		var eop models.OperationKill
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = s.handleKill(eop)
		}
	default:
		err = models.ErrInvalidOpCode
	}

	return
}

func (s *session) handleExec(op models.OperationExec) (err error) {
	if op.Args.Code == "" {
		return models.ErrEmptyCode
	}

	cSpn := make(chan string, 1)
	cStdOut := make(chan []byte)
	cStdErr := make(chan []byte)
	cStop := make(chan bool, 1)

	var runId string

	go func() {
		for {
			select {
			case <-cStop:
				return
			case runId = <-cSpn:
				err = s.Send(models.Event{
					Code:  models.EventSpawn,
					Nonce: op.Nonce,
					Data: models.DataSpawn{
						DataRunId: models.DataRunId{
							RunId: runId,
						},
					},
				})
			case p := <-cStdOut:
				err = s.Send(models.Event{
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
				err = s.Send(models.Event{
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
				logrus.WithError(err).Error("Failed sending event")
				if err = s.SendError(err, op.Nonce); err != nil {
					logrus.WithError(err).Error("Failed sending error event")
					return
				}
			}
		}
	}()

	execTime := util.MeasureTime(func() {
		err = s.manager.RunInSandbox(&op.Args, cSpn, cStdOut, cStdErr, cStop)
	})

	if err != nil {
		cStop <- false
	}

	err = s.Send(models.Event{
		Code:  models.EventStop,
		Nonce: op.Nonce,
		Data: models.DataStop{
			DataRunId: models.DataRunId{
				RunId: runId,
			},
			ExecTimeMS: int(execTime.Milliseconds()),
		},
	})

	return
}

func (s *session) handleKill(op models.OperationKill) (err error) {
	ok, err := s.manager.KillAndCleanUp(op.Args.RunId)
	if err != nil {
		return
	}
	if !ok {
		err = models.ErrSandboxNotRunning
	}
	return
}
