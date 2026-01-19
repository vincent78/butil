package thread

import (
	"sync/atomic"
	"time"
)

type functinType func(any) error

// Worker is the actual executor who runs the tasks,
// it starts a goroutine that accepts tasks and
// performs function calls.
type Worker struct {
	// pool who owns this routine.
	pool *Pool
	// task is a job should be done.
	task chan functinType
	// recycleTime will be update when putting a routine back into queue.
	recycleTime time.Time

	input chan any
}

// run starts a goroutine to repeat the process
// that performs the function calls.
func (w *Worker) run() {

	go func() {
		//监听任务列表，一旦有任务立马取出运行
		count := 1
		var input any
		var f functinType
		for count <= 2 {
			select {
			case str_temp, ok := <-w.input:
				if !ok {
					return
				}
				count++
				input = str_temp
			case f_temp, ok := <-w.task:
				if !ok {
					//如果接收到关闭
					atomic.AddInt32(&w.pool.running, -1)
					close(w.task)
					return
				}
				count++
				f = f_temp
			}
		}
		err := f(input)
		if err != nil {
			//fmt.Println("执行任务失败")
		}
		//回收复用
		w.pool.putWorker(w)
		return
	}()
}

// stop this routine.
func (w *Worker) stop() {
	w.sendTask(nil)
	close(w.input)
}

// sendTask sends a task to this routine.
func (w *Worker) sendTask(task functinType) {
	w.task <- task
}

func (w *Worker) sendarg(input any) {
	w.input <- input
}
