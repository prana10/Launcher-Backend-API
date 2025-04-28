package main

import (
	_ "github.com/lib/pq" 
	_ "launcherbackend_api/internal/delivery/http/docs" 
	"launcherbackend_api/internal/app"
)

func main() {
	app.StartApp()
}

