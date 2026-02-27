package thread

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

var (
	poolWorksNonblocking     *ants.Pool
	poolWorksBlocking        *ants.Pool
	poolOnceNonblocking      sync.Once
	poolOnceBlocking         sync.Once
	poolWorksNonblockingSize = 100
)

// GetPoolWorks 获取工作池（延迟初始化）
// nonblocking=true：非阻塞提交；false：阻塞提交
func GetPoolWorks(nonblocking bool) *ants.Pool {
	if nonblocking {
		poolOnceNonblocking.Do(func() {
			var err error
			poolWorksNonblocking, err = ants.NewPool(poolWorksNonblockingSize,
				ants.WithNonblocking(true),
				ants.WithPreAlloc(true),
				ants.WithPanicHandler(func(p any) {
					fmt.Printf("ants panic: %v\n", p)
				}),
			)
			if err != nil {
				panic(fmt.Sprintf("failed to create ants pool: %v", err))
			}
		})
		return poolWorksNonblocking
	}

	poolOnceBlocking.Do(func() {
		var err error
		poolWorksBlocking, err = ants.NewPool(poolWorksNonblockingSize,
			ants.WithPreAlloc(true),
			ants.WithPanicHandler(func(p any) {
				fmt.Printf("ants panic: %v\n", p)
			}),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create blocking ants pool: %v", err))
		}
	})
	return poolWorksBlocking
}
