package store

// Store interface
type Store interface {
	Get(id string) *Account
	Set(id string, account *Account)
	GetAll() []*Account
}

