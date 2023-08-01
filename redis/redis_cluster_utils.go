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

func ClusterRedisSet(key string, value interface{}, expire int) error {
	if expire > 0 {

		err := redisClusterDb.Do("SET", key, value, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisClusterDb.Do("SET", key, value).Err()
		if err != nil {
			log.Errorf("RedisSet Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}
func ClusterRedisKeyExists(key string) (bool, error) {
	ok, err := redisClusterDb.Do("EXISTS", key).Bool()
	return ok, err
}

func ClusterRedisGet(key string) (string, error) {
	value, err := redisClusterDb.Do("GET", key).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func ClusterRedisGetResult(key string) (interface{}, error) {
	v, err := redisClusterDb.Do("GET", key).Result()
	if err == redis.Nil {
		return v, nil
	}
	return v, err
}

func ClusterRedisGetInt(key string) (int, error) {
	v, err := redisClusterDb.Do("GET", key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func ClusterRedisGetInt64(key string) (int64, error) {
	v, err := redisClusterDb.Do("GET", key).Int64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func ClusterRedisGetUint64(key string) (uint64, error) {
	v, err := redisClusterDb.Do("GET", key).Uint64()
	if err == redis.Nil {
		return 0, nil
	}
	return v, err
}

func ClusterRedisGetFloat64(key string) (float64, error) {
	v, err := redisClusterDb.Do("GET", key).Float64()
	if err == redis.Nil {
		return 0.0, nil
	}
	return v, err
}

func ClusterRedisExpire(key string, expire int) error {
	err := redisClusterDb.Do("EXPIRE", key, expire).Err()
	if err != nil {
		log.Errorf("RedisExpire Error!", key, "Details:", err.Error())
		return err
	}

	return nil
}

func ClusterRedisPTTL(key string) (int, error) {
	ttl, err := redisClusterDb.Do("PTTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func ClusterRedisTTL(key string) (int, error) {
	ttl, err := redisClusterDb.Do("TTL", key).Int()
	if err != nil {
		return -1, err
	}

	return ttl, nil
}

func ClusterRedisSetJson(key string, value interface{}, expire int) error {
	jsonData, _ := json.Marshal(value)
	if expire > 0 {
		err := redisClusterDb.Do("SET", key, jsonData, "PX", expire).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	} else {
		err := redisClusterDb.Do("SET", key, jsonData).Err()
		if err != nil {
			log.Errorf("RedisSetJson Error! key:", key, "Details:", err.Error())
			return err
		}
	}

	return nil
}

func ClusterRedisGetJson(key string) ([]byte, error) {
	value, err := redisClusterDb.Do("GET", key).String()
	if err != nil {
		return nil, nil
	}

	return []byte(value), nil
}

func ClusterRedisDel(key string) error {
	err := redisClusterDb.Do("DEL", key).Err()
	if err != nil {
		log.Errorf("RedisDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func ClusterRedisHGet(key, field string) (string, error) {
	value, err := redisClusterDb.Do("HGET", key, field).String()
	if err != nil {
		return "", nil
	}

	return value, nil
}

func ClusterRedisHSet(key, field, value string) error {
	err := redisClusterDb.Do("HSET", key, field, value).Err()
	if err != nil {
		log.Errorf("RedisHSet Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func ClusterRedisHDel(key, field string) error {
	err := redisClusterDb.Do("HDEL", key, field).Err()
	if err != nil {
		log.Errorf("RedisHDel Error!", key, "field:", field, "Details:", err.Error())
	}
	return err
}

func ClusterRedisZAdd(key, member, score string) error {
	err := redisClusterDb.Do("ZADD", key, score, member).Err()
	if err != nil {
		log.Errorf("RedisZAdd Error!", key, "member:", member, "score:", score, "Details:", err.Error())
	}
	return err
}

func ClusterRedisZRank(key, member string) (int, error) {
	rank, err := redisClusterDb.Do("ZRANK", key, member).Int()
	if err == redis.Nil {
		return -1, nil
	}

	if err != nil {
		log.Errorf("RedisZRank Error!", key, "member:", member, "Details:", err.Error())
		return -1, nil
	}

	return rank, err
}

func ClusterRedisZRange(key string, start, stop int) (values []string, err error) {
	values, err = redisClusterDb.ZRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func ClusterRedisZRangeWithScores(key string, start, stop int) (values []redis.Z, err error) {
	values, err = redisClusterDb.ZRangeWithScores(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisZRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func ClusterRedisZRem(key, member string) error {
	err := redisClusterDb.Do("ZREM", key, member).Err()
	if err != nil {
		log.Errorf("RedisZRem Error!", key, "member:", member, "Details:", err.Error())
	}
	return err
}

func ClusterRedisRPUSH(key string, member string) (err error) {
	err = redisClusterDb.Do("RPUSH", key, member).Err()
	if err != nil {
		log.Errorf("RedisRPUSH Error!", key, member, "Details:", err.Error())
		return
	}

	return
}

func ClusterRedisBLPOP(timeout time.Duration, keys ...string) (value []string, err error) {
	value, err = redisClusterDb.BLPop(timeout, keys...).Result()
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

func ClusterRedisLLEN(key string) (value int64, err error) {
	value, err = redisClusterDb.LLen(key).Result()
	if err != nil {
		log.Errorf("RedisLLEN Error!", key, "Details:", err.Error())
		return
	}

	return
}

func ClusterRedisLRange(key string, start, stop int) (values []string, err error) {
	values, err = redisClusterDb.LRange(key, int64(start), int64(stop)).Result()
	if err != nil {
		log.Errorf("RedisLRange Error!", key, "start:", start, "stop:", stop, "Details:", err.Error())
		return
	}

	return
}

func ClusterRedisKeys(pattern string) (keys []string, err error) {
	keys, err = redisClusterDb.Keys(pattern).Result()
	if err != nil {
		log.Errorf("RedisKeys Error!", pattern, "Details:", err.Error())
		return
	}

	return
}

// ClusterRedisListAllValuesWithPrefix will take in a key prefix and return the value of all the keys that contain that prefix
func ClusterRedisListAllValuesWithPrefix(prefix string) (map[string]string, error) {
	// Grab all the keys with the prefix
	keys, err := getClusterKeys(fmt.Sprintf("%s*", prefix))
	if err != nil {
		return nil, err
	}

	// We will now iterate through all of the values to
	values, err := getClusterKeyAndValuesMap(keys, prefix)

	return values, nil
}

// getClusterKeys will take a certain prefix that the keys share and return a list of all the keys
func getClusterKeys(prefix string) ([]string, error) {
	var allKeys []string
	var cursor uint64
	count := int64(10) // count specifies how many keys should be returned in every Scan call

	for {
		var keys []string
		var err error
		keys, cursor, err = redisClusterDb.Scan(cursor, prefix, count).Result()
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

// getClusterKeyAndValuesMap generates a [string]string map structure that will associate an ID with the token_util value stored in Redis
func getClusterKeyAndValuesMap(keys []string, prefix string) (map[string]string, error) {
	values := make(map[string]string)
	for _, key := range keys {
		value, err := redisClusterDb.Do("GET", key).String()
		if err != nil {
			return nil, errors.New("error retrieving " + prefix + " keys")
		}

		// Strip off the prefix from the key so that we save the key to the user ID
		strippedKey := strings.Split(key, prefix)
		values[strippedKey[1]] = value
	}

	return values, nil
}

func ClusterRedisBatchDel(key ...string) error {
	err := redisClusterDb.Del(key...).Err()
	if err != nil {
		log.Errorf("RedisBatchDel Error! key:", key, "Details:", err.Error())
	}
	return err
}

func ClusterRedisMset(pairs ...interface{}) error {
	err := redisClusterDb.MSet(pairs...).Err()
	if err != nil {
		log.Errorf("RedisMset Error! pairs:", pairs, "Details:", err.Error())
	}
	return err
}
