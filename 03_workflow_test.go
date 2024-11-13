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
	None
)

func TestSuccess(t *testing.T) {
	var (
		init = wf.NewStage(Init, func(ctx wf.Context) error {
			println("this is Init")
			ctx.Set(Init, "success")
			return nil
		})
		create = wf.NewStage(
			Create,
			func(ctx wf.Context) error {
				fmt.Printf("this is Create, Init stage is %v, Update stage still is %v\n", ctx.Get(Init), ctx.Get(Update))
				ctx.Set(Create, "success")
				return nil
			},
			wf.WithDependOn([]WorkType{Init}),
		)
		update = wf.NewStage(
			Update,
			func(ctx wf.Context) error {
				fmt.Printf("this is Update, Create stage is %v\n", ctx.Get(Create))
				ctx.Child().Set(Update, "success")
				return nil
			},
			wf.WithDependOn([]WorkType{Create}),
		)
		delete = wf.NewStage(
			Delete,
			func(ctx wf.Context) error {
				fmt.Printf("this is Delete, can not get Update stage: %v\n", ctx.Get(Update))
				return nil
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
			func(ctx wf.Context) error {
				println("this is Init")
				return nil
			},
		)
		init2 = wf.NewStage(
			Init,
			func(ctx wf.Context) error {
				println("this is Init2")
				return nil
			},
		)
	)
	builder, err := wf.NewBuilder(init, init2)
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

func TestNoStage(t *testing.T) {
	var (
		init = wf.NewStage(
			Init,
			func(ctx wf.Context) error {
				println("this is Init")
				return nil
			},
		)
		init2 = wf.NewStage(
			Create,
			func(ctx wf.Context) error {
				println("this is Init2")
				return nil
			},
			wf.WithDependOn([]WorkType{Init, None}),
		)
	)
	builder, err := wf.NewBuilder(init, init2)
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

func TestCircular(t *testing.T) {
	var (
		init = wf.NewStage(Init, func(ctx wf.Context) error {
			println("this is Init")
			return nil
		})
		create = wf.NewStage(
			Create,
			func(ctx wf.Context) error {
				println("this is Create")
				return nil
			},
			wf.WithDependOn([]WorkType{Init}),
		)
		update = wf.NewStage(
			Update,
			func(ctx wf.Context) error {
				println("this is Update")
				return nil
			},
			wf.WithDependOn([]WorkType{Create, Delete}),
		)
		delete = wf.NewStage(
			Delete,
			func(ctx wf.Context) error {
				println("this is Delete")
				return nil
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
