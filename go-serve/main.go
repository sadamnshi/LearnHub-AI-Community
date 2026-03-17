//TIP <p>To run your code, right-click the code and select <b>Run</b>.</p> <p>Alternatively, click
// the <icon src="AllIcons.Actions.Execute"/> icon in the gutter and select the <b>Run</b> menu item from here.</p>

package main

import (
	"log"

	"gin_demo/databases"
	"gin_demo/router"
)

func main() {
	// 1. 初始化数据库并且自动迁移数据库创建数据表
	databases.InitDatabase()
	databases.InitRedis()
	// 2. 初始化路由并启动（CORS 已在 SetupRouter 内部配置）
	r := router.SetupRouter()

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
