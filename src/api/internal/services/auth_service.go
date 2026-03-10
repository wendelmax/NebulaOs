package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/wendelmax/nebulaos/src/api/domain"
	"golang.org/x/crypto/bcrypt"
)

type InternalIdentityManager struct {
	userRepo domain.UserRepository
	jwtKey   []byte
}

func NewInternalIdentityManager(userRepo domain.UserRepository, jwtKey string) *InternalIdentityManager {
	return &InternalIdentityManager{
		userRepo: userRepo,
		jwtKey:   []byte(jwtKey),
	}
}

func (m *InternalIdentityManager) Authenticate(ctx context.Context, username, password string) (string, error) {
	user, err := m.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":                  user.ID,
		"username":             user.Username,
		"tenant_id":            user.TenantID,
		"must_change_password": user.MustChangePassword,
		"exp":                  time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(m.jwtKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

func (m *InternalIdentityManager) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.jwtKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, _ := claims["sub"].(string)
		user, err := m.userRepo.GetByID(ctx, userID)
		if err != nil {
			return nil, fmt.Errorf("user not found: %w", err)
		}
		return user, nil
	}

	return nil, errors.New("invalid token claims")
}

func (m *InternalIdentityManager) ChangePassword(ctx context.Context, userID, oldPassword, newPassword, email string) error {
	user, err := m.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(oldPassword)); err != nil {
		return errors.New("invalid old password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.PasswordHash = string(hashedPassword)
	user.MustChangePassword = false
	if email != "" {
		user.Email = email
	}

	return m.userRepo.Update(ctx, user)
}

func (m *InternalIdentityManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
