package main

import (
	"nbarbey.fr/url-shortener/urlshortener"
	"os"
)

func main() {
	var app *urlshortener.Application
	if os.Getenv("DB_TYPE") == "memory" {
		app = urlshortener.NewInMemoryApplication()
	} else {
		app = urlshortener.NewPGpplication()
	}
	_ = app.Start()
}
