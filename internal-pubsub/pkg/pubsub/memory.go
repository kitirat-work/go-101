package pubsub

import (
	"context"
	"errors"
	"sync"
	"time"
)

// Bus คือสัญญาใช้งาน Pub/Sub
type Bus interface {
	// Subscribe สมัครรับอีเวนต์ตาม topic
	// buffer ถ้า <= 0 จะใช้ค่าจาก Options.DefaultBuffer
	Subscribe(topic Topic, buffer int) Subscription

	// Publish ส่งอีเวนต์ (เคารพ ctx เมื่อ DeliveryMode เป็น Block/Timeout)
	Publish(ctx context.Context, topic Topic, data any) error

	// Close ปิด bus และปิด chan ของทุก subscriber (idempotent)
	Close(ctx context.Context) error
}

type memoryBus struct {
	mu        sync.RWMutex
	topics    map[Topic]map[*memSub]struct{}
	closed    bool
	opts      Options
	closeOnce sync.Once
}

func New(opts Options) Bus {
	if opts.DefaultBuffer <= 0 {
		opts.DefaultBuffer = DefaultOptions().DefaultBuffer
	}
	return &memoryBus{
		topics: make(map[Topic]map[*memSub]struct{}),
		opts:   opts,
	}
}

func (b *memoryBus) Subscribe(topic Topic, buffer int) Subscription {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		// สร้าง sub ว่างที่ปิดแล้วจะทำให้ range ออกทันที
		ch := make(chan Event)
		close(ch)
		return &memSub{bus: b, topic: topic, ch: ch, closed: true}
	}

	if buffer <= 0 {
		buffer = b.opts.DefaultBuffer
	}

	sub := &memSub{
		bus: b, topic: topic,
		ch: make(chan Event, buffer),
	}
	if b.topics[topic] == nil {
		b.topics[topic] = make(map[*memSub]struct{})
	}
	b.topics[topic][sub] = struct{}{}
	return sub
}

func (b *memoryBus) Publish(ctx context.Context, topic Topic, data any) error {
	b.mu.RLock()
	if b.closed {
		b.mu.RUnlock()
		return ErrClosed
	}
	// snapshot subscribers เพื่อหลีกเลี่ยง hold lock นานเกินไปตอนส่ง
	subs := make([]*memSub, 0, len(b.topics[topic]))
	for s := range b.topics[topic] {
		subs = append(subs, s)
	}
	mode := b.opts.DeliveryMode
	timeout := time.Duration(b.opts.DeliveryTimeoutMs) * time.Millisecond
	b.mu.RUnlock()

	ev := Event{Topic: topic, Data: data}

	for _, s := range subs {
		switch mode {
		case DeliveryDrop:
			select {
			case s.ch <- ev:
			default:
				// ทิ้ง ไม่บล็อก
			}
		case DeliveryTimeout:
			if timeout <= 0 {
				timeout = 100 * time.Millisecond
			}
			timer := time.NewTimer(timeout)
			select {
			case s.ch <- ev:
				if !timer.Stop() {
					<-timer.C
				}
			case <-timer.C:
				// ทิ้งเพราะหมดเวลา
			case <-ctx.Done():
				if !timer.Stop() {
					<-timer.C
				}
				return ctx.Err()
			}
		default: // DeliveryBlock (เคารพ ctx)
			select {
			case s.ch <- ev:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
	}
	return nil
}

func (b *memoryBus) Close(ctx context.Context) error {
	var err error
	b.closeOnce.Do(func() {
		b.mu.Lock()
		defer b.mu.Unlock()
		if b.closed {
			return
		}
		b.closed = true
		// ปิดทุก subscription channel
		for _, set := range b.topics {
			for s := range set {
				s.closeNoLock()
			}
		}
		// ล้าง map เพื่อช่วย GC
		b.topics = make(map[Topic]map[*memSub]struct{})
	})
	// เคารพ ctx เฉย ๆ แม้การปิดจะ instant
	if ctx.Err() != nil && !errors.Is(ctx.Err(), context.Canceled) {
		err = ctx.Err()
	}
	return err
}
