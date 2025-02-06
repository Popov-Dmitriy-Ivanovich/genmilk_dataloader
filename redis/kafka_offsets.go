package redis

import (
	"context"
	"strconv"
	"sync"

	"github.com/redis/go-redis/v9"
)

type RedisQuerier struct {
	RedisDB *redis.Client
	Ctx     context.Context
}

var redisQuerier RedisQuerier
var redisQuerierMtx sync.Mutex
var initialized bool

const (
	LACTATION_LOAD_OFFSET_KEY  = "LactationLoadOffset"
	COW_LOAD_OFFSET_KEY        = "CowLoadOffset"
	GENETIC_LOAD_OFFSET_KEY    = "GeneticLoadOffset"
	CHECK_MILK_LOAD_OFFSET_KEY = "CheckMilkLoadOffset"
	DAILY_MILK_LOAD_OFFSET_KEY = "DailyMilkLoadOffset"
)

func GetRedisQuerier() RedisQuerier {
	if !initialized {
		redisQuerierMtx.Lock()
		defer redisQuerierMtx.Unlock()
		if initialized {
			return redisQuerier
		}
		res := RedisQuerier{}
		res.RedisDB = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
		res.Ctx = context.Background()
		redisQuerier = res
		initialized = true
	}
	return redisQuerier
}

func (rq RedisQuerier) getUintKey(key string) (uint64, error) {
	val, err := rq.RedisDB.Get(rq.Ctx, key).Result()
	if err != nil {
		return 0, err
	}
	valUint, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		return 0, err
	}
	return valUint, nil
}

func (rq RedisQuerier) GetLactationKafkaOffset() (uint64, error) {
	return rq.getUintKey(LACTATION_LOAD_OFFSET_KEY)
}
