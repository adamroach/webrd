package auth_test

import (
	"testing"

	"github.com/adamroach/webrd/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestDarwinPasswordChecker_CurrentUser(t *testing.T) {
	checker := auth.NewPasswordChecker()

	username := checker.CurrentUser()
	t.Log("Current user:", username)
	assert.NotEmpty(t, username, "CurrentUser should return a non-empty username")
}

func TestDarwinPasswordChecker_CheckPassword(t *testing.T) {
	checker := auth.NewPasswordChecker()

	// Test with invalid credentials
	username := "invalid_user"
	password := "invalid_password"
	result := checker.CheckPassword(username, password)
	assert.False(t, result, "CheckPassword should return false for invalid credentials")

	// Note: Testing with valid credentials is beyond the scope of this unit test
}
