package command

import (
	"context"
	"log/slog"

	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"gorm.io/gorm"
)

type CreateCategory struct {
	logger       *slog.Logger
	domainDb     *gorm.DB
	repoCategory repogeneric.Category
}

type RequestCreateCategory struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type ResultCreateCategory struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

func NewCreateCategoryHandler(
	logger *slog.Logger,
	domainDb *gorm.DB,
	repoCategory repogeneric.Category,
) *CreateCategory {
	return &CreateCategory{
		logger:       logger,
		domainDb:     domainDb,
		repoCategory: repoCategory,
	}
}

func (h *CreateCategory) Handle(
	ctx context.Context,
	request RequestCreateCategory,
) (ResultCreateCategory, error) {
	// 1. Prepare the entity
	newCategory := &entity.Category{
		Name: request.Name,
		Slug: request.Slug,
	}

	// 2. Insert into database using the repository
	if err := h.repoCategory.Create(h.domainDb, newCategory); err != nil {
		// Logged inside the repository already
		return ResultCreateCategory{}, customerror.NewInternalErr("Could not create category. Slug might already exist.")
	}

	// 3. Map back to Result DTO
	result := ResultCreateCategory{
		ID:   newCategory.ID,
		Name: newCategory.Name,
		Slug: newCategory.Slug,
	}

	return result, nil
}
