package go_workflow

import "context"

type stage[nt nameType] struct {
	name nt
	deps []nt
	desc string
	run  func(ctx context.Context) bool
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

func (s *stage[nt]) Run(ctx context.Context) bool {
	return s.run(ctx)
}

func NewStage[nt nameType](
	name nt,
	run func(ctx context.Context) bool,
	opts ...stageOpt[nt]) Stage[nt] {
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
