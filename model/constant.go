package model

const (
	FAILURE = 9999
	SUCCESS = 0
)

const (
	availability   = true
	unavailability = false
)

//程序的模式

const (
	CLIENT = iota
	MASTER
	SLAVER
	ANALYSE
)

// RunStatus 运行状态
type RunStatus int

const (
	// RunStatusStopped 停止
	RunStatusStopped RunStatus = 0
	// RunStatusPrepare 准备中
	RunStatusPrepare RunStatus = 1
	// RunStatusRunning 当前运行中
	RunStatusRunning RunStatus = 2
	// RunStatusIdle 当前空闲中
	RunStatusIdle RunStatus = 9
)

var (
	Identify  = ""
	Version   = ""
	GitCommit = ""
	GitTag    = ""
	BuildDate = ""
)

const (
	ErrorCodeConsoleWSClient = 1000 + iota
)

var ErrorConsoleParamsFormat = func(msg string) *ErrorModel {
	return NewErrModelByStr(ErrorCodeConsoleWSClient, "the params is not right: [%v]", msg)
}
