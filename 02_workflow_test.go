package go_workflow_test

import (
	"context"
	"fmt"
	"testing"

	wf "github.com/arknights-w/go-workflow"
)

type WorkType int

const (
	Init WorkType = iota
	Create
	Update
	Delete
)

func TestSuccess(t *testing.T) {
	var (
		init = wf.NewStage(Init, func(ctx context.Context) bool {
			println("this is Init")
			return true
		})
		create = wf.NewStage(
			Create,
			func(ctx context.Context) bool {
				println("this is Create")
				return true
			},
			wf.WithDependOn([]WorkType{Init}),
		)
		update = wf.NewStage(
			Update,
			func(ctx context.Context) bool {
				println("this is Update")
				return true
			},
			wf.WithDependOn([]WorkType{Create}),
		)
		delete = wf.NewStage(
			Delete,
			func(ctx context.Context) bool {
				println("this is Delete")
				return true
			},
			wf.WithDependOn([]WorkType{Update}),
		)
	)
	builder, err := wf.NewBuilder(init, create, delete, update)
	if err != nil {
		fmt.Printf("1 err: %v\n", err)
		return
	}
	workflow, err := builder.Build()
	if err != nil {
		fmt.Printf("2 err: %v\n", err)
		return
	}
	workflow.Work(context.Background())
}

func TestDuplicate(t *testing.T) {
	var (
		init = wf.NewStage(
			Init,
			func(ctx context.Context) bool {
				println("this is Init")
				return true
			},
		)
		init2 = wf.NewStage(
			Init,
			func(ctx context.Context) bool {
				println("this is Init2")
				return true
			},
		)
	)
	builder, err := wf.NewBuilder(init, init2)
	if err != nil {
		fmt.Printf("1 err: %v\n", err)
		return
	}
	workflow, err := builder.Build()
	if err != nil {
		fmt.Printf("2 err: %v\n", err)
		return
	}
	workflow.Work(context.Background())
}

func TestCircular(t *testing.T) {
	var (
		init = wf.NewStage(Init, func(ctx context.Context) bool {
			println("this is Init")
			return true
		})
		create = wf.NewStage(
			Create,
			func(ctx context.Context) bool {
				println("this is Create")
				return true
			},
			wf.WithDependOn([]WorkType{Init}),
		)
		update = wf.NewStage(
			Update,
			func(ctx context.Context) bool {
				println("this is Update")
				return true
			},
			wf.WithDependOn([]WorkType{Create, Delete}),
		)
		delete = wf.NewStage(
			Delete,
			func(ctx context.Context) bool {
				println("this is Delete")
				return true
			},
			wf.WithDependOn([]WorkType{Update}),
		)
	)
	builder, err := wf.NewBuilder(init, create, delete, update)
	if err != nil {
		t.Fatalf("NewBuilder err: %v\n", err)
		return
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
		return
	}
	workflow.Work(context.Background())
}
