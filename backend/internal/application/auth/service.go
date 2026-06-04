package auth

import (
	"context"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/security"
)

type Service struct {
	users  domain.UserRepository
	hasher security.PasswordHasher
	tokens security.TokenManager
}

type Session struct {
	User         domain.User `json:"user"`
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
}

func NewService(users domain.UserRepository, hasher security.PasswordHasher, tokens security.TokenManager) Service {
	return Service{users: users, hasher: hasher, tokens: tokens}
}

func (s Service) Register(ctx context.Context, email string, name string, password string) (Session, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" || len(password) < 8 {
		return Session{}, domain.ErrInvalidInput
	}
	hash, err := s.hasher.Hash(password)
	if err != nil {
		return Session{}, err
	}
	user, err := s.users.Create(ctx, domain.User{
		ID:           uuid.New(),
		Email:        email,
		Name:         strings.TrimSpace(name),
		PasswordHash: hash,
	})
	if err != nil {
		return Session{}, err
	}
	return s.createSession(ctx, user)
}

func (s Service) Login(ctx context.Context, email string, password string) (Session, error) {
	user, err := s.users.FindByEmail(ctx, strings.ToLower(strings.TrimSpace(email)))
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return Session{}, domain.ErrUnauthorized
		}
		return Session{}, err
	}
	if !s.hasher.Compare(user.PasswordHash, password) {
		return Session{}, domain.ErrUnauthorized
	}
	return s.createSession(ctx, user)
}

func (s Service) Refresh(ctx context.Context, refreshToken string) (Session, error) {
	claims, err := s.tokens.ParseRefresh(refreshToken)
	if err != nil {
		return Session{}, domain.ErrUnauthorized
	}
	hash := s.tokens.HashToken(refreshToken)
	stored, err := s.users.FindRefreshToken(ctx, hash)
	if err != nil {
		return Session{}, domain.ErrUnauthorized
	}
	if stored.RevokedAt != nil {
		return Session{}, domain.ErrUnauthorized
	}
	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return Session{}, err
	}
	if err := s.users.RevokeRefreshToken(ctx, stored.ID); err != nil {
		return Session{}, err
	}
	return s.createSession(ctx, user)
}

func (s Service) Logout(ctx context.Context, refreshToken string) error {
	hash := s.tokens.HashToken(refreshToken)
	stored, err := s.users.FindRefreshToken(ctx, hash)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil
		}
		return err
	}
	return s.users.RevokeRefreshToken(ctx, stored.ID)
}

func (s Service) createSession(ctx context.Context, user domain.User) (Session, error) {
	pair, err := s.tokens.IssuePair(ctx, user.ID)
	if err != nil {
		return Session{}, err
	}
	if err := s.users.StoreRefreshToken(ctx, domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    user.ID,
		TokenHash: s.tokens.HashToken(pair.RefreshToken),
		ExpiresAt: pair.RefreshExpiry,
	}); err != nil {
		return Session{}, err
	}
	user.PasswordHash = ""
	return Session{User: user, AccessToken: pair.AccessToken, RefreshToken: pair.RefreshToken}, nil
}
