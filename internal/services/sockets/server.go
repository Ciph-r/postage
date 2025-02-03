package sockets

import "net/http"

func NewServer() *http.Server {
	return &http.Server{
		Addr:    ":80", // TODO: set addr with a configurable var.
		Handler: nil,   // TODO: set handler with socket router.
	}
}
