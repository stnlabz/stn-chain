package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

// Block holds threats.
type Block struct {
    Index     int
    Timestamp int64
    Threats   []Threat
    PrevHash  string
    Hash      string
}

// ComputeHash covers Index, Timestamp, PrevHash and all Threat.Hashes.
func (b *Block) ComputeHash() string {
    data := fmt.Sprintf("%d|%d|%s", b.Index, b.Timestamp, b.PrevHash)
    for _, t := range b.Threats {
        data += "|" + t.Hash
    }
    sum := sha256.Sum256([]byte(data))
    return hex.EncodeToString(sum[:])
}

// NewBlock builds a block, computes its hash.
func NewBlock(index int, prevHash string, threats []Threat) *Block {
    b := &Block{
        Index:     index,
        Timestamp: time.Now().Unix(),
        Threats:   threats,
        PrevHash:  prevHash,
    }
    b.Hash = b.ComputeHash()
    return b
}
