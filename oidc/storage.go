package oidc

import (
	"context"
	"errors"

	"github.com/ory/fosite"
)

// MemoryUserRelation :
type MemoryUserRelation struct {
	Username string
	Password string
}

// MemoryStore :
type MemoryStore struct {
	Clients        map[string]*fosite.DefaultClient
	AuthorizeCodes map[string]fosite.Requester
	IDSessions     map[string]fosite.Requester
	AccessTokens   map[string]fosite.Requester
	Implicit       map[string]fosite.Requester
	RefreshTokens  map[string]fosite.Requester
	Users          map[string]MemoryUserRelation
	// In-memory request ID to token signatures
	AccessTokenRequestIDs  map[string]string
	RefreshTokenRequestIDs map[string]string
}

// NewMemoryStore :
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		Clients:        make(map[string]*fosite.DefaultClient),
		AuthorizeCodes: make(map[string]fosite.Requester),
		IDSessions:     make(map[string]fosite.Requester),
		AccessTokens:   make(map[string]fosite.Requester),
		Implicit:       make(map[string]fosite.Requester),
		RefreshTokens:  make(map[string]fosite.Requester),
		Users:          make(map[string]MemoryUserRelation),
		AccessTokenRequestIDs:  make(map[string]string),
		RefreshTokenRequestIDs: make(map[string]string),
	}
}

// NewExampleStore :
func NewExampleStore() *MemoryStore {
	return &MemoryStore{
		IDSessions: make(map[string]fosite.Requester),
		Clients: map[string]*fosite.DefaultClient{
			"my-client": {
				ID:            "my-client",
				Secret:        []byte(`$2a$10$IxMdI6d.LIRZPpSfEwNoeu4rY3FhDREsxFJXikcgdRRAStxUlsuEO`), // = "foobar"
				RedirectURIs:  []string{"http://localhost:3846/callback"},
				ResponseTypes: []string{"id_token", "code", "token"},
				GrantTypes:    []string{"implicit", "refresh_token", "authorization_code", "password", "client_credentials"},
				Scopes:        []string{"fosite", "openid", "photos", "offline"},
			},
		},
		Users: map[string]MemoryUserRelation{
			"peter": {
				// This store simply checks for equality, a real storage implementation would obviously use
				// a hashing algorithm for encrypting the user password.
				Username: "peter",
				Password: "secret",
			},
		},
		AuthorizeCodes:         map[string]fosite.Requester{},
		Implicit:               map[string]fosite.Requester{},
		AccessTokens:           map[string]fosite.Requester{},
		RefreshTokens:          map[string]fosite.Requester{},
		AccessTokenRequestIDs:  map[string]string{},
		RefreshTokenRequestIDs: map[string]string{},
	}
}

// CreateOpenIDConnectSession :
func (s *MemoryStore) CreateOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) error {
	s.IDSessions[authorizeCode] = requester
	return nil
}

// GetOpenIDConnectSession :
func (s *MemoryStore) GetOpenIDConnectSession(_ context.Context, authorizeCode string, requester fosite.Requester) (fosite.Requester, error) {
	cl, ok := s.IDSessions[authorizeCode]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return cl, nil
}

// DeleteOpenIDConnectSession :
func (s *MemoryStore) DeleteOpenIDConnectSession(_ context.Context, authorizeCode string) error {
	delete(s.IDSessions, authorizeCode)
	return nil
}

// GetClient :
func (s *MemoryStore) GetClient(_ context.Context, id string) (fosite.Client, error) {
	cl, ok := s.Clients[id]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return cl, nil
}

// CreateAuthorizeCodeSession :
func (s *MemoryStore) CreateAuthorizeCodeSession(_ context.Context, code string, req fosite.Requester) error {
	s.AuthorizeCodes[code] = req
	return nil
}

// GetAuthorizeCodeSession :
func (s *MemoryStore) GetAuthorizeCodeSession(_ context.Context, code string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.AuthorizeCodes[code]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

// DeleteAuthorizeCodeSession :
func (s *MemoryStore) DeleteAuthorizeCodeSession(_ context.Context, code string) error {
	delete(s.AuthorizeCodes, code)
	return nil
}

// CreateAccessTokenSession :
func (s *MemoryStore) CreateAccessTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.AccessTokens[signature] = req
	s.AccessTokenRequestIDs[req.GetID()] = signature
	return nil
}

// GetAccessTokenSession :
func (s *MemoryStore) GetAccessTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.AccessTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

// DeleteAccessTokenSession :
func (s *MemoryStore) DeleteAccessTokenSession(_ context.Context, signature string) error {
	delete(s.AccessTokens, signature)
	return nil
}

// CreateRefreshTokenSession :
func (s *MemoryStore) CreateRefreshTokenSession(_ context.Context, signature string, req fosite.Requester) error {
	s.RefreshTokens[signature] = req
	s.RefreshTokenRequestIDs[req.GetID()] = signature
	return nil
}

// GetRefreshTokenSession :
func (s *MemoryStore) GetRefreshTokenSession(_ context.Context, signature string, _ fosite.Session) (fosite.Requester, error) {
	rel, ok := s.RefreshTokens[signature]
	if !ok {
		return nil, fosite.ErrNotFound
	}
	return rel, nil
}

// DeleteRefreshTokenSession :
func (s *MemoryStore) DeleteRefreshTokenSession(_ context.Context, signature string) error {
	delete(s.RefreshTokens, signature)
	return nil
}

// CreateImplicitAccessTokenSession :
func (s *MemoryStore) CreateImplicitAccessTokenSession(_ context.Context, code string, req fosite.Requester) error {
	s.Implicit[code] = req
	return nil
}

// Authenticate :
func (s *MemoryStore) Authenticate(_ context.Context, name string, secret string) error {
	rel, ok := s.Users[name]
	if !ok {
		return fosite.ErrNotFound
	}
	if rel.Password != secret {
		return errors.New("Invalid credentials")
	}
	return nil
}

// RevokeRefreshToken :
func (s *MemoryStore) RevokeRefreshToken(ctx context.Context, requestID string) error {
	if signature, exists := s.RefreshTokenRequestIDs[requestID]; exists {
		s.DeleteRefreshTokenSession(ctx, signature)
		s.DeleteAccessTokenSession(ctx, signature)
	}
	return nil
}

// RevokeAccessToken :
func (s *MemoryStore) RevokeAccessToken(ctx context.Context, requestID string) error {
	if signature, exists := s.AccessTokenRequestIDs[requestID]; exists {
		s.DeleteAccessTokenSession(ctx, signature)
	}
	return nil
}
