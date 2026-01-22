// File: blockchain/chain.go
package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

const chainFile = "data/chain_data.json"

var (
	Chain      []*Block
	ChainMutex sync.Mutex
)

func InitGenesis() {
	if loadChainFromDisk() {
		return
	}
	gen := NewBlock(0, "genesis", nil)
	Chain = []*Block{gen}
	saveChainToDisk()
}

// FIXED: Only one declaration of LatestBlock
func LatestBlock() *Block {
	ChainMutex.Lock()
	defer ChainMutex.Unlock()
	return Chain[len(Chain)-1]
}

func MineThreats() *Block {
	prev := LatestBlock()
	threats := LoadThreats()
	
	// Create the new block
	newBlock := NewBlock(prev.Index+1, prev.Hash, threats)
	
	ChainMutex.Lock()
	Chain = append(Chain, newBlock)
	ChainMutex.Unlock()
	
	saveChainToDisk()
	log.Printf("[BLOCKCHAIN] Mined block #%d with %d threats", newBlock.Index, len(threats))
	return newBlock
}

func AppendBlock(b *Block) error {
	ChainMutex.Lock()
	defer ChainMutex.Unlock()
	last := Chain[len(Chain)-1]

	if b.Index != last.Index+1 {
		return fmt.Errorf("invalid index %d, expected %d", b.Index, last.Index+1)
	}
	if b.PrevHash != last.Hash {
		return fmt.Errorf("invalid prev hash")
	}

	// Validate Hash (ASIC-Resistance check)
	var valid bool
	if b.Index < 2 {
		valid = (b.ComputeHash() == b.Hash)
	} else {
		valid = (b.ComputeArgonHash() == b.Hash)
	}

	if !valid {
		return fmt.Errorf("hash validation failed")
	}

	Chain = append(Chain, b)
	saveChainToDisk()
	return nil
}

// Standard file helpers below
func saveChainToDisk() {
	data, _ := json.MarshalIndent(Chain, "", "  ")
	ioutil.WriteFile(chainFile, data, 0644)
}

func loadChainFromDisk() bool {
	data, err := ioutil.ReadFile(chainFile)
	if err != nil {
		return false
	}
	return json.Unmarshal(data, &Chain) == nil
}

// ProcessExternalBlock handles blocks solved by stn-stratumd workers
func ProcessExternalBlock(index int, prevHash string, solution string) error {
	ChainMutex.Lock()
	defer ChainMutex.Unlock()

	last := Chain[len(Chain)-1]

	// 1. Context Validation
	if index != last.Index+1 {
		return fmt.Errorf("stale work: expected index %d", last.Index+1)
	}
	if prevHash != last.Hash {
		return fmt.Errorf("parent mismatch: expected %s", last.Hash)
	}

	// 2. Create Candidate Block
	threats := LoadThreats() // Pull any threats queued while mining
	newBlock := &Block{
		Index:     index,
		Timestamp: time.Now().Unix(),
		Threats:   threats,
		PrevHash:  prevHash,
		Hash:      solution, // The Argon2id hash found by the miner
	}

	// 3. PoW Validation
	// Note: The miner must have solved the hash using the Argon2id parameters 
	// (64MB, 1 pass, 4 threads) to match b.ComputeArgonHash()
	if newBlock.ComputeArgonHash() != solution {
		return fmt.Errorf("invalid proof of work solution")
	}

	// 4. Commit
	Chain = append(Chain, newBlock)
	saveChainToDisk()
	log.Printf("[STRATUM] Block #%d accepted from external miner", index)
	return nil
}
