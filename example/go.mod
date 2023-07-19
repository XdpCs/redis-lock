module github.com/XdpCs/redis-lock/example

go 1.17

replace github.com/XdpCs/redis-lock => ../

require (
	github.com/XdpCs/redis-lock v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.0.5
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)
