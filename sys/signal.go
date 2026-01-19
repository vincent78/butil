package sys

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Signal struct {
	sigChan    chan os.Signal
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewSignal() *Signal {
	s := &Signal{
		sigChan:    make(chan os.Signal),
		ctx:        nil,
		cancelFunc: nil,
	}
	s.ctx, s.cancelFunc = context.WithCancel(context.Background())

	return s
}

func (s *Signal) notify() {
	signal.Notify(s.sigChan, syscall.SIGINT, syscall.SIGTERM)
}

func (s *Signal) Waiter() error {
	s.notify()

	for sig := range s.sigChan {
		switch sig {
		case syscall.SIGINT:
			//fmt.Println("control signal int:", s)
			return nil
		case syscall.SIGTERM:
			//fmt.Println("control signal term:", s)
			return nil
		}
	}

	return nil
}

func (s *Signal) Cannel(duration ...time.Duration) {
	s.cancelFunc()

	if len(duration) > 0 {
		time.Sleep(duration[0])
	}
}

func (s *Signal) GetCtx() context.Context {
	return s.ctx
}
