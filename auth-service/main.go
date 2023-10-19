package main

import (
	_ "github.com/mfajri11/family-catering-micro-service/auth-service/config"
	"github.com/mfajri11/family-catering-micro-service/auth-service/internal/app"
)

func main() {
	app.Run()
}
