package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

const (
	chainFile   = "data/chain_data.json"
	backupFile  = "data/chain_data.json.bak"
)

var (
	Chain      []*Block
	ChainMutex sync.Mutex
)

func InitGenesis() {
	// 1. Force create data directory if missing
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		log.Println("[DISK] Data directory missing. Creating...")
		os.Mkdir("data", 0755)
	}

	// 2. Load existing ledger
	if loadChainFromDisk() {
		log.Printf("[DISK] Ledger loaded: %d blocks detected", len(Chain))
		return
	}

	// 3. Fallback to Genesis
	log.Println("[DISK] No ledger found. Generating Genesis Block...")
	gen := NewBlock(0, "genesis", nil)
	Chain = []*Block{gen}
	saveChainToDisk()
}

func LatestBlock() *Block {
	if len(Chain) == 0 { return nil }
	return Chain[len(Chain)-1]
}

func MineThreats() *Block {
	prev := LatestBlock()
	threats := LoadThreats()
	
	newBlock := NewBlock(prev.Index+1, prev.Hash, threats)
	
	ChainMutex.Lock()
	Chain = append(Chain, newBlock)
	ChainMutex.Unlock()
	
	saveChainToDisk()
	log.Printf("[BLOCKCHAIN] Mined block #%d and saved to disk", newBlock.Index)
	return newBlock
}

func saveChainToDisk() {
	data, err := json.MarshalIndent(Chain, "", "  ")
	if err != nil {
		log.Printf("[ERROR] Failed to encode chain: %v", err)
		return
	}

	// Rotate main file to backup before writing new data (Safety First)
	if _, err := os.Stat(chainFile); err == nil {
		os.Rename(chainFile, backupFile)
	}

	err = ioutil.WriteFile(chainFile, data, 0644)
	if err != nil {
		log.Printf("[CRITICAL] Write failed: %v", err)
	} else {
		log.Println("[DISK] Persistence Successful.")
	}
}

func loadChainFromDisk() bool {
	// Try main file first, then backup
	data, err := ioutil.ReadFile(chainFile)
	if err != nil {
		data, err = ioutil.ReadFile(backupFile)
		if err != nil {
			return false
		}
		log.Println("[DISK] Main ledger corrupted or missing. Restored from backup.")
	}
	return json.Unmarshal(data, &Chain) == nil
}
