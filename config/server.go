package config

type Server interface {
	Reload() error
}
