package databases

import (
	"fmt"
	"gin_demo/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 创造全局gorm.DB变量,为了在持久层能够使用DB进行对数据库的操作
var DB *gorm.DB

// 初始化数据库连接等等
func InitDatabase() {
	dsn := getDSN()
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = db
	AutoMigrate()
}

// 获取数据库连接串
func getDSN() string {
	host := mustGetEnv("PG_HOST")
	port := mustGetEnv("PG_PORT")
	user := mustGetEnv("PG_USER")
	password := mustGetEnv("PG_PASSWORD")
	dbname := mustGetEnv("PG_DB")

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

func AutoMigrate() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Category{},
		&models.ChatMessage{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
