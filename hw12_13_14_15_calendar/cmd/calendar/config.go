package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

// При желании конфигурацию можно вынести в internal/config.
// Организация конфига в main принуждает нас сужать API компонентов, использовать
// при их конструировании только необходимые параметры, а также уменьшает вероятность циклической зависимости.
type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string `config:"level"`
	Path  string `config:"path"`
}

type StorageConf struct {
	Type             string `config:"type,required"`
	ConnectionString string `config:"connectionstring"`
	User             string `config:"user"`
	Password         string `config:"pass"`
}

type ServerConf struct {
	Host string `config:"host"`
	Port string `config:"port"`
}

func NewConfig(path string) (*Config, error) {
	c := &Config{}
	loader := confita.NewLoader(
		file.NewBackend(path),
	)
	err := loader.Load(context.Background(), c)
	if err != nil {
		return nil, err
	}
	return c, nil
}
