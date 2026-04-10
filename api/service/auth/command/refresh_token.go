package command

import (
	"context"
	"log/slog"

	"github.com/thanyakwaonueng/shopgo/api/repository/custom"
	"github.com/thanyakwaonueng/shopgo/lib/jwt"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type RefreshToken struct {
	logger           *slog.Logger
	jwtManager       jwt.Manager
	domainDb         *gorm.DB
	repoUserRoleById repocustom.UserRoleById
}

type RequestRefreshToken struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type ResultRefreshToken struct {
	AccessToken    string `json:"access_token"`
	AccessTokenExp int64  `json:"access_token_exp"`
}

func NewRefreshToken(
	logger *slog.Logger,
	jwtManager jwt.Manager,
	domainDb *gorm.DB,
	repoUserRoleById repocustom.UserRoleById,
) *RefreshToken {
	return &RefreshToken{
		logger:           logger,
		jwtManager:       jwtManager,
		domainDb:         domainDb,
		repoUserRoleById: repoUserRoleById,
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

	// 2. Quick DB check using Custom Repository
	userRole, err := r.repoUserRoleById.Execute(r.domainDb, refreshClaims.UserId)
	if err != nil {
		// If user doesn't exist or DB is down, we deny the rotation
		return ResultRefreshToken{}, customerror.NewInternalErr("User account no longer active")
	}

	// 3. Generate New Pair (Rotation)
	newTokens, err := r.jwtManager.GenerateLoginToken(
		refreshClaims.UserId,
		userRole,
		false,
	)
	if err != nil {
		return ResultRefreshToken{}, customerror.NewInternalErr("Failed to rotate tokens")
	}

	return ResultRefreshToken{
		AccessToken:    newTokens.AccessToken,
		AccessTokenExp: newTokens.AtExpire,
	}, nil
}
