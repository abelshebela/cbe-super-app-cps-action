package main

import (
	"context"

	"cbe-super-app-cps-action/initiator"
)

func main() {
	initiator.Init(context.Background())
}
