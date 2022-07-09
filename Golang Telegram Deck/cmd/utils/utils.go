package utils

import (
	"crypto/rand"
	"fmt"
)

func GeneratedID() string {
	x := make([]byte, 16)
	rand.Read(x)
	return fmt.Sprintf("%x", x)
}
