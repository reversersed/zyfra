package main

import (
	"github.com/reversersed/zyfra/internal/app"
)

func main() {
	app := app.New()
	if app == nil {
		return
	}

	app.Run()
}
