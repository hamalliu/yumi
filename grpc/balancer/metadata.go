package balancer

const (
	// AttributesKey 从resolver中获取参数的key
	AttributesKey = "balancer"
)

// MetaFromResolver ...
type MetaFromResolver struct {
	Color  string `json:"color"`
	Weight uint64 `json:"weight"`
}
