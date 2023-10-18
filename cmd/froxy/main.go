package main

import (
	"froxy/init/args"
	"froxy/internal/http/reverse"
	"log"
)

func main() {
	var err error

	secure, port, target := args.InitArgs()

	switch target {
	case nil:
		log.Println("forward proxy not implemented yet")
	default:
		if secure {
			log.Println("secure reverse proxy not implemented yet")
		} else {
			err = reverse.ReverseProxy(port, *target)
		}
	}
	log.Fatal(err)
}
