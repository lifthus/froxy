package main

import (
	"goroxy/config"
	"os"
)

func main() {
	config.ReadConfig(os.Args[1:])
}
