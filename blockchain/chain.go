// File: blockchain/chain.go
// Version 1.0
package blockchain

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "sync"
)

const chainFile = "data/chain_data.json"

// Chain is the in-memory ledger of blocks.
var (
    Chain      []*Block
    ChainMutex sync.Mutex
)

// InitGenesis initializes the chain from disk or with a genesis block.
func InitGenesis() {
    log.Println("[BLOCKCHAIN] InitGenesis started")
    if loadChainFromDisk() {
        log.Println("[BLOCKCHAIN] Loaded existing chain from disk")
        return
    }
    log.Println("[BLOCKCHAIN] No existing chain on disk, creating genesis block")
    gen := NewBlock(0, "genesis", nil)
    Chain = []*Block{gen}
    saveChainToDisk()
}

// LatestBlock returns the most recent block on the chain.
func LatestBlock() *Block {
    return Chain[len(Chain)-1]
}

// MineThreats loads pending threats, creates a new block, appends to the chain, persists, and returns the new block.
func MineThreats() *Block {
    prev := LatestBlock()
    threats := LoadThreats()
    newBlock := NewBlock(prev.Index+1, prev.Hash, threats)
    ChainMutex.Lock()
    Chain = append(Chain, newBlock)
    ChainMutex.Unlock()
    saveChainToDisk()
    log.Printf("[BLOCKCHAIN] Mined block #%d with %d threats", newBlock.Index, len(threats))
    return newBlock
}

// AppendBlock validates and appends an incoming block to the chain.
func AppendBlock(b *Block) error {
    ChainMutex.Lock()
    defer ChainMutex.Unlock()
    last := Chain[len(Chain)-1]
    if b.Index != last.Index+1 {
        return fmt.Errorf("invalid index %d, expected %d", b.Index, last.Index+1)
    }
    if b.PrevHash != last.Hash {
        return fmt.Errorf("invalid prev hash, got %s, expected %s", b.PrevHash, last.Hash)
    }
    if b.ComputeHash() != b.Hash {
        return fmt.Errorf("invalid hash, computed %s, got %s", b.ComputeHash(), b.Hash)
    }
    Chain = append(Chain, b)
    saveChainToDisk()
    return nil
}

func loadChainFromDisk() bool {
    data, err := ioutil.ReadFile(chainFile)
    if err != nil {
        return false
    }
    var loaded []*Block
    if err := json.Unmarshal(data, &loaded); err != nil {
        return false
    }
    Chain = loaded
    return true
}

func saveChainToDisk() {
    data, err := json.MarshalIndent(Chain, "", "  ")
    if err != nil {
        return
    }
    ioutil.WriteFile(chainFile, data, 0644)
}
