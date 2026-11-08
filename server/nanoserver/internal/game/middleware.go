package game

import (
	"github.com/lonng/nano/pipeline"
	"github.com/lonng/nano/session"
)

type GameEventProcess struct{}

func NewGameEventProcess() *GameEventProcess {
	return &GameEventProcess{}
}

func (gep *GameEventProcess) OutBound(s *session.Session, msg *pipeline.Message) error {
	// logger.Debugf("GameEventProcess OutBound msg: %+v, session: %+v", msg, s)
	p, err := playerWithSession(s)
	if err != nil {
		logger.Errorf("GameEventProcess OutBound playerWithSession error: %v", err)
		return err
	}

	p.eventMgr.ProcessEvents()

	return nil
}

func (gep *GameEventProcess) InBound(s *session.Session, msg *pipeline.Message) error {
	// logger.Debugf("GameEventProcess InBound msg: %+v, session: %+v", msg, s)
	return nil
}
