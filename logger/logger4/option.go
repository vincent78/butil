package logger4

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	formatConsole = "console"
	formatJSON    = "json"
)

type Level zapcore.Level

type options struct {
	level         Level
	encoding      string
	disableCaller bool
	callerSkip    int
	stacktrace    Level

	isSave     bool
	fileConfig *fileOptions

	hooks []func(zapcore.Entry) error
}

func defaultOptions() *options {
	opts := &options{
		encoding:   formatConsole,
		stacktrace: Level(zapcore.ErrorLevel),
		callerSkip: 1,
	}
	opts.level = Level(zapcore.DebugLevel)
	return opts
}

func (o *options) apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the logger options.
type Option func(*options)

// WithLevel setting the log level
func WithLevel(name string) Option {
	return func(o *options) {
		ln := strings.ToLower(name)
		l, err := zapcore.ParseLevel(ln)
		if err != nil {
			panic(err)
		}
		o.level = Level(l)
	}
}

// WithFormat set the output log format, console or json
func WithFormat(format string) Option {
	return func(o *options) {
		if strings.ToLower(format) == formatJSON {
			o.encoding = formatJSON
		}
	}
}

// WithCallerSkip setting the log caller
func WithCaller(caller bool, skip int) Option {
	return func(o *options) {
		o.disableCaller = caller
		o.callerSkip = skip
	}
}

// WithStacktraceLevel setting the log stacktraceLevel
func WithStacktraceLevel(name string) Option {
	return func(o *options) {
		l, err := zapcore.ParseLevel(strings.ToLower(name))
		if err != nil {
			panic(err)
		}
		o.stacktrace = Level(l)
	}
}

// WithHooks set the log hooks
func WithHooks(hooks ...func(zapcore.Entry) error) Option {
	return func(o *options) {
		o.hooks = hooks
	}
}

/*****************************************************************************
 *
 * logger file options
 *
 *****************************************************************************/

// WithSave save log to file
func WithSave(isSave bool, opts ...FileOption) Option {
	return func(o *options) {
		if isSave {
			o.isSave = true
			fo := defaultFileOptions()
			fo.apply(opts...)
			o.fileConfig = fo
		}
	}
}

type fileOptions struct {
	filename      string
	maxSize       int
	maxBackups    int
	maxAge        int
	isCompression bool
	isLocalTime   bool
}

func defaultFileOptions() *fileOptions {
	return &fileOptions{
		filename:    "out.log",
		maxSize:     50, // maximum file size (MB)
		maxBackups:  20, // maximum number of old files
		maxAge:      30, // maximum number of days for old documents
		isLocalTime: true,
	}
}

func (o *fileOptions) apply(opts ...FileOption) {
	for _, opt := range opts {
		opt(o)
	}
}

// FileOption set the file options.
type FileOption func(*fileOptions)

// WithFileName set log filename
func WithFileName(filename string) FileOption {
	return func(f *fileOptions) {
		if filename != "" {
			f.filename = filename
		}
	}
}

// WithFileMaxSize set maximum file size (MB)
func WithFileMaxSize(maxSize int) FileOption {
	return func(f *fileOptions) {
		if maxSize > 0 {
			f.maxSize = maxSize
		}
	}
}

// WithFileMaxBackups set maximum number of old files
func WithFileMaxBackups(maxBackups int) FileOption {
	return func(f *fileOptions) {
		if f.maxBackups > 0 {
			f.maxBackups = maxBackups
		}
	}
}

// WithFileMaxAge set maximum number of days for old documents
func WithFileMaxAge(maxAge int) FileOption {
	return func(f *fileOptions) {
		if f.maxAge > 0 {
			f.maxAge = maxAge
		}
	}
}

// WithFileIsCompression set whether to compress log files
func WithFileIsCompression(isCompression bool) FileOption {
	return func(f *fileOptions) {
		f.isCompression = isCompression
	}
}

// WithLocalTime set whether to use local time
func WithLocalTime(isLocalTime bool) FileOption {
	return func(f *fileOptions) {
		f.isLocalTime = isLocalTime
	}
}
