package datastructure

import (
	"fmt"
	"testing"
)

func TestFnv32Key1(t *testing.T) {
	u := fnv32("test-key-1")
	index := uint(u) % uint(SHARD_COUNT)
	fmt.Println(index)
}

func TestFnv32Key2(t *testing.T) {
	u := fnv32("test-key-2")
	index := uint(u) % uint(SHARD_COUNT)
	fmt.Println(index)
}

func TestFnv32Key3(t *testing.T) {
	u := fnv32("test-key")
	index := uint(u) % uint(SHARD_COUNT)
	fmt.Println(index)
}

func TestConcurrentMap(t *testing.T) {
	concurrentMap := New[string]()
	concurrentMap.Set("test-key-1", "test-key-1")
	fmt.Println(concurrentMap.Get("test-key-1"))
	concurrentMap.Remove("test-key-1")
	fmt.Println(concurrentMap.Get("test-key-1"))
}
