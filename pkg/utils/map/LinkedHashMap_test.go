package _map

import (
	"testing"
)

// TestNewLinkedHashMap_InsertOrder 测试插入顺序模式
func TestNewLinkedHashMap_InsertOrder(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	})

	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}

	// 验证插入顺序（使用 Keys 迭代器）
	expected := []string{"a", "b", "c"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected order %v, got %v", expected, actual)
			break
		}
	}
}

// TestNewLinkedHashMap_Empty 测试空 Map 创建
func TestNewLinkedHashMap_Empty(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil)

	if m.Size() != 0 {
		t.Errorf("Expected empty map, got size %d", m.Size())
	}

	if m.head != nil || m.tail != nil {
		t.Error("Expected nil head and tail for empty map")
	}
}

// TestNewLinkedHashMap_WithLRU 测试 LRU 模式
func TestNewLinkedHashMap_WithLRU(t *testing.T) {
	m := NewLinkedHashMap[string, int](
		[]Entry[string, int]{
			{"a", 1},
			{"b", 2},
			{"c", 3},
		},
		WithAccessOrder[string, int](),
	)

	if !m.accessOrder {
		t.Error("Expected accessOrder to be true")
	}

	// 访问 "a"，应该移到尾部
	m.Get("a")

	// 验证顺序：b, c, a（使用 Keys 迭代器）
	expected := []string{"b", "c", "a"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected LRU order %v, got %v", expected, actual)
			break
		}
	}
}

// TestSet_UpdateValue 测试更新已有键的值
func TestSet_UpdateValue(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil)
	m.Set("a", 1)
	m.Set("a", 10) // 更新

	val, ok := m.Get("a")
	if !ok {
		t.Error("Expected key 'a' to exist")
	}
	if val != 10 {
		t.Errorf("Expected value 10, got %d", val)
	}
	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}
}

// TestSet_InsertOrderPreserved 测试插入顺序保持不变
func TestSet_InsertOrderPreserved(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil)
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("a", 10) // 更新，但不改变顺序（插入顺序模式）

	expected := []string{"a", "b"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected order %v, got %v", expected, actual)
			break
		}
	}
}

// TestGet_NonExistentKey 测试获取不存在的键
func TestGet_NonExistentKey(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil)

	val, ok := m.Get("nonexistent")
	if ok {
		t.Error("Expected ok to be false for non-existent key")
	}
	if val != 0 {
		t.Errorf("Expected zero value, got %d", val)
	}
}

// TestRemove 测试删除节点
func TestRemove(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	})

	// 删除中间节点
	removed := m.Remove("b")
	if !removed {
		t.Error("Expected Remove to return true")
	}
	if m.Size() != 2 {
		t.Errorf("Expected size 2, got %d", m.Size())
	}
	if m.Has("b") {
		t.Error("Expected key 'b' to be removed")
	}

	// 验证顺序：a, c（使用 Keys 迭代器）
	expected := []string{"a", "c"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected order %v, got %v", expected, actual)
			break
		}
	}
}

// TestRemove_HeadAndTail 测试删除头尾节点
func TestRemove_HeadAndTail(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	})

	// 删除头节点
	m.Remove("a")
	if m.head == nil {
		t.Fatal("Expected head to not be nil after removing first element")
	}
	if m.head.key != "b" {
		t.Errorf("Expected new head to be 'b', got '%s'", m.head.key)
	}

	// 删除尾节点
	m.Remove("c")
	if m.tail == nil {
		t.Fatal("Expected tail to not be nil after removing last element")
	}
	if m.tail.key != "b" {
		t.Errorf("Expected new tail to be 'b', got '%s'", m.tail.key)
	}

	if m.Size() != 1 {
		t.Errorf("Expected size 1, got %d", m.Size())
	}
}

// TestRemove_NonExistentKey 测试删除不存在的键
func TestRemove_NonExistentKey(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil)

	removed := m.Remove("nonexistent")
	if removed {
		t.Error("Expected Remove to return false for non-existent key")
	}
}

// TestHas 测试 Has 方法
func TestHas(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
	})

	if !m.Has("a") {
		t.Error("Expected Has('a') to return true")
	}
	if m.Has("b") {
		t.Error("Expected Has('b') to return false")
	}
}

// TestClear 测试清空 Map
func TestClear(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
	})

	m.Clear()

	if m.Size() != 0 {
		t.Errorf("Expected size 0 after clear, got %d", m.Size())
	}
	if m.head != nil || m.tail != nil {
		t.Error("Expected nil head and tail after clear")
	}
	if len(m.data) != 0 {
		t.Errorf("Expected empty data map after clear, got length %d", len(m.data))
	}
}

