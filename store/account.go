package store

import "math/big"

type Account struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Balance *big.Float `json:"balance"`
}
