package main

import (
	"apple-findmy-to-mqtt/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := bootstrap.RootApp.Execute(); err != nil {
		panic(err)
	}
}
