package tests

import (
	"testing"

	"github.com/fathimasithara01/tradeverse/pkg/auth"
)

func TestGenerateAndValidateToken(t *testing.T) {
	userID := uint(1)
	email := "test@example.com"
	role := "admin"
	roleID := uint(10)
	secret := "test-secret-key"

	// Generate token
	token, err := auth.GenerateJWT(userID, email, role, roleID, secret)
	if err != nil {
		t.Fatalf("failed to generate JWT: %v", err)
	}

	// Validate token
	claims, err := auth.ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("failed to validate JWT: %v", err)
	}

	// Assertions
	if claims.UserID != userID {
		t.Errorf("expected UserID %d, got %d", userID, claims.UserID)
	}
	if claims.Email != email {
		t.Errorf("expected Email %s, got %s", email, claims.Email)
	}
	if claims.Role != role {
		t.Errorf("expected Role %s, got %s", role, claims.Role)
	}
	if claims.RoleID != roleID {
		t.Errorf("expected RoleID %d, got %d", roleID, claims.RoleID)
	}
}
