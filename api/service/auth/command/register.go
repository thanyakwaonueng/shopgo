package command

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/thanyakwaonueng/shopgo/lib/database/entity"
	"github.com/thanyakwaonueng/shopgo/lib/jwt"
	repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/lib/util/customerror"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Register struct {
	logger      *slog.Logger
	jwtManager  jwt.Manager
	domainDb    *gorm.DB
	repoUser    repogeneric.User 
}

type RequestRegister struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResultRegister struct {
	User            UserResponse `json:"user"`
	AccessToken     string       `json:"access_token"`
	RefreshToken    string       `json:"refresh_token"`
	AccessTokenExp  int64        `json:"access_token_exp"`
	RefreshTokenExp int64        `json:"refresh_token_exp"`
}

type UserResponse struct {
	Id    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

func NewRegister(
	logger *slog.Logger,
	jwtManager jwt.Manager,
	domainDb *gorm.DB,
	repoUser repogeneric.User, 
) *Register {
	return &Register{
		logger:     logger,
		jwtManager: jwtManager,
		domainDb:   domainDb,
		repoUser:   repoUser,
	}
}

func (r *Register) Handle(
	ctx context.Context,
	request RequestRegister,
) (ResultRegister, error) {

	// 1. Check if user already exists using repoUser.Search
	existingUser, err := r.repoUser.Search(r.domainDb, map[string]interface{}{
		"email": request.Email,
	}, "")

	if err != nil {
		return ResultRegister{}, customerror.NewInternalErr("Database error during check")
	}

	if existingUser != nil {
		// Requirement: Error if email exists
		return ResultRegister{}, customerror.NewInternalErr("Email already registered")
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return ResultRegister{}, customerror.NewInternalErr("Failed to process password")
	}

	var result ResultRegister

	// 3. The Transaction
	err = r.domainDb.Transaction(func(tx *gorm.DB) error {

		// 3a. Prepare User entity
		newUser := &entity.User{
			ID:           uuid.New(),
			Email:        request.Email,
			PasswordHash: string(hashedPassword),
			Name:         request.Name,
			Role:         entity.RoleCustomer,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// 3b. Create User record using repoUser.Create
		if err := r.repoUser.Create(tx, newUser); err != nil {
			return customerror.NewInternalErr("Database save failed")
		}

		// 4. Generate Tokens
		loginToken, err := r.jwtManager.GenerateLoginToken(
			newUser.ID,
			string(newUser.Role),
			false,
		)
		if err != nil {
			return customerror.NewInternalErr("Token generation failed")
		}

		// 5. Build Result
		result = ResultRegister{
			User: UserResponse{
				Id:    newUser.ID,
				Name:  newUser.Name,
				Email: newUser.Email,
				Role:  string(newUser.Role),
			},
			AccessToken:     loginToken.AccessToken,
			AccessTokenExp:  loginToken.AtExpire,
			RefreshToken:    loginToken.RefreshToken,
			RefreshTokenExp: loginToken.RtExpire,
		}

		return nil
	})

	if err != nil {
		// Use the rescue pattern for the transaction error
		customErr := customerror.UnmarshalError(err)
		if customErr.Message == "" {
			customErr.Message = err.Error()
		}
		return ResultRegister{}, customErr
	}

	return result, nil
}
