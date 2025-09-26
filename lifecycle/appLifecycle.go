package lifecycle

import (
	"context"

	"github.com/vincent78/butil/logger/logger1"
	"github.com/vincent78/butil/sys"
)

type AppLifecycleInitModel struct {
	Finished chan bool
}

func init() {
	go func() {
		for {
			select {
			case m := <-AppLifecycleInit:
				AppInit(m)
			case <-AppLifecyclePrepared:
				AppPrepare()
			case ctx := <-AppLifecycleActived:
				AppPause(ctx)
			case ctx := <-AppLifecyclePaused:
				AppActive(ctx)
			case ctx := <-AppLifecycleDestroy:
				AppDestroy(ctx)
			}
		}
	}()
}

func AppInit(model *AppLifecycleInitModel) {
	logger1.Info("---- appInit ----")
	AppBaseInit()
	sys.ManualGC()
	close(model.Finished)
}

func AppPrepare() {
	logger1.Info("---- appPrepare ----")
	if AppBasePrepared() {
		sys.ManualGC()
		logger1.Info("App Has Prepared")
		ctx := context.Background()
		ctx = context.WithValue(ctx, "result", true)
		AppLifecycleWorkBegin <- ctx
	}
}

func AppPause(ctx context.Context) {
	logger1.Info("---- appPause ----")
	AppBasePause()
	sys.ManualGC()
}

func AppActive(ctx context.Context) {
	logger1.Info("---- appActive ----")
	AppBaseActive()
	sys.ManualGC()
}

func AppDestroy(ctx context.Context) {
	logger1.Info("---- appDestroy ----")
	AppBaseDestory()
}
