package main

import (
	"notification-service/internal/context"
)

func main() {
	appCtx := context.NewAppContext()
	appCtx.Init()

	defer appCtx.Close()
}
