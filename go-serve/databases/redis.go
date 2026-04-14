package databases

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端，其他包通过 databases.RDB 使用
var RDB *redis.Client

// mustGetEnv 获取环境变量，缺失时终止启动
func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("关键环境变量缺失: %s", key)
	}
	return value
}

// mustGetEnvInt 获取整数类型环境变量，缺失或格式错误时终止启动
func mustGetEnvInt(key string) int {
	value := mustGetEnv(key)
	intVal, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("环境变量 %s 不是有效的整数: %v", key, err)
	}
	return intVal
}

// InitRedis 初始化 Redis 连接，启动时调用一次
func InitRedis() {
	// 第一步：从环境变量读取 Redis 配置
	addr := mustGetEnv("REDIS_ADDR")
	password := os.Getenv("REDIS_PASSWORD") // 密码允许为空
	db := mustGetEnvInt("REDIS_DB")
	poolSize := mustGetEnvInt("REDIS_POOL_SIZE")
	minIdleConns := mustGetEnvInt("REDIS_MIN_IDLE_CONNS")

	// 第二步：创建客户端 + 连接池配置
	RDB = redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     password,
		DB:           db,
		PoolSize:     poolSize,
		MinIdleConns: minIdleConns,
	})

	// 第三步：主动 Ping 检验连接（快速失败原则）
	ctx := context.Background()
	Result, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}
	fmt.Println("redis 连接成功:" + Result)
}
