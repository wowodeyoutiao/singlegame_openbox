package async

import (
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

// func pcall(fn func()) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logrus.Errorf("async/pcall: Error=%v\nStack:\n%s", err, debug.Stack())
// 		}
// 	}()

// 	fn()
// }

// func Run(fn func()) {
// 	go pcall(fn)
// }

func Run(fn func()) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				logrus.Errorf("goroutine panic at %s: %v\nStack:\n%s", time.Now().Format("2006-01-02 15:04:05"), r, debug.Stack())
			}
		}()
		fn()
	}()
}
