package go_workflow

type stage[nt nameType] struct {
	name nt
	deps []nt
	desc string
	run  func(ctx Context) error
}

func (s *stage[nt]) Name() nt {
	return s.name
}

func (s *stage[nt]) DependOn() []nt {
	return s.deps
}

func (s *stage[nt]) Desc() string {
	return s.desc
}

func (s *stage[nt]) Run(ctx Context) error {
	return s.run(ctx)
}

func NewStage[nt nameType](
	name nt,
	run func(ctx Context) error,
	opts ...stageOpt[nt],
) Stage[nt] {
	one := &stage[nt]{
		name: name,
		run:  run,
	}
	for _, opt := range opts {
		opt(one)
	}
	return one
}

type stageOpt[nt nameType] func(s *stage[nt])

func WithDependOn[nt nameType](deps []nt) stageOpt[nt] {
	return func(s *stage[nt]) {
		s.deps = deps
	}
}

func WithDesc[nt nameType](desc string) stageOpt[nt] {
	return func(s *stage[nt]) {
		s.desc = desc
	}
}
