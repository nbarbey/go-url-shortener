package main

import "nbarbey.fr/url-shortener/urlshortener"

func main() {
	_ = urlshortener.NewPGpplication().Start()
}