// TestAll 测试 All 迭代器（替代 ForEach）
func TestAll(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	})

	keys := make([]string, 0)
	values := make([]int, 0)

	for key, val := range m.Entries() {
		keys = append(keys, key)
		values = append(values, val)
	}

	expectedKeys := []string{"a", "b", "c"}
	expectedValues := []int{1, 2, 3}

	if len(keys) != len(expectedKeys) {
		t.Errorf("Expected %d keys, got %d", len(expectedKeys), len(keys))
	}

	for i := range expectedKeys {
		if keys[i] != expectedKeys[i] {
			t.Errorf("Expected key %s at index %d, got %s", expectedKeys[i], i, keys[i])
		}
		if values[i] != expectedValues[i] {
			t.Errorf("Expected value %d at index %d, got %d", expectedValues[i], i, values[i])
		}
	}
}

// TestKeys 测试 Keys 迭代器
func TestKeys(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"x", 10},
		{"y", 20},
		{"z", 30},
	})

	var count int
	for key := range m.Keys() {
		if !m.Has(key) {
			t.Errorf("Key %s should exist", key)
		}
		count++
	}

	if count != 3 {
		t.Errorf("Expected 3 keys, got %d", count)
	}
}

// TestValues 测试 Values 迭代器
func TestValues(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
	})

	sum := 0
	for val := range m.Values() {
		sum += val
	}

	if sum != 6 {
		t.Errorf("Expected sum 6, got %d", sum)
	}
}

// TestEarlyExit 测试提前退出迭代
func TestEarlyExit(t *testing.T) {
	m := NewLinkedHashMap[string, int]([]Entry[string, int]{
		{"a", 1},
		{"b", 2},
		{"c", 3},
		{"d", 4},
	})

	count := 0
	for key := range m.Keys() {
		count++
		if key == "b" {
			break // 提前退出
		}
	}

	if count != 2 {
		t.Errorf("Expected to iterate 2 items, got %d", count)
	}
}

// TestLRU_EvictionOrder 测试 LRU 淘汰顺序
func TestLRU_EvictionOrder(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil, WithAccessOrder[string, int]())
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// 访问 "a"，移到尾部
	m.Get("a")

	// 此时顺序应该是: b, c, a
	// 最少使用的是 "b"（在头部）

	expected := []string{"b", "c", "a"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected LRU order %v, got %v", expected, actual)
			break
		}
	}
}

// TestLRU_UpdateMovesToTail 测试 LRU 模式下更新值会移到尾部
func TestLRU_UpdateMovesToTail(t *testing.T) {
	m := NewLinkedHashMap[string, int](nil, WithAccessOrder[string, int]())
	m.Set("a", 1)
	m.Set("b", 2)
	m.Set("c", 3)

	// 更新 "a" 的值，应该移到尾部
	m.Set("a", 10)

	// 顺序应该是: b, c, a
	expected := []string{"b", "c", "a"}
	actual := make([]string, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected order %v, got %v", expected, actual)
			break
		}
	}

	// 验证值已更新
	val, _ := m.Get("a")
	if val != 10 {
		t.Errorf("Expected value 10, got %d", val)
	}
}

// TestMixedTypes 测试不同类型
func TestMixedTypes(t *testing.T) {
	// 字符串 -> 任意类型
	m := NewLinkedHashMap[string, any]([]Entry[string, any]{
		{"name", "John"},
		{"age", 30},
		{"active", true},
	})

	if m.Size() != 3 {
		t.Errorf("Expected size 3, got %d", m.Size())
	}

	name, ok := m.Get("name")
	if !ok || name != "John" {
		t.Errorf("Expected name 'John', got %v", name)
	}

	age, ok := m.Get("age")
	if !ok || age != 30 {
		t.Errorf("Expected age 30, got %v", age)
	}
}

// TestIntKey 测试整数键
func TestIntKey(t *testing.T) {
	m := NewLinkedHashMap[int, string]([]Entry[int, string]{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	})

	val, ok := m.Get(2)
	if !ok || val != "two" {
		t.Errorf("Expected 'two', got %v", val)
	}

	expected := []int{1, 2, 3}
	actual := make([]int, 0)
	for key := range m.Keys() {
		actual = append(actual, key)
	}

	for i, exp := range expected {
		if i >= len(actual) || actual[i] != exp {
			t.Errorf("Expected order %v, got %v", expected, actual)
			break
		}
	}
}

// BenchmarkSet 性能测试：插入
func BenchmarkSet(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}
}

// BenchmarkGet 性能测试：查找
func BenchmarkGet(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Get(i % 1000)
	}
}

// BenchmarkRemove 性能测试：删除
func BenchmarkRemove(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < b.N; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		m.Remove(i)
	}
}

// BenchmarkAll 性能测试：All 迭代器
func BenchmarkAll(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for key, val := range m.Entries() {
			_ = key
			_ = val
		}
	}
}

// BenchmarkKeys 性能测试：Keys 迭代器
func BenchmarkKeys(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for key := range m.Keys() {
			_ = key
		}
	}
}

// BenchmarkValues 性能测试：Values 迭代器
func BenchmarkValues(b *testing.B) {
	m := NewLinkedHashMap[int, int](nil)
	for i := 0; i < 1000; i++ {
		m.Set(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for val := range m.Values() {
			_ = val
		}
	}
}
