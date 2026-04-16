package framework

import (
	"context"
	"framework/server"
	"framework/store"
)

const frameworkContextName = "frameworkContext"

type frameworkContext struct {
	values  store.Store[any, any]
	servers store.Store[string, server.Component[server.Server]]
}

func newFrameworkContext() *frameworkContext {
	return &frameworkContext{
		values:  store.NewConcurrent[any, any](),
		servers: store.NewConcurrent[string, server.Component[server.Server]](),
	}
}

func ContextWithComponent(ctx context.Context, component server.Component[server.Server]) context.Context {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		frameworkCtx = newFrameworkContext()
		ctx = context.WithValue(ctx, frameworkContextName, frameworkCtx)
	}
	frameworkCtx.servers.Put(component.Name(), component)
	return ctx
}

func ComponentFromContext[T server.Server](ctx context.Context, name string) (c server.Component[T]) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}
	v, ok := frameworkCtx.servers.Get(name)
	if !ok {
		return
	}

	s, ok := v.Server().(T)
	if !ok {
		return
	}

	return server.NewServerComponent(v.Name(), s)
}

func ComponentsFromContext[T server.Server](ctx context.Context) (cs []server.Component[T]) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}

	frameworkCtx.servers.Foreach(func(name string, c server.Component[server.Server]) {
		var s T
		s, ok = c.Server().(T)
		if ok {
			cs = append(cs, server.NewServerComponent(c.Name(), s))
		}
	})
	return
}

func ContextWithValue[Key comparable, Value any](ctx context.Context, k Key, v Value) context.Context {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		frameworkCtx = newFrameworkContext()
		ctx = context.WithValue(ctx, frameworkContextName, frameworkCtx)
	}

	frameworkCtx.values.Put(k, v)
	return ctx
}

func ValueFromContext[Value any](ctx context.Context, k any) (v Value) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}
	r, ok := frameworkCtx.values.Get(k)
	if !ok {
		return
	}
	return r.(Value)
}
