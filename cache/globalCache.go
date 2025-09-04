package cache

import (
	"os"
	"sync"
	"time"

	"github.com/muesli/cache2go"
)

var once sync.Once

var defaultCachePool *cache2go.CacheTable
var defaultDuration = time.Second * 60 * 3

const ENV_GO_CACHE_POOL_NAME = "env_go_default_cache_pool_name"
const ENV_GO_CACHE_DEFAULT_DURATION = "env_go_default_cache_duration"

func init() {
	n, r := os.LookupEnv(ENV_GO_CACHE_POOL_NAME)
	if !r {
		n = "DEFAULT_CAHCE_POOL_NAME"
	}
	defaultCachePool = cache2go.Cache(n)

	d, r := os.LookupEnv(ENV_GO_CACHE_DEFAULT_DURATION)
	if !r {
		defaultDuration, _ = time.ParseDuration(d)
	}

}

// 保存数据，永远不过期
func Save(key string, value interface{}) {
	SaveWithLifeSpan(key, value, 0)
}

// 保存数据，指定时间后过期
func SaveWithLifeSpan(key string, value interface{}, lifeDuration time.Duration) {
	defaultCachePool.Add(key, lifeDuration, value)
	defaultCachePool.Flush()
}

func Delete(key string) error {
	_, err := defaultCachePool.Delete(key)
	if err == nil {
		defaultCachePool.Flush()
	}
	return err
}

func Get(key string) interface{} {
	val, err := defaultCachePool.Value(key)
	if err == nil {
		return val
	} else {
		return nil
	}
}

func Exist(key string) bool {
	return defaultCachePool.Exists(key)
}
