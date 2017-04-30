package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"time"
)

var ShortFormat = "2006-01-02"

type Transactions []*Transaction
type Budgets map[string]float64

func (a Transactions) Len() int           { return len(a) }
func (a Transactions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Transactions) Less(i, j int) bool { return a[i].GetDate().After(a[j].GetDate()) }

type Ledger struct {
	version      int                `json:"version"`
	Balance      float64            `json:"balance"`
	Budgets      map[string]float64 `json:"budgets"`
	Transactions Transactions       `json:"transactions"`
}

type Transaction struct {
	Number      int       `json:"number"`
	ID          string    `json:"id"`
	Amount      float64   `json:"amount"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
	UserDate    time.Time `json:"user-date"`
	Type        string    `json:"type"`
	Category    string    `json:"category"`
	Balance     float64   `json:"balance"`
}

func (t *Transaction) GetDate() time.Time {
	if !t.UserDate.After(time.Time{}) {
		return t.Date
	}
	return t.UserDate
}

func LoadLedger() (*Ledger, error) {
	filepath := "/Users/guillaume/.config/ledger.json"
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return &Ledger{1, 0, make(map[string]float64, 0), make([]*Transaction, 0)}, nil
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
	previousBalance := l.Balance
	for k, v := range t {
		v.Number = k + 1
		v.Balance = previousBalance
		previousBalance = previousBalance - v.Amount
		t[k] = v
	}
	j, err := json.Marshal(l)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath, j, 770)
}
func (l *Ledger) Add(ID string, date time.Time, txtype string, description string, amount float64) bool {
	t := l.Get(ID)
	if t == nil {
		transaction := Transaction{
			ID:          ID,
			Date:        date,
			UserDate:    date,
			Description: description,
			Amount:      amount,
			Type:        txtype,
			Category:    "",
		}
		l.Transactions = append(l.Transactions, &transaction)
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

//TODO rename in GetById
func (l *Ledger) Get(ID string) *Transaction {
	for _, t := range l.Transactions {
		if t.ID == ID {
			return t
		}
	}
	return nil

}

// TODO remove this
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

func (l *Ledger) UpdateTransaction(number int, tx *Transaction) error {

	for k, t := range l.Transactions {
		if t.Number == number {
			l.Transactions[k] = tx
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

func (f *Filter) IsFiltered(transaction *Transaction) bool {
	// fmt.Printf("%v %v\n", transaction.Date, f.StartDate)
	if f.StartDate.After(time.Time{}) {
		if transaction.GetDate().Before(f.StartDate) {
			return true
		}
	}
	if f.EndDate.After(time.Time{}) {
		if transaction.GetDate().After(f.EndDate) {
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
func (l *Ledger) GetTransaction(number int) (*Transaction, error) {
	for _, v := range l.Transactions {
		if v.Number == number {
			return v, nil
		}
	}
	return nil, fmt.Errorf("%s not found")
}

func (l *Ledger) GetTransactions(f *Filter) Transactions {
	t := Transactions(l.Transactions)
	sort.Sort(t)
	result := []*Transaction{}
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
	l.Budgets[newName] = l.Budgets[oldName]
	delete(l.Budgets, oldName)
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
