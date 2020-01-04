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
	err1 := Uid.LoadH24FromRedis(newClient, "UUID:UID:24")
	err2 := Uid.LoadH32FromRedis(newClient, "UUID:UID:32")
	if err1 != nil {
		log.Error("初始化UUID错误:%v",err1.Error())
		return err1
	}
	if err2 != nil {
		log.Error("初始化UUID错误:%v",err2.Error())
		return err2
	}
	return nil
}

func GenUid() uint64 {
	return Uid.Next()
}

func GenStringUUID() string {
	return strconv.FormatUint(GenUid(), 10)
}
