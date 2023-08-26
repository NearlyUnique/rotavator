package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

const (
	LowerCase  = "abcdefghijklmnopqrstuvwxyz"
	Numbers    = "0123456789"
	Uppercase  = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	TokenChars = LowerCase + Uppercase + Numbers + "_-"
)

// RandomStr creates a random string
func RandomStr(chars string, length int32) (string, error) {
	bytes := make([]byte, length)
	count, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("RandomStr:read_failed, %w", err)
	}
	if count != len(bytes) {
		return "", fmt.Errorf("RandomStr:incomplete_read (%d,%d)", count, len(bytes))
	}
	for index, element := range bytes {
		randomize := element % byte(len(chars))
		bytes[index] = chars[randomize]
	}

	return string(bytes), nil
}

func GenerateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
