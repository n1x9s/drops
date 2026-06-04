package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessTTL     time.Duration
	refreshTTL    time.Duration
	now           func() time.Time
}

type TokenPair struct {
	AccessToken   string
	RefreshToken  string
	AccessExpiry  time.Time
	RefreshExpiry time.Time
}

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func NewTokenManager(accessSecret string, refreshSecret string, accessTTL time.Duration, refreshTTL time.Duration) TokenManager {
	return TokenManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessTTL:     accessTTL,
		refreshTTL:    refreshTTL,
		now:           time.Now,
	}
}

func (m TokenManager) IssuePair(_ context.Context, userID uuid.UUID) (TokenPair, error) {
	now := m.now().UTC()
	accessExpiry := now.Add(m.accessTTL)
	refreshExpiry := now.Add(m.refreshTTL)

	access, err := m.sign(userID, accessExpiry, m.accessSecret, "access")
	if err != nil {
		return TokenPair{}, err
	}
	refresh, err := m.sign(userID, refreshExpiry, m.refreshSecret, "refresh")
	if err != nil {
		return TokenPair{}, err
	}
	return TokenPair{
		AccessToken:   access,
		RefreshToken:  refresh,
		AccessExpiry:  accessExpiry,
		RefreshExpiry: refreshExpiry,
	}, nil
}

func (m TokenManager) ParseAccess(tokenString string) (Claims, error) {
	return m.parse(tokenString, m.accessSecret, "access")
}

func (m TokenManager) ParseRefresh(tokenString string) (Claims, error) {
	return m.parse(tokenString, m.refreshSecret, "refresh")
}

func (m TokenManager) HashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func (m TokenManager) sign(userID uuid.UUID, expiresAt time.Time, secret []byte, tokenType string) (string, error) {
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.NewString(),
			Subject:   userID.String(),
			Audience:  []string{"secondbrain"},
			Issuer:    "secondbrain-api",
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(m.now().UTC()),
			NotBefore: jwt.NewNumericDate(m.now().UTC()),
		},
	}
	claims.RegisteredClaims.Issuer = fmt.Sprintf("secondbrain-%s", tokenType)
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
}

func (m TokenManager) parse(tokenString string, secret []byte, tokenType string) (Claims, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", token.Header["alg"])
		}
		return secret, nil
	}, jwt.WithAudience("secondbrain"), jwt.WithIssuer(fmt.Sprintf("secondbrain-%s", tokenType)))
	if err != nil {
		return Claims{}, err
	}
	if !token.Valid {
		return Claims{}, errors.New("invalid token")
	}
	return claims, nil
}
