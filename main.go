package main

import (
	"log"

	"transmission-otel/config"
	"transmission-otel/fetcher"
	otelexporter "transmission-otel/otel"
)

func main() {
	if _, err := config.Load(); err != nil {
		log.Fatal(err)
	}

	if !config.C.Debug {
		otelexporter.Start()
	}

	fetcher.Start()
}
