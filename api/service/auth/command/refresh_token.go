package command

import (
	"context"
	"log/slog"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/jwt"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"

	"gorm.io/gorm"
)

type RefreshToken struct {
	logger                  *slog.Logger
	jwtManager              jwt.Manager
	domainDb                *gorm.DB
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ResultRefreshToken struct {
	AccessToken     string    `json:"access_token"`
	AccessTokenExp  int64     `json:"access_token_exp"`
}

func NewRefreshToken(
	logger *slog.Logger,
	jwtManager jwt.Manager,
	domainDb *gorm.DB,
) *RefreshToken {
	return &RefreshToken{
		logger:                  logger,
		jwtManager:              jwtManager,
		domainDb:                domainDb,
	}
}

func (r *RefreshToken) Handle(
    ctx context.Context,
    request RequestRefreshToken,
) (ResultRefreshToken, error) {
    // 1. Validate the Refresh Token and extract claims
    refreshClaims, err := r.jwtManager.ExtractRefreshToken(request.RefreshToken)
    if err != nil {
        r.logger.Error("Refresh failed: invalid signature or expired", "error", err)
        return ResultRefreshToken{}, customerror.NewInternalErr("Session expired, please login again")
    }

    // 2. Optional: Quick DB check (Recommended by mr.germini)
    // Even in stateless apps, it's good to check if the user still exists
    // before giving them a fresh 15-minute Access Token.
    var userRole string
    err = r.domainDb.Model(&entity.User{}).
        Select("role").
        Where("id = ?", refreshClaims.UserId).
        First(&userRole).Error

    if err != nil {
        return ResultRefreshToken{}, customerror.NewInternalErr("User account no longer active")
    }

    // 3. Generate New Pair (Rotation)
    // We pass the role we just fetched to ensure the token is up-to-date
    newTokens, err := r.jwtManager.GenerateLoginToken(
        refreshClaims.UserId,
        userRole,
        false, // Or implement the duration-check logic from before
    )
    if err != nil {
        return ResultRefreshToken{}, customerror.NewInternalErr("Failed to rotate tokens")
    }

    return ResultRefreshToken{
        AccessToken:     newTokens.AccessToken,
        AccessTokenExp:  newTokens.AtExpire,
    }, nil
}
