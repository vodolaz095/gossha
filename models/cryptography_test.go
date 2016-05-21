package models

import (
	//	"fmt"
	"testing"
)

func TestHash(t *testing.T) {
	if Hash("хуй") != "d05661c81df53f54e9973d8ccdbb0666cd91925d89b24abfad58a1073d3f0a2e" {
		t.Error("Wrong hash!")
	}
}

func TestSalt(t *testing.T) {
	salt1, _ := GenSalt()
	salt2, _ := GenSalt()
	if salt1 == salt2 {
		t.Error("We got the same salt twice!!!")
	}
}
