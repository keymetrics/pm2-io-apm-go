package structures

import (
	"crypto/rand"
	"encoding/hex"
	"os"

	"github.com/google/uuid"
)

// Config from user code
type Config struct {
	PublicKey       string
	PrivateKey      string
	Name            string
	ServerName      string
	Hostname        string
	Node            *string
	Proxy           string
	ProcessUniqueID string
}

// InitNames with random values
func (config *Config) InitNames() {
	realHostname, err := os.Hostname()
	if err != nil {
		config.Hostname = randomHex(5)
	} else {
		config.Hostname = realHostname
	}
	if config.ServerName == "" {
		config.ServerName = config.Hostname
	}
	config.ProcessUniqueID = uuid.Must(uuid.NewRandom()).String()
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return hex.EncodeToString(bytes)
}
