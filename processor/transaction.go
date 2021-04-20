package processor

import (
	"errors"
	"fmt"
	"math/big"
	"time"
)

type TransactionStatus int

const(
    Pending TransactionStatus = iota
    Success
	Invalid
	Timeout
)

type Transaction struct {
	From string `json:"from"`
	To string `json:"to"`
	Amount *big.Float `json:"amount"`

	Status TransactionStatus `json:"-"`
	Error string `json:"-"`
}

func NewTransaction(from string, to string, amount float64) *Transaction {
	return &Transaction{
		From: from,
		To: to,
		Amount: big.NewFloat(amount),
	}
}

func (t *Transaction) Wait(timeout time.Duration) error {
	started := time.Now()

	for {
		if t.Status == Success {
			break
		}

		if t.Status == Invalid {
			return errors.New(t.Error)
		}

		if timeout > 0 {
			now :=  time.Now()
			diff := now.Sub(started)
			if diff > timeout {
				return errors.New(fmt.Sprintf("wait timeout of %s exceeded, transaction still pending", timeout.String()))
			}
		}
	}

	return nil
}
