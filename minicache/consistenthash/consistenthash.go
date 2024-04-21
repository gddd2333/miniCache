package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// 定义函数类型，采用依赖注入，允许替换成自定义，，默认crc32
type Hash func(data []byte) uint32

// 一致性哈希的主数据结构
type Map struct {
	hash     Hash           //func
	replicas int            //虚拟节点倍数
	keys     []int          //sorted  哈希环
	hashMap  map[int]string //虚拟节点和真实节点的映射
}

// 构造函数，允许自定义replicas和哈希算法
func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// 添加真是节点的方法，允许传入0个或多个真实节点的名称
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 对每个真实节点key，创建repliccas个虚拟节点
		for i := 0; i < m.replicas; i++ {
			// 名称是 序号i加真实key
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			// m.hash计算虚拟节点的哈希值，添加到环上
			m.keys = append(m.keys, hash)
			// map建立虚拟节点和真实节点的映射
			m.hashMap[hash] = key
		}
	}
	// 换上的哈希值排序
	sort.Ints(m.keys)
}

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	// 计算节点的哈希值
	hash := int(m.hash([]byte(key)))
	// 顺时针找到第一个匹配的虚拟节点下标idx，并从keys中获取对应的哈希值
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 通过hashMap映射得到真实节点

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
