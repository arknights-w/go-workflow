package go_workflow

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/arknights-w/go-workflow/tools"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/go-echarts/go-echarts/v2/types"
)

type WorkflowBuilder[nt nameType] struct {
	stages map[nt]Stage[nt]
}

func NewBuilder[nt nameType](stages ...Stage[nt]) (*WorkflowBuilder[nt], error) {
	// 构建工作流
	builder := &WorkflowBuilder[nt]{
		stages: make(map[nt]Stage[nt]),
	}
	// stage add
	for _, stage := range stages {
		if _, ok := builder.stages[stage.Name()]; ok {
			return nil, ErrDupStage.WithDesc(fmt.Sprintf("duplicate stage: %v", stage.Name()))
		}
		builder.stages[stage.Name()] = stage
	}
	return builder, nil
}

func (builder *WorkflowBuilder[nt]) AddStage(stage Stage[nt]) error {
	if _, ok := builder.stages[stage.Name()]; ok {
		return ErrDupStage.WithDesc(fmt.Sprintf("duplicate stage: %v", stage.Name()))
	}
	builder.stages[stage.Name()] = stage
	return nil
}

func (builder *WorkflowBuilder[nt]) Build() (Workflow[nt], error) {
	// 1. 构建依赖关系
	edges := make(map[nt][]nt)
	for _, stage := range builder.stages {
		if _, ok := edges[stage.Name()]; !ok {
			edges[stage.Name()] = nil
		}
		for _, dep := range stage.DependOn() {
			// 检查依赖的阶段是否存在
			if _, ok := builder.stages[dep]; !ok {
				return nil, ErrNoStage.WithDesc(fmt.Sprintf("no stage: %v", dep))
			}
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
		sortedStages = append(sortedStages, builder.stages[name])
	}

	// 5. 构建 printer
	prtGraph := builder.buildPrinter()

	return &workflow[nt]{
		stages:   sortedStages,
		printG: prtGraph,
	}, nil
}

func (builder *WorkflowBuilder[nt]) buildPrinter() *charts.Graph {
	var (
		prtEdge   = make([]opts.GraphLink, 0, len(builder.stages)+1)
		prtNode   = make([]opts.GraphNode, 0, len(builder.stages)+1)
		graph     = charts.NewGraph()
		degreeMap = make(map[string]int)
	)
	// 1. build node and edge
	for _, stage := range builder.stages {
		stageName := fmt.Sprint(stage.Name())
		degreeMap[stageName] = len(stage.DependOn())
		prtNode = append(prtNode, opts.GraphNode{
			Name:       stageName,
			SymbolSize: 50,
			Tooltip:    &opts.Tooltip{Formatter: types.FuncStr(stage.Desc())},
		})
		for _, dep := range stage.DependOn() {
			prtEdge = append(prtEdge, opts.GraphLink{
				Source: fmt.Sprint(dep),
				Target: stageName,
			})
		}
	}
	// 2. build start node
	prtNode = append(prtNode, opts.GraphNode{
		Name:       "Σ graph start",
		Tooltip:    &opts.Tooltip{Formatter: "avoid same name"},
		SymbolSize: 50,
		Fixed:      opts.Bool(true),
		X:          200,
		Y:          200,
	})
	for name, degree := range degreeMap {
		if degree == 0 {
			prtEdge = append(prtEdge, opts.GraphLink{
				Source: "Σ graph start",
				Target: name,
			})
		}
	}
	// 3. build graph
	graph.AddSeries("", prtNode, prtEdge,
		charts.WithGraphChartOpts(opts.GraphChart{
			EdgeSymbol: []string{"circle", "arrow"},
			Force:      &opts.GraphForce{Repulsion: 1000, EdgeLength: 100},
			Draggable:  opts.Bool(true),
		}),
		charts.WithLabelOpts(opts.Label{
			Show: opts.Bool(true),
		}),
	).SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{
			Width:     "90vw",
			Height:    "90vh",
			PageTitle: "workflow dependency graph",
		}),
	)
	return graph
}

type workflow[nt nameType] struct {
	stages   []Stage[nt]
	printG *charts.Graph
}

func (w *workflow[nt]) Work(ctx context.Context) (err error) {
	context := NewContext(ctx)
	for _, stage := range w.stages {
		if err = stage.Run(context); err != nil {
			return
		}
	}
	return
}

func (w *workflow[nt]) GetStage(name nt) Stage[nt] {
	for idx := range w.stages {
		if w.stages[idx].Name() != name {
			return w.stages[idx]
		}
	}
	return nil
}

func (w *workflow[nt]) Print(_path string, name string) error {
	fp := path.Join(_path, name)
	if !strings.HasSuffix(fp, ".html") {
		fp = fp + ".html"
	}
	return os.WriteFile(
		fp,
		[]byte(w.printG.RenderContent()),
		0755,
	)
}
