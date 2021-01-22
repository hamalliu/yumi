package entity

import (
	"fmt"

	"yumi/pkg/codec"
	"yumi/pkg/codes"
	"yumi/pkg/login"
	"yumi/pkg/sessions"
	"yumi/pkg/status"
	"yumi/pkg/strfmt"
	"yumi/pkg/types"
)

// UserStatus ...
type UserStatus struct {
	Disabled bool
}

// UserAttribute ...
type UserAttribute struct {
	// 用户id，唯一
	UserID string
	// 用户密码
	Password string
	// 用户名
	UserName string
	//电话号码
	PhoneNumber string
	// 注册时间
	RegisteTime types.Timestamp

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
	if len(u.attr.UserID) < 6 || len(u.attr.UserID) > 60 {
		return status.New(codes.InvalidArgument, "密码长度在6-60之间")
	}
	err = strfmt.RegexpUser(u.attr.UserID)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error())
	}

	// 2. 密码强度
	if len(u.attr.Password) < 10 || len(u.attr.Password) > 30 {
		return fmt.Errorf("密码字数在10~30之间")
	}
	err = strfmt.RegexpPassword(u.attr.Password)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error())
	}

	return nil
}

// BcryptPassword 使用bcrypt 格式化密码
func (u *User) BcryptPassword() (err error) {
	pwd, err := login.GetBcrypt().GenerateFromPassword([]byte(u.attr.Password))
	if err != nil {
		err := status.Internal().WithDetails(err.Error())
		return err
	}
	u.attr.Password = string(pwd)
	return nil
}

// VerifyPassword 验证密码
func (u *User) VerifyPassword(password string) (err error) {
	pass, err := login.GetBcrypt().VarifyPassword([]byte(password), []byte(u.attr.Password))
	if err != nil {
		return status.Internal().WithDetails(err.Error())
	}

	if !pass {
		return status.New(codes.Unauthenticated, "密码错误")
	}

	return nil
}

// Session 构建session
func (u *User) Session(store sessions.Store, userID, password, client string) (string, error) {
	sess, err := sessions.NewSession(store, u.attr.UserID, client, 1)
	err = sess.Save()
	if err != nil {
		return "", status.Internal().WithDetails(err.Error())
	}

	sess.AddFlash(sess.ID, "session_id")
	sess.AddFlash(userID, "user_id")

	// secureKey = Md5(sessionID+userID+password)
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
		return status.New(codes.Unauthenticated, "请重新登录").WithDetails("secure_key 泄漏")
	}
	srvSecureKey := vars[0].(string)

	if secureKey != srvSecureKey {
		return status.New(codes.Unauthenticated, "请重新登录").WithDetails("secureKey error:" + secureKey)
	}

	return nil
}
