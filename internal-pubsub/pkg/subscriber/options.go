package subscriber

import "log"

type Options struct {
	// ถ้า true: เมื่อ handler คืน error ให้หยุด Run() (รีเทิร์น error กลับไป)
	// ถ้า false: แค่ log แล้วอ่านข้อความถัดไป
	StopOnError bool

	// Logger (ค่าเริ่มต้นใช้ log.Default())
	Logger *log.Logger
}

func DefaultOptions() Options {
	return Options{
		StopOnError: false,
		Logger:      log.Default(),
	}
}

type Option func(*Options)

func WithStopOnError(stop bool) Option {
	return func(o *Options) { o.StopOnError = stop }
}

func WithLogger(l *log.Logger) Option {
	return func(o *Options) {
		if l != nil {
			o.Logger = l
		}
	}
}
