package jwt

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/environment"
)

type LoginToken struct {
	AccessToken  string
	RefreshToken string
	AtExpire     int64
	RtExpire     int64
	UserId       uuid.UUID
}

type (
	ClaimUpdatePassword struct {
		UserId uuid.UUID `json:"user_id"`
		jwt.RegisteredClaims
	}

	// Requirement: Only user_id and role
	ClaimLoginAccess struct {
		UserId uuid.UUID `json:"user_id"`
		Role   string    `json:"role"`
		jwt.RegisteredClaims
	}

	// Germini told me this: Only user_id (Refresh tokens don't usually need roles)
	ClaimLoginRefresh struct {
		UserId uuid.UUID `json:"user_id"`
		jwt.RegisteredClaims
	}
)

type Manager interface {
	GenerateLoginToken(userId uuid.UUID, role string, rememberMe bool) (LoginToken, error)
	ExtractAccessToken(tokenStr string) (ClaimLoginAccess, error)
	GetAccessTokenFromContext(c *fiber.Ctx) (token string, err error)
}

type manager struct {
	loginConfig loginConfig
	logger      *slog.Logger
}

type loginConfig struct {
	AccessExpMins  int
	RefreshExpMins int
	AccessSecret   string
	RefreshSecret  string
}

func New(logger *slog.Logger) Manager {
	return &manager{
		loginConfig: loginConfig{
			AccessExpMins:       environment.GetInt(environment.LoginAccessExpMinsKey),
			RefreshExpMins:      environment.GetInt(environment.LoginRefreshExpMinsKey),
			AccessSecret:             environment.GetString(environment.LoginAccessSecretKey),
			RefreshSecret:            environment.GetString(environment.LoginRefreshSecretKey),
		},
		logger: logger,
	}
}

func (m *manager) GenerateLoginToken(userId uuid.UUID, role string, rememberMe bool) (LoginToken, error) {
	now := time.Now()

	if userId == uuid.Nil {
		return LoginToken{}, errors.New("invalid user id")
	}

	accessSecret := m.loginConfig.AccessSecret
	refreshSecret := m.loginConfig.RefreshSecret

	if accessSecret == "" || refreshSecret == "" {
		return LoginToken{}, errors.New("jwt secrets are missing")
	}

	accessExp := now.Add(time.Minute * time.Duration(m.loginConfig.AccessExpMins))
	accessClaims := ClaimLoginAccess{
		UserId:     userId,
        Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	accessToken, err := m.createJwt(accessSecret, accessClaims)
	if err != nil {
		return LoginToken{}, err
	}

	refreshExp := now.Add(time.Minute * time.Duration(m.loginConfig.RefreshExpMins))
	refreshClaims := ClaimLoginRefresh{
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	refreshToken, err := m.createJwt(refreshSecret, refreshClaims)
	if err != nil {
		return LoginToken{}, err
	}

	return LoginToken{
		UserId:       userId,
		AccessToken:  accessToken,
		AtExpire:     accessExp.Unix(),
		RefreshToken: refreshToken,
		RtExpire:     refreshExp.Unix(),
	}, nil
}

func (m *manager) ExtractAccessToken(tokenStr string) (ClaimLoginAccess, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &ClaimLoginAccess{}, m.validateToken(m.loginConfig.AccessSecret))
	if err != nil {
		return ClaimLoginAccess{}, err
	}

	claims, ok := token.Claims.(*ClaimLoginAccess)
	if !ok {
		return ClaimLoginAccess{}, errors.New("invalid claims")
	}

	return *claims, nil
}

func (m *manager) GetAccessTokenFromContext(c *fiber.Ctx) (string, error) {
	bearToken := c.Get("Authorization")
	if len(bearToken) > 7 && strings.ToUpper(bearToken[0:6]) == "BEARER" {
		return bearToken[7:], nil
	}
	return "", errors.New("invalid token format")
}

func (m *manager) createJwt(secret string, claims jwt.Claims) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

func (m *manager) validateToken(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	}
}
