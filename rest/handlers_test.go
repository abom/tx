package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/abom/tx/processor"
	"github.com/abom/tx/store"
)

func GetHandlers() *Handlers {
	storage := store.NewMemoryStore()
	storage.Set("a", &store.Account{
		ID:      "a",
		Name:    "a test",
		Balance: big.NewFloat(100),
	})
	storage.Set("b", &store.Account{
		ID:      "b",
		Name:    "b test",
		Balance: big.NewFloat(100),
	})

	p := processor.NewProcessor(storage)
	p.Start()

	duration, _ := time.ParseDuration("10s")
	handlers := NewHandlers(storage, p, duration)

	return handlers
}

func GetServer() (*Handlers, *httptest.Server) {
	handlers := GetHandlers()
	server := httptest.NewServer(GetRouter(handlers))
	return handlers, server
}

func doRequest(method string, endpoint string, body string, server *httptest.Server) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v%s%s", server.URL, VERSION, endpoint)
	req, _ := http.NewRequest(method, url, strings.NewReader(body))
	client := &http.Client{}

	return client.Do(req)
}

func doSingleRequest(method string, url string, body string) (*http.Response, error) {
	handlers, server := GetServer()

	defer server.Close()
	defer handlers.Processor.Stop()

	return doRequest(method, url, body, server)
}

func TestGetAccounts(t *testing.T) {
	resp, err := doSingleRequest("GET", "/accounts", "")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d, got %d", http.StatusOK, resp.StatusCode)
	}

	var accounts []*store.Account
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &accounts)
	accountLength := len(accounts)
	if accountLength != 2 {
		t.Errorf("unexpected no of accounts returned: %d, wanted: 2", accountLength)
	}
}

func TestGetAccount(t *testing.T) {
	resp, err := doSingleRequest("GET", "/account/a", "")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code: %d, got %d", http.StatusOK, resp.StatusCode)
	}

	// test not found error
	notFoundResp, notFoundErr := doSingleRequest("GET", "/account/xyz", "")
	if notFoundErr != nil {
		t.Error(notFoundErr)
	}

	if notFoundResp.StatusCode != http.StatusNotFound {
		t.Errorf("unexpected status code: %d, got %d", http.StatusNotFound, notFoundResp.StatusCode)
	}
}

func TestTransfer(t *testing.T) {
	tx := processor.NewTransaction("a", "b", 1)

	s, _ := json.Marshal(tx)

	handlers, server := GetServer()
	// do 100 request with amount of 1
	for i := 0; i < 100; i++ {
		resp, err := doRequest("POST", "/transfer", string(s), server)
		if err != nil {
			t.Error(err)
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("unexpected status code: %d, got %d", http.StatusOK, resp.StatusCode)
		}
	}

	resp, err := doRequest("POST", "/transfer", string(s), server)
	if err != nil {
		t.Error(err)
	}

	// should get insufficient funds
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code: %d, got %d", http.StatusBadRequest, resp.StatusCode)
	}

	defer handlers.Processor.Stop()
	defer server.Close()
}
