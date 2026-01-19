package concurrentUtil

import "sync"

func SyncMapSize(m *sync.Map) (size int) {
	m.Range(func(key, value any) bool {
		size++
		return true
	})
	return size
}
