package game

import (
	"nanoserver/protocol/protobuf/pb"
	"sync"
)

type Stats struct {
	lock sync.RWMutex

	Attrs map[pb.Attribute]int32
}

// 清空
func (s *Stats) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	for key := range s.Attrs {
		delete(s.Attrs, key)
	}
}

// 查询
func (s *Stats) getStats(k pb.Attribute) int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	st, ok := s.Attrs[k]
	if ok {
		return st
	} else {
		return 0
	}
}

// 累加
func (s *Stats) setStats(k pb.Attribute, v int32) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if k <= pb.Attribute_INVALID_Attribute || k >= pb.Attribute_MAX_Attribute {
		return
	}

	_, ok := s.Attrs[k]
	if ok {
		s.Attrs[k] += v
	} else {
		s.Attrs[k] = v
	}
}

// 替换
func (s *Stats) changeStats(k pb.Attribute, v int32) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if k <= pb.Attribute_INVALID_Attribute || k >= pb.Attribute_MAX_Attribute {
		return
	}

	if v < 0 {
		v = 0
	}
	s.Attrs[k] = v
}

// 拷贝
func (s *Stats) copyStats() map[pb.Attribute]int32 {
	s.lock.RLock()
	defer s.lock.RUnlock()

	attrs := make(map[pb.Attribute]int32)
	for k, v := range s.Attrs {
		attrs[k] = v
	}
	return attrs
}

func (s *Stats) Range(f func(key pb.Attribute, value int32) bool) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for k, v := range s.Attrs {
		if !f(k, v) {
			break
		}
	}
}
