package database

import (
	"context"
	"crypto/tls"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Host      string
	Port      int
	Password  string
	DB        int
	TLSConfig *tls.Config
}

func LoadRedisConfig() RedisConfig {
	host := os.Getenv("REDIS_HOST")
	if host == "" {
		host = "localhost"
	}
	portStr := os.Getenv("REDIS_PORT")
	if portStr == "" {
		if os.Getenv("REDIS_TLS") == "true" {
			portStr = "6380"
		} else {
			portStr = "6379"
		}
	}
	port, _ := strconv.Atoi(portStr)

	password := os.Getenv("REDIS_PASSWORD")

	dbStr := os.Getenv("REDIS_DB")
	db, err := strconv.Atoi(dbStr)
	if err != nil || db < 0 {
		db = 0
	}

	var tlsCfg *tls.Config
	if os.Getenv("REDIS_TLS") == "true" {
		tlsCfg = &tls.Config{
			MinVersion: tls.VersionTLS12,
		}
	}

	return RedisConfig{Host: host, Port: port, Password: password, DB: db, TLSConfig: tlsCfg}
}

var Ctx = context.Background()
var Rdb *redis.Client

func ConnectRedis() {
	config := LoadRedisConfig()

	addr := config.Host + ":" + strconv.Itoa(config.Port)

	opt := &redis.Options{
		Addr:      addr,
		Password:  config.Password,
		DB:        config.DB,
		TLSConfig: config.TLSConfig,
	}

	Rdb = redis.NewClient(opt)

	_, err := Rdb.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	log.Println("Connected to Redis successfully")
}
