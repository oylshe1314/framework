package component

import "context"

type Component interface {
	Init(ctx context.Context) error
	Start() error
	Close() error

	Name() string
}

type ComponentsOption map[string]string

func (cso ComponentsOption) Get(component string) string {
	c, ok := cso[component]
	if !ok {
		return ""
	}
	return c
}
