package redis

import (
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/juju/errgo"
)

type RedisClient struct {
	Pool *redis.Pool
}

func NewRedisClient(addr string) *RedisClient {
	return &RedisClient{newPool(addr)}
}

func (r *RedisClient) RemBackend(domain, backendName string) error {
	backendAddr, err := r.get(backendName)
	if err != nil {
		return errgo.Mask(err)
	}

	Verbosef("INFO: remove backendName '%s' with backendAddr '%s' from domain '%s'", backendName, backendAddr, domain)

	if err := r.del(backendName); err != nil {
		return errgo.Mask(err)
	}

	if err := r.srem(domain, backendAddr); err != nil {
		return errgo.Mask(err)
	}

	return nil
}

// AddBackend adds backends to redis. If the backend is already set for this
// domain, we overwrite it using the given value to always have the newest one
// used.
func (r *RedisClient) AddBackend(domain, backendName, backendAddr string) error {
	Verbosef("INFO: add backendName '%s' with backendAddr '%s' to domain '%s'", backendName, backendAddr, domain)

	if err := r.set(backendName, backendAddr); err != nil {
		return errgo.Mask(err)
	}

	if err := r.sadd(domain, backendAddr); err != nil {
		return errgo.Mask(err)
	}

	return nil
}

//------------------------------------------------------------------------------
// private

func newPool(addr string) *redis.Pool {
	return &redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 100,

		// Maximum number of connections allocated by the pool at a given time.
		// When zero, there is no limit on the number of connections in the pool.
		MaxActive: 100,

		// Close connections after remaining idle for this duration. If the value
		// is zero, then idle connections are not closed. Applications should set
		// the timeout to a value less than the server's timeout.
		IdleTimeout: 240 * time.Second,

		// Dial is an application supplied function for creating and configuring a
		// connection
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}

			return c, err
		},
	}
}

func (r *RedisClient) get(key string) (string, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	data, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", errgo.Mask(err)
	}

	return data, nil
}

func (r *RedisClient) set(key, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SET", key, value); err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (r *RedisClient) sadd(key, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SADD", key, value); err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (r *RedisClient) srem(key, value string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("SREM", key, value); err != nil {
		return errgo.Mask(err)
	}

	return nil
}

func (r *RedisClient) exists(key string) (bool, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	data, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false, errgo.Mask(err)
	}

	return data, nil
}

func (r *RedisClient) del(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	if _, err := conn.Do("DEL", key); err != nil {
		return errgo.Mask(err)
	}

	return nil
}
