package logger

import (
	"strings"

	"go.uber.org/zap/zapcore"
)

const (
	FormatConsole = "console"
	formatJSON    = "json"
)

type Level zapcore.Level

type Options struct {
	Level         Level
	Encoding      string
	DisableCaller bool
	CallerSkip    int
	Stacktrace    Level
	IsSave        bool
	FileConfig    *fileOptions
	Hooks         []func(zapcore.Entry) error
}

func DefaultOptions() *Options {
	opts := &Options{
		Encoding:   FormatConsole,
		Stacktrace: Level(zapcore.ErrorLevel),
		CallerSkip: 1,
	}
	opts.Level = Level(zapcore.DebugLevel)
	return opts
}

func (o *Options) Apply(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

// Option set the logger options.
type Option func(*Options)

// WithLevel setting the log level
func WithLevel(name string) Option {
	return func(o *Options) {
		ln := strings.ToLower(name)
		l, err := zapcore.ParseLevel(ln)
		if err != nil {
			panic(err)
		}
		o.Level = Level(l)
	}
}

// WithFormat set the output log format, console or json
func WithFormat(format string) Option {
	return func(o *Options) {
		if strings.ToLower(format) == formatJSON {
			o.Encoding = formatJSON
		}
	}
}

// WithCallerSkip setting the log caller
func WithCaller(caller bool, skip int) Option {
	return func(o *Options) {
		o.DisableCaller = caller
		o.CallerSkip = skip
	}
}

// WithStacktraceLevel setting the log stacktraceLevel
func WithStacktraceLevel(name string) Option {
	return func(o *Options) {
		l, err := zapcore.ParseLevel(strings.ToLower(name))
		if err != nil {
			panic(err)
		}
		o.Stacktrace = Level(l)
	}
}

// WithHooks set the log hooks
func WithHooks(hooks ...func(zapcore.Entry) error) Option {
	return func(o *Options) {
		o.Hooks = hooks
	}
}

/*****************************************************************************
 *
 * logger file options
 *
 *****************************************************************************/

// WithSave save log to file
func WithSave(isSave bool, opts ...FileOption) Option {
	return func(o *Options) {
		if isSave {
			o.IsSave = true
			fo := defaultFileOptions()
			fo.apply(opts...)
			o.FileConfig = fo
		}
	}
}

type fileOptions struct {
	Filename      string
	MaxSize       int
	MaxBackups    int
	MaxAge        int
	IsCompression bool
	isLocalTime   bool
}

func defaultFileOptions() *fileOptions {
	return &fileOptions{
		Filename:    "out.log",
		MaxSize:     50, // maximum file size (MB)
		MaxBackups:  20, // maximum number of old files
		MaxAge:      30, // maximum number of days for old documents
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
			f.Filename = filename
		}
	}
}

// WithFileMaxSize set maximum file size (MB)
func WithFileMaxSize(maxSize int) FileOption {
	return func(f *fileOptions) {
		if maxSize > 0 {
			f.MaxSize = maxSize
		}
	}
}

// WithFileMaxBackups set maximum number of old files
func WithFileMaxBackups(maxBackups int) FileOption {
	return func(f *fileOptions) {
		if f.MaxBackups > 0 {
			f.MaxBackups = maxBackups
		}
	}
}

// WithFileMaxAge set maximum number of days for old documents
func WithFileMaxAge(maxAge int) FileOption {
	return func(f *fileOptions) {
		if f.MaxAge > 0 {
			f.MaxAge = maxAge
		}
	}
}

// WithFileIsCompression set whether to compress log files
func WithFileIsCompression(isCompression bool) FileOption {
	return func(f *fileOptions) {
		f.IsCompression = isCompression
	}
}

// WithLocalTime set whether to use local time
func WithLocalTime(isLocalTime bool) FileOption {
	return func(f *fileOptions) {
		f.isLocalTime = isLocalTime
	}
}
