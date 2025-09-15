package subscriber

import (
	"context"
	"fmt"
	"internal-pubsub/pkg/pubsub"
	"runtime/debug"
	"sync"
)

// Handler ฟังก์ชันประมวลผล event แบบทีละรายการ
type Handler func(ctx context.Context, ev pubsub.Event) error

// Subscriber ดูแล lifecycle ของ subscription
type Subscriber struct {
	bus    pubsub.Bus
	topic  pubsub.Topic
	sub    pubsub.Subscription
	opts   Options
	closed bool
	mu     sync.Mutex
}

// New สร้าง Subscriber และ subscribe ทันที (buffer <= 0 ใช้ค่า default ของ bus)
func New(bus pubsub.Bus, topic pubsub.Topic, buffer int, optFns ...Option) *Subscriber {
	opts := DefaultOptions()
	for _, f := range optFns {
		f(&opts)
	}
	return &Subscriber{
		bus:   bus,
		topic: topic,
		sub:   bus.Subscribe(topic, buffer),
		opts:  opts,
	}
}

// Close ยกเลิกการรับข้อความ (idempotent)
func (s *Subscriber) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.sub.Unsubscribe()
	s.closed = true
}

// Run จะอ่านจาก channel แบบ sequential และเรียก handler ทีละ event
// จะจบเมื่อ ctx.Done() หรือ channel ถูกปิด หรือ handler คืน error (เมื่อ StopOnError == true)
func (s *Subscriber) Run(ctx context.Context, handler Handler) error {
	// ห่อ handler ให้ recover panic เสมอ
	runHandler := func(ctx context.Context, ev pubsub.Event) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("handler panic: %v", r)
				s.opts.Logger.Printf("[subscriber] panic recovered topic=%s err=%v\n%s",
					s.topic, err, string(debug.Stack()))
			}
		}()
		return handler(ctx, ev)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case ev, ok := <-s.sub.C():
			if !ok {
				// channel ปิดจาก Unsubscribe() หรือ bus.Close()
				return nil
			}
			if err := runHandler(ctx, ev); err != nil {
				s.opts.Logger.Printf("[subscriber] handler error topic=%s err=%v", s.topic, err)
				if s.opts.StopOnError {
					return err
				}
				// ถ้าไม่หยุด ให้ continue อ่านตัวถัดไป
			}
		}
	}
}
