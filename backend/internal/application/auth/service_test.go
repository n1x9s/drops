package auth

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/n1x9s/second-brain/backend/internal/domain"
	"github.com/n1x9s/second-brain/backend/internal/infrastructure/security"
)

func TestServiceRegisterLoginRefreshLogout(t *testing.T) {
	repo := newMemoryUserRepo()
	service := NewService(repo, security.NewPasswordHasher(), security.NewTokenManager("access", "refresh", time.Minute, time.Hour))
	ctx := context.Background()

	session, err := service.Register(ctx, "USER@example.com", "User", "password123")
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	if session.User.Email != "user@example.com" {
		t.Fatalf("email = %s", session.User.Email)
	}

	login, err := service.Login(ctx, "user@example.com", "password123")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if login.AccessToken == "" || login.RefreshToken == "" {
		t.Fatal("expected tokens")
	}

	refreshed, err := service.Refresh(ctx, login.RefreshToken)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if refreshed.RefreshToken == login.RefreshToken {
		t.Fatal("refresh token should rotate")
	}

	if err := service.Logout(ctx, refreshed.RefreshToken); err != nil {
		t.Fatalf("logout: %v", err)
	}
	if _, err := service.Refresh(ctx, refreshed.RefreshToken); err == nil {
		t.Fatal("revoked token should not refresh")
	}
}

type memoryUserRepo struct {
	mu     sync.Mutex
	users  map[uuid.UUID]domain.User
	email  map[string]uuid.UUID
	tokens map[string]domain.RefreshToken
}

func newMemoryUserRepo() *memoryUserRepo {
	return &memoryUserRepo{
		users:  map[uuid.UUID]domain.User{},
		email:  map[string]uuid.UUID{},
		tokens: map[string]domain.RefreshToken{},
	}
}

func (r *memoryUserRepo) Create(_ context.Context, user domain.User) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.email[user.Email]; ok {
		return domain.User{}, domain.ErrConflict
	}
	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = user.CreatedAt
	r.users[user.ID] = user
	r.email[user.Email] = user.ID
	return user, nil
}

func (r *memoryUserRepo) FindByEmail(_ context.Context, email string) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	id, ok := r.email[email]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return r.users[id], nil
}

func (r *memoryUserRepo) FindByID(_ context.Context, id uuid.UUID) (domain.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	user, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return user, nil
}

func (r *memoryUserRepo) StoreRefreshToken(_ context.Context, token domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	token.CreatedAt = time.Now().UTC()
	r.tokens[token.TokenHash] = token
	return nil
}

func (r *memoryUserRepo) FindRefreshToken(_ context.Context, tokenHash string) (domain.RefreshToken, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	token, ok := r.tokens[tokenHash]
	if !ok || token.ExpiresAt.Before(time.Now().UTC()) {
		return domain.RefreshToken{}, domain.ErrNotFound
	}
	return token, nil
}

func (r *memoryUserRepo) RevokeRefreshToken(_ context.Context, tokenID uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	for hash, token := range r.tokens {
		if token.ID == tokenID {
			now := time.Now().UTC()
			token.RevokedAt = &now
			r.tokens[hash] = token
			return nil
		}
	}
	return domain.ErrNotFound
}
