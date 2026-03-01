package registry

import (
	reg "github.com/vincent78/butil/bus/core/registry"
)

var (
	handlerReg reg.Registry[NewHandler] = new(handlerRegistry)
	loggerReg  reg.Registry[NewLogger]  = new(loggerRegistry)
	testReg    reg.Registry[NewTest]    = new(testRegistry)
)

func HandlerRegistry() reg.Registry[NewHandler] {
	return handlerReg
}

func TestRegistry() reg.Registry[NewTest] {
	return testReg
}

func LoggerRegistry() reg.Registry[NewLogger] {
	return loggerReg
}
