package conf

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/wethedevelop/account/model"
	"github.com/wethedevelop/account/util"
)

// 初始化配置和链接数据库等操作
func Init() {
	_ = godotenv.Load()

	// 连接数据库
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DB_NAME"),
	)
	util.Log().Debug("dsn: %s", dsn)
	model.Database(dsn)
}
