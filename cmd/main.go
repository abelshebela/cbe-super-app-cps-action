package main

import (
	"log"
	"github.com/CBE-Super-App/cbe-super-app-cps-action/initiator"
)

func main() {
	if err := initiator.Initiator(); err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
}
