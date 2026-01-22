package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
    "golang.org/x/crypto/argon2" // Run: go get golang.org/x/crypto/argon2
)

type Block struct {
    Index     int
    Timestamp int64
    Threats   []Threat
    PrevHash  string
    Hash      string
}

// ComputeHash (Legacy SHA256) - Keep for Block 0 and 1
func (b *Block) ComputeHash() string {
    data := fmt.Sprintf("%d|%d|%s", b.Index, b.Timestamp, b.PrevHash)
    for _, t := range b.Threats {
        data += "|" + t.Hash
    }
    sum := sha256.Sum256([]byte(data))
    return hex.EncodeToString(sum[:])
}

// ComputeArgonHash (ASIC Resistant) - Mandatory for Block 2+
func (b *Block) ComputeArgonHash() string {
    data := fmt.Sprintf("%d|%d|%s", b.Index, b.Timestamp, b.PrevHash)
    for _, t := range b.Threats {
        data += "|" + t.Hash
    }
    // Salt is static for chain-wide consistency
    salt := []byte("stn-sovereign-quad-salt")
    // Argon2id: 1 pass, 64MB memory, 4 threads (Tuned for Pi 5)
    hash := argon2.IDKey([]byte(data), salt, 1, 64*1024, 4, 32)
    return hex.EncodeToString(hash)
}

func NewBlock(index int, prevHash string, threats []Threat) *Block {
    b := &Block{
        Index:     index,
        Timestamp: time.Now().Unix(),
        Threats:   threats,
        PrevHash:  prevHash,
    }
    // Logic: Genesis and initial testing use SHA256; transition at Block 2
    if index >= 2 {
        b.Hash = b.ComputeArgonHash()
    } else {
        b.Hash = b.ComputeHash()
    }
    return b
}

// GetHeaderHex formats the block for the RPC bridge
func (b *Block) GetHeaderHex() string {
    header := fmt.Sprintf("%08x%s", b.Index, b.PrevHash)
    return hex.EncodeToString([]byte(header))
}
