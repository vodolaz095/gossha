package models

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"time"
)

func randSeq(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Hash makes irreversible sha256 hash of string
func Hash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

// GenSalt generates random string
func GenSalt() (string, error) {
	return Hash(randSeq(64)), nil
}
