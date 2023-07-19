package main

import (
	"context"
	"time"

	redislock "github.com/XdpCs/redis-lock"
	"github.com/redis/go-redis/v9"
)

func main() {
	// init context
	ctx := context.Background()
	// init redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: ":6379",
	})
	// close redis client
	defer rdb.Close()
	// flush redis
	_ = rdb.FlushDB(ctx).Err()
	// init redislock client
	client := redislock.NewDefaultClient(rdb)
	// try lock with default parameter
	mutex, err := client.TryLock(ctx, "XdpCs", -1)
	if err != nil {
		panic(err)
	}

	defer func(mutex *redislock.Mutex, ctx context.Context) {
		// unlock mutex
		err := mutex.Unlock(ctx)
		if err != nil {
			panic(err)
		}
	}(mutex, ctx)
	time.Sleep(time.Second * 30)
}
