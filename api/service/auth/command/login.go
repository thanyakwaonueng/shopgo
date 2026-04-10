package command

import (
	"context"
	"log/slog"

	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/jwt"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Login struct {
	logger     *slog.Logger
	jwtManager jwt.Manager
	domainDb   *gorm.DB
	repoUser   repogeneric.User
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResultLogin struct {
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	AccessTokenExp  int64  `json:"access_token_exp"`
	RefreshTokenExp int64  `json:"refresh_token_exp"`
}

func NewLogin(
	logger *slog.Logger,
	jwtManager jwt.Manager,
	domainDb *gorm.DB,
	repoUser repogeneric.User,
) *Login {
	return &Login{
		logger:     logger,
		jwtManager: jwtManager,
		domainDb:   domainDb,
		repoUser:   repoUser,
	}
}

func (l *Login) Handle(
	ctx context.Context,
	request RequestLogin,
) (ResultLogin, error) {

	// 1. Check if user exists using repoUser.Search
	existingUser, err := l.repoUser.Search(l.domainDb, map[string]interface{}{
		"email": request.Email,
	}, "")

	if err != nil {
		l.logger.Error("Database error during login check", "error", err)
		return ResultLogin{}, customerror.NewInternalErr("Database error")
	}

	if existingUser == nil {
		// Security Tip: Generic error for both email and password issues
		return ResultLogin{}, customerror.NewInternalErr("Invalid email or password")
	}

	// 2. Verify Password (Compare hash with plain text)
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(request.Password))
	if err != nil {
		return ResultLogin{}, customerror.NewInternalErr("Invalid email or password")
	}

	// 3. Generate Tokens
	loginToken, err := l.jwtManager.GenerateLoginToken(
		existingUser.ID,
		string(existingUser.Role),
		false, // rememberMe
	)
	if err != nil {
		l.logger.Error("Failed to generate login tokens", "error", err)
		return ResultLogin{}, customerror.NewInternalErr("Token generation failed")
	}

	// 4. Build Result
	result := ResultLogin{
		AccessToken:     loginToken.AccessToken,
		AccessTokenExp:  loginToken.AtExpire,
		RefreshToken:    loginToken.RefreshToken,
		RefreshTokenExp: loginToken.RtExpire,
	}

	return result, nil
}
