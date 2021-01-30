package models

import (
	//	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	if Hash([]byte("хуй")) != "d05661c81df53f54e9973d8ccdbb0666cd91925d89b24abfad58a1073d3f0a2e" {
		t.Error("Wrong hash!")
	}
}
