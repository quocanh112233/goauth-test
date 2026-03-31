package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/quocanh112233/goauth-test/gin/internal/model"
	"github.com/quocanh112233/goauth-test/gin/internal/pkg/jwt"
	"github.com/quocanh112233/goauth-test/gin/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Signup(ctx context.Context, name, email, phone, password string) error
	Login(ctx context.Context, email, password string) (string, string, string, error)
	Logout(ctx context.Context, refreshToken string) error
	RefreshAccessToken(ctx context.Context, refreshToken string) (string, error)
	GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error)
}

type authService struct {
	userRepo    repository.UserRepository
	sessionRepo repository.SessionRepository
	jwtSecret   string
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		jwtSecret:   jwtSecret,
	}
}

func (s *authService) Signup(ctx context.Context, name, email, phone, password string) error {
	// Check existing
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return errors.New("email already exists")
	}
	existingPhone, _ := s.userRepo.FindByPhone(ctx, phone)
	if existingPhone != nil {
		return errors.New("phone number already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	user := &model.User{
		Name:     name,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		Role:     "user",
		Provider: "local",
	}

	return s.userRepo.Create(ctx, user)
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return "", "", "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", "", "", errors.New("invalid email or password")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.Hex(), user.Role, s.jwtSecret)
	if err != nil {
		return "", "", "", err
	}

	refreshToken, err := s.generateRandomToken()
	if err != nil {
		return "", "", "", err
	}

	session := &model.Session{
		UserID:       user.ID,
		RefreshToken: refreshToken,
		ExpiredAt:    time.Now().Add(7 * 24 * time.Hour), // 7 days
	}
	if err := s.sessionRepo.Create(ctx, session); err != nil {
		return "", "", "", err
	}

	return accessToken, refreshToken, user.Role, nil
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
}

func (s *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (string, error) {
	session, err := s.sessionRepo.FindByRefreshToken(ctx, refreshToken)
	if err != nil || session == nil {
		return "", errors.New("invalid refresh token")
	}

	if time.Now().After(session.ExpiredAt) {
		_ = s.sessionRepo.DeleteByRefreshToken(ctx, refreshToken)
		return "", errors.New("refresh token expired")
	}

	user, err := s.userRepo.FindByID(ctx, session.UserID)
	if err != nil || user == nil {
		return "", errors.New("user not found")
	}

	return jwt.GenerateAccessToken(user.ID.Hex(), user.Role, s.jwtSecret)
}

func (s *authService) GetUserByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *authService) generateRandomToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
