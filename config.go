package main

type Config struct {
	WorkDir string
	Port    uint16
}

func (c *Config) workDir() string {
	if c == nil {
		return "/var/sitehub"
	}
	return c.WorkDir
}

var config = &Config{"/var/sitehub", 80}
