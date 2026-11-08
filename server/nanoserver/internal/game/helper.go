package game

import (
	"runtime"
	"strings"

	"nanoserver/pkg/errutil"

	"github.com/lonng/nano/session"
)

func playerWithSession(s *session.Session) (*Player, error) {
	p, ok := s.Value(kCurPlayer).(*Player)
	if !ok {
		return nil, errutil.ErrPlayerNotFound
	}
	return p, nil
}

func stack() string {
	buf := make([]byte, 10000)
	n := runtime.Stack(buf, false)
	buf = buf[:n]

	s := string(buf)

	// skip nano frames lines
	const skip = 7
	count := 0
	index := strings.IndexFunc(s, func(c rune) bool {
		if c != '\n' {
			return false
		}
		count++
		return count == skip
	})
	return s[index+1:]
}
