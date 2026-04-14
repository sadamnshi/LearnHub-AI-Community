package main

import (
	"log"

	"gin_demo/databases"
	"gin_demo/router"

	"github.com/joho/godotenv"
)

func main() {
	// 0. 加载 .env 环境变量文件（本地开发）
	// 如果 .env 不存在，则使用系统环境变量（生产部署）
	if err := godotenv.Load("config/.env"); err != nil {
		log.Println("📝 未找到 config/.env 文件，将使用系统环境变量")
	}

	// 1. 初始化数据库并且自动迁移数据库创建数据表
	databases.InitDatabase()
	databases.InitRedis()
	// 2. 初始化路由并启动（CORS 已在 SetupRouter 内部配置）
	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
