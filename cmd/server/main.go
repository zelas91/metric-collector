package main

import (
	"github.com/zelas91/metric-collector/internal/server"
)

func main() {
	if err := server.Run("8080"); err != nil {
		panic(err)
	}

}
