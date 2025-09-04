package logger4

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/vincent78/butil/config"
	logger "github.com/vincent78/butil/logger/logger2/logger"
	"github.com/vincent78/butil/utils/fileUtil"

	"go.uber.org/zap/zapcore"
)

func printInfo() {
	defer func() {
		recover()
	}()

	Debug("this is debug")
	Debugf("this is debugf %d", 2)
	Info("this is info")
	Infof("this is infof %d", 2)
	Warn("this is warn")
	Warnf("this is warnf %d", 2)
	Error("this is error")
	Errorf("this is errorf %d", 2)
	WithFields(Int("key", 2)).Info("this is info")
	//Fatal("this is fatal")
	//Fatalf("this is fatal %d", 2)
	_ = Sync()

	type people struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	p := &people{"Mr Zhang", 11}
	ps := []people{{"Mr Zhang", 11}, {"Mr Li", 12}}
	pMap := map[string]*people{"123": p, "456": p}
	Info("this is info object", Any("object1", p), Any("object2", ps), Any("object3", pMap)) // this sentence cannot be printed using debug

	Panic("this is panic")
}

func TestInit(t *testing.T) {
	type args struct {
		opts []Option
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "terminal console debug",
			args:    args{},
			wantErr: false,
		},
		{
			name: "terminal json info",
			args: args{[]Option{
				WithFormat("json"), WithLevel("info"),
			}},
			wantErr: false,
		},
		{
			name: "terminal json warn",
			args: args{[]Option{
				WithFormat("json"), WithLevel("warn"),
			}},
			wantErr: false,
		},
		{
			name: "with hooks info",
			args: args{[]Option{
				WithFormat("json"),
				WithLevel("info"),
				WithHooks(func(entry zapcore.Entry) error {
					if strings.Contains(entry.Message, "this is error") {
						fmt.Println("it contains error message")
					}
					return nil
				}),
			}},
			wantErr: false,
		},
		{
			name: "file json debug",
			args: args{[]Option{
				WithFormat("json"), WithLevel("unknown"),
				WithSave(
					true,
					WithFileName(os.TempDir()+"/testLog/my.log"),
					WithFileMaxSize(5),
					WithFileMaxBackups(5),
					WithFileMaxAge(10),
					WithFileIsCompression(true),
					WithLocalTime(true),
				),
			}},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Init(tt.args.opts...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Init() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			printInfo()
		})
	}

	time.Sleep(time.Second)
	_ = os.RemoveAll("my.log")
}

func BenchmarkString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("this is info", String("string", "hello golang"))
	}
}

func BenchmarkInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark type int", Int("int", i))
	}
}

func BenchmarkAny(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Info("benchmark type any", Any(fmt.Sprintf("object_%d", i), map[string]int{"Mr Zhang": 11}))
	}
}

func Test_getLevelSize(t *testing.T) {
	_ = getLevelSize(levelDebug)
	_ = getLevelSize(levelInfo)
	_ = getLevelSize(levelWarn)
	_ = getLevelSize(levelWarn)
	_ = getLevelSize(levelError)
	_ = getLevelSize("unknown")

	defaultLogger = nil
	_ = GetWithSkip(5)
	_ = Get()
}

func Test_InitTestConfigByFile(t *testing.T) {
	path := fileUtil.GetCurrentProjectPath()
	file := fileUtil.Join(path, "config", "test.yaml")
	cfg := &config.CmdConfig{}
	_ = config.Parse(file, cfg)
	dfl := cfg.Logger["test"]
	l, _ := Init(
		WithLevel(dfl.Level),
		WithFormat(dfl.Format),
		WithDisableCaller(dfl.DisableCaller),
		WithSave(
			dfl.IsSave,
			WithFileName(dfl.LogFileConfig.Filename),
			WithFileMaxSize(dfl.LogFileConfig.MaxSize),
			WithFileMaxBackups(dfl.LogFileConfig.MaxBackups),
			WithFileMaxAge(dfl.LogFileConfig.MaxAge),
			WithFileIsCompression(dfl.LogFileConfig.IsCompression),
		),
	)
	l.Info("this is info")
	l.Debug("this is debug")
}

func Test_InitConfigByFile(t *testing.T) {
	path := fileUtil.GetCurrentProjectPath()
	file := fileUtil.Join(path, "config", "test.yaml")
	cfgs := &config.CmdConfig{}
	_ = config.Parse(file, cfgs)
	InitLoggerByConfig(cfgs.Logger)
	SetDefaultLogger("default")
	logger.Info("this is info")
	logger.Debug("this is debug")
}
