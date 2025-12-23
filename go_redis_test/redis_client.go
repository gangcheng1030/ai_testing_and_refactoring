package go_redis_test

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// RedisClient Redis 客户端结构体
type RedisClient struct {
	client *redis.Client
}

// RedisConfig Redis 配置
type RedisConfig struct {
	Addr     string // Redis 地址，格式: host:port
	Password string // Redis 密码，如果没有密码则为空字符串
	DB       int    // Redis 数据库编号，默认为 0
}

// NewRedisClient 创建一个新的 Redis 客户端
// 参数:
//
//	config: Redis 配置信息
//
// 返回:
//
//	RedisClient 实例
func NewRedisClient(config RedisConfig) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisClient{
		client: rdb,
	}
}

// Get 从 Redis 获取指定 key 的值
// 参数:
//
//	ctx: 上下文对象，用于超时控制和取消操作
//	key: 要获取的键名
//
// 返回:
//
//	string: 键对应的值
//	error: 如果键不存在，返回 redis.Nil 错误；其他错误返回相应的错误信息
func (rc *RedisClient) Get(ctx context.Context, key string) (string, error) {
	return rc.client.Get(ctx, key).Result()
}

// Set 向 Redis 设置键值对
// 参数:
//
//	ctx: 上下文对象
//	key: 键名
//	value: 值
//	expiration: 过期时间，0 表示永不过期
//
// 返回:
//
//	error: 设置失败时返回错误
func (rc *RedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return rc.client.Set(ctx, key, value, expiration).Err()
}

// Del 删除指定的 key
// 参数:
//
//	ctx: 上下文对象
//	keys: 要删除的键名列表
//
// 返回:
//
//	int64: 实际删除的键数量
//	error: 删除失败时返回错误
func (rc *RedisClient) Del(ctx context.Context, keys ...string) (int64, error) {
	return rc.client.Del(ctx, keys...).Result()
}

// Exists 检查 key 是否存在
// 参数:
//
//	ctx: 上下文对象
//	keys: 要检查的键名列表
//
// 返回:
//
//	int64: 存在的键数量
//	error: 检查失败时返回错误
func (rc *RedisClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return rc.client.Exists(ctx, keys...).Result()
}

// Close 关闭 Redis 连接
// 返回:
//
//	error: 关闭失败时返回错误
func (rc *RedisClient) Close() error {
	return rc.client.Close()
}

// Ping 测试 Redis 连接是否正常
// 参数:
//
//	ctx: 上下文对象
//
// 返回:
//
//	error: 连接异常时返回错误，正常返回 nil
func (rc *RedisClient) Ping(ctx context.Context) error {
	return rc.client.Ping(ctx).Err()
}
