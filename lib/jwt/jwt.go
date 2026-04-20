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
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
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
    ExtractRefreshToken(tokenStr string) (ClaimLoginRefresh, error)
}

type manager struct {
	loginConfig loginConfig
	logger      *slog.Logger
}

type loginConfig struct {
	AccessExpMins  int
	RefreshExpMins int
	RefreshExpMinsRememberMe int
	AccessSecret   string
	RefreshSecret  string
}

func New(logger *slog.Logger) Manager {
	return &manager{
		loginConfig: loginConfig{
			AccessExpMins:       environment.GetInt(environment.LoginAccessExpMinsKey),
			RefreshExpMins:      environment.GetInt(environment.LoginRefreshExpMinsKey),
            RefreshExpMinsRememberMe: environment.GetInt(environment.LoginRefreshExpMinsRememberMeKey),
			AccessSecret:             environment.GetString(environment.LoginAccessSecretKey),
			RefreshSecret:            environment.GetString(environment.LoginRefreshSecretKey),
		},
		logger: logger,
	}
}

func (m *manager) GenerateLoginToken(userId uuid.UUID, role string, rememberMe bool) (LoginToken, error) {
	now := time.Now()

	if userId == uuid.Nil {
        errMsg := "invalid parameter to generate token"
        m.logger.Error(errMsg)
        return LoginToken{}, errors.New(errMsg)
	}

	accessSecret := m.loginConfig.AccessSecret
	refreshSecret := m.loginConfig.RefreshSecret

	if accessSecret == "" || refreshSecret == "" {
        errMsg := "token secret from environment is empty"
        m.logger.Error(errMsg)
        return LoginToken{}, errors.New(errMsg)
	}
    
    // set LoginAccessClaims struct for access token
	accessExp := now.Add(time.Minute * time.Duration(m.loginConfig.AccessExpMins))
	accessClaims := ClaimLoginAccess{
		UserId:     userId,
        Role:       role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}
    
    // create access token
	accessToken, err := m.createJwt(accessSecret, accessClaims)
	if err != nil {
        m.logger.Error("Cannot create access token", customerror.LogErrorKey, err.Error())
		return LoginToken{}, err
	}
    
    // set LoginRefreshClaims struct for refresh token
    // Extend refresh token expiration if remember me is enabled
    var refreshExpMins int
    if rememberMe {
        refreshExpMins = m.loginConfig.RefreshExpMinsRememberMe
        m.logger.Info("generating refresh token with remember me", "exp_mins", refreshExpMins)
    } else {
        refreshExpMins = m.loginConfig.RefreshExpMins
    }
    refreshExp := now.Add(time.Duration(refreshExpMins * int(time.Minute)))
	refreshClaims := ClaimLoginRefresh{
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

    // create refresh token
	refreshToken, err := m.createJwt(refreshSecret, refreshClaims)
	if err != nil {
        m.logger.Error("Cannot create refresh token", customerror.LogErrorKey, err.Error())
		return LoginToken{}, err
	}

	token := LoginToken{
		UserId:       userId,
		AccessToken:  accessToken,
		AtExpire:     accessExp.Unix(),
		RefreshToken: refreshToken,
		RtExpire:     refreshExp.Unix(),
	}

    return token, nil
}

func (m *manager) ExtractRefreshToken(tokenStr string) (ClaimLoginRefresh, error) {
	secret := environment.GetString(environment.LoginRefreshSecretKey)
	if secret == "" {
		errMsg := "refresh token secret from environment is empty"
		m.logger.Error(errMsg)
		return ClaimLoginRefresh{}, errors.New(errMsg)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &ClaimLoginRefresh{}, m.validateToken(secret))
	if err != nil {
        m.logger.Error("Cannot parse token", customerror.LogErrorKey, err.Error())
		return ClaimLoginRefresh{}, err
	}

	claims, ok := token.Claims.(*ClaimLoginRefresh)
	if !ok {
		err = errors.New("unknown claims type, cannot proceed")
        m.logger.Error("Cannot get claims", customerror.LogErrorKey, err.Error())
		return ClaimLoginRefresh{}, err
	}

	return *claims, nil
}

func (m *manager) ExtractAccessToken(tokenStr string) (ClaimLoginAccess, error) {
    secret := environment.GetString(environment.LoginAccessSecretKey)
	if secret == "" {
		errMsg := "refresh token secret from environment is empty"
		m.logger.Error(errMsg)
		return ClaimLoginAccess{}, errors.New(errMsg)
	}

	token, err := jwt.ParseWithClaims(tokenStr, &ClaimLoginAccess{}, m.validateToken(m.loginConfig.AccessSecret))
	if err != nil {
        m.logger.Error("Cannot parse token", customerror.LogErrorKey, err.Error())
		return ClaimLoginAccess{}, err
	}

	claims, ok := token.Claims.(*ClaimLoginAccess)
	if !ok {
		err = errors.New("unknown claims type, cannot proceed")
        m.logger.Error("Cannot get claims", customerror.LogErrorKey, err.Error())
		return ClaimLoginAccess{}, errors.New("invalid claims")
	}

	return *claims, nil
}

func (m *manager) GetAccessTokenFromContext(c *fiber.Ctx) (token string, err error) {
	var tokenstr string
	bearToken := c.Get("Authorization")

	if bearToken == "" {
		errMsg := "token is empty"
		m.logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

	if len(bearToken) > 7 && strings.ToUpper(bearToken[0:6]) == "BEARER" {
		tokenstr = bearToken[7:]
		return tokenstr, nil
	} else {
		errMsg := "invalid authorize token header"
		m.logger.Error(errMsg)
		return "", errors.New(errMsg)
	}

}

func (m *manager) createJwt(secret string, claims jwt.Claims) (string, error) {
	token, err := jwt.
            NewWithClaims(jwt.SigningMethodHS256, claims).
            SignedString([]byte(secret))
  	if err != nil {
		m.logger.Error("jwt token Cannot be generated", customerror.LogErrorKey, err.Error())
		return token, err
	}

	return token, nil 
}

func (m *manager) validateToken(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			m.logger.Error("Cannot parse token", customerror.LogErrorKey, err.Error())
			return nil, err
		}

		return []byte(secret), nil
	}
}
