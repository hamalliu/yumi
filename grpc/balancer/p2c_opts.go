package balancer

// P2cOption ...
type P2cOption struct {
	f func(*p2cOptions)
}

type p2cOptions struct {
	Color string
}

// SetColor ...
func SetColor(color string) P2cOption {
	p := P2cOption{}
	p.f = func(po *p2cOptions) {
		po.Color = color
	}
	return p
}

