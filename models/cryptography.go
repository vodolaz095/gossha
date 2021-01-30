package models

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
)

func randSeq(n int) (b []byte, err error) {
	b = make([]byte, n)
	_, err = rand.Read(b)
	return
}

// Hash makes irreversible sha256 hash of string
func Hash(input []byte) string {
	return fmt.Sprintf("%x", sha256.Sum256(input))
}

// GenSalt generates random string
func GenSalt() (salt string, err error) {
	x, err := randSeq(64)
	if err != nil {
		return
	}
	return Hash(x), nil
}
