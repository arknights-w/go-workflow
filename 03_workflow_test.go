package go_workflow_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	wf "github.com/arknights-w/go-workflow"
)

type WorkType string

const (
	Init   WorkType = "init"
	Create WorkType = "create"
	Update WorkType = "update"
	Delete WorkType = "delete"
	None   WorkType = "none"
)

func TestSuccess(t *testing.T) {
	stages := []wf.Stage[WorkType]{
		wf.NewStage(Init,
			func(ctx wf.Context) error {
				println("this is Init")
				ctx.Set(Init, "success")
				return nil
			},
		), wf.NewStage(Create,
			func(ctx wf.Context) error {
				fmt.Printf("this is Create, Init stage is %v, Update stage still is %v\n", ctx.Get(Init), ctx.Get(Update))
				ctx.Set(Create, "success")
				return nil
			},
			wf.WithDependOn([]WorkType{Init}),
		), wf.NewStage(Update,
			func(ctx wf.Context) error {
				fmt.Printf("this is Update, Create stage is %v\n", ctx.Get(Create))
				ctx.Child().Set(Update, "success")
				return nil
			},
			wf.WithDependOn([]WorkType{Create}),
		), wf.NewStage(Delete,
			func(ctx wf.Context) error {
				fmt.Printf("this is Delete, can not get Update stage: %v\n", ctx.Get(Update))
				return nil
			},
			wf.WithDependOn([]WorkType{Create, Update}),
		),
	}
	builder, err := wf.NewBuilder(stages...)
	if err != nil {
		t.Fatalf("NewBuilder err: %v\n", err)
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
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
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
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
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
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
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
	}
	workflow.Work(context.Background())
}

func TestPrint(t *testing.T) {
	stages := []wf.Stage[WorkType]{}
	for i := 1; i < 10; i++ {
		str_i := strconv.Itoa(i)
		stages = append(stages, wf.NewStage(
			WorkType("init_"+str_i),
			func(ctx wf.Context) error {
				println("this is Init", str_i)
				return nil
			},
			wf.WithDesc[WorkType]("初始化 "+str_i),
		))
	}
	for i := 1; i < 10; i++ {
		str_i := strconv.Itoa(i)
		str_sub_i := strconv.Itoa(i - 1)
		if i == 1 {
			stages = append(stages, wf.NewStage(
				WorkType("stage "+str_i),
				func(ctx wf.Context) error {
					println("this is stage", str_i)
					return nil
				},
			))
		} else {
			stages = append(stages, wf.NewStage(
				WorkType("stage "+str_i),
				func(ctx wf.Context) error {
					println("this is stage", str_i)
					return nil
				},
				wf.WithDependOn([]WorkType{WorkType("stage " + str_sub_i)}),
			))
		}
	}

	builder, err := wf.NewBuilder(stages...)
	if err != nil {
		t.Fatalf("NewBuilder err: %v\n", err)
	}
	workflow, err := builder.Build()
	if err != nil {
		t.Fatalf("Build err: %v\n", err)
	}
	workflow.Work(context.Background())
	workflow.Print("", "test.html")
}
