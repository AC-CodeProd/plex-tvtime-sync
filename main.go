package main

import (
	"plex-tvtime-sync/bootstrap"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()
	if err := bootstrap.RootApp.Execute(); err != nil {
		panic(err)
	}
}
