package manager

type Queue struct {
	PoolSize int
	PoolChan chan any
}

func NewQueue(size int) *Queue {
	return &Queue{
		PoolSize: size,
		PoolChan: make(chan any, size),
	}
}

func (tq *Queue) Init(size int) *Queue {
	tq.PoolSize = size
	tq.PoolChan = make(chan any, size)
	return tq
}

func (tq *Queue) Push(i any) bool {
	if len(tq.PoolChan) == tq.PoolSize {
		return false
	}
	tq.PoolChan <- i
	return true
}

func (tq *Queue) PushSlice(s []any) {
	for _, i := range s {
		tq.Push(i)
	}
}

func (tq *Queue) Pull() any {
	return <-tq.PoolChan
}

// 二次使用Queue实例时，根据容量需求进行高效转换
func (tq *Queue) Exchange(num int) (add int) {
	last := len(tq.PoolChan)

	if last >= num {
		add = 0
		return
	}

	if tq.PoolSize < num {
		var pool []any
		for i := 0; i < last; i++ {
			pool = append(pool, <-tq.PoolChan)
		}
		// 重新定义、赋值
		tq.Init(num).PushSlice(pool)
	}

	add = num - last
	return
}
