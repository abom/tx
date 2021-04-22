package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/abom/tx/processor"
	"github.com/abom/tx/rest"
	"github.com/abom/tx/store"
)

func main() {
	path := flag.String("path", "", "path to accounts mock file in json format")
	timeout := flag.String("timeout", "10s", "transaction timeout e.g. 10s")
	host := flag.String("host", "127.0.0.1", "server host, defaults to '127.0.0.1'")
	port := flag.String("port", "8000", "server port, defaults to '8000'")

	flag.Parse()

	if *path == "" {
		flag.Usage()
		return
	}

	storage, err := store.NewMemoryStoreFromFile(*path)
	if err != nil {
		log.Printf("couldn't load accounts from '%s'", *path)
		log.Fatal(err)
	} else {
		log.Printf("accounts load successfully from '%s'", *path)
	}

	queueProcessor := processor.NewProcessor(storage)
	queueProcessor.Start()

	transactionTimeout, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Fatal(err)
	}

	handlers := rest.NewHandlers(storage, queueProcessor, transactionTimeout)
	router := rest.GetRouter(handlers)
	addr := fmt.Sprintf("%s:%s", *host, *port)
	server := &http.Server{
		Handler: router,
		Addr:    addr,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server is ready at %s", addr)

	defer queueProcessor.Stop()
}
