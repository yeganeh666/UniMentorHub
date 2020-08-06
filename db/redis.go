package db

import (
	"Atrovan_Q1/config"
	"Atrovan_Q1/services/models"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

// context for using in redis.
var ctx = context.Background()

const lessonsKey = "lessons"

// ConnectRedis initialize the redis connection.
func ConnectRedis() (*RedisCli, error) {
	if redisInstance == nil {
		redisInstance = &RedisCli{}
		redis := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", config.EnvGetStr("REDIS_DATABASE_HOST", "localhost"), config.EnvGetInt("REDIS_PORT", 6379)),
			Password: config.EnvGetStr("REDIS_PASSWORD", ""),
			DB:       0, // use default DB
		})
		_, err := redis.Ping(ctx).Result()
		if err != nil {
			return nil, err
		}
		redisInstance.conn = redis
		return redisInstance, nil
	}
	return redisInstance, nil
}

// SetValue set data into redis.
func (redisCli *RedisCli) SetValue(key string, data interface{}, ttl time.Duration) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return redisCli.conn.Set(ctx, key, b, ttl).Err()
}

// GetValue get data from redis.
func (redisCli *RedisCli) GetValue(key string) (interface{}, error) {
	return redisCli.conn.Get(ctx, key).Result()
}

// AddLesson add data to redis.
func (redisCli *RedisCli) AddLesson(lessons []models.Lesson) error {
	for _, l := range lessons {
		response, _ := json.Marshal(&l)
		_, err := redisCli.conn.ZAdd(ctx, lessonsKey, &redis.Z{Score: float64(l.Code), Member: response}).Result()
		if err != nil {
			return err
		}
	}
	return nil
}

// GetLessons get data from redis.
func (redisCli *RedisCli) GetLessons(from, to int) ([]string, error) {
	return redisCli.conn.ZRevRange(ctx, lessonsKey, int64(from), int64(to)).Result()
}
