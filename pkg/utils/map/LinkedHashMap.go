package _map

import "iter"

type Entry[K comparable, V any] struct {
	Key K
	Val V
}

type Option[K comparable, V any] func(*LinkedHashMap[K, V])

// WithAccessOrder 设置访问顺序模式（LRU）
func WithAccessOrder[K comparable, V any]() Option[K, V] {
	return func(m *LinkedHashMap[K, V]) {
		m.accessOrder = true
	}
}

// Node 双向链表节点（对应 LinkedHashMap 内部 Entry）
type Node[K comparable, V any] struct {
	key  K
	val  V
	prev *Node[K, V]
	next *Node[K, V]
}

// LinkedHashMap 复刻 Java LinkedHashMap
type LinkedHashMap[K comparable, V any] struct {
	data        map[K]*Node[K, V] // 哈希表
	head, tail  *Node[K, V]       // 双向链表头尾
	size        int
	accessOrder bool // false=插入顺序(默认)，true=访问顺序LRU
}

// NewLinkedHashMap 默认插入顺序（等价 new LinkedHashMap<>()）
func NewLinkedHashMap[K comparable, V any](entries []Entry[K, V], opts ...Option[K, V]) *LinkedHashMap[K, V] {
	m := &LinkedHashMap[K, V]{
		data:        make(map[K]*Node[K, V]),
		accessOrder: false, // 默认插入顺序
	}

	for _, entry := range entries {
		m.Put(entry.Key, entry.Val)
	}

	// 应用所有选项
	for _, opt := range opts {
		opt(m)
	}

	return m
}

// Put 对应 Java put()，JS set()
func (m *LinkedHashMap[K, V]) Put(key K, val V) {
	// 已存在key：更新值；accessOrder=true则移到队尾
	if node, ok := m.data[key]; ok {
		node.val = val
		if m.accessOrder {
			m.moveToTail(node)
		}
		return
	}

	// 新建节点，尾部插入
	newNode := &Node[K, V]{key: key, val: val}
	m.data[key] = newNode
	m.addToTail(newNode)
	m.size++
}

// Get 对应 Java get()，JS get()
func (m *LinkedHashMap[K, V]) Get(key K) (V, bool) {
	node, ok := m.data[key]
	if !ok {
		var zero V
		return zero, false
	}
	// 访问顺序模式：get后移到尾部
	if m.accessOrder {
		m.moveToTail(node)
	}
	return node.val, true
}

// Remove 对应 Java remove()，JS delete()
func (m *LinkedHashMap[K, V]) Remove(key K) bool {
	node, ok := m.data[key]
	if !ok {
		return false
	}
	delete(m.data, key)
	m.removeNode(node)
	m.size--
	return true
}

// ContainsKey 对应 Java containsKey()，JS has()
func (m *LinkedHashMap[K, V]) ContainsKey(key K) bool {
	_, ok := m.data[key]
	return ok
}

// Size 对应 size()，JS size
func (m *LinkedHashMap[K, V]) Size() int {
	return m.size
}

// Clear 清空
func (m *LinkedHashMap[K, V]) Clear() {
	m.data = make(map[K]*Node[K, V])
	m.head, m.tail = nil, nil
	m.size = 0
}

// Seq2 返回一个迭代器，支持 range 遍历（Go 1.23+ Seq2）
func (m *LinkedHashMap[K, V]) Seq2() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		cur := m.head
		for cur != nil {
			if !yield(cur.key, cur.val) {
				return
			}
			cur = cur.next
		}
	}
}

// Seq 返回键的迭代器（Go 1.23+ Seq）
func (m *LinkedHashMap[K, V]) Seq() iter.Seq[V] {
	return func(yield func(V) bool) {
		cur := m.head
		for cur != nil {
			if !yield(cur.val) {
				return
			}
			cur = cur.next
		}
	}
}

// Keys 返回键的迭代器（Go 1.23+ Seq）
func (m *LinkedHashMap[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		cur := m.head
		for cur != nil {
			if !yield(cur.key) {
				return
			}
			cur = cur.next
		}
	}
}

// Values 返回值的迭代器（Go 1.23+ Seq）
func (m *LinkedHashMap[K, V]) Values() iter.Seq[V] {
	return m.Seq()
}

// 内部：尾部添加节点
func (m *LinkedHashMap[K, V]) addToTail(node *Node[K, V]) {
	if m.tail == nil {
		m.head = node
		m.tail = node
		return
	}
	node.prev = m.tail
	m.tail.next = node
	m.tail = node
}

// 内部：移除节点
func (m *LinkedHashMap[K, V]) removeNode(node *Node[K, V]) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		m.head = node.next // 删头节点
	}
	if node.next != nil {
		node.next.prev = node.prev
	} else {
		m.tail = node.prev // 删尾节点
	}
	node.prev, node.next = nil, nil
}

// 内部：移到尾部（LRU访问顺序）
func (m *LinkedHashMap[K, V]) moveToTail(node *Node[K, V]) {
	if node == m.tail {
		return
	}
	m.removeNode(node)
	m.addToTail(node)
}
