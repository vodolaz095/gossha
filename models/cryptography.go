package models

import (
	"crypto/sha256"
	"fmt"
)

// Hash makes irreversible sha256 hash of string
func Hash(input []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(input))
}
