package utils

import (
	"github.com/go-redis/redis"
	"github.com/hero1s/gotools/log"
	"github.com/hero1s/gotools/utils/uuid"
	"strconv"
)

var (
	Uid *uuid.UUID
)

func InitUUID(redisHost, password string) error {
	newClient := func() (redis.Cmdable, bool, error) {
		return redis.NewClient(&redis.Options{
			Addr:     redisHost,
			Password: password,
		}), true, nil
	}
	Uid = uuid.NewUUID("uid")
	err := Uid.LoadH24FromRedis(newClient, "UUID:UID:24")
	if err != nil {
		log.Error("初始化UUID错误:%v",err.Error())
		return err
	}
	Uid.Renew32 = func() error {//只设置函数不加计数
		return Uid.LoadH32FromRedis(newClient, "UUID:UID:32")
	}

	return nil
}

func GenUid() uint64 {
	return Uid.Next()
}

func GenStringUUID() string {
	return strconv.FormatUint(GenUid(), 10)
}
