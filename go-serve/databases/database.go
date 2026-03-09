package databases

import (
	"fmt"
	"gin_demo/models"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// 创造全局gorm.DB变量
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
	host := getEnv("PG_HOST", "localhost")
	port := getEnv("PG_PORT", "5432")
	user := getEnv("PG_USER", "postgres")
	password := getEnv("PG_PASSWORD", "123456")
	dbname := getEnv("PG_DB", "postgres")
	fmt.Println("HOST:", host)
	fmt.Println("PORT:", port)
	fmt.Println("USER:", user)
	fmt.Println("PASS:", password)
	fmt.Println("DBNAME:", dbname)

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
}

// 获取环境变量中的数据库连接信息
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

func AutoMigrate() {
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Category{},
		&models.Tag{}); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}
