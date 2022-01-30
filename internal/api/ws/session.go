package ws

import (
	"encoding/json"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/ranna-go/ranna/internal/sandbox"
	"github.com/ranna-go/ranna/internal/static"
	"github.com/ranna-go/ranna/internal/util"
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
}

func newSession(ctn di.Container) (s *session) {
	s = sessionPool.Get().(*session)
	s.conn = nil
	s.manager = ctn.Get(static.DiSandboxManager).(sandbox.Manager)
	return
}

func (s *session) Close() {
	logrus.
		WithField("addr", s.conn.RemoteAddr().String()).
		Debug("websocket connection closed")
	sessionPool.Put(s)
}

func (s *session) Handdler() fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		logrus.
			WithField("addr", c.RemoteAddr().String()).
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
				s.SendError(ErrInvalidMessageType, 0)
			}
			go func() {
				if err, nonce = s.HandleOp(msg); err != nil {
					s.SendError(err, nonce)
				}
			}()
		}
	})
}

func (s *session) Send(v Event) (err error) {
	err = s.conn.WriteJSON(v)
	return
}

func (s *session) SendError(err error, nonce int) error {
	return s.Send(Event{
		Code:  EventError,
		Nonce: nonce,
		Data:  err.Error(),
	})
}

func (s *session) HandleOp(msg []byte) (err error, nonce int) {
	var op Operation
	if err = json.Unmarshal(msg, &op); err != nil {
		return
	}

	var event Event
	event.Nonce = op.Nonce
	nonce = op.Nonce

	switch op.Op {
	case OpPing:
		event.Code = EventPong
		event.Data = "Pong!"
		err = s.Send(event)
	case OpExec:
		var eop OperationExec
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = s.handleExec(eop)
		}
	case OpKill:
		var eop OperationKill
		err = json.Unmarshal(msg, &eop)
		if err == nil {
			err = s.handleKill(eop)
		}
	default:
		err = ErrInvalidOpCode
	}

	return
}

func (s *session) handleExec(op OperationExec) (err error) {
	if op.Args.Code == "" {
		return ErrEmptyCode
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
				s.Send(Event{
					Code:  EventSpawn,
					Nonce: op.Nonce,
					Data: DataSpawn{
						DataRunId: DataRunId{
							RunId: runId,
						},
					},
				})
			case p := <-cStdOut:
				s.Send(Event{
					Code:  EventLog,
					Nonce: op.Nonce,
					Data: DataLog{
						DataRunId: DataRunId{
							RunId: runId,
						},
						StdOut: string(p),
					},
				})
			case p := <-cStdErr:
				s.Send(Event{
					Code:  EventLog,
					Nonce: op.Nonce,
					Data: DataLog{
						DataRunId: DataRunId{
							RunId: runId,
						},
						StdErr: string(p),
					},
				})
			}
		}
	}()

	execTime := util.MeasureTime(func() {
		err = s.manager.RunInSandbox(&op.Args, cSpn, cStdOut, cStdErr, cStop)
	})

	if err != nil {
		cStop <- false
	}

	err = s.Send(Event{
		Code:  EventStop,
		Nonce: op.Nonce,
		Data: DataStop{
			DataRunId: DataRunId{
				RunId: runId,
			},
			ExecTimeMS: int(execTime.Milliseconds()),
		},
	})

	return
}

func (s *session) handleKill(op OperationKill) (err error) {
	ok, err := s.manager.KillAndCleanUp(op.Args.RunId)
	if err != nil {
		return
	}
	if !ok {
		err = ErrSandboxNotRunning
	}
	return
}
