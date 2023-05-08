package main

import (
	"log"

	"PasManagerGophKeeper/internal/router"
)

func main() {
	rout := router.InitServer()
	err := rout.Router()

	if err != nil {
		log.Fatal("Router:", err)
	}
}
