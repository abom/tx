package processor

import (
	"fmt"

	"github.com/abom/tx/store"
)

const MAX_TX = 100

type Processor struct {
	id           string
	transactions chan *Transaction
	running      bool

	storage *store.MemoryStore
}

// get a new processor
func NewProcessor(storage *store.MemoryStore) *Processor {
	var processor *Processor

	processor = &Processor{
		transactions: make(chan *Transaction, MAX_TX),
		storage:      storage,
	}

	return processor
}

// submit a transaction to current processor
func (p *Processor) Submit(t *Transaction) {
	t.Status = Pending
	p.transactions <- t
}

// commit a transaction
func (p *Processor) commit(t *Transaction) {
	// do some validations and set proper status accordingly
	if t.From == t.To {
		t.Status = Invalid
		t.Error = "cannot transfer to the same account"
		return
	}

	accountFrom := p.storage.Get(t.From)
	accountTo := p.storage.Get(t.To)

	if accountFrom == nil {
		t.Status = Invalid
		t.Error = fmt.Sprintf("account '%s' cannot be found", t.From)
	} else {
		if accountFrom.Balance.Cmp(t.Amount) == -1 {
			t.Status = Invalid
			t.Error = "insufficient balance"
		} else {
			if accountTo == nil {
				t.Status = Invalid
				t.Error = fmt.Sprintf("account '%s' cannot be found", t.To)
			}
		}
	}

	// after validations, if it's still pending, apply the transaction
	if t.Status == Pending {
		accountFrom.Balance.Sub(accountFrom.Balance, t.Amount)
		accountTo.Balance.Add(accountTo.Balance, t.Amount)
		t.Status = Success
	}
}

// start the processor handler
func (p *Processor) Start() {
	p.running = true

	go p.handler()
}

// processor handler, which should run in a goroutine
// it pulls transaction from the buffered channel and commits them in order
func (p *Processor) handler() {

	for {
		if !p.running {
			break
		}

		tx := <-p.transactions
		p.commit(tx)
	}
}

func (p *Processor) Stop() {
	p.running = false
}
