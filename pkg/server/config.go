package server

import (
	"fmt"
)

type Config struct {
	Host string
	Port int
}

func (c Config) Address() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
