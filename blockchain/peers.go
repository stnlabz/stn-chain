// File: blockchain/peers.go
// Version: 1.0
package blockchain

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "sync"
)

var (
    peerMutex sync.Mutex
    peers     = make(map[string]struct{})
)

// PeersHandler handles GET and POST for /peers.
// GET returns JSON list of peer URLs.
// POST accepts {"url":"http://other:8333"} to register a peer.
func PeersHandler(w http.ResponseWriter, r *http.Request) {
    peerMutex.Lock()
    defer peerMutex.Unlock()

    switch r.Method {
    case http.MethodGet:
        list := make([]string, 0, len(peers))
        for u := range peers {
            list = append(list, u)
        }
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(list); err != nil {
            http.Error(w, "failed to encode peers list", http.StatusInternalServerError)
        }

    case http.MethodPost:
        var req struct{ URL string `json:"url"` }
        if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
            http.Error(w, "invalid JSON", http.StatusBadRequest)
            return
        }
        peers[req.URL] = struct{}{}
        w.WriteHeader(http.StatusCreated)

    default:
        http.Error(w, "use GET or POST", http.StatusMethodNotAllowed)
    }
}

// BroadcastBlock sends the given block JSON to each peer's /block endpoint asynchronously.
func BroadcastBlock(block Block) {
    peerMutex.Lock()
    list := make([]string, 0, len(peers))
    for u := range peers {
        list = append(list, u)
    }
    peerMutex.Unlock()

    data, err := json.Marshal(block)
    if err != nil {
        fmt.Printf("[peers] marshal error: %v\n", err)
        return
    }
    for _, url := range list {
        go func(u string) {
            resp, err := http.Post(u+"/block", "application/json", bytes.NewReader(data))
            if err != nil {
                fmt.Printf("[peers] failed to broadcast to %s: %v\n", u, err)
                return
            }
            defer resp.Body.Close()
            body, _ := io.ReadAll(resp.Body)
            fmt.Printf("[peers] sent block to %s: %s\n", u, string(body))
        }(url)
    }
}
