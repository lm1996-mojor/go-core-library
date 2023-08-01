package redis

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lm1996-mojor/go-core-library/log"

	"github.com/go-redis/redis"
)

func RedisSet(key string, value interface{}, expire int) error {
	if expire > 0 {
		err := redisSingleDb.Do("SET", key, value, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisSingleDb.Do("SET", key, value).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}
func RedisKeyExists(key string) (bool, error) {
	ok, err := redisSingleDb.Do("EXISTS", key).Bool()
	return ok, err
}

func RedisGet(key string) (string, error) {
	value, err := redisSingleDb.Do("GET", key).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func RedisGetResult(key string) (interface{}, error) {
	v, err := redisSingleDb.Do("GET", key).Result()
	if err == redis.Nil {
		return v, nil
	}
	return v, err
}

func RedisGetInt(key string) (int, error) {
	v, err := redisSingleDb.Do("GET", key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetInt64(key string) (int64, error) {
	v, err := redisSingleDb.Do("GET", key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetUint64(key string) (uint64, error) {
	v, err := redisSingleDb.Do("GET", key).Uint64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func RedisGetFloat64(key string) (float64, error) {
	v, err := redisSingleDb.Do("GET", key).Float64()
	if err == redis.Nil {
		return 0.0, nil
	}
	return v, err
}

func RedisExpire(key string, expire int) error {
	err := redisSingleDb.Do("EXPIRE", key, expire).Err()
	if err != nil {
		log.Errorf("RedisExpire Error!", key, "Details:", err.Error())
		return err
	}

	return nil
}

func RedisPTTL(key string) (int, error) {
	ttl, err := redisSingleDb.Do("PTTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func RedisTTL(key string) (int, error) {
	ttl, err := redisSingleDb.Do("TTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func RedisSetJson(key string, value interface{}, expire int) error {
	jsonData, _ := json.Marshal(value)
	if expire > 0 {
		err := redisSingleDb.Do("SET", key, jsonData, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisSingleDb.Do("SET", key, jsonData).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func RedisGetJson(key string) ([]byte, error) {
	value, err := redisSingleDb.Do("GET", key).String()
	if err != nil {
		return nil, nil
	}

	return []byte(value), nil
}

func RedisDel(key string) error {
	err := redisSingleDb.Do("DEL", key).Err()
	if err != nil {
		log.Errorf("RedisDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func RedisHGet(key, field string) (string, error) {
	value, err := redisSingleDb.Do("HGET", key, field).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func RedisHSet(key, field, value string) error {
	err := redisSingleDb.Do("HSET", key, field, value).Err()
	if err != nil {
		log.Errorf("RedisHSet Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func RedisHDel(key, field string) error {
	err := redisSingleDb.Do("HDEL", key, field).Err()
	if err != nil {
		log.Errorf("RedisHDel Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func RedisZAdd(key, member, score string) error {
	err := redisSingleDb.Do("ZADD", key, score, member).Err()
	if err != nil {
		log.Errorf("RedisZAdd Error!", key, "member:", member, "score:", score, "Details:", err.Error())
	}
	return err
}

func RedisZRank(key, member string) (int, error) {
	rank, err := redisSingleDb.Do("ZRANK", key, member).Int()
	if err == redis.Nil {
		return -1, nil
	}

	if err != nil {
		log.Errorf("RedisZRank Error!", key, "member:", member, "Details:", err.Error())
		return -1, nil
	}

	return rank, err
}

func RedisZRange(key string, start, stop int) (values []string, err error) {
	values, err = redisSingleDb.ZRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisZRangeWithScores(key string, start, stop int) (values []redis.Z, err error) {
	values, err = redisSingleDb.ZRangeWithScores(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisZRem(key, member string) error {
	err := redisSingleDb.Do("ZREM", key, member).Err()
	if err != nil {
		log.Errorf("RedisZRem Error!", key, "member:", member, "Details:", err.Error())
	}
	return err
}

func RedisRPUSH(key string, member string) (err error) {
	err = redisSingleDb.Do("RPUSH", key, member).Err()
	if err != nil {
		log.Errorf("RedisRPUSH Error!", key, member, "Details:", err.Error())
		return
	}

	return
}

func RedisBLPOP(timeout time.Duration, keys ...string) (value []string, err error) {
	value, err = redisSingleDb.BLPop(timeout, keys...).Result()
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

func RedisLLEN(key string) (value int64, err error) {
	value, err = redisSingleDb.LLen(key).Result()
	if err != nil {
		log.Errorf("RedisLLEN Error!", key, "Details:", err.Error())
		return
	}

	return
}

func RedisLRange(key string, start, stop int) (values []string, err error) {
	values, err = redisSingleDb.LRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisLRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func RedisKeys(pattern string) (keys []string, err error) {
	keys, err = redisSingleDb.Keys(pattern).Result()
	if err != nil {
		log.Errorf("RedisKeys Error!", pattern, "Details:", err.Error())
		return
	}

	return
}

// RedisListAllValuesWithPrefix will take in a key prefix and return the value of all the keys that contain that prefix
func RedisListAllValuesWithPrefix(prefix string) (map[string]string, error) {
	// Grab all the keys with the prefix
	keys, err := getKeys(fmt.Sprintf("%s*", prefix))
	if err != nil {
		return nil, err
	}

	// We will now iterate through all of the values to
	values, err := getKeyAndValuesMap(keys, prefix)

	return values, nil
}

// getClusterKeys will take a certain prefix that the keys share and return a list of all the keys
func getKeys(prefix string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	count := int64(10) // count specifies how many keys should be returned in every Scan call

	for {
		var keys []string
		var err error
		keys, cursor, err = redisSingleDb.Scan(cursor, prefix, count).Result()
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
func getKeyAndValuesMap(keys []string, prefix string) (map[string]string, error) {
	values := make(map[string]string)
	for _, key := range keys {
		value, err := redisSingleDb.Do("GET", key).String()
		if err != nil {
			return nil, errors.New("error retrieving " + prefix + " keys")
		}

		// Strip off the prefix from the key so that we save the key to the user ID
		strippedKey := strings.Split(key, prefix)
		values[strippedKey[1]] = value
	}

	return values, nil
}

func RedisBatchDel(key ...string) error {
	err := redisSingleDb.Del(key...).Err()
	if err != nil {
		log.Errorf("RedisBatchDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func RedisMset(pairs ...interface{}) error {
	err := redisSingleDb.MSet(pairs...).Err()
	if err != nil {
		log.Errorf("RedisMset Error! pairs:", pairs, "Details:", err.Error())
	}
	return err
}
