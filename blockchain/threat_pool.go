package blockchain

import (
    "fmt"
    "sync"
)

var (
    pendingThreats []Threat
    muThreats      sync.Mutex
)

// AddThreat enqueues a validated threat for the next block.
func AddThreat(t Threat) error {
    if t.Hash != t.ComputeHash() {
        return fmt.Errorf("hash mismatch")
    }
    muThreats.Lock()
    pendingThreats = append(pendingThreats, t)
    muThreats.Unlock()
    return nil
}

// LoadThreats atomically returns all queued threats and clears the pool.
func LoadThreats() []Threat {
    muThreats.Lock()
    defer muThreats.Unlock()
    out := pendingThreats
    pendingThreats = nil
    return out
}
