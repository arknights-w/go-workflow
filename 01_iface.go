package go_workflow

import "context"

// nameType 是一个可比较的类型，用于定义工作流和阶段的名称。
type nameType = comparable

/**
 * Workflow 接口定义了一个工作流，它包含了一个工作方法和获取阶段的方法。
 *
 * @type_param nt 工作流名称的类型，必须是可比较的。
 */
type Workflow[nt nameType] interface {
	// Work 方法执行工作流的主要逻辑。
	Work(ctx context.Context)

	// GetStage 方法根据名称获取工作流中的某个阶段。
	GetStage(nt) Stage[nt]
}

/**
 * Stage 接口定义了工作流中的一个阶段，它包含了阶段的名称、依赖、描述和执行方法。
 *
 * @type_param nt 阶段名称的类型，必须是可比较的。
 */
type Stage[nt nameType] interface {
	// Name 方法返回阶段的名称。
	Name() nt

	// Rely 方法返回当前阶段所依赖的阶段名称列表。
	DependOn() []nt

	// Desc 方法返回阶段的描述信息。
	Desc() string

	// Run 方法执行当前阶段的逻辑，并返回是否继续执行下一个阶段。
	Run(ctx context.Context) (isContinue bool)
}
