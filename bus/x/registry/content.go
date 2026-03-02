package registry

import (
	"github.com/vincent78/butil/bus/core/logger"
	reg "github.com/vincent78/butil/bus/core/registry"
)

var (
	handlerReg reg.Registry[NewHandler]           = new(handlerRegistry)
	loggerReg  reg.Registry[logger.ILoggerMethod] = new(loggerRegistry)
	testReg    reg.Registry[NewTest]              = new(testRegistry)
)

func HandlerRegistry() reg.Registry[NewHandler] {
	return handlerReg
}

func TestRegistry() reg.Registry[NewTest] {
	return testReg
}

func LoggerRegistry() reg.Registry[logger.ILoggerMethod] {
	return loggerReg
}
