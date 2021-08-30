package main

import (
	"github.com/ozonva/ova-purchase-api/internal/server"
)

func main() {
	server := server.NewServer(81, 8181)
	server.Run()
}
