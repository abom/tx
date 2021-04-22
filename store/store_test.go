package store

import (
	"math/big"
	"testing"
)

func TestStore(t *testing.T) {
	accountsStore := NewMemoryStore()

	account := &Account{
		ID:      "aaa",
		Name:    "aaa",
		Balance: big.NewFloat(100),
	}

	accountsStore.Set("aaa", account)
	if accountsStore.Get("aaa") == nil {
		t.Error("cannot get item (id: aaa)")
	}
}
