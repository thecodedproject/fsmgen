package main

import (
	internal "github.com/thecodedproject/fsmgen/internal"
	log "log"
)

func main() {

	err := internal.Generate()
	if err != nil {
		log.Fatal(err.Error())
	}
}

