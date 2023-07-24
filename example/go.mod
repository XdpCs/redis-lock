module github.com/XdpCs/redis-lock/example

go 1.18

replace github.com/XdpCs/redis-lock => ../

require (
	github.com/XdpCs/redis-lock v0.0.0-20230719094903-e79ff7e15277
	github.com/redis/go-redis/v9 v9.0.5
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
