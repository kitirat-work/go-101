package pubsub

import "sync"

// Subscription ให้ฝั่งผู้รับใช้ดึงอีเวนต์ผ่าน C และยกเลิกด้วย Unsubscribe
type Subscription interface {
	C() <-chan Event
	Unsubscribe()
}

type memSub struct {
	bus    *memoryBus
	topic  Topic
	ch     chan Event
	closed bool
	once   sync.Once
}

func (s *memSub) C() <-chan Event { return s.ch }

func (s *memSub) Unsubscribe() {
	s.once.Do(func() {
		s.bus.mu.Lock()
		defer s.bus.mu.Unlock()
		if s.closed {
			return
		}
		if subs, ok := s.bus.topics[s.topic]; ok {
			delete(subs, s)
			if len(subs) == 0 {
				delete(s.bus.topics, s.topic)
			}
		}
		close(s.ch)
		s.closed = true
	})
}

func (s *memSub) closeNoLock() {
	if s.closed {
		return
	}
	close(s.ch)
	s.closed = true
}
