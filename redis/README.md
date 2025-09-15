# Go + Redis Best Practices Example

โปรเจกต์ตัวอย่างการใช้งาน **Redis** ใน **Go** ด้วยไลบรารีทางการ `github.com/redis/go-redis/v9` โดยเน้นแนวปฏิบัติที่ดี (best practices):
- แยกแพ็กเกจ `internal/redisclient` สำหรับ config + สร้าง client ครั้งเดียว (singleton) เพื่อ **reuse** connection pool
- ใช้ `context` + **timeouts** สำหรับทุกคำสั่งสำคัญ
- ตั้งค่า **connection pool** และ **timeouts** ผ่าน `redis.Options`
- Health check ด้วย `PING`
- ตัวอย่าง **pipeline** ลด RTT
- ตัวอย่าง **Pub/Sub** แบบสั้น
- ปิดทรัพยากรอย่าง **graceful** ตอน shutdown

> แนะนำให้อ่าน:  
> - *go-redis guide* และ *Connect your Go app to Redis* เพื่อดูการติดตั้ง/เชื่อมต่อขั้นพื้นฐาน :contentReference[oaicite:5]{index=5}  
> - *Pipelines and transactions* สำหรับการ batch คำสั่งแบบมีประสิทธิภาพ :contentReference[oaicite:6]{index=6}  
> - *Connection pools and multiplexing* เพื่อเข้าใจกลไกของ pooling/การคืน connection :contentReference[oaicite:7]{index=7}  
> - รายการ API/ตัวเลือก `Options` บน pkg.go.dev (อัปเดตตามเวอร์ชัน v9) :contentReference[oaicite:8]{index=8}

## โครงสร้าง
cmd/app/main.go // entrypoint
internal/redisclient/* // config + new client + ping + close
internal/cache/user_repo.go// ตัวอย่างเลเยอร์ cache

makefile
Copy code

## การตั้งค่า
คัดลอก `.env.example` เป็น `.env` แล้วปรับค่าให้เหมาะสม:
```env
REDIS_ADDR=127.0.0.1:6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_POOL_SIZE=20
REDIS_MIN_IDLE=5
REDIS_DIAL_TIMEOUT=5s
REDIS_READ_TIMEOUT=2s
REDIS_WRITE_TIMEOUT=2s
REDIS_CONN_MAX_IDLE=30m
REDIS_CONN_MAX_LIFETIME=2h
```

## การรัน

1. ติดตั้ง Redis (หรือต่อ Redis Cloud) แล้วเปิดใช้งาน
2. ติดตั้งไลบรารี:
    
    ```bash
    go mod tidy
    
    ```
    
3. รันโปรแกรม:
    
    ```bash
    go run ./cmd/app
    
    ```
    

จะเห็น log: `redis: ping OK`, ค่า `GET`, และผลลัพธ์ pipeline + ข้อความ Pub/Sub

## Best Practices ย่อ

- **Reuse client** ตัวเดียวต่อโปรเซส (go-redis มี pool ในตัว) แทนการสร้าง client ต่อ request [Redis](https://redis.io/docs/latest/develop/clients/pools-and-muxing/?utm_source=chatgpt.com)
- กำหนด **timeouts** (dial/read/write) ให้เหมาะกับระบบ ป้องกัน hang ยาว ๆ และช่วยคุม latency โดยรวม [Go Packages](https://pkg.go.dev/github.com/redis/go-redis/v9?utm_source=chatgpt.com)
- ใช้ **pipeline** เมื่อต้องยิงหลายคำสั่งติดกัน ลด round-trip time (RTT) [Redis](https://redis.io/docs/latest/develop/clients/go/transpipe/?utm_source=chatgpt.com)
- ใส่ **health check** (PING) ตอนเริ่มและ monitoring runtime เป็นระยะ (แล้วแต่ระบบ)
- วาง **key naming** ให้สื่อความหมาย/มี namespace เช่น `user:<id>`
- จัดชั้น **abstraction**: สร้าง repo/service หุ้มไปเหนือ client เพื่อลดการผูกติด vendor
- จัดการ **graceful shutdown**: ปิด client คืนทรัพยากร
- (ถ้าใช้เยอะ) **ปรับ pool** (`PoolSize/MinIdleConns`) ให้พอดีกับโหลดและ CPU/คอร์ที่มี

> หมายเหตุเวอร์ชัน: ใช้ v9.7.x (ดู release notes) และอาจมี option ใหม่/เปลี่ยนชื่อเล็กน้อยในอนาคต ควรเช็ค changelog สม่ำเสมอ
>

---

## หมายเหตุเชิงสถาปัตยกรรมสั้น ๆ (สำหรับโปรดักชัน)

- **Singleton client + DI**: ผูก `*redis.Client` เข้ากับเลเยอร์ repo/service ผ่าน constructor เพื่อให้เทส/ม็อคง่าย  
- **Observability**: ครอบคำสั่งด้วย metric (latency, errors, pool stats) และ tracing (OpenTelemetry)  
- **Resilience**: retry แบบมี backoff เฉพาะ error ประเภท network/timeout (ระวังไม่ให้เขียนซ้ำที่ไม่ idempotent)  
- **Security**: ใช้ TLS/ACL/Role ตามแวดล้อม (Redis Cloud/OSS) และอย่า log ค่า secret  
- **Key lifecycle**: ใส่ TTL กับคีย์ cache ทุกครั้ง อย่าเก็บถาวรโดยไม่จำเป็น  

---

