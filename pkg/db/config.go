package db

import (
	"fmt"

	"github.com/benchkram/bobc/pkg/rnd"
)

type Config struct {
	Host string
	Port string

	User     string
	Password string

	Name string

	UseSSL bool
}

func (c Config) ForTesting() Config {
	c.Name = nameTesting + "_" + rnd.RandStringBytesMaskImprSrc(7)
	return c
}

func (c Config) ConnectString() string {
	return fmt.Sprintf(connectStr,
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Password,
	)
}

func NewConfig() Config {
	return Config{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Name:     name,
		UseSSL:   false,
	}
}
