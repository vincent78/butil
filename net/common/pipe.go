package common

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/go-gost/core/common/bufpool"
)

const (
	// tcpWaitTimeout implements a TCP half-close timeout.
	tcpWaitTimeout = 10 * time.Second
	bufferSize     = 64 * 1024
)

var (
	ErrUnsupported = errors.New("unsupported")
)

type CloseRead interface {
	CloseRead() error
}

type CloseWrite interface {
	CloseWrite() error
}

type setReadDeadline interface {
	SetReadDeadline(t time.Time) error
}

func Pipe(ctx context.Context, rw1, rw2 io.ReadWriteCloser) error {
	wg := sync.WaitGroup{}
	wg.Add(2)

	ch := make(chan error, 2)

	go func() {
		defer wg.Done()
		if err := pipeBuffer(rw1, rw2, bufferSize/2); err != nil {
			ch <- err
		}
	}()
	go func() {
		defer wg.Done()
		if err := pipeBuffer(rw2, rw1, bufferSize/2); err != nil {
			ch <- err
		}
	}()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		return nil
	}

	select {
	case err := <-ch:
		return err
	default:
	}

	return nil
}

func pipeBuffer(dst io.ReadWriteCloser, src io.ReadWriteCloser, bufferSize int) error {
	buf := bufpool.Get(bufferSize)
	defer bufpool.Put(buf)

	_, err := io.CopyBuffer(dst, src, buf)

	// Do the upload/download side TCP half-close.
	if cr, ok := src.(CloseRead); ok {
		cr.CloseRead()
	}

	if cw, ok := dst.(CloseWrite); ok {
		if e := cw.CloseWrite(); e == ErrUnsupported {
			dst.Close()
		} else {
			// Set TCP half-close timeout.
			SetReadDeadline(dst, time.Now().Add(tcpWaitTimeout))
		}
	} else {
		dst.Close()
	}

	return err
}

func SetReadDeadline(rw io.ReadWriter, t time.Time) error {
	if v, _ := rw.(setReadDeadline); v != nil {
		return v.SetReadDeadline(t)
	}
	return nil
}
