package main

import (
	"log"

	"PasManagerGophKeeper/internal/router"
)

func main() {
	rout := router.InitServer()
	err := rout.StartServer()

	if err != nil {
		log.Fatal("StartServer:", err)
	}
}
