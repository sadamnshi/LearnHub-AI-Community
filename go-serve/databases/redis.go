package databases

import (
	"context"
	"fmt"
	"log"

	"github.com/redis/go-redis/v9"
)

// RDB 全局 Redis 客户端，其他包通过 databases.RDB 使用
var RDB *redis.Client

// InitRedis 初始化 Redis 连接，启动时调用一次
func InitRedis() {
	// 第一步：创建客户端 + 连接池配置
	// 此时不发任何网络请求，只是在内存中初始化配置对象和连接池
	RDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis 服务地址
		Password: "",               // 无密码则留空
		DB:       0,                // 使用默认数据库 0

		// 连接池配置（可选，按需调整）
		PoolSize:     10, // 最大连接数，默认 10 * CPU核心数
		MinIdleConns: 2,  // 最少保持 2 个空闲连接，减少冷启动延迟
	})

	// 第二步：主动 Ping —— 这才是第一次真正建立 TCP 连接的时刻
	// 流程：TCP 三次握手 → 发送 PING 命令 → 收到 PONG → 连接归还连接池复用
	// 如果 Redis 未启动或地址错误，这里会立即报错，程序启动失败（快速失败原则）
	ctx := context.Background()
	Result, err := RDB.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Redis 连接失败: %v", err)
	}
	fmt.Println("redis 连接成功:" + Result)

}
