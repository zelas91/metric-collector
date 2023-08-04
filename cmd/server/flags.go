package main

import "flag"

var addr *string

func init() {
	addr = flag.String("a", "localhost:8080", "endpoint start server")
}

type Config struct {
	Addr string
}

func NewConfig() *Config {

	flag.Parse()
	return &Config{
		Addr: *addr,
	}
}
