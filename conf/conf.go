package conf

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wethedevelop/account/model"
)

// 初始化配置和链接数据库等操作
func Init() {
	_ = godotenv.Load()

	// 连接数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("KAE_MYSQL_USER"),
		os.Getenv("KAE_MYSQL_PASSWORD"),
		os.Getenv("KAE_MYSQL_HOST"),
		os.Getenv("KAE_MYSQL_PORT"),
		os.Getenv("KAE_MYSQL_DB_NAME"),
	)
	model.Database(dsn)
}
