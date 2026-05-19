package _map

import "iter"

// Map 定义哈希映射的完整行为，Go标准命名
type Map[K comparable, V any] interface {
	Put(k K, v V) V
	Get(k K) (V, bool)
	Remove(k K) V
	RemoveMatch(k K, oldVal V) bool

	Size() int
	IsEmpty() bool
	Clear()
	PutAll(other Map[K, V])

	// GetOrDefault Java8+ 便捷方法
	GetOrDefault(k K, def V) V
	PutIfAbsent(k K, v V) V
	Replace(k K, newVal V) (V, bool)
	ReplaceMatch(k K, old, new V) bool

	// Seq2 Go 1.23+ 标准迭代器（替代 entrySet/keySet/values）
	Seq2() iter.Seq2[K, V]
	Seq() iter.Seq[V]
	Values() iter.Seq[V]
	Keys() iter.Seq[K]
}
