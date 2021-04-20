package rest

import (
	"encoding/json"
	"log"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"github.com/abom/tx/store"
	"github.com/abom/tx/processor"
)

const VERSION = "1"

type Handlers struct {
	Store *store.MemoryStore
	Processor *processor.Processor

	TransactionTimeout time.Duration
}


// return a new handlers given the store and processor
func NewHandlers(s *store.MemoryStore, p *processor.Processor, t time.Duration) *Handlers {
	return &Handlers{
		Store: s,
		Processor: p,
		TransactionTimeout: t,
	}
}

func writeJson(status int, w http.ResponseWriter, value interface{}) error {
	w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)

	res, err := json.Marshal(value)
	if err != nil {
		return err
	}
	w.Write(res)
	return nil
}

func writeStatus(status int, w http.ResponseWriter, message string) error {
	return writeJson(status, w, map[string]interface{}{
		"message": message,
		"status": status,
	})
}

func (h *Handlers) AccountsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: implement pagination
	writeJson(200, w, h.Store.GetAll())
}

// account handler, will only try to get an account of <id>
func (h *Handlers) AccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	account := h.Store.Get(id)
	if account == nil {
		writeStatus(404, w, "account cannot be found")
	} else {
		writeJson(200, w, account)
	}
}

// transfer handler, which takes a Transaction as `json` in the body
func (h *Handlers) TransferHandler(w http.ResponseWriter, r *http.Request) {
	tx := &processor.Transaction{}

	bytes, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(bytes, tx)
	if err != nil {
		log.Printf("Invalid transaction: %s (%s)", string(bytes), err)

		writeStatus(400, w, "Invalid transaction format")
	} else {
		// if it's a valid transaction, then submit to the processor
		log.Printf("Got a new transaction: %s", string(bytes))
		h.Processor.Submit(tx)
		// wait for the transaction to be committed
		err = tx.Wait(h.TransactionTimeout)
		if err != nil {
			writeStatus(400, w, err.Error())
		} else {
			writeStatus(200, w, "Transaction completed successfully")
		}
	}
}
