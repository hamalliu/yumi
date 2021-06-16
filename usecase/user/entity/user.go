package entity

import (
	"errors"
	"fmt"

	"yumi/pkg/codec"
	"yumi/pkg/login"
	"yumi/pkg/sessions"
	"yumi/pkg/regexputility"
	"yumi/pkg/types"
	"yumi/status"
)

// UserStatus ...
type UserStatus struct {
	Disabled bool
}

// UserAttribute ...
type UserAttribute struct {
	// 用户id，唯一
	UserID string `bson:"user_id"`
	// 用户uuid，唯一
	UserUUID string `bson:"user_uuid"`
	// 用户密码
	Password string `bson:"password"`
	// 用户名
	UserName string `bson:"user_name"`
	//电话号码
	PhoneNumber string `bson:"phone_number"`
	// 注册时间
	RegisteTime types.Timestamp `bson:"registe_time"`

	// 用户状态
	UserStatus
}

// User ...
type User struct {
	attr *UserAttribute
}

// NewUser ...
func NewUser(attr *UserAttribute) *User {
	return &User{attr: attr}
}

// LawEnforcement 执法：检查当前数据是否合乎业务规定
func (u *User) LawEnforcement() (err error) {
	// 1. 用户名格式
	if !regexputility.RegexpUser(u.attr.UserID) {
		return status.FailedPrecondition().WithMessage(status.UserFmtIncorrect)
	}

	// 2. 密码强度
	if regexputility.RegexpPassword(u.attr.Password) {
		return status.FailedPrecondition().WithMessage(status.PasswordFmtIncorrect)
	}

	return nil
}

// BcryptPassword 使用bcrypt 格式化密码
func (u *User) BcryptPassword() (err error) {
	pwd, err := login.GetBcrypt().GenerateFromPassword([]byte(u.attr.Password))
	if err != nil {
		err := status.Internal().WithDetails(err)
		return err
	}
	u.attr.Password = string(pwd)
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) (err error) {
	pass, err := login.GetBcrypt().VarifyPassword([]byte(password), []byte(u.attr.Password))
	if err != nil {
		return status.Internal().WithDetails(err)
	}

	if !pass {
		return status.FailedPrecondition().WithMessage(status.PasswordIncorrect)
	}

	return nil
}

// Session 构建session
func (u *User) Session(store sessions.Store, userID, password, client string) (string, error) {
	sess, err := sessions.NewSession(store, u.attr.UserID, client, 1)
	if err != nil {
		return "", status.Internal().WithDetails(err)
	}
	err = sess.Save()
	if err != nil {
		return "", status.Internal().WithDetails(err)
	}

	sess.AddFlash(sess.ID, "session_id")
	sess.AddFlash(userID, "user_id")

	cnt := fmt.Sprintf("%s&%s&%s", sess.ID, userID, password)
	secureKey := codec.Md5([]byte(cnt))
	sess.AddFlash(secureKey, "secure_key")
	return sess.ID, nil
}

// Authenticate 认证登录状态
func (u *User) Authenticate(store sessions.Store, sessionID, secureKey string) error {
	sess, err := sessions.GetSession(store, sessionID)
	if err != nil {
		return err
	}

	vars := sess.GetValues("secure_key")
	if len(vars) == 0 {
		return status.Unauthenticated().WithDetails(errors.New("secure_key 泄漏"))
	}
	srvSecureKey := vars[0].(string)

	if secureKey != srvSecureKey {
		return status.Unauthenticated().WithDetails(errors.New("secureKey error:" + secureKey))
	}

	return nil
}
