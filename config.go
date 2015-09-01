package main

import (
	"os"
)

type Config struct {
	WorkDir string
	Port    string
}

func (c *Config) workDir() string {
	if c == nil {
		return "/var/sitehub"
	}
	return c.WorkDir
}

func (c *Config) port() string {
	if c == nil {
		return "80"
	}
	return c.Port
}

var config = &Config{
	os.Getenv("SH_WORKDIR"),
	os.Getenv("SH_PORT"),
}
