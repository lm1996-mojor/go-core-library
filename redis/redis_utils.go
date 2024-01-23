package redis

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kataras/iris/v12"
	"github.com/lm1996-mojor/go-core-library/log"

	"github.com/redis/go-redis/v9"
)

func RedisSet(ctx iris.Context, key string, value interface{}, expire int) error {
	if expire > 0 {
		err := redisDb.Do(ctx, "SET", key, value, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisDb.Do(ctx, "SET", key, value).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}
func RedisKeyExists(ctx iris.Context, key string) (bool, error) {
	ok, err := redisDb.Do(ctx, "EXISTS", key).Bool()
	return ok, err
}

func RedisGet(ctx iris.Context, key string) (string, error) {
	value, err := redisDb.Do(ctx, "GET", key).Result()
	if err != nil {
		return "", nil
	}

	return value.(string), nil
}

func RedisGetResult(ctx iris.Context, key string) (interface{}, error) {
	v, err := redisDb.Do(ctx, "GET", key).Result()
	if err == redis.Nil {
		return v, nil
	}
	return v, err
}

func RedisGetInt(ctx iris.Context, key string) (int, error) {
	v, err := redisDb.Do(ctx, "GET", key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetInt64(ctx iris.Context, key string) (int64, error) {
	v, err := redisDb.Do(ctx, "GET", key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetUint64(ctx iris.Context, key string) (uint64, error) {
	v, err := redisDb.Do(ctx, "GET", key).Uint64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetFloat64(ctx iris.Context, key string) (float64, error) {
	v, err := redisDb.Do(ctx, "GET", key).Float64()
	if err == redis.Nil {
		return 0.0, nil
	}
	return v, err
}

func RedisExpire(ctx iris.Context, key string, expire int) error {
	err := redisDb.Do(ctx, "EXPIRE", key, expire).Err()
	if err != nil {
		log.Errorf("RedisExpire Error!", key, "Details:", err.Error())
		return err
	}

	return nil
}

func RedisPTTL(ctx iris.Context, key string) (int, error) {
	ttl, err := redisDb.Do(ctx, "PTTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func RedisTTL(ctx iris.Context, key string) (int, error) {
	ttl, err := redisDb.Do(ctx, "TTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func RedisSetJson(ctx iris.Context, key string, value interface{}, expire int) error {
	jsonData, _ := json.Marshal(value)
	if expire > 0 {
		err := redisDb.Do(ctx, "SET", key, jsonData, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisDb.Do(ctx, "SET", key, jsonData).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func RedisGetJson(ctx iris.Context, key string) ([]byte, error) {
	value, err := redisDb.Do(ctx, "GET", key).Result()
	if err != nil {
		return nil, nil
	}

	return []byte(value.(string)), nil
}

func RedisDel(ctx iris.Context, key string) error {
	err := redisDb.Do(ctx, "DEL", key).Err()
	if err != nil {
		log.Errorf("RedisDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func RedisHGet(ctx iris.Context, key, field string) (string, error) {
	value, err := redisDb.Do(ctx, "HGET", key, field).Result()
	if err != nil {
		return "", nil
	}

	return value.(string), nil
}

func RedisHSet(ctx iris.Context, key, field, value string) error {
	err := redisDb.Do(ctx, "HSET", key, field, value).Err()
	if err != nil {
		log.Errorf("RedisHSet Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func RedisHDel(ctx iris.Context, key, field string) error {
	err := redisDb.Do(ctx, "HDEL", key, field).Err()
	if err != nil {
		log.Errorf("RedisHDel Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func RedisZAdd(ctx iris.Context, key, member, score string) error {
	err := redisDb.Do(ctx, "ZADD", key, score, member).Err()
	if err != nil {
		log.Errorf("RedisZAdd Error!", key, "member:", member, "score:", score, "Details:", err.Error())
	}
	return err
}

func RedisZRank(ctx iris.Context, key, member string) (int, error) {
	rank, err := redisDb.Do(ctx, "ZRANK", key, member).Int()
	if err == redis.Nil {
		return -1, nil
	}

	if err != nil {
		log.Errorf("RedisZRank Error!", key, "member:", member, "Details:", err.Error())
		return -1, nil
	}

	return rank, err
}

func RedisZRange(ctx iris.Context, key string, start, stop int) (values []string, err error) {
	values, err = redisDb.ZRange(ctx, key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisZRangeWithScores(ctx iris.Context, key string, start, stop int) (values []redis.Z, err error) {
	values, err = redisDb.ZRangeWithScores(ctx, key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisZRem(ctx iris.Context, key, member string) error {
	err := redisDb.Do(ctx, "ZREM", key, member).Err()
	if err != nil {
		log.Errorf("RedisZRem Error!", key, "member:", member, "Details:", err.Error())
	}
	return err
}

func RedisRPUSH(ctx iris.Context, key string, member string) (err error) {
	err = redisDb.Do(ctx, "RPUSH", key, member).Err()
	if err != nil {
		log.Errorf("RedisRPUSH Error!", key, member, "Details:", err.Error())
		return
	}

	return
}

func RedisBLPOP(ctx iris.Context, timeout time.Duration, keys ...string) (value []string, err error) {
	value, err = redisDb.BLPop(ctx, timeout, keys...).Result()
	if err == redis.Nil {
		err = nil
		return
	}

	if err != nil {
		log.Errorf("BLPop Error!", keys, timeout, "Details:", err.Error())
		return
	}
	return
}

func RedisLLEN(ctx iris.Context, key string) (value int64, err error) {
	value, err = redisDb.LLen(ctx, key).Result()
	if err != nil {
		log.Errorf("RedisLLEN Error!", key, "Details:", err.Error())
		return
	}

	return
}

func RedisLRange(ctx iris.Context, key string, start, stop int) (values []string, err error) {
	values, err = redisDb.LRange(ctx, key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisLRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisKeys(ctx iris.Context, pattern string) (keys []string, err error) {
	keys, err = redisDb.Keys(ctx, pattern).Result()
	if err != nil {
		log.Errorf("RedisKeys Error!", pattern, "Details:", err.Error())
		return
	}

	return
}

// RedisListAllValuesWithPrefix will take in a key prefix and return the value of all the keys that contain that prefix
func RedisListAllValuesWithPrefix(ctx iris.Context, prefix string) (map[string]string, error) {
	// Grab all the keys with the prefix
	keys, err := getKeys(ctx, fmt.Sprintf("%s*", prefix))
	if err != nil {
		return nil, err
	}

	// We will now iterate through all of the values to
	values, err := getKeyAndValuesMap(ctx, keys, prefix)

	return values, nil
}

// getClusterKeys will take a certain prefix that the keys share and return a list of all the keys
func getKeys(ctx iris.Context, prefix string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	count := int64(10) // count specifies how many keys should be returned in every Scan call

	for {
		var keys []string
		var err error
		keys, cursor, err = redisDb.Scan(ctx, cursor, prefix, count).Result()
		if err != nil {

			return nil, errors.New("error retrieving " + prefix + " keys")
		}

		allKeys = append(allKeys, keys...)

		if cursor == 0 {
			break
		}

	}

	return allKeys, nil
}

// getKeyAndValuesMap generates a [string]string map structure that will associate an ID with the token_util value stored in Redis
func getKeyAndValuesMap(ctx iris.Context, keys []string, prefix string) (map[string]string, error) {
	values := make(map[string]string)
	for _, key := range keys {
		value, err := redisDb.Do(ctx, "GET", key).Result()
		if err != nil {
			return nil, errors.New("error retrieving " + prefix + " keys")
		}

		// Strip off the prefix from the key so that we save the key to the user ID
		strippedKey := strings.Split(key, prefix)
		values[strippedKey[1]] = value.(string)
	}

	return values, nil
}

func RedisBatchDel(ctx iris.Context, key ...string) error {
	err := redisDb.Del(ctx, key...).Err()
	if err != nil {
		log.Errorf("RedisBatchDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func RedisMset(ctx iris.Context, pairs ...interface{}) error {
	err := redisDb.MSet(ctx, pairs...).Err()
	if err != nil {
		log.Errorf("RedisMset Error! pairs:", pairs, "Details:", err.Error())
	}
	return err
}

// RedisPublish 消息发布
func RedisPublish(ctx iris.Context, channelName string, message interface{}) int64 {
	result, err := redisDb.Publish(ctx, channelName, message).Result()
	if err != nil {
		log.Error("发布消息至redis中错误：" + err.Error())
		panic(err)
	}
	return result
}

// RedisPSubscribe 匹配获取指定通道中的订阅消息
func RedisPSubscribe(ctx context.Context, channelName string) *redis.PubSub {
	return redisDb.PSubscribe(ctx, channelName)
}
