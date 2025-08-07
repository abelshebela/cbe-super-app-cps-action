package main

import (
	"context"

	"cbe-super-app-member-auth/initiator"
)

func main() {
	initiator.Init(context.Background())
}
