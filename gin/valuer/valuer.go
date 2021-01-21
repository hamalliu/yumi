package valuer

//Key ...
type Key string

const (
	//KeyUser ...
	KeyUser Key = "user"
)

//Valuer ...
type Valuer struct {
	Value interface{}
}

//User ...
type User struct {
	UserID   string
	UserName string
}

// SwitchKey ...
func SwitchKey(key string) Key {
	switch key {
	case "user":
		return KeyUser
	default:
		return ""
	}
}

//User ...
func (v Valuer) User() User {
	val, ok := v.Value.(User)
	if ok {
		return val
	}

	return User{}
}

//String ...
func (v Valuer) String() string {
	val, ok := v.Value.(string)
	if ok {
		return val
	}
	return ""
}

//Float64 ...
func (v Valuer) Float64() float64 {
	val, ok := v.Value.(float64)
	if ok {
		return val
	}
	return 0
}

//Float32 ...
func (v Valuer) Float32() float32 {
	val, ok := v.Value.(float32)
	if ok {
		return val
	}
	return 0
}

//Int ...
func (v Valuer) Int() int {
	val, ok := v.Value.(int)
	if ok {
		return val
	} else {
		return 0
	}
}

//Int64 ...
func (v Valuer) Int64() int64 {
	val, ok := v.Value.(int64)
	if ok {
		return val
	}
	return 0
}
