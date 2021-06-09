package model

import (
	"strconv"
	"time"

	"github.com/wethedevelop/account/cache"
	"github.com/wethedevelop/account/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	Account        string
	PasswordDigest string
	Nickname       string
	Status         string
	Avatar         string `gorm:"size:1000"`
}

const (
	// PassWordCost 密码加密难度
	PassWordCost = 12
	// Active 激活用户
	Active string = "active"
	// Inactive 未激活用户
	Inactive string = "inactive"
	// Suspend 被封禁用户
	Suspend string = "suspend"
)

// GetUser 用ID获取用户
func GetUser(ID interface{}) (User, error) {
	var user User
	result := DB.First(&user, ID)
	return user, result.Error
}

// SetPassword 设置密码
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// CheckPassword 校验密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

// 创建用户
func (user *User) Create() error {
	err := DB.Create(user).Error
	return err
}

// 保存
func (item *User) Save() error {
	err := DB.Save(item).Error
	return err
}

// UserID 返回string版的uid
func (user *User) UserID() string {
	return strconv.Itoa(int(user.ID))
}

// MakeToken 生成token
func (user *User) MakeToken() (string, int64, error) {
	// 移动端生成token, 2周自动过期
	token := util.RandStringRunes(15)
	exp := 14 * 24 * time.Hour
	tokenExpire := time.Now().Add(exp).Unix()
	if err := cache.SaveUserToken(token, user.UserID(), exp); err != nil {
		return "", 0, err
	}
	return token, tokenExpire, nil
}

// CheckRegistered 同用户名检测
func CheckRegistered(account string) (bool, error) {
	var count int64 = 0
	err := DB.Model(&User{}).Where("account = ?", account).Count(&count).Error
	return count > 0, err
}
