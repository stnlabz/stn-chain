// File: main.go
// Version 1.4
package main

import (
    "encoding/json"
    "log"
    "net/http"

    "dw-chain/blockchain"
)

func main() {
    // Initialize genesis and load chain
    blockchain.InitGenesis()

    // HTTP endpoints for threat index functionality
    http.HandleFunc("/threat", handleThreat)
    http.HandleFunc("/mine",   handleMine)
    http.HandleFunc("/chain",  handleChain)
    http.HandleFunc("/peers",  blockchain.PeersHandler)
    http.HandleFunc("/block",  handleBlockReceive)

    log.Println("[HTTP] Starting server on :8333")
    log.Fatal(http.ListenAndServe(":8333", nil))
}

func handleThreat(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "use POST", http.StatusMethodNotAllowed)
        return
    }
    var t blockchain.Threat
    if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
        http.Error(w, "bad JSON", http.StatusBadRequest)
        return
    }
    if err := blockchain.AddThreat(t); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("threat queued"))
}

func handleMine(w http.ResponseWriter, r *http.Request) {
    newBlk := blockchain.MineThreats()
    blockchain.BroadcastBlock(*newBlk)
    w.Write([]byte("new block mined"))
}

func handleChain(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(blockchain.Chain)
}

func handleBlockReceive(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "use POST", http.StatusMethodNotAllowed)
        return
    }
    var blk blockchain.Block
    if err := json.NewDecoder(r.Body).Decode(&blk); err != nil {
        http.Error(w, "invalid JSON", http.StatusBadRequest)
        return
    }
    if err := blockchain.AppendBlock(&blk); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("block accepted"))
}
