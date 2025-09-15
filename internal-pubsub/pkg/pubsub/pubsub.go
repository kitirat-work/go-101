package pubsub

// Topic คือชื่อหัวข้อสำหรับกระจายอีเวนต์
type Topic string

// Event คือข้อมูลอีเวนต์ที่ส่งข้ามโมดูล
type Event struct {
	Topic Topic
	Data  any
}

// DeliveryMode กำหนดกลยุทธ์การส่ง
type DeliveryMode int

const (
	// DeliveryBlock: ส่งแบบบล็อกจนกว่าจะส่งเข้า chan ของผู้รับได้ (เคารพ ctx)
	DeliveryBlock DeliveryMode = iota
	// DeliveryDrop: ถ้า chan ผู้รับเต็ม ให้ทิ้ง (ไม่บล็อก)
	DeliveryDrop
	// DeliveryTimeout: พยายามส่งจนถึง timeout ที่กำหนดไว้ใน Options แล้วค่อยทิ้ง
	DeliveryTimeout
)

// Options ปรับแต่งพฤติกรรมของ Bus
type Options struct {
	DefaultBuffer     int          // ความจุ chan เริ่มต้นของ subscriber (ถ้าไม่ระบุ)
	DeliveryMode      DeliveryMode // โหมดการส่ง
	DeliveryTimeoutMs int          // ใช้เมื่อ DeliveryTimeout (>0)
}

// DefaultOptions ค่าปริยาย
func DefaultOptions() Options {
	return Options{
		DefaultBuffer:     1,
		DeliveryMode:      DeliveryBlock,
		DeliveryTimeoutMs: 0,
	}
}
