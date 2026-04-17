package framework

import (
	"context"

	"github.com/oylshe1314/framework/server"
	"github.com/oylshe1314/framework/store"
)

const frameworkContextName = "_frameworkContext"

type frameworkContext struct {
	values  store.Store[any, any]
	servers store.Store[string, server.Server]
}

func newFrameworkContext() *frameworkContext {
	return &frameworkContext{
		values:  store.NewConcurrent[any, any](),
		servers: store.NewConcurrent[string, server.Server](),
	}
}

func ContextWithComponent(ctx context.Context, name string, svr server.Server) context.Context {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		frameworkCtx = newFrameworkContext()
		ctx = context.WithValue(ctx, frameworkContextName, frameworkCtx)
	}
	frameworkCtx.servers.Put(name, svr)
	return ctx
}

func ServerFromContext[T server.Server](ctx context.Context, name string) (s T) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}
	v, ok := frameworkCtx.servers.Get(name)
	if !ok {
		return
	}

	s, ok = v.(T)
	if !ok {
		return
	}

	return s
}

func ServersFromContext[T server.Server](ctx context.Context) (ss []T) {
	frameworkCtx, ok := ctx.Value(frameworkContextName).(*frameworkContext)
	if !ok {
		return
	}

	frameworkCtx.servers.Foreach(func(name string, svr server.Server) {
		var s T
		s, ok = svr.(T)
		if ok {
			ss = append(ss, s)
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
