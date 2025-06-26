package cache

import (
	"slices"
	"sync"
	"time"
)

type MemoryCache struct {
	syncBlockNumber       int64 //  当前同步的block编号
	syncBlockNumberLock   sync.RWMutex
	energyAddressPool     map[string][]string // 能量地址的缓冲池 [能量地址][能量发放地址]
	energyAddressPoolLock sync.RWMutex
	addressActiveTime     map[string]int64 // 地址激活时间
	addressActiveTimeLock sync.RWMutex
	extend                map[string]MemeryCacheValue // 通用的缓存
	extendLock            sync.RWMutex
}

type MemeryCacheValue struct {
	Value     interface{}
	Timestamp int64
}

var (
	localMemoryCache *MemoryCache
	initOnce         sync.Once
)

func initMemoryCache() {
	localMemoryCache = &MemoryCache{}
	localMemoryCache.syncBlockNumber = 0
	localMemoryCache.syncBlockNumberLock = sync.RWMutex{}
	localMemoryCache.energyAddressPool = make(map[string][]string)
	localMemoryCache.energyAddressPoolLock = sync.RWMutex{}
	localMemoryCache.addressActiveTime = make(map[string]int64)
	localMemoryCache.addressActiveTimeLock = sync.RWMutex{}
	localMemoryCache.extend = make(map[string]MemeryCacheValue)
	localMemoryCache.extendLock = sync.RWMutex{}
}

func GetLocalMemoryCache() *MemoryCache {
	if localMemoryCache == nil {
		initOnce.Do(func() {
			initMemoryCache()
		})
	}
	return localMemoryCache
}

func (cache *MemoryCache) Get(key string) (interface{}, bool) {
	cache.extendLock.RLock()
	defer cache.extendLock.RUnlock()
	if cache.extend != nil {
		if val, ok := cache.extend[key]; ok {
			if val.Timestamp > (time.Now().UTC().UnixNano() / 1e6) {
				return val.Value, true
			} else {
				delete(cache.extend, key)
			}
			return val, true
		}
	}
	return nil, false
}

func getNowMillSecond() int64 {
	return time.Now().UnixNano() / 1e6
}

func (cache *MemoryCache) Put(key string, val interface{}) error {
	// 默认30分钟后失效
	return cache.PutWithTimeout(key, val, 30*60*1000*getNowMillSecond())
}

// PutWithTimeout
// 参数  out: 毫秒数
func (cache *MemoryCache) PutWithTimeout(key string, val interface{}, out int64) error {
	cache.extend[key] = MemeryCacheValue{
		Value:     val,
		Timestamp: out + getNowMillSecond(),
	}
	return nil
}

func (cache *MemoryCache) PutEnergyAddr(addr, otherAddr string) {
	if cache.energyAddressPool != nil {
		if p, ok := cache.energyAddressPool[addr]; ok {
			if !slices.Contains(p, otherAddr) {
				if len(otherAddr) == 0 {
					cache.energyAddressPool[addr] = make([]string, 0)
				} else {
					cache.energyAddressPool[addr] = append(p, otherAddr)
				}
			}
		} else {
			cache.energyAddressPool[addr] = []string{otherAddr}
		}
	}
}

func (cache *MemoryCache) ExistEnergyAddr(addr string) bool {
	_, ok := cache.energyAddressPool[addr]
	return ok
}

func (cache *MemoryCache) IncrBlockNumber() int64 {
	cache.syncBlockNumber = cache.syncBlockNumber + 1
	return cache.syncBlockNumber
}

func (cache *MemoryCache) GetBlockNumber() int64 {
	return cache.syncBlockNumber
}

func (cache *MemoryCache) SetBlockNumber(blockNumber int64) {
	cache.syncBlockNumberLock.Lock()
	cache.syncBlockNumber = blockNumber
	cache.syncBlockNumberLock.Unlock()
}

func (cache *MemoryCache) GetAddressActiveTime(address string) int64 {
	return cache.addressActiveTime[address]
}

func (cache *MemoryCache) PutAddressActiveTime(address string, ts int64) {
	cache.addressActiveTimeLock.Lock()
	cache.addressActiveTime[address] = ts
	cache.addressActiveTimeLock.Unlock()
}
