package work

import (
	"context"
	"github.com/vincent78/butil/model"
)

type WorkStatusType int

const (
	Created  WorkStatusType = iota // 已创建
	Prepared                       // 已就绪
	Running                        // 运行中
	Idle                           // 空闲中
	Fault                          // 故障中
)

type Worker struct {
	Id       string               // 唯一标识
	Context  context.Context      // worker的上下文
	JobCtx   context.Context      // 当前工作的上下文
	Status   WorkStatusType       // 当前状态
	Desc     string               // 描述
	JobCh    chan context.Context // 工作队列（输入）
	OutCh    chan string          // 运行情况（日志）队列 （输出）
	ResultCh chan model.RespModel // 运行结果（输出）
}

func NewWorker(ctx context.Context, id string) *Worker {

	w := &Worker{
		Id:       id,
		Context:  ctx,
		JobCtx:   nil,
		Status:   Created,
		JobCh:    make(chan context.Context),
		OutCh:    make(chan string),
		ResultCh: make(chan model.RespModel),
	}

	go func() {
		for {
			select {
			case jCtx := <-w.JobCh:
				w.doWork(jCtx)
				//case <-w.JobCh
			}
		}
	}()
	return w
}

func (w *Worker) doWork(jCtx context.Context) {
}

func (w *Worker) CancelWork() {

}

func (w *Worker) GetStatus() WorkStatusType {
	return w.Status
}
