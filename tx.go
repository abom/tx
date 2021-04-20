package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/abom/tx/store"
	"github.com/abom/tx/processor"
	"github.com/abom/tx/rest"
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
		log.Fatal(err)
	}

	queueProcessor := processor.NewProcessor(storage)
	queueProcessor.Start()

	transactionTimeout, err := time.ParseDuration(*timeout)
	if err != nil {
		log.Fatal(err)
	}

	handlers := rest.NewHandlers(storage, queueProcessor, transactionTimeout)

	router := mux.NewRouter()
	apiRouter := router.PathPrefix(fmt.Sprintf("/api/v%s", rest.VERSION)).Subrouter()

	apiRouter.HandleFunc("/accounts", handlers.AccountsHandler)
	apiRouter.HandleFunc("/account/{id}", handlers.AccountHandler)
	apiRouter.HandleFunc("/transfer", handlers.TransferHandler).Methods("POST")

	server := &http.Server{
        Handler:      router,
        Addr:         fmt.Sprintf("%s:%s", *host, *port),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	defer queueProcessor.Stop()
}
