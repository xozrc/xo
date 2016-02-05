package common

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

type RedisConfig struct {
	Host        string `json:"host"`
	Port        int32  `json:"port"`
	MaxIdle     int32  `json:"max_idle"`
	IdleTimeout int32  `json:"idle_timeout"`
	MaxActive   int32  `json:"max_active"`
}

func NewRedisConf(m json.RawMessage) (conf *RedisConfig, err error) {
	conf = &RedisConfig{}
	err = json.Unmarshal(m, conf)
	return
}

func RedisPoolForCfg(cfg *RedisConfig) (pool *redis.Pool, err error) {
	dataSourceName := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	//connect
	pool = &redis.Pool{
		MaxIdle:     int(cfg.MaxIdle),
		IdleTimeout: 240 * time.Second,

		Dial: func() (redis.Conn, error) {
			fmt.Println("redis connect " + dataSourceName)
			c, err := redis.Dial("tcp", dataSourceName)
			if err != nil {
				panic(err.Error())
			}

			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return

}

func NewRedisService(m json.RawMessage) (rs *RedisService, err error) {
	conf, err1 := NewRedisConf(m)
	if err1 != nil {
		err = err1
		return
	}
	rs = &RedisService{RedisCfg: conf}
	return
}

type RedisService struct {
	RedisCfg  *RedisConfig
	RedisPool *redis.Pool
}

func (r *RedisService) Init() {
	pool, err := RedisPoolForCfg(r.RedisCfg)
	dataSourceName := fmt.Sprintf("%s:%d", r.RedisCfg.Host, r.RedisCfg.Port)
	if err != nil {
		panic("init redis service failed:" + err.Error())
	}
	r.RedisPool = pool
	logger.Printf("connect  %s success\n", dataSourceName)

}

func (r *RedisService) AfterInit() {

}

func (r *RedisService) BeforeDestroy() {

}

func (r *RedisService) Destroy() {

}
