package go_workflow

import "context"

type ctx struct {
	kv     map[any]any
	parent Context
}

func (c *ctx) Get(key any) any {
	if v, ok := c.kv[key]; ok {
		return v
	}
	return c.parent.Get(key)
}

func (c *ctx) Set(key, value any) {
	c.kv[key] = value
}

func (c *ctx) Child() Context {
	return &ctx{kv: make(map[any]any), parent: c}
}

type sourceCtx struct {
	ctx
	parent context.Context
}

func (c *sourceCtx) Get(key any) any {
	if v, ok := c.kv[key]; ok {
		return v
	}
	return c.parent.Value(key)
}

func NewContext(parent context.Context) Context {
	return &sourceCtx{
		ctx: ctx{
			kv:     make(map[any]any),
			parent: nil,
		},
		parent: parent,
	}
}
