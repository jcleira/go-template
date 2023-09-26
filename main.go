package main

import (
	"log"

	"github.com/jcleira/go-template/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute(), Err: %v", err)
	}
}
