package lifecycle

import (
	"gitee.com/vincent78/gcutil/logger"
	"gitee.com/vincent78/gcutil/sys"
	"github.com/urfave/cli/v2"
)

type AppLifecycleInitModel struct {
	Ctx      *cli.Context
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
				AppPause(ctx.(*cli.Context))
			case ctx := <-AppLifecyclePaused:
				AppActive(ctx.(*cli.Context))
			case ctx := <-AppLifecycleDestroy:
				AppDestroy(ctx.(*cli.Context))
			}
		}
	}()
}

func AppInit(model *AppLifecycleInitModel) {
	logger.Info("---- appInit ----")
	AppBaseInit()
	sys.ManualGC()
	close(model.Finished)
}

func AppPrepare() {
	logger.Info("---- appPrepare ----")
	if AppBasePrepared() {
		sys.ManualGC()
		logger.Info("App Has Prepared")
		AppLifecycleWorkBegin <- true
	}
}

func AppPause(ctx *cli.Context) {
	logger.Info("---- appPause ----")
	AppBasePause()
	sys.ManualGC()
}

func AppActive(ctx *cli.Context) {
	logger.Info("---- appActive ----")
	AppBaseActive()
	sys.ManualGC()
}

func AppDestroy(ctx *cli.Context) {
	logger.Info("---- appDestroy ----")
	AppBaseDestory()
}
