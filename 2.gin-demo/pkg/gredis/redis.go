package gredis

import (
	"demo/2.gin-demo/pkg/setting"
	"encoding/json"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
)

/*
    @Create by GoLand
    @Author: hong
    @Time: 2018/9/26 17:30 
    @File: redis.go    
*/

var RedisConn *redis.Pool

func Setup() error {
	RedisConn = &redis.Pool{
		MaxIdle:     setting.RedisSetting.MaxIdle,
		MaxActive:   setting.RedisSetting.MaxActive,
		IdleTimeout: setting.RedisSetting.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", setting.RedisSetting.Host)
			if err != nil {
				return nil, err
			}
			if setting.RedisSetting.Password != "" {
				if _, err := c.Do("AUTH", setting.RedisSetting.Password); err != nil {
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
	return nil
}

func Set(key string, data interface{}, time int) error {
	conn := RedisConn.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("redis set -> json marshal fail: [%v]", err)
		return err
	}
	_, err = conn.Do("SET", key, value)
	if err != nil {
		log.Fatalf("redis set -> set fail: [%v]", err)
		return err
	}

	_, err = conn.Do("EXPIRE", key, time)
	if err != nil {
		log.Fatalf("redis set -> expire fail: [%v]", err)
		return err
	}
	return nil
}

func Exists(key string) bool {
	conn := RedisConn.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Fatalf("redis exists -> do fail: %v", err)
		return false
	}
	return exists
}

func Get(key string) ([]byte, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", key))
	if err != nil {
		log.Fatalf("redis get -> do fail: %v", err)
		return nil, err
	}
	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := RedisConn.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := RedisConn.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		log.Fatalf("redis likeDeletes -> do fail: %v", err)
		return err
	}
	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			log.Fatalf("redis likeDeletes -> delete [%v] fail: %v", key, err)
			return err
		}
	}
	return nil
}
