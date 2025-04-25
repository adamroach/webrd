package auth_test

import (
	"os"
	"testing"

	"github.com/adamroach/webrd/pkg/auth"
)

func TestNewHmacSecret_FileDoesNotExist(t *testing.T) {
	filename := "test_secret_file"
	defer os.Remove(filename) // Clean up after the test

	hmacSecret, err := auth.NewHmacSecret(filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hmacSecret == nil {
		t.Fatal("expected HmacSecret instance, got nil")
	}

	if len(hmacSecret.Get().([]byte)) != 64 {
		t.Fatalf("expected secret length to be 64, got %d", len(hmacSecret.Get().([]byte)))
	}

	// Check if the file was created
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatal("expected secret file to be created, but it does not exist")
	}
}

func TestNewHmacSecret_FileExists(t *testing.T) {
	filename := "test_secret_file"
	defer os.Remove(filename) // Clean up after the test

	// Create a file with a valid 64-byte secret
	secret := make([]byte, 64)
	for i := range secret {
		secret[i] = byte(i)
	}
	err := os.WriteFile(filename, secret, 0600)
	if err != nil {
		t.Fatalf("failed to create test secret file: %v", err)
	}

	hmacSecret, err := auth.NewHmacSecret(filename)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hmacSecret == nil {
		t.Fatal("expected HmacSecret instance, got nil")
	}

	if string(hmacSecret.Get().([]byte)) != string(secret) {
		t.Fatal("expected secret to match the content of the file")
	}
}

func TestNewHmacSecret_EmptyFilename(t *testing.T) {
	_, err := auth.NewHmacSecret("")
	if err == nil {
		t.Fatal("expected an error for empty filename, got nil")
	}
}

func TestNewHmacSecret_InvalidFileLength(t *testing.T) {
	filename := "test_invalid_secret_file"
	defer os.Remove(filename) // Clean up after the test

	// Create a file with an invalid secret length
	invalidSecret := make([]byte, 32) // Less than 64 bytes
	err := os.WriteFile(filename, invalidSecret, 0600)
	if err != nil {
		t.Fatalf("failed to create test invalid secret file: %v", err)
	}

	_, err = auth.NewHmacSecret(filename)
	if err == nil {
		t.Fatal("expected an error for invalid secret file length, got nil")
	}
}
