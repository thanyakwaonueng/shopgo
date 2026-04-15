package query

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type GetUsers struct {
	logger   *slog.Logger
	domainDb *gorm.DB
	repoUser repogeneric.User
}

type RequestGetUsers struct {
	Page  int
	Limit int
	Q     string
}

type ResultGetUsers struct {
	Items []UserItem `json:"items"`
	Total int64      `json:"total"`
}

type UserItem struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt string    `json:"created_at"`
}

func NewGetUsersHandler(logger *slog.Logger, domainDb *gorm.DB, repoUser repogeneric.User) *GetUsers {
	return &GetUsers{logger: logger, domainDb: domainDb, repoUser: repoUser}
}

func (h *GetUsers) Handle(ctx context.Context, request RequestGetUsers) (ResultGetUsers, error) {
	condition := make(map[string]interface{})
	var queryStr string
	var queryArgs []interface{}

	if request.Q != "" {
		queryStr = "name ILIKE ? OR email ILIKE ?"
		queryArgs = append(queryArgs, "%"+request.Q+"%", "%"+request.Q+"%")
	}

	total, err := h.repoUser.Count(h.domainDb, condition, queryStr, queryArgs)
	if err != nil {
		return ResultGetUsers{}, customerror.NewInternalErr("Failed to retrieve user count")
	}

	offset := (request.Page - 1) * request.Limit
	users, err := h.repoUser.ListWithPagination(h.domainDb, condition, queryStr, queryArgs, "created_at DESC", offset, request.Limit)
	if err != nil {
		return ResultGetUsers{}, customerror.NewInternalErr("Database error while fetching users")
	}

	items := make([]UserItem, len(users))
	for i, u := range users {
		items[i] = UserItem{
			ID:        u.ID,
			Email:     u.Email,
			Name:      u.Name,
			Role:      string(u.Role),
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return ResultGetUsers{Items: items, Total: total}, nil
}
