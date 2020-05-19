package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/dmitriy-vas/p2p/cmd"
	_ "github.com/dmitriy-vas/p2p/handler"
)

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	log.Println("Shutting down...")
}
