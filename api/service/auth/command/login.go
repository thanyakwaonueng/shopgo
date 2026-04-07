package command

import (
    "context"
    "errors"
    "log/slog"

    "github.com/thanyakwaonueng/shopgo/lib/database/entity"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)


type Login struct {
    logger                      *slog.Logger
    jwtManager                   jwt.Manager
    domainDb                     *gorm.DB
}


type RequestLogin struct {
    Email           string `json:"email" validate:"required"`
    Password        string `json:"password" validate:"required"`
}


type ResultLogin struct {
	AccessToken     string       `json:"access_token"`
	RefreshToken    string       `json:"refresh_token"`
	AccessTokenExp  int64        `json:"access_token_exp"`
	RefreshTokenExp int64        `json:"refresh_token_exp"`
}


func NewLogin(
    logger *slog.Logger,
    jwtManager jwt.Manager,
    domainDb *gorm.DB,
) *Login {
    return &Login{
        logger:                         logger,
        jwtManager:                     jwtManager,
        domainDb:                       domainDb,
    }
}


func (l *Login) Handle(
	ctx context.Context,
	request RequestLogin,
) (ResultLogin, error) {


	// 1. Check if user already exists
	var existingUser entity.User
	err := l.domainDb.Where("email = ?", request.Email).First(&existingUser).Error

    if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Security Tip: Use a generic error so attackers don't know if the email exists
			return ResultLogin{}, customerror.NewInternalErr("Invalid email or password")
		}
		l.logger.Error("Database error during login", "error", err)
		return ResultLogin{}, customerror.NewInternalErr("Database error")
	}

    // 2. Verify Password (Compare hash with plain text)
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(request.Password))
	if err != nil {
		return ResultLogin{}, customerror.NewInternalErr("Invalid email or password")
	}
    
    // 3. Generate Tokens using your updated clean JWT Manager
	// Note: We removed tenantId and clientId as per your requirement
	loginToken, err := l.jwtManager.GenerateLoginToken(
		existingUser.ID,
		string(existingUser.Role),
		false, // rememberMe
	)
	if err != nil {
		l.logger.Error("Failed to generate login tokens", "error", err)
		return ResultLogin{}, customerror.NewInternalErr("Token generation failed")
	}

	// 4. Build ResultLogin (Slim format: only tokens and expiration)
	result := ResultLogin{
		AccessToken:     loginToken.AccessToken,
		AccessTokenExp:  loginToken.AtExpire,
		RefreshToken:    loginToken.RefreshToken,
		RefreshTokenExp: loginToken.RtExpire,
	}

	return result, nil
}
