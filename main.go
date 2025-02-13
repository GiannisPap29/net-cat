package main

import (
	"fmt"
	"log"
	"os"

	"net-cat/internal/config" // constables
	"net-cat/internal/server"
	"net-cat/internal/utils" // just atoi
)

func main() {
	port := config.DefaultPort
	if len(os.Args) == 2 {
		port = utils.Atoi(os.Args[1])
	} else if len(os.Args) > 2 {
		fmt.Println(config.UsageMessage)
		return
	}

	s := server.NewServer(port)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
