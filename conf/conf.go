package conf

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wethedevelop/account/model"
)

// 初始化配置和链接数据库等操作
func Init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 连接数据库
	model.Database(os.Getenv("MYSQL_DSN"))
}
