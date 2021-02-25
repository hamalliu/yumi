package internal

// BuildOptions build options.
type BuildOptions struct {
	Subset     func([]*Instance, int) []*Instance
	SubsetSize int
}

// BuildOption build option interface.
type BuildOption interface {
	Apply(*BuildOptions)
}

type funcOpt struct {
	f func(*BuildOptions)
}

func (f *funcOpt) Apply(opt *BuildOptions) {
	f.f(opt)
}

// Subset Subset option.
func Subset(defaultSize int) BuildOption {
	return &funcOpt{f: func(opt *BuildOptions) {
		opt.SubsetSize = defaultSize
		opt.Subset = subset
	}}
}
