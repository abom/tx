package processor

import (
	"math/big"
	"testing"
	"time"

	"github.com/abom/tx/store"
)

func TestTransactionTimeout(t *testing.T) {
	tx := NewTransaction("a", "b", 100)

	duration, _ := time.ParseDuration("0.2s")
	err := tx.Wait(duration)
	if err == nil {
		t.Error("no timeout error was returned")
	}
}

func TestProcessor(t *testing.T) {
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

	p := NewProcessor(storage)

	p.Start()

	tx := &Transaction{
		From:   "a",
		To:     "b",
		Amount: big.NewFloat(100),
	}

	p.Submit(tx)
	tx.Wait(0)

	aBalance := storage.Get("a").Balance
	if aBalance.Cmp(big.NewFloat(0.0)) != 0 {
		t.Errorf("a balance should be 0, got %f", aBalance)
	}

	bBalance := storage.Get("b").Balance
	if bBalance.Cmp(big.NewFloat(200)) != 0 {
		t.Errorf("b balance should be 200, got %f", bBalance)
	}

	defer p.Stop()
}
