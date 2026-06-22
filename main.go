package main

import (
	"log"

	"git.aads.cloud/aad/transmission-otel/config"
	"git.aads.cloud/aad/transmission-otel/fetcher"
	otelexporter "git.aads.cloud/aad/transmission-otel/otel"
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
