package main

import (
	"github.com/zelas91/metric-collector/internal/server"
)

func main() {
	conf := NewConfig()
	if err := server.Run(conf.Addr); err != nil {
		panic(err)
	}

}
