package lru

import "container/list"

// 这里的cache只是lru中的，不是并发安全的
type Cache struct {
	maxBytes int64
	nbytes   int64
	// 标准库中的双向链表
	ll    *list.List
	cache map[string]*list.Element
	// 某条记录被移除时的回调函数，可以为 nil
	OnEvicted func(key string, value Value)
}

// 节点淘汰的时候，需要同时删除map中的kv对，为链表节点中的数据
type entry struct {
	key   string
	value Value
}

// 用Len计算用了多大内存
type Value interface {
	Len() int
}

// 方便实例化 Cache，实现 New()函数
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get从map中查找节点指针，移动节点到队首，返回值和成功与否
func (c *Cache) Get(key string) (value Value, ok bool) {
	// 查询的key存在的话，移动到队首
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	// key不存在，占用内存0，并且返回false
	return
}

// Remove0ldest()删除队尾的节点和map中的kv并更新内存
func (c *Cache) Remove0ldest() {
	// 获取队尾的节点，为空返回nil
	ele := c.ll.Back()
	if ele != nil {
		//取尾部节点，删除
		c.ll.Remove(ele)
		//拿到尾部节点的map kv对，删除
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//更新当前所用内存
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

func (c *Cache) Add(key string, value Value) {
	//操作完成后更新nbytes，使小于maxBytes
	// key已经存在了，移动到队首,修改value(相同则没影响)
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		// key不存在，&entry{key,value}新加到队首，增加map
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.Remove0ldest()
	}
}

// Len()获取添加了多少条数据，用于测试
func (c *Cache) Len() int {
	return c.ll.Len()
}
