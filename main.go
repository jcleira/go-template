package main

import (
	"log"

	"github.com/jcleira/alerts/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("cmd.Execute(), Err: %v", err)
	}
}
