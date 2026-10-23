package app

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 12

var (
	ErrEmailTaken   = errors.New("email taken")
	ErrCredentials  = errors.New("invalid credentials")
	ErrUnauthorized = errors.New("unauthorized")
	ErrTaskNotFound = errors.New("task not found")
	ErrRewardLocked = errors.New("reward locked")
	ErrTaskDone     = errors.New("task done")
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Coins        int
	CreatedAt    time.Time
}

type Session struct {
	User         User
	AccessToken  string
	RefreshToken string
}

type Refresh struct {
	UserID    string
	ExpiresAt time.Time
	Revoked   bool
}

type Store interface {
	CreateUser(ctx context.Context, email, passwordHash string) (User, error)
	UserByEmail(ctx context.Context, email string) (User, error)
	UserByID(ctx context.Context, id string) (User, error)
	SaveRefresh(ctx context.Context, userID, tokenHash string, expiresAt time.Time) error
	UseRefresh(ctx context.Context, tokenHash, nextHash string, expiresAt time.Time) (string, error)
	RevokeRefresh(ctx context.Context, userID, tokenHash string) error
	CreateTask(ctx context.Context, userID, title, description string, reward int) (Task, error)
	ListTasks(ctx context.Context, userID, status string, limit, offset int) ([]Task, error)
	TaskByID(ctx context.Context, userID, taskID string) (Task, error)
	UpdateTask(ctx context.Context, userID, taskID string, patch TaskPatch) (Task, error)
	DeleteTask(ctx context.Context, userID, taskID string) error
	CompleteTask(ctx context.Context, userID, taskID string) (Task, error)
	ListTransactions(ctx context.Context, userID string, limit, offset int) ([]Transaction, error)
}

type Service struct {
	store      Store
	secret     []byte
	issuer     string
	accessTTL  time.Duration
	refreshTTL time.Duration
	cost       int
	now        func() time.Time
	dummy      []byte
}

func New(store Store, secret []byte, issuer string, accessTTL, refreshTTL time.Duration) *Service {
	return newService(store, secret, issuer, accessTTL, refreshTTL, bcryptCost)
}

func NewForTest(store Store, secret []byte, issuer string, accessTTL, refreshTTL time.Duration, cost int) *Service {
	return newService(store, secret, issuer, accessTTL, refreshTTL, cost)
}

func newService(store Store, secret []byte, issuer string, accessTTL, refreshTTL time.Duration, cost int) *Service {
	dummy, err := bcrypt.GenerateFromPassword([]byte("dummy-password"), cost)
	if err != nil {
		panic(err)
	}
	return &Service{
		store:      store,
		secret:     append([]byte(nil), secret...),
		issuer:     issuer,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		cost:       cost,
		now:        time.Now,
		dummy:      dummy,
	}
}

func (s *Service) Register(ctx context.Context, email, password string) (Session, error) {
	email = normalizeEmail(email)
	hash, err := bcrypt.GenerateFromPassword([]byte(password), s.cost)
	if err != nil {
		return Session{}, err
	}
	user, err := s.store.CreateUser(ctx, email, string(hash))
	if err != nil {
		return Session{}, err
	}
	return s.issue(ctx, user)
}

func (s *Service) Login(ctx context.Context, email, password string) (Session, error) {
	email = normalizeEmail(email)
	user, err := s.store.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			_ = bcrypt.CompareHashAndPassword(s.dummy, []byte(password))
			return Session{}, ErrCredentials
		}
		return Session{}, err
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return Session{}, ErrCredentials
	}
	return s.issue(ctx, user)
}

func (s *Service) Refresh(ctx context.Context, refresh string) (Session, error) {
	next, nextHash, exp, err := s.newRefresh()
	if err != nil {
		return Session{}, err
	}
	userID, err := s.store.UseRefresh(ctx, hashToken(refresh), nextHash, exp)
	if err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotFound) {
			return Session{}, ErrUnauthorized
		}
		return Session{}, err
	}
	user, err := s.store.UserByID(ctx, userID)
	if err != nil {
		return Session{}, ErrUnauthorized
	}
	access, err := s.accessToken(user.ID)
	if err != nil {
		return Session{}, err
	}
	return Session{User: user, AccessToken: access, RefreshToken: next}, nil
}

func (s *Service) Logout(ctx context.Context, access, refresh string) error {
	userID, err := s.UserID(access)
	if err != nil {
		return ErrUnauthorized
	}
	if err := s.store.RevokeRefresh(ctx, userID, hashToken(refresh)); err != nil {
		if errors.Is(err, ErrUnauthorized) || errors.Is(err, ErrNotFound) {
			return ErrUnauthorized
		}
		return err
	}
	return nil
}

func (s *Service) Profile(ctx context.Context, access string) (User, error) {
	userID, err := s.UserID(access)
	if err != nil {
		return User{}, ErrUnauthorized
	}
	user, err := s.store.UserByID(ctx, userID)
	if err != nil {
		return User{}, ErrUnauthorized
	}
	user.PasswordHash = ""
	return user, nil
}

func (s *Service) UserID(access string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(access, claims, func(token *jwt.Token) (any, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, ErrUnauthorized
		}
		return s.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(s.issuer),
		jwt.WithAudience(s.issuer),
		jwt.WithExpirationRequired())
	if err != nil || !token.Valid || claims.Subject == "" {
		return "", ErrUnauthorized
	}
	return claims.Subject, nil
}

func (s *Service) issue(ctx context.Context, user User) (Session, error) {
	refresh, refreshHash, exp, err := s.newRefresh()
	if err != nil {
		return Session{}, err
	}
	if err := s.store.SaveRefresh(ctx, user.ID, refreshHash, exp); err != nil {
		return Session{}, err
	}
	access, err := s.accessToken(user.ID)
	if err != nil {
		return Session{}, err
	}
	user.PasswordHash = ""
	return Session{User: user, AccessToken: access, RefreshToken: refresh}, nil
}

func (s *Service) accessToken(userID string) (string, error) {
	now := s.now()
	claims := jwt.RegisteredClaims{
		Subject:   userID,
		Issuer:    s.issuer,
		Audience:  jwt.ClaimStrings{s.issuer},
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(s.accessTTL)),
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.secret)
}

func (s *Service) newRefresh() (string, string, time.Time, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", time.Time{}, err
	}
	token := base64.RawURLEncoding.EncodeToString(buf)
	return token, hashToken(token), s.now().Add(s.refreshTTL), nil
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func normalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

var ErrNotFound = errors.New("not found")
