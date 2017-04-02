package service

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	Budgets []Budget
}

type Budget struct {
	Category string  `json:"category"`
	Value    float64 `json:"value"`
}

func GetBudgets() []Budget {
	v := viper.GetStringMap("budgets")
	b := make([]Budget, len(v))
	i := 0
	for k, v := range v {
		f, _ := v.(float64)
		b[i] = Budget{
			Category: k,
			Value:    f,
		}
		i++
	}
	return b
}
