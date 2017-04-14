package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

type Transactions []Transaction
type Budgets map[string]float64

func (a Transactions) Len() int           { return len(a) }
func (a Transactions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Transactions) Less(i, j int) bool { return a[i].Date.After(a[j].Date) }

type Ledger struct {
	version      int                `json:"version"`
	Budgets      map[string]float64 `json:"budgets"`
	Transactions Transactions       `json:"transactions"`
}

type Transaction struct {
	Number      int       `json:"number"`
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
}

func LoadLedger() (*Ledger, error) {
	filepath := "/Users/guillaume/.config/ledger.json"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return &Ledger{1, make(map[string]float64, 0), make([]Transaction, 0)}, nil
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
	t := Transactions(l.Transactions)
	sort.Sort(t)
	for k, v := range t {
		v.Number = k + 1
		t[k] = v
	}
	j, _ := json.Marshal(l)
	return ioutil.WriteFile(filepath, j, 770)
}
func (l *Ledger) Add(ID string, date time.Time, txtype string, description string, amount float64) bool {
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
		return true
	} else {
		t.ID = ID
		t.Date = date
		t.Type = txtype
		t.Description = description
		t.Amount = amount
		return false
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

func (l *Ledger) SetCategory(number int, category string) error {

	for k, t := range l.Transactions {
		if t.Number == number {
			t.Category = category
			l.Transactions[k] = t
			return nil
		}
	}
	return fmt.Errorf("Transaction %d not found", number)
}

type Filter struct {
	StartDate time.Time
	EndDate   time.Time
	Category  string
}

func NewFilter() *Filter {
	return &Filter{}
}

func (f *Filter) IsFiltered(transaction Transaction) bool {
	// fmt.Printf("%v %v\n", transaction.Date, f.StartDate)
	if f.StartDate.After(time.Time{}) {
		if transaction.Date.Before(f.StartDate) {
			return true
		}
	}
	if f.EndDate.After(time.Time{}) {
		if transaction.Date.After(f.EndDate) {
			return true
		}
	}
	if f.Category != "" {
		if transaction.Category != f.Category {
			return true
		}

	}
	return false
}

func (l *Ledger) GetTransactions(f *Filter) Transactions {
	t := Transactions(l.Transactions)
	sort.Sort(t)
	result := []Transaction{}
	for _, tx := range t {
		if !f.IsFiltered(tx) {
			result = append(result, tx)
		}
	}
	return result
}

func (l *Ledger) RenameCategory(oldName, newName string) {
	for k, tx := range l.Transactions {
		if tx.Category == oldName {
			tx.Category = newName
			l.Transactions[k] = tx
		}
	}
}

func (l *Ledger) DeleteBudget(category string) error {
	delete(l.Budgets, category)
	return nil
}

func (l *Ledger) GetBudgets() Budgets {
	return l.Budgets
}

func (l *Ledger) SetBudget(category string, value float64) error {
	l.Budgets[category] = value
	return nil
}
