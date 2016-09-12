package services

import (
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/garyburd/redigo/redis"
)

var (
	RedisPool *redis.Pool
)

func PanicOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%v: %v", msg, err)
		panic(fmt.Sprintf("%v: %v", msg, err))
	}
}

func newRedisPool(server, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func RedisSetup() {
	serverUrl := os.Getenv("REDIS_URL")
	password := os.Getenv("REDIS_PASSWORD")
	RedisPool = newRedisPool(serverUrl, password)
	c := RedisPool.Get()
	defer c.Close()

	pong, err := redis.String(c.Do("PING"))
	PanicOnError(err, "Cannot ping Redis")
	log.Infof("Redis Ping: %s", pong)
}
