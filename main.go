package main

import (
	"context"

	"github.com/dejandjenic/go-gin-sample/application/setup"
)

func main() {
	r, cleanup := setup.GinHandler("")
	defer cleanup(context.Background())
	r.Run(":8080")
}

// @BasePath /api/v1
