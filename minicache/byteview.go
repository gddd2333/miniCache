package minicache

// 缓存值的抽象和封装
// 只读的数据结构ByteView表示缓存值

// 存储真实的缓存值，byte类型能够支持任意类型的存储（字符串、图片等）
type ByteView struct {
	b []byte
}

// 实现Len() int 方法，lru.Cache中，要求被缓存的对象实现Value接口
func (v ByteView) Len() int {
	return len(v.b)
}

// b作为只读数据，返回b的拷贝
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

// []byte转化为字符串返回，必要时拷贝
func (v ByteView) String() string {
	return string(v.b)
}

// 拷贝b
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
