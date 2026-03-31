package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	destination := os.Getenv("DESTINATION")
	if destination == "" {
		log.Fatal().Msg("DESTINATION must be set")
	}

	log.Info().Str("destination", destination).Msg("Starting client")
	client := &http.Client{}

	go serveProbes()

	for {
		<-time.After(5 * time.Second)
		go callHTTP(destination, client)
	}
}

func serveProbes() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to start probe server")
	}
}

func callHTTP(destination string, client *http.Client) {
	log.Info().Msg("Calling endpoint")
	url := fmt.Sprintf("%s", destination)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create request")
		return
	}

	do, err := client.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to perform request")
		return
	}
	defer func() {
		_ = do.Body.Close()
	}()

	log.Info().Str("status", do.Status).Msg("Received response from endpoint")
}
