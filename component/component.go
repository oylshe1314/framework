package component

import "context"

type Component interface {
	Name() string

	Init(ctx context.Context) error
	Start() error
	Close() error
}

type ComponentsOption map[string]string

func (cso ComponentsOption) Get(component string) string {
	c, ok := cso[component]
	if !ok {
		return ""
	}
	return c
}
