package sys

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"

	"github.com/vincent78/butil/logger/logger4"
)

var cpuprofile = flag.String("cpuprofile", "", "Where to write CPU profile")
var memprofile = flag.String("memprofile", "", "Where to write MEM profile")
var meminterval = flag.Int("meminterval", 10, "the interval of get mem info ,default is 10s")

// RunCPUProfile Run starts up stuff at the beginning of a main function, and returns a
// function to defer until the function completes.  It should be used like this:
//
//	func main() {
//	  defer sys.RunCPUProfile()()
//	  ... stuff ...
//	}
func RunCPUProfile() func() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("could not open cpu profile file %q", *cpuprofile)
		}
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			return nil
		}
		return func() {
			pprof.StopCPUProfile()
			_ = f.Close()
		}
	}
	return func() {}
}

func RunMemProfile() func() {
	flag.Parse()
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatalf("could not open mem profile file %q", *memprofile)
		}
		runtime.GC()
		err = pprof.WriteHeapProfile(f)
		if err != nil {
			return nil
		}
		tmi := time.Duration(*meminterval)
		ticker := time.NewTicker(tmi * time.Second)
		ch := make(chan bool)
		go func(ticker *time.Ticker, f *os.File) {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					_ = pprof.WriteHeapProfile(f)
				case stop := <-ch:
					if stop {
						return
					}
				}
			}
		}(ticker, f)
		return func() {
			ticker.Stop()
			_ = pprof.WriteHeapProfile(f)
			ch <- true
			_ = f.Close()
		}
	}
	return func() {}
}

// 计算运行f前后内存信息
func MemoryShow(f func(), log logger4.Logger) {
	// 启用内存分配统计
	runtime.MemProfileRate = 1 // 每分配1次就记录一次
	runtime.GC()               // 强制进行垃圾回收

	// 获取内存统计信息
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	log.Info("memory show: begin",
		logger4.Any("totalAlloc", mem.TotalAlloc),
		logger4.Any("heapInuse", mem.HeapInuse),
		logger4.Any("stackInuse", mem.StackInuse),
	)

	if f != nil {
		f()
	}

	runtime.ReadMemStats(&mem)
	log.Info("memory show: end",
		logger4.Any("totalAlloc", mem.TotalAlloc),
		logger4.Any("heapInuse", mem.HeapInuse),
		logger4.Any("stackInuse", mem.StackInuse),
	)
}
