module github.com/wethedevelop/account

go 1.16

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/joho/godotenv v1.3.0
	github.com/onsi/ginkgo v1.16.4 // indirect
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/wethedevelop/proto v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	google.golang.org/grpc v1.38.0
	gorm.io/driver/mysql v1.1.0
	gorm.io/gorm v1.21.10
)

replace github.com/wethedevelop/proto => /Users/chengnan/go16/src/github.com/wethedevelop/proto
