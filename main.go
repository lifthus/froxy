package main

import (
	"fmt"
	"froxy/config"
	"froxy/proxy"
	"log"
	"os"
)

func main() {
	conf, err := config.ReadConfig(os.Args[1:])
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("configurations set")

	go func() {
		err = proxy.Start(*conf)
		if err != nil {
			log.Println(err)
			return
		}
	}()

	// any input to exit
	var input string
	fmt.Scanln(&input)
}
