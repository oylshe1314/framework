package config

type Server interface {
	Load(dir string) error
}
