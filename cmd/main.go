package main

import (
	"context"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/initiator"
)

func main() {
	initiator.Init(context.Background())
}
