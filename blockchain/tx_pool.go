// tx_pool.go
package blockchain

import (
	"encoding/json"
	"os"
	"time"
)

//const dataDir = "data/"

func LoadTransactions() []Transaction {
	data, _ := os.ReadFile(dataDir+"transactions.json")
	var txs []Transaction
	json.Unmarshal(data, &txs)
	return txs
}

func SaveTransactions(txs []Transaction) {
	data, _ := json.MarshalIndent(txs, "", "  ")
	os.WriteFile(dataDir+"transactions.json", data, 0644)
}

func NewTransaction(sender, recipient string, amount float64, signature string) Transaction {
	return Transaction{
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		Timestamp: time.Now().Unix(),
		Signature: signature,
	}
}