package core

import (
	"context"
	"strings"
	"time"

	"siteLetterJob/config"
	"siteLetterJob/mdata"
	"siteLetterJob/utils"

	"github.com/go-redis/redis/v8"
)

var redisDb *RedisDb

const RedisNil = redis.Nil

type RedisDb struct {
	wPool redis.UniversalClient
}

func InitRedis() error {
	rw := config.GetConfig().RedisCore
	rwPool, err := initRedis(rw.Host, rw.Auth, rw.Master, rw.PoolSize)
	if err != nil {
		return err
	}

	redisDb = &RedisDb{
		wPool: rwPool,
	}
	return nil
}

func initRedis(host, auth, master string, poolSize int) (redis.UniversalClient, error) {
	auth = utils.GetRealString(config.GetConfig().DBSecretKey, auth)
	options := &redis.UniversalOptions{
		Addrs:              strings.Split(host, ","), // redis地址
		MaxRedirects:       0,                        // 放弃前最大重试次数,默认是不重试失败的命令,默认是3次
		ReadOnly:           false,                    // 在从库上打开只读命令
		RouteByLatency:     false,                    // 允许将只读命令路由到最近的主节点或从节点,自动启用只读
		RouteRandomly:      false,                    // 允许将只读命令路由到随机主节点或从节点。 它自动启用只读。
		Password:           auth,
		MaxRetries:         2,
		MinRetryBackoff:    8 * time.Millisecond,
		MaxRetryBackoff:    512 * time.Millisecond,
		DialTimeout:        5 * time.Second,
		ReadTimeout:        10 * time.Second,
		WriteTimeout:       20 * time.Second,
		PoolSize:           poolSize,
		MaxConnAge:         6 * time.Minute,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        5 * time.Minute,
		IdleCheckFrequency: 1 * time.Minute, //空闲连接检查频率
	}
	//哨兵模式
	if len(master) > 0 {
		options.SentinelPassword = auth
		options.MasterName = master
	}
	redisPool := redis.NewUniversalClient(options)
	_, err := redisPool.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	return redisPool, nil
}

func GetKey(key string) (string, error) {
	value, err := redisDb.wPool.Get(context.Background(), key).Result()
	if err != nil && err != redis.Nil {
		return "", err
	}
	return value, nil
}

func GetKeyBytes(key string) ([]byte, error) {
	return redisDb.wPool.Get(context.Background(), key).Bytes()
}

