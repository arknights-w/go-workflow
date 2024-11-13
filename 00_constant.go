package go_workflow

import "fmt"

var _ error = (*WorkflowErr)(nil)

type WorkflowErr struct {
	code int
	msg  string
	desc string
}

func (e *WorkflowErr) Error() string {
	return fmt.Sprintf("{code: %d, msg: \"%s\", desc: \"%s\"}", e.code, e.msg, e.desc)
}

func (e *WorkflowErr) WithDesc(desc string) *WorkflowErr {
	var newErr = *e
	newErr.desc = desc
	return &newErr
}

func (e *WorkflowErr) Code() int {
	return e.code
}

func (e *WorkflowErr) Msg() string {
	return e.msg
}

func (e *WorkflowErr) Desc() string {
	return e.desc
}

var (
	// workflow build err form 10001 to 20000

	// 阶段名称重复
	ErrDupStage = &WorkflowErr{code: 10001, msg: "duplicate stage"}
	// 循环依赖
	ErrCircDep = &WorkflowErr{code: 10002, msg: "circular dependency"}
	// 阶段不存在
	ErrNoStage = &WorkflowErr{code: 10003, msg: "stage not found"}

	// workflow run err form 20001 to 30000

	// stage err form 30001 to 40000
)
