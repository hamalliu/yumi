package login

import "golang.org/x/crypto/bcrypt"

// Bcrypt 包装bcrypt算法
type Bcrypt struct{}

var _bcrypt Bcrypt

// GetBcrypt 获取bcrypt对象
func GetBcrypt() *Bcrypt {
	return &_bcrypt
}

// GenerateFromPassword 根据用户输入的密码，生成经bcrypt算法加密后的密码。
func (b *Bcrypt) GenerateFromPassword(password []byte) (bcryptPassword []byte, err error) {
	bcryptPassword, err = bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	return
}

// VarifyPassword 验证密码
func (b *Bcrypt) VarifyPassword(bcryptPassword, password []byte) (pass bool, err error) {
	err = bcrypt.CompareHashAndPassword(bcryptPassword, password)
	if err != nil {
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
