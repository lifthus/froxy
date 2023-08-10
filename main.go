package main

import (
	"goroxy/config"
	"log"
	"os"
)

func main() {
	conf, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(conf)
}
