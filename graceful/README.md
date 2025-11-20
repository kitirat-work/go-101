# Graceful Shutdown Library

Library สำหรับจัดการการปิดเซิร์ฟเวอร์อย่างสวยงาม (Graceful Shutdown) ใน Go

## แนวคิด (Concept)

เมื่อเซิร์ฟเวอร์ได้รับสัญญาณให้หยุดทำงาน (เช่น SIGTERM, SIGINT) การปิดเซิร์ฟเวอร์แบบทันทีอาจทำให้:
- Request ที่กำลังประมวลผลอยู่ถูกตัดกลางคัน
- ข้อมูลสูญหาย
- Transaction ไม่สมบูรณ์
- Client ได้รับ error

**Graceful Shutdown** คือการปิดเซิร์ฟเวอร์อย่างสวยงามโดย:
1. หยุดรับ request ใหม่
2. รอให้ request ที่กำลังประมวลผลเสร็จสิ้น
3. ปิดการเชื่อมต่อและทำความสะอาดทรัพยากร
4. ปิดเซิร์ฟเวอร์อย่างปลอดภัย

## หลักการทำงาน (How It Works)

Library นี้ใช้หลักการทำงาน 5 ขั้นตอน:

```
1. Signal Handling
   ├─> ดักจับ OS signals (SIGTERM, SIGINT, etc.)
   └─> สร้าง context ที่จะ cancel เมื่อได้รับ signal

2. Server Start
   ├─> รันเซิร์ฟเวอร์ใน goroutine
   └─> ส่ง error กลับมาผ่าน channel

3. Wait for Signal or Error
   ├─> รอ signal จาก OS
   ├─> หรือรอ error จากเซิร์ฟเวอร์
   └─> ถ้าเซิร์ฟเวอร์พังเองก็ออกทันที

4. Graceful Shutdown
   ├─> สร้าง timeout context
   ├─> เรียก StopServer function
   └─> รอให้ shutdown เสร็จ

5. Timeout Protection
   ├─> ถ้า shutdown เสร็จภายใน timeout → ออกปกติ
   └─> ถ้า timeout → force exit
```

## การออกแบบ (Design)

### 1. Configuration-Based Design

ใช้ `GracefulShutdownConfig` struct เป็น configuration object ที่รวมทุกอย่างไว้:

```go
type GracefulShutdownConfig struct {
    ShutdownTimeoutSeconds int                            // เวลา timeout
    Signals                []os.Signal                    // signals ที่ต้องการดัก
    StartServer            func() error                   // ฟังก์ชันเริ่มเซิร์ฟเวอร์
    StopServer             func(ctx context.Context) error // ฟังก์ชันหยุดเซิร์ฟเวอร์
    Logger                 ILogger                        // logger
}
```

**ข้อดี:**
- ยืดหยุ่น: ผู้ใช้กำหนดวิธีการ start/stop เองได้
- ทดสอบง่าย: inject dependencies ผ่าน config
- ไม่ couple กับ framework เฉพาะ: ใช้กับ http.Server, gRPC, หรืออะไรก็ได้

### 2. Interface-Based Logging

```go
type ILogger interface {
    Println(v ...interface{})
    Printf(format string, v ...interface{})
}
```

**ข้อดี:**
- ผู้ใช้สามารถใช้ logger ของตัวเองได้ (zap, logrus, etc.)
- ทดสอบง่ายด้วย mock logger
- ไม่บังคับต้องใช้ `log.Default()`

### 3. Context-Based Control Flow

ใช้ `context.Context` เป็นหัวใจสำคัญในการควบคุม:

```go
// Signal context: cancel เมื่อได้รับ signal
ctx, stop := signal.NotifyContext(context.Background(), signals...)

// Timeout context: ป้องกัน shutdown ค้างนานเกินไป
shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
```

**ข้อดี:**
- เป็น Go idiom ที่ standard
- Propagate cancellation ได้
- Control timeout ได้แม่นยำ

### 4. Channel-Based Communication

ใช้ channel สำหรับสื่อสารระหว่าง goroutines:

```go
srvErr := make(chan error, 1)  // รับ error จากเซิร์ฟเวอร์
done := make(chan struct{})    // signal เมื่อ shutdown เสร็จ
```

**ข้อดี:**
- Thread-safe โดยธรรมชาติ
- Non-blocking communication
- รองรับ `select` statement

### 5. Defensive Programming

มีการป้องกันและ fallback หลายจุด:

```go
// Default logger ถ้าไม่ได้กำหนด
if logger == nil {
    logger = newLogger()
}

// Default signals
if len(cfg.Signals) == 0 {
    cfg.Signals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}
}

// Default timeout
if timeout <= 0 {
    timeout = 30 * time.Second
}

// เพิกเฉย http.ErrServerClosed (normal shutdown)
if err != nil && !errors.Is(err, http.ErrServerClosed) {
    srvErr <- err
}
```

## ตัวอย่างการใช้งาน

### HTTP Server

```go
package main

import (
    "context"
    "net/http"
    "time"
    
    "your-module/graceful"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        time.Sleep(2 * time.Second) // simulate work
        w.Write([]byte("Hello, World!"))
    })
    
    srv := &http.Server{
        Addr:    ":8080",
        Handler: mux,
    }
    
    graceful.Run(graceful.GracefulShutdownConfig{
        ShutdownTimeoutSeconds: 30,
        StartServer: func() error {
            return srv.ListenAndServe()
        },
        StopServer: func(ctx context.Context) error {
            return srv.Shutdown(ctx)
        },
    })
}
```

### Custom Logger

```go
import "go.uber.org/zap"

logger, _ := zap.NewProduction()
defer logger.Sync()

graceful.Run(graceful.GracefulShutdownConfig{
    Logger: logger.Sugar(),
    // ... other config
})
```

### Custom Signals

```go
graceful.Run(graceful.GracefulShutdownConfig{
    Signals: []os.Signal{syscall.SIGUSR1, syscall.SIGUSR2},
    // ... other config
})
```

## Flow Diagram

```
    [Start]
       |
       v
[ดักจับ Signals] ──────────────┐
       |                        |
       v                        |
[Start Server]                  |
   (goroutine)                  |
       |                        |
       v                        |
   [Select]<───────────────────┘
     /   \
    /     \
[Signal] [Error]
   |        |
   |        └──> [Log & Exit]
   |
   v
[Stop Server]
(with timeout)
   |
   v
 [Select]
  /    \
[Done] [Timeout]
  |       |
  v       v
[Exit] [Force Exit]
```

## ข้อควรระวัง

1. **StartServer ต้อง block**: ฟังก์ชัน `StartServer` ควร block จนกว่าเซิร์ฟเวอร์จะหยุด
2. **StopServer ต้องเคารพ context**: ต้องหยุดการทำงานเมื่อ context timeout
3. **Timeout ควรเพียงพอ**: กำหนด `ShutdownTimeoutSeconds` ให้เพียงพอกับ request ที่ยาวที่สุด
4. **http.ErrServerClosed เป็นเรื่องปกติ**: error นี้เกิดขึ้นเมื่อ shutdown ปกติ ไม่ต้องกังวล