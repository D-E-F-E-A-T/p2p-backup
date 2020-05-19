package main

import (
	"log"
	"os"
	"os/signal"

	_ "github.com/dmitriy-vas/node/cmd"
	"github.com/dmitriy-vas/node/wire"
)

func main() {
	go wire.Serve()
	sigChan := make(chan os.Signal, 1)
	defer close(sigChan)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
	log.Println("Shutting down...")
}
