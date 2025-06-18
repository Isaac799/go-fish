package bridge

import (
	"testing"
)

func TestRandomID_Unique(t *testing.T) {
	n := 1000
	m := make(map[string]bool, n)
	for _ = range n {
		s := RandomID()
		if _, exists := m[s]; exists {
			t.Fatal("unique id recreated")
		}
		m[s] = true
	}
}
