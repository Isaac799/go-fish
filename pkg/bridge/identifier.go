package bridge

import (
	"crypto/rand"
	"encoding/base64"
)

// RandomID provides a random strung to be used as an ID
// for an element
func RandomID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		// realistically unreachable
		panic("fail gen random ID " + err.Error())
	}
	// RawURLEncoding better for html ID than RawStdEncoding
	return base64.RawURLEncoding.EncodeToString(b)
}
