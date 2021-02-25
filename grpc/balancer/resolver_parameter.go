package balancer

const (
	// AttributesKey 从resolver中获取参数的key
	AttributesKey = "balancer"
)

// ParamsFromResolver ...
type ParamsFromResolver struct {
	Color  string `json:"color"`
	Weight int    `json:"weight"`
}
