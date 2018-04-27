package model

import (
    "fmt"
    "crypto/md5"
    "gin-blog/config"
    "strconv"
    "time"
)


// User Table
type User struct {
    BaseModel
    Email       string  `gorm:"unique_index;default:null"`  //邮箱
    Telephone   string  `gorm:"unique_index;default:null"`  //手机号码
    Password    string  `gorm:"default:null"`               //密码
    VerifyState string  `gorm:"default:'0'"`                //邮箱验证状态
    IsAdmin     bool                                        //是否是管理员
    AvatarUrl   string                                      //头像链接
    NickName    string                                      //昵称
    LockState   bool    `gorm:"default:'0'"`                //锁定状态

}

// CheclPassword 验证密码是否正确
func (user User) CheckPassword(password string) bool {
    if password == "" || user.Password == "" {
        return false
    }
    return user.EncryptPassword(password, user.Salt()) == user.Password
}

// Salt 每个用户都有一个不同的盐
func (user User) Salt() string {
    var userSalt string
    if user.Password == "" {
        userSalt = strconv.Itoa(int(time.Now().Unix()))
    } else {
        userSalt = user.Password[0:6]
    }
    return userSalt
}

// EncryptPassword 给密码加密
func (user User) EncryptPassword(password, salt string) (hash string) {
    password = fmt.Sprintf("%x", md5.Sum([]byte(password)))
    hash = salt + password + config.ServerConfig.PassSalt
    hash = salt + fmt.Sprintf("%x", md5.Sum([]byte(hash)))
    return
}

// Lock User 锁定用户
func (user *User) Lock() error {
    return DB.Model(user).Update(map[string]interface{}{
        "lock_state": user.LockState,
    }).Error
}

// Get User
func GetUserById(userId uint64) (*User, error) {

    var newUser User
    err := DB.First(&newUser, "id = ?", userId).Error
    return &newUser, err
}