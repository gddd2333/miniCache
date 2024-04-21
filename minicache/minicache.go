package minicache

import (
	"fmt"
	"log"
	pb "minicache/minicachepb"
	"minicache/singleflight"
	"sync"
)

// 缓存的命名空间，每个group有唯一的name可以区分
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     PeerPicker
	// singleflight确保多个相同的请求只向节点客户端请求一次
	loader *singleflight.Group
}

// get方法的接口
type Getter interface {
	Get(key string) ([]byte, error)
}

// 定义函数类型
type GetterFunc func(key string) ([]byte, error)

// 实现函数类型的Get方法
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

// 创建一个新的group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}
	groups[name] = g
	return g
}

// 获取特定名称为name的group,只读锁
func GetGroup(name string) *Group {
	mu.Lock()
	g := groups[name]
	mu.Unlock()
	return g
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}
	if v, ok := g.mainCache.get(key); ok {
		log.Println("[miniCache] hit")
		return v, nil
	}
	return g.load(key)
}

// RegisterPeers方法将实现了peerpick接口的HTTPPOOL注入到group中
func (g *Group) RegisterPeers(peers PeerPicker) {
	if g.peers != nil {
		panic("RegisterPeerPicker called more than once")
	}
	g.peers = peers
}

// 缓存不存在，是否从远程加载，还是从本地加载
func (g *Group) load(key string) (value ByteView, err error) {
	// 用singleflight包裹，确保多个相同请求只发送一次，防止缓存击穿
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}
				log.Println("[miniCache] failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

// 调用Get方法，从peer客户端查询key
func (g *Group) getFromPeer(peer PeerGetter, key string) (ByteView, error) {
	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &pb.Response{}
	err := peer.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}

// 回调
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

// 回调后写入cache
func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
