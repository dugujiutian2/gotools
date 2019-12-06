package common

import (
	"fmt"
	"github.com/hero1s/gotools/cache"
	"github.com/hero1s/gotools/log"
	"math/rand"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"
)

type AccessLimitConf struct {
	Frequency  int64 `json:"frequency"`
	ExpireTime int64 `json:"expire_time"`
}

var Al AccessLimitConf //访问限速配置

func FilterEmoji(content string) string {
	newContent := ""
	for _, value := range content {
		_, size := utf8.DecodeRuneInString(string(value))
		if size <= 3 {
			newContent += string(value)
		} else {
			newContent += "*"
		}
	}
	return newContent
}

//判断是否为纯数字
func IsNumber(number string) bool {
	for _, v := range number {
		if '9' < v || v < '0' {
			return false
		}
	}
	return true
}

func RandomString(length int64) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyz"
	var result string
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var i int64
	for i = 0; i < length; i++ {
		result = result + string(str[r.Intn(len(str))])
	}
	return result
}

func RandomNum() int64 {
	num1 := []int64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := rand.New(rand.NewSource(time.Now().UnixNano() + rand.Int63n(10000)*5))
	return num1[r.Intn(len(num1))]
}

//根据时间戳,获取星座
func Constellation(tt int64) string {
	t := time.Unix(tt, 0).Format("0102")
	d, _ := strconv.ParseInt(t, 10, 64)
	if d >= 321 && d <= 419 {
		return "白羊座"
	}
	if d >= 420 && d <= 520 {
		return "金牛座"
	}
	if d >= 521 && d <= 621 {
		return "双子座"
	}
	if d >= 622 && d <= 722 {
		return "巨蟹座"
	}
	if d >= 723 && d <= 822 {
		return "狮子座"
	}
	if d >= 823 && d <= 922 {
		return "处女座"
	}
	if d >= 923 && d <= 1023 {
		return "天秤座"
	}
	if d >= 1024 && d <= 1122 {
		return "天蝎座"
	}

	if d >= 1123 && d <= 1221 {
		return "射手座"
	}
	if d >= 1222 || d <= 119 {
		return "魔羯座"
	}
	if d >= 120 && d <= 218 {
		return "水平座"
	}
	if d >= 219 && d <= 320 {
		return "双鱼座"
	}

	return "水平座"
}

// 过来结构体空字段，转换json字段的map
func ChangeStructPointToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		if !f.IsNil() {
			data[t.Field(i).Tag.Get("json")] = f.Interface()
		}
	}
	return data
}
func ChangeStructToJsonMap(p interface{}) map[string]interface{} {
	data := map[string]interface{}{}
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	count := v.NumField()
	for i := 0; i < count; i++ {
		f := v.Field(i)
		data[t.Field(i).Tag.Get("json")] = f.Interface()
	}
	return data
}

/*
 *desc:用于访问次数限制
 *@key:需要以什么来做标识做访问次数限制的标志
 *@frequency: 次数
 *@expireTime:多少秒超时
 */
func AccessLimit(cache cache.Cache, key string, frequency int64, expireTime int64) bool {
	ok := cache.IsExist(key)
	if !ok { //doesn't exist, set a key and expire time
		cache.Put(key, 1, time.Duration(expireTime)*time.Second)
		return true
	}

	f := cache.GetInt64(key)
	if f < frequency {
		cache.Incr(key)
		return true
	}
	log.Debug(fmt.Sprintf("触及访问限速,key:%v,frequency:%v,expireTime:%v", key, frequency, expireTime))
	return false
}
