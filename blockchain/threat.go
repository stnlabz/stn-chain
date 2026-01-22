package blockchain

import (
    "crypto/sha256"
    "encoding/hex"
    "fmt"
    "time"
)

// Threat represents a single threat intelligence report.
type Threat struct {
    ID        string // unique ID, e.g. UUID
    Domain    string // source domain
    Type      string // e.g. "malware", "phishing"
    Severity  int    // 1–10 scale
    Timestamp int64  // unix seconds
    Hash      string // SHA256 of the serialized fields
}

// ComputeHash serializes the core fields and returns hex SHA256.
func (t *Threat) ComputeHash() string {
    payload := fmt.Sprintf("%s|%s|%s|%d|%d",
        t.ID, t.Domain, t.Type, t.Severity, t.Timestamp,
    )
    sum := sha256.Sum256([]byte(payload))
    return hex.EncodeToString(sum[:])
}

// NewThreat constructs a Threat with current timestamp and correct Hash.
func NewThreat(id, domain, typ string, severity int) *Threat {
    t := &Threat{
        ID:        id,
        Domain:    domain,
        Type:      typ,
        Severity:  severity,
        Timestamp: time.Now().Unix(),
    }
    t.Hash = t.ComputeHash()
    return t
}
