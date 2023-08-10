package main

import (
	"goroxy/config"
	"goroxy/proxy"
	"log"
	"os"
)

func main() {
	conf, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		log.Println(err)
		return
	}
	err = proxy.Start(*conf)
	if err != nil {
		log.Println(err)
		return
	}
}
