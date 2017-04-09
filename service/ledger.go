package service

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type Ledger struct {
	Transactions []Transaction
}
type Transaction struct {
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:description`
	Date        time.Time `json:"date"`
	Type        string    `json:type`
	Category    string    `json:category`
}

func LoadLedger() (*Ledger, error) {
	filepath := "/Users/guillaume/.config/ledger.json"
	if _, err := os.Stat("/path/to/whatever"); os.IsNotExist(err) {
		return &Ledger{make([]Transaction, 0)}, nil
	}
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		return nil, e
	}
	var ledger Ledger
	json.Unmarshal(file, &ledger)
	return &ledger, nil
}
func (l *Ledger) Save() error {
	filepath := "/Users/guillaume/.config/ledger.json"
	j, _ := json.Marshal(l)
	return ioutil.WriteFile(filepath, j, 770)
}
func (l *Ledger) Add(ID string, date time.Time, txtype string, description string, amount float64) {
	t := l.Get(ID)
	if t == nil {
		transaction := Transaction{
			ID:          ID,
			Date:        date,
			Description: description,
			Amount:      amount,
			Type:        txtype,
			Category:    "",
		}
		l.Transactions = append(l.Transactions, transaction)
	} else {
		t.ID = ID
		t.Date = date
		t.Type = txtype
		t.Description = description
		t.Amount = amount
	}
}

func (l *Ledger) Get(ID string) *Transaction {
	for _, t := range l.Transactions {
		if t.ID == ID {
			return &t
		}
	}
	return nil

}
