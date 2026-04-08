package query

import (
    "log/slog"
    "errors"
	"context"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetMe struct {
    logger                      *slog.Logger
    domainDb                     *gorm.DB
}

type RequestGetMe struct {
	UserId uuid.UUID `json:"id"`
}

type ResultGetMe struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}


func NewGetMeHandler(
    logger *slog.Logger,
    domainDb *gorm.DB,
) *GetMe {
    return &GetMe{
        logger:                         logger,
        domainDb:                       domainDb,
    }
}

func (h *GetMe) Handle(
    ctx context.Context, 
    request RequestGetMe,
) (ResultGetMe, error) {
	var user entity.User

	// Direct GORM query to find the user by ID
	// We use First() which returns ErrRecordNotFound if the ID doesn't exist
    err := h.domainDb.Where("id = ?", request.UserId).First(&user).Error
	if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ResultGetMe{}, customerror.NewInternalErr("User profile not found")
        }
        return ResultGetMe{}, customerror.NewInternalErr("Database error")
	}

	// Map the database entity to our Result struct
	result := ResultGetMe{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	}
    return result, nil
}
