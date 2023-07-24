# redis-lock

![GitHub watchers](https://img.shields.io/github/watchers/XdpCs/redis-lock?style=social)
![GitHub stars](https://img.shields.io/github/stars/XdpCs/redis-lock?style=social)
![GitHub forks](https://img.shields.io/github/forks/XdpCs/redis-lock?style=social)
![GitHub last commit](https://img.shields.io/github/last-commit/XdpCs/redis-lock?style=flat-square)
![GitHub repo size](https://img.shields.io/github/repo-size/XdpCs/redis-lock?style=flat-square)
![GitHub license](https://img.shields.io/github/license/XdpCs/redis-lock?style=flat-square)

Distributed lock based on [redis](https://redis.io/docs/manual/patterns/distributed-locks/).

redis-lock supports watchdog mechanism in [redisson](https://github.com/redisson/redisson).

## install

`go get`

```shell
go get -u github.com/XdpCs/redis-lock
```

`go mod`

```shell
require github.com/XdpCs/redis-lock latest
```

## example

Error handling is simplified to panic for shorter example.

You can run this program in this [directory](./example/main.go).

```go
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
	client, err := redislock.NewDefaultClient(rdb)
	if err != nil {
		panic(err)
	}
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
```

## License

redis-lock is under the [MIT](LICENSE). Please refer to LICENSE for more information.