package service

import (
	"encoding/json"
	"io/ioutil"
)

type Budgets map[string]float64

type Config struct {
	Budgets map[string]float64 `json:"budgets"`
}

func LoadConfig() (*Config, error) {
	filepath := "/Users/guillaume/.config/gobud.json"
	file, e := ioutil.ReadFile(filepath)
	if e != nil {
		return nil, e
	}
	var config Config
	json.Unmarshal(file, &config)
	return &config, nil
}

func (c *Config) Save() error {
	filepath := "/Users/guillaume/.config/gobud.json"
	j, _ := json.Marshal(c)
	return ioutil.WriteFile(filepath, j, 644)
}

func (c *Config) Delete(category string) error {
	delete(c.Budgets, category)
	return nil
}

func (c *Config) GetBudgets() Budgets {
	return c.Budgets
}

func (c *Config) SetBudget(category string, value float64) error {
	c.Budgets[category] = value
	return nil
}
