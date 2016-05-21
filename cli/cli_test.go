package cli

import (
	"testing"
)

func TestGreet(t *testing.T) {
	greet := Greet()
	if greet == "" {
		t.Error("Empty greet!")
	}
}
