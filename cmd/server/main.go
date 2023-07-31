package main

import "github.com/zelas91/metric-collector/internal/server"

func main() {
	err := server.Run("8080")
	if err != nil {
		panic(err)
	}
}
