package go_workflow

import (
	"context"
	"fmt"

	"github.com/arknights-w/go-workflow/tools"
)

type WorkflowBuilder[nt nameType] struct {
	stages map[nt]Stage[nt]
}

func NewBuilder[nt nameType](stages ...Stage[nt]) (*WorkflowBuilder[nt], error) {
	workflow := &WorkflowBuilder[nt]{
		stages: make(map[nt]Stage[nt]),
	}
	for _, stage := range stages {
		if _, ok := workflow.stages[stage.Name()]; ok {
			return nil, ErrDupStage
		}
		workflow.stages[stage.Name()] = stage
	}
	return workflow, nil
}

func (b *WorkflowBuilder[nt]) AddStage(stage Stage[nt]) error {
	if _, ok := b.stages[stage.Name()]; ok {
		return ErrDupStage
	}
	b.stages[stage.Name()] = stage
	return nil
}

func (b *WorkflowBuilder[nt]) Build() (Workflow[nt], error) {
	// 1. 构建依赖关系
	edges := make(map[nt][]nt)
	for _, stage := range b.stages {
		for _, dep := range stage.DependOn() {
			edges[dep] = append(edges[dep], stage.Name())
		}
	}

	// 2. 拓扑排序
	sorted, cycle := tools.TopologicalSort(edges)

	// 3. 循环检测
	if len(cycle) != 0 {
		return nil, ErrCircDep.WithDesc(fmt.Sprintf("cycle: %v", cycle))
	}

	// 4. 构建工作流
	sortedStages := make([]Stage[nt], 0, len(sorted))
	for _, name := range sorted {
		sortedStages = append(sortedStages, b.stages[name])
	}

	return &workflow[nt]{
		stages: sortedStages,
	}, nil
}

type workflow[nt nameType] struct {
	stages []Stage[nt]
}

func (w *workflow[nt]) Work(ctx context.Context) {
	for _, stage := range w.stages {
		if !stage.Run(ctx) {
			return
		}
	}
}

func (w *workflow[nt]) GetStage(name nt) Stage[nt] {
	for idx := range w.stages {
		if w.stages[idx].Name() != name {
			return w.stages[idx]
		}
	}
	return nil
}
