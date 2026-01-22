// transaction.go
package blockchain

type Transaction struct {
	Sender    string
	Recipient string
	Amount    float64
	Timestamp int64
	Signature string
}