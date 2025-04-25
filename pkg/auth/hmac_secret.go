package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"os"
)

type HmacSecret struct {
	secret [64]byte
}

func NewHmacSecret(filename string) (*HmacSecret, error) {
	h := HmacSecret{}
	if filename == "" {
		return nil, errors.New("filename cannot be empty")
	}

	// If the file is not found, generate a new secret and write it to the file
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			log.Printf("Secret file '%s' does not exist; creating a new secret", filename)
			_, err := rand.Read(h.secret[:])
			if err != nil {
				return nil, fmt.Errorf("failed to generate random secret: %w", err)
			}
			err = os.WriteFile(filename, h.secret[:], 0600)
			if err != nil {
				return nil, fmt.Errorf("failed to write secret to file: %w", err)
			}
			return &h, nil
		}
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if len(data) != len(h.secret) {
		return nil, fmt.Errorf("secret file must be 64 bytes")
	}
	copy(h.secret[:], data)
	return &h, nil
}

func (h *HmacSecret) Get() any {
	return h.secret[:]
}
