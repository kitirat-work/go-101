package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/joho/godotenv"

	"go-redis-best-practice/internal/cache"
	"go-redis-best-practice/internal/redisclient"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_ = godotenv.Load() // optional

	cfg := redisclient.LoadConfig()
	rdb := redisclient.New(cfg)
	defer func() {
		_ = redisclient.Close()
	}()

	ctx := context.Background()

	// 1) Health check
	must(redisclient.PingWithTimeout(ctx, rdb, 1*time.Second))
	log.Println("redis: ping OK")

	// 2) ใช้งานผ่านเลเยอร์ repo
	userRepo := cache.NewUserRepo(rdb)
	must(userRepo.CacheUser(ctx, "42", `{"name":"Ada"}`, 10*time.Minute))
	val, err := userRepo.GetUser(ctx, "42")
	must(err)
	log.Println("redis GET user:42 =", val)

	// 3) Pipeline example
	_, _ = rdb.Set(ctx, "user:1", `{"name":"Alan"}`, 0).Result()
	_, _ = rdb.Set(ctx, "user:2", `{"name":"Edsger"}`, 0).Result()
	data, err := userRepo.PipelineExample(ctx, []string{"1", "2", "nope"})
	must(err)
	log.Printf("pipeline results: %+v\n", data)

	// 4) (ทางเลือก) Pub/Sub demo แบบสั้น
	go func() {
		sub := rdb.Subscribe(ctx, "demo:news")
		defer sub.Close()
		for msg := range sub.Channel() {
			log.Printf("[SUB] %s: %s\n", msg.Channel, msg.Payload)
			break
		}
	}()
	time.Sleep(200 * time.Millisecond)
	_ = rdb.Publish(ctx, "demo:news", "hello from publisher").Err()

	// 5) Graceful shutdown (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Println("shutting down...")
}
