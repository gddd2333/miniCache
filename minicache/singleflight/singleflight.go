package singleflight

import "sync"

// 防止缓存击穿，提高性能
// 将多次相同的请求，只执行一次

// call表示正在进行中，或已经结束的请求。使用sync.WaitGroup锁避免重复进入
type call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

// Group是singleflight的主数据结构，管理不同key和call
type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	//防止g被并发读写
	g.mu.Lock()
	//延迟实例化，节省内存
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	if c, ok := g.m[key]; ok {
		// 请求正在执行，等待后，直接拿请求的结果
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}
	// 请求没有被执行过，new call，Add锁加一
	c := new(call)
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()
	// 执行，执行完成后done，锁减一，拿到结果
	c.val, c.err = fn()
	c.wg.Done()
	// 执行完成后删除call的map
	g.mu.Lock()
	delete(g.m, key)
	g.mu.Unlock()

	return c.val, c.err
}
