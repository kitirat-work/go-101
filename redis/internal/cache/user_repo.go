package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserRepo struct {
	rdb *redis.Client
}

func NewUserRepo(rdb *redis.Client) *UserRepo {
	return &UserRepo{rdb: rdb}
}

func (r *UserRepo) CacheUser(ctx context.Context, id string, payload string, ttl time.Duration) error {
	// SET key value EX ttl  (atomic)
	return r.rdb.Set(ctx, fmt.Sprintf("user:%s", id), payload, ttl).Err()
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (string, error) {
	return r.rdb.Get(ctx, fmt.Sprintf("user:%s", id)).Result()
}

func (r *UserRepo) PipelineExample(ctx context.Context, ids []string) (map[string]string, error) {
	// ตัวอย่าง pipeline ลด RTT ตามเอกสาร Redis (batch commands)
	pipe := r.rdb.Pipeline()
	cmds := make([]*redis.StringCmd, 0, len(ids))
	for _, id := range ids {
		cmds = append(cmds, pipe.Get(ctx, fmt.Sprintf("user:%s", id)))
	}
	_, err := pipe.Exec(ctx) // runs the batch
	if err != nil && err != redis.Nil {
		return nil, err
	}
	out := make(map[string]string, len(cmds))
	for i, cmd := range cmds {
		if val, err := cmd.Result(); err == nil {
			out[ids[i]] = val
		}
	}
	return out, nil
}
