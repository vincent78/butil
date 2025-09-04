package sys

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/vincent78/butil/logger/logger1"
)

func BlockBySignal() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		s := <-c
		logger1.Info("discovery get a signal %s", s.String())
		switch s {
		case os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM:
			logger1.Info("application quit !!!")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

func BlockBySignalAndCancel(cancel context.CancelFunc, duration time.Duration) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM)
	for {
		s := <-c
		logger1.Info("discovery get a signal %s", s.String())
		switch s {
		case os.Interrupt, os.Kill, syscall.SIGQUIT, syscall.SIGTERM:
			if cancel != nil {
				cancel()
			}
			time.Sleep(duration)
			logger1.Info("application quit !!!")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

type actionFunc func() error

// 守护方法
func runLoop(action actionFunc) error {
	errCh := make(chan error)
	// run the proxy in goroutine
	go func() {
		err := action()
		errCh <- err
	}()

	for {
		err := <-errCh
		if err != nil {
			return err
		}
	}
}

func SignalContext() context.Context {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-sigs
		logger1.Info("--- receive system signal ---")
		cancel()
		time.Sleep(time.Second * 2)
		os.Exit(1)
	}()
	return ctx
}
