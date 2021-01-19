package entity

import (
	"errors"
	"fmt"
	"regexp"

	"yumi/pkg/codes"
	"yumi/pkg/login"
	"yumi/pkg/sessions"
	"yumi/pkg/status"
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
	err = regexpUser(u.attr.UserID)
	if err != nil {
		return status.New(codes.InvalidArgument, err.Error())
	}

	// 2. 密码强度
	if len(u.attr.Password) < 10 || len(u.attr.Password) > 30 {
		return fmt.Errorf("密码字数在10~30之间")
	}
	err = regexpPassword(u.attr.Password)
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
		return errors.New("密码错误")
	}

	return nil
}

func regexpUser(str string) error {
	re := regexp.MustCompile(`\w+`)
	if string(re.Find([]byte(str))) != str {
		return fmt.Errorf("要求字母，数字或下划线")
	}

	return nil
}

func regexpPassword(str string) error {
	specailChars := "[~`!@#\\$%\\^&\\*\\(\\)_\\+\\-=\\{\\}\\[\\]\\|\\\\:;\"'<>\\,\\./\\?]"
	alphabet := "[A-Za-z]"
	number := "[0-9]"
	//包含特殊字符和字母
	specialCharAndAlphabet := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, alphabet, alphabet, specailChars)
	re := regexp.MustCompile(specialCharAndAlphabet)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含特殊字符和数字
	specialCharAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", specailChars, number, number, specailChars)
	re = regexp.MustCompile(specialCharAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	//包含数字和字母
	alphabetAndNumber := fmt.Sprintf(".*%s.*%s.*|.*%s.*%s.*", alphabet, number, number, alphabet)
	re = regexp.MustCompile(alphabetAndNumber)
	if string(re.Find([]byte(str))) == str {
		return nil
	}

	return fmt.Errorf("密码必须包含特殊字符，数字，字母中的两种以上")
}

// Session ...
func (u *User) Session(store sessions.Store, client string) (string, error) {
	loginLen, err := sessions.GetSessionLen(store, u.attr.UserID, client)
	if err != nil {
		return "", err
	}
	if loginLen > 0 {
		sessions.DeleteSession(store, u.attr.UserID, client, 1)
	}
	sess := sessions.NewSession(store, u.attr.UserID, client)
	err = sess.Save()
	if err != nil {
		return "", err
	}

	return sess.ID, nil
}