// SetNotExpireKV 设置不过期的 key
func SetNotExpireKV(key, value string) error {
	err := redisDb.wPool.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetNotExpireKVInterface 设置不过期的 key
func SetNotExpireKVInterface(key string, value interface{}) error {
	err := redisDb.wPool.Set(context.Background(), key, value, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetExpireKV 设置过期的 key
func SetExpireKV(key, value string, expire time.Duration) error {
	err := redisDb.wPool.Set(context.Background(), key, value, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

func GetTTL(key string) time.Duration {
	return redisDb.wPool.TTL(context.Background(), key).Val()
}

// SetExpireKey 设置 key 过期
func SetExpireKey(key string, expire time.Duration) error {
	err := redisDb.wPool.Expire(context.Background(), key, expire).Err()
	if err != nil {
		return err
	}
	return nil
}

// SetNX 设置 key, value 以及过期时间
func SetNX(key string, value string, expire time.Duration) (bool, error) {
	flag, err := redisDb.wPool.SetNX(context.Background(), key, value, expire).Result()
	if err != nil {
		return false, err
	}
	return flag, nil
}

// DelKey 删除 redis 的key
func DelKey(key string) error {
	return redisDb.wPool.Del(context.Background(), key).Err()
}

// KeyExist 判断某一个key 是否存在
func KeyExist(keys string) (bool, error) {

	count, err := redisDb.wPool.Exists(context.Background(), keys).Result()
	if err != nil {
		return false, err
	}
	if count > 0 {
		return true, nil
	}
	return false, nil
}

// HSet 设置 hash
func HSet(key, field string, value interface{}) error {
	return redisDb.wPool.HSet(context.Background(), key, field, value).Err()
}

// HMSet 批量存储 hash
func HMSet(key string, fields map[string]interface{}) error {
	if len(fields) < 1 {
		return nil
	}

	err := redisDb.wPool.HMSet(context.Background(), key, fields).Err()
	if err != nil {
		return err
	}

	return nil
}

// HGet 获取单个 hash
func HGet(key, field string) (string, error) {
	return redisDb.wPool.HGet(context.Background(), key, field).Result()
}

func HKeys(key string) ([]string, error) {
	return redisDb.wPool.HKeys(context.Background(), key).Result()
}

// HMGet 批量获取 hash
func HMGet(key string, fields ...string) ([]interface{}, error) {
	res, err := redisDb.wPool.HMGet(context.Background(), key, fields...).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HGetAll 获取 hash 全部值
func HGetAlL(key string) (map[string]string, error) {
	res, err := redisDb.wPool.HGetAll(context.Background(), key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

// HScan 获取 hash 键值树
func HScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return redisDb.wPool.HScan(context.Background(), key, cursor, match, count).Result()
}

func SScan(key string, cursor uint64, match string, count int64) ([]string, uint64, error) {
	return redisDb.wPool.SScan(context.Background(), key, cursor, match, count).Result()
}

func HLen(key string) (int, error) {
	res, err := redisDb.wPool.HLen(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return int(res), nil
}

// HDel 删除 hash key
func HDel(key string, fields ...string) error {
	err := redisDb.wPool.HDel(context.Background(), key, fields...).Err()
	if err != nil {
		return err
	}

	return nil
}

// RPush 在名称为key的list尾添加一个值为value的元素
func RPush(key string, values ...interface{}) error {
	return redisDb.wPool.RPush(context.Background(), key, values...).Err()
}

// LPush 在名称为key的list头添加一个值为value的 元素
func LPush(key string, values ...interface{}) error {
	return redisDb.wPool.LPush(context.Background(), key, values...).Err()
}

// Publish 在名称为key的list头添加一个值为value的 元素
func Publish(channel string, values interface{}) error {
	return redisDb.wPool.Publish(context.Background(), channel, values).Err()
}

// LLen 返回名称为key的list的长度
func LLen(key string) (int64, error) {
	return redisDb.wPool.LLen(context.Background(), key).Result()
}

// LRange 返回名称为key的list中start至end之间的元素, start为0, end为-1 则是获取所有 list key
func LRange(key string, start, end int64) ([]string, error) {
	return redisDb.wPool.LRange(context.Background(), key, start, end).Result()
}

// LSet 给名称为key的list中index位置的元素赋值
func LSet(key string, index int64, value interface{}) error {
	return redisDb.wPool.LSet(context.Background(), key, index, value).Err()
}

// LRem 删除count个key的list中值为value的元素
func LRem(key string, count int64, value interface{}) error {
	return redisDb.wPool.LRem(context.Background(), key, count, value).Err()
}

// ZAdd 有序集合中增加一个成员
func ZAdd(key, member string, score float64) error {
	z := redis.Z{
		Score:  score,
		Member: member,
	}
	_, err := redisDb.wPool.ZAdd(context.Background(), key, &z).Result()
	if err != nil {
		return err
	}

	return nil
}

// ZCount  有序集合中 min-max中的成员数量
func ZCount(key, min, max string) (int64, error) {
	count, err := redisDb.wPool.ZCount(context.Background(), key, min, max).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZCARD 获取中元素的数量
func ZCARD(key string) (int64, error) {
	count, err := redisDb.wPool.ZCard(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ZRange 通过索引区间返回有序集合成指定区间内的成员
func ZRange(key string, start, stop int64) ([]string, error) {
	arr, err := redisDb.wPool.ZRange(context.Background(), key, start, stop).Result()
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

// ZRangeByScore 通过索引区间返回有序集合成指定区间内的成员
func ZRangeByScore(key string, min, max string) ([]string, error) {
	opt := redis.ZRangeBy{
		Min: min,
		Max: max,
	}
	arr, err := redisDb.wPool.ZRangeByScore(context.Background(), key, &opt).Result()
	if err != nil {
		return []string{}, err
	}

	return arr, nil
}

func ZRem(key string, members ...string) error {
	return redisDb.wPool.ZRem(context.Background(), key, members).Err()
}

func HGetBytesByField(key, filed string) ([]byte, error) {
	return redisDb.wPool.HGet(context.Background(), key, filed).Bytes()
}

func SIsMember(key, field string) (bool, error) {
	return redisDb.wPool.SIsMember(context.Background(), key, field).Result()
}
func Incr(key string) {
	redisDb.wPool.Incr(context.Background(), key)
}

func IncrBy(key string, value int64) (int64, error) {
	return redisDb.wPool.IncrBy(context.Background(), key, value).Result()
}

func IncrWithResult(key string) (int64, error) {
	return redisDb.wPool.Incr(context.Background(), key).Result()
}

func DecrWithResult(key string) (int64, error) {
	return redisDb.wPool.Decr(context.Background(), key).Result()
}

func SMembers(key string) ([]string, error) {
	return redisDb.wPool.SMembers(context.Background(), key).Result()
}

func SAdd(key string, members ...interface{}) (int64, error) {
	return redisDb.wPool.SAdd(context.Background(), key, members...).Result()
}

func SRem(key string, members ...interface{}) (int64, error) {
	return redisDb.wPool.SRem(context.Background(), key, members...).Result()
}

func LIndex(key string, index int64) (string, error) {
	return redisDb.wPool.LIndex(context.Background(), key, index).Result()
}

// SetSscan 集合读取
func SetSscan(key string, match string, perCount int64) ([]string, error) {
	var (
		cursor = uint64(0)
		data   []string
	)
	for {
		keys, retCursor, err := redisDb.wPool.SScan(context.Background(), key, cursor, match, perCount).Result()
		if err != nil {
			return data, err
		}
		if len(keys) == 0 {
			break
		}
		data = append(data, keys...)
		if retCursor == 0 {
			break
		}
		cursor = retCursor
	}
	return data, nil
}

// 获取键值,如不存在 则获取func 存入到键中
func GetOrSet(key string, f func() (interface{}, error), expire time.Duration) ([]byte, error) {
	result, err := redisDb.wPool.Get(context.Background(), key).Bytes()
	if err != nil || len(result) == 0 {
		data, err := f()
		if err == nil {
			var value []byte
			value, err = mdata.Cjson.Marshal(data)
			if err != nil {
				return nil, err
			}
			err = redisDb.wPool.Set(context.Background(), key, value, expire).Err()
			if err != nil {
				return nil, err
			}
			return value, nil
		}
		return nil, err
	}
	return result, nil
}
