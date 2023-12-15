package testdata

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("hello world")
	stop()
}

func stop() {
	os.Exit(1)
}
