package main

import (
	"cbe-super-app-budget/initiator"
	"context"
)

func main() {
	initiator.Init(context.Background())
}
