package framework

import (
	"context"
	"framework/store"
)

const frameworkContextName = "frameworkContext"

type frameworkContext struct {
	values  store.Store[any, any]
	servers store.Store[string, Server]
}

func newFrameworkContext() *frameworkContext {
	return &frameworkContext{
		values:  store.NewConcurrent[any, any](),
		servers: store.NewConcurrent[string, Server](),
	}
}

func ContextWithServer(ctx context.Context, server Server) context.Context {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		frameworkCtx = newFrameworkContext()
		ctx = context.WithValue(ctx, frameworkContextName, frameworkCtx)
	}
	frameworkCtx.servers.Put(server.Name(), server)
	return ctx
}

func ServerFromContext[T Server](ctx context.Context, name string) (t T) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}
	s, ok := frameworkCtx.values.Get(name)
	if !ok {
		return
	}
	return s.(T)
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
