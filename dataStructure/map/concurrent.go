package datastructure

// https://github.com/orcaman/concurrent-map
import "sync"

var SHARD_COUNT = 32

type ConcurrentMap[V any] []*ConcurrentMapShared[V]

type ConcurrentMapShared[V any] struct {
	Items        map[string]V
	sync.RWMutex // 读写锁
}

func New[V any]() ConcurrentMap[V] {
	m := make(ConcurrentMap[V], SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		m[i] = &ConcurrentMapShared[V]{
			Items: make(map[string]V),
		}
	}
	return m
}

func (m ConcurrentMap[V]) GetShard(key string) *ConcurrentMapShared[V] {
	return m[uint(fnv32(key))%uint(SHARD_COUNT)]
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	keyLength := len(key)
	for i := 0; i < keyLength; i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}

func (m ConcurrentMap[V]) Set(key string, value V) {
	shard := m.GetShard(key)
	shard.Lock()
	shard.Items[key] = value
	shard.Unlock()
}

func (m ConcurrentMap[V]) Get(key string) (V, bool) {
	shard := m.GetShard(key)
	shard.RLock()
	val, ok := shard.Items[key]
	shard.RUnlock()
	return val, ok
}

func (m ConcurrentMap[V]) Remove(key string) {
	shard := m.GetShard(key)
	shard.Lock()
	delete(shard.Items, key)
	shard.Unlock()
}
