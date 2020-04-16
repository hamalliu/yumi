package ecodes

type Status struct {
	Code int    //状态码
	Desc string //状态描述
}

func (s Status) Error() string {
	return s.Desc
}

func NewStatus(code int, desc string) func() Status {
	return func() Status {
		return Status{Code: code, Desc: desc}
	}
}
