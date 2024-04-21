package minicache

import pb "minicache/minicachepb"

// PeerGetter对应http客户端

// PeerPicker根据传入的key选择相应的节点PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter接口中的Get方法用于从对应的group查找缓存
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
