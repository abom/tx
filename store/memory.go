package store

import (
	"encoding/json"
	"io/ioutil"
)

type MemoryStore struct {
	accounts map[string]*Account
}

// return a new empty MemoryStore
func NewMemoryStore() *MemoryStore {
	var store MemoryStore
	store.accounts = make(map[string]*Account)
	return &store
}

// return a new MemoryStore from file, the file should contain
// a list of accounts as json (with from, to and balance fields)
func NewMemoryStoreFromFile(path string) (*MemoryStore, error) {
	bytes, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	memoryStore := NewMemoryStore()
	var accounts []*Account
	json.Unmarshal(bytes, &accounts)

	for _, account := range accounts {
		memoryStore.Set(account.ID, account)
	}

	return memoryStore, nil
}

// get an account by id
func (s *MemoryStore) Get(id string) *Account {
	return s.accounts[id]
}

// set an account by id
func (s *MemoryStore) Set(id string, account *Account) {
	s.accounts[id] = account
}

// get a list of all accounts
func (s *MemoryStore) GetAll() []*Account {
	all := make([]*Account, 0, len(s.accounts))

	for _, account := range s.accounts {
		all = append(all, account)
	}

	return all
}
