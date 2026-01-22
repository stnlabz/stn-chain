// Replace your main.go with this logic to re-enable the "Heartbeat"
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time" // Added for the ticker
	"stn-chain/blockchain"
)

func main() {
	blockchain.InitGenesis()

	// Re-injecting the internal mining tick
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		for range ticker.C {
			log.Println("[MINER] Internal Tick: Sealing current state...")
			// Triggering internal mine without needing a curl command
			newBlk := blockchain.MineThreats()
			if newBlk != nil {
				log.Printf("[MINER] Success: Block #%d added", newBlk.Index)
			}
		}
	}()

	http.HandleFunc("/threat", handleThreat)
	http.HandleFunc("/mine",   handleMine)
	http.HandleFunc("/chain",  handleChain)
	http.HandleFunc("/peers",  blockchain.PeersHandler)
	http.HandleFunc("/block",  handleBlockReceive)
	
	// FIXED: Actual RPC Bridge for stratumd
	http.HandleFunc("/rpc", handleRPC)

	log.Println("[HTTP] Server hot on :8333")
	log.Fatal(http.ListenAndServe(":8333", nil))
}

func handleRPC(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		return
	}

	if request["method"] == "getwork" {
		// This provides the work template stn-stratumd is looking for
		work := blockchain.LatestBlock().GetWorkFormat() 
		json.NewEncoder(w).Encode(map[string]interface{}{
			"result": work,
			"id":     request["id"],
			"error":  nil,
		})
	}
}
