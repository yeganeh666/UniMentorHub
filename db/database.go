package db

import (
	"log"

	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
)

// RedisCli for creating instance of redis.
type RedisCli struct {
	conn *redis.Client
}

// GormCli for creating instance of gorm(postgres).
type GormCli struct {
	conn *gorm.DB
}

var redisInstance *RedisCli = nil

var gormInstance *GormCli = nil

// InitConnections initialize Database Connection
func InitConnections() {
	_, err := connectGorm()
	if err != nil {
		log.Fatalf("postgres connection error %v", err)
	}

	_, err = ConnectRedis()
	if err != nil {
		log.Fatalf("redis connection error %v", err)
	}
}
