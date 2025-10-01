package main

import (
	"server/internal/web"
)

func main() {
	handler := web.NewGameHandler()
	router := web.NewRouter(handler)
	router.Run(":8080")
}
