package sys

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"
	"time"
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
		pprof.WriteHeapProfile(f)
		return func() {
			pprof.StopCPUProfile()
			f.Close()
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
		pprof.WriteHeapProfile(f)
		tmi := time.Duration(*meminterval)
		ticker := time.NewTicker(tmi * time.Second)
		ch := make(chan bool)
		go func(ticker *time.Ticker, f *os.File) {
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					pprof.WriteHeapProfile(f)
				case stop := <-ch:
					if stop {
						return
					}
				}
			}
		}(ticker, f)
		return func() {
			ticker.Stop()
			pprof.WriteHeapProfile(f)
			ch <- true
			f.Close()
		}
	}
	return func() {}
}
