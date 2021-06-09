package serializer

import (
	"github.com/wethedevelop/account/model"
	pb "github.com/wethedevelop/proto/auth"
)

// BuildUser 序列化用户
func BuildUser(user model.User) *pb.User {
	return &pb.User{
		Id:       int64(user.ID),
		Nickname: user.Nickname,
	}
}
