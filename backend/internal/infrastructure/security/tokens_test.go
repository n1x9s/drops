package security

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTokenManagerIssueAndParse(t *testing.T) {
	manager := NewTokenManager("access-secret", "refresh-secret", time.Minute, time.Hour)
	userID := uuid.New()

	pair, err := manager.IssuePair(context.Background(), userID)
	if err != nil {
		t.Fatalf("issue pair: %v", err)
	}
	accessClaims, err := manager.ParseAccess(pair.AccessToken)
	if err != nil {
		t.Fatalf("parse access: %v", err)
	}
	if accessClaims.UserID != userID {
		t.Fatalf("access user id = %s, want %s", accessClaims.UserID, userID)
	}
	if _, err := manager.ParseRefresh(pair.AccessToken); err == nil {
		t.Fatal("access token should not parse as refresh token")
	}
	if manager.HashToken(pair.RefreshToken) == pair.RefreshToken {
		t.Fatal("token hash should not equal token")
	}
}
