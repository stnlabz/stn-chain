// balances.go
package blockchain

import (
	"encoding/json"
	"os"
)

type Balances map[string]float64

func LoadBalances() Balances {
	file, err := os.ReadFile(dataDir+"balances.json")
	if err != nil {
		return make(Balances)
	}
	var b Balances
	json.Unmarshal(file, &b)
	return b
}

func SaveBalances(b Balances) {
	data, _ := json.MarshalIndent(b, "", "  ")
	os.WriteFile(dataDir+"balances.json", data, 0644)
}