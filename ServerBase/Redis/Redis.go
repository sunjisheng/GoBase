package Redis

import (
	"context"
	"github.com/go-redis/redis"
)


var redis_instances []*redis.Client
func InitRedisInstance(maxDbCount uint32) {
	redis_instances = make([]*redis.Client, maxDbCount)
}

func Redis_Instance(index uint32) *redis.Client {
	return redis_instances[index]
}

func Redis_Init(index uint32, addr string,pwd string)  bool{
	rdb := redis.NewClient(&redis.Options{Addr: addr,	Password: pwd, 	DB: 0, })
	ctx := context.Background()
	_, err:= rdb.Ping(ctx).Result()
	if err == nil {
		redis_instances[index] = rdb
		return true
	} else {
		return false
	}
}
