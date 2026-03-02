package bus

import (
	"context"
	"testing"

	"github.com/vincent78/butil/bus/core/handler"
	"github.com/vincent78/butil/bus/core/test"
	"github.com/vincent78/butil/bus/x/metadata"
	"github.com/vincent78/butil/bus/x/registry"
	_ "github.com/vincent78/butil/bus/x/test"
	"github.com/vincent78/butil/utils/timeUtil"
)

/*
	这是一个非常经典的 “延迟加载（Lazy Loading）” 或 “插件/驱动注册模式”。在 Go 语言（尤其是开发框架或驱动库，如 database/sql）中非常常见。
	简单来说，你在 init 函数里做的不是“调用函数”，而是**“存入说明书”**。只有当程序真正需要这个对象时，它才会照着说明书去“生产”。
	这个过程通常分为三个阶段：注册期、存储期、激活期。

	阶段 A：注册期 (Registration)
	在 Go 程序启动时，所有的 init() 函数会先于 main() 执行。
	动作：你将 NewXXX 这个函数本身（作为变量或接口）存入 sync.Map。
	注意：此时 NewXXX 并没有被执行，它只是像一个“函数指针”一样被存在了内存里。

	阶段 B：存储期 (Storage)
	容器：sync.Map。
	内容：Key 是一个字符串（比如 "mysql" 或 "redis"），Value 是 NewXXX 函数。
	状态：此时系统开销极小，因为你只是存了一个函数的地址。

	阶段 C：激活期 (Loading/Factory)
	当你的业务代码调用类似 Get("XXX") 的方法时：
	程序从 sync.Map 中根据 Key 找到对应的 NewXXX 函数。
	此时才真正**调用（Invoke）**该函数：instance := creator()。
	返回生成的实例。
*/

func Test_RegisterTest(t *testing.T) {
	// 必须在import中添加  _ "github.com/vincent78/butil/bus/x/test"
	// 这样才能让init方法被调用。

	if f := registry.TestRegistry().Get("test"); f != nil {
		//执行注册的方法（注：这里才是真正初始化的地方）
		md1 := metadata.NewMetadata(map[string]any{
			"test": "test",
			"now":  timeUtil.NowUtcStr(),
		})
		obj := f(test.WithMetadata(md1))
		md2 := metadata.NewMetadata(map[string]any{
			"xxx": "xxx",
		})
		_ = obj.Init(md2)
		_ = obj.Write(context.Background(), "test")
	}
}

func Test_RegisterHttpHandler(t *testing.T) {
	var opts []handler.Option
	if f := registry.HandlerRegistry().Get("http"); f != nil {
		_ = f(opts...)
	}
}

func Test_RegisterLogger(t *testing.T) {

}
