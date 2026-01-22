package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"stn-chain/blockchain"
)

func main() {
	blockchain.InitGenesis()

	// 1. RE-INJECTED: Internal Mining Tick (Heartbeat)
	// This ensures the node mines even without manual intervention.
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			newBlk := blockchain.MineThreats()
			if newBlk != nil {
				log.Printf("[MINER] Autonomous Block #%d Sealed", newBlk.Index)
			}
		}
	}()

	// 2. HTTP Route Handling
	http.HandleFunc("/threat", handleThreat)
	http.HandleFunc("/mine",   handleMine)
	http.HandleFunc("/chain",  handleChain)
	http.HandleFunc("/block",  handleBlockReceive)
	http.HandleFunc("/rpc",    handleRPC) // Stratum Bridge

	log.Println("[STN] Sovereignty Active on :8333")
	log.Fatal(http.ListenAndServe(":8333", nil))
}

// 3. RE-INJECTED: Missing Handlers
func handleThreat(w http.ResponseWriter, r *http.Request) {
	var t blockchain.Threat
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Bad JSON", 400)
		return
	}
	blockchain.AddThreat(t)
	w.Write([]byte("Threat Queued"))
}

func handleMine(w http.ResponseWriter, r *http.Request) {
	newBlk := blockchain.MineThreats()
	w.Write([]byte("Mined Block Index: " + string(rune(newBlk.Index))))
}

func handleChain(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(blockchain.Chain)
}

func handleBlockReceive(w http.ResponseWriter, r *http.Request) {
	var blk blockchain.Block
	if err := json.NewDecoder(r.Body).Decode(&blk); err != nil {
		return
	}
	blockchain.AppendBlock(&blk)
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
    var req map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        return
    }

    if req["method"] == "getwork" {
        last := blockchain.LatestBlock()
        
        // Render the response for stratumd
        response := map[string]interface{}{
            "result": map[string]string{
                "data":   last.GetHeaderHex(),
                "target": "00000ffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
            },
            "id":    req["id"],
            "error": nil,
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    }
}
