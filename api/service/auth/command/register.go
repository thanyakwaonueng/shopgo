package command

import (
    "context"
    "errors"
    "log/slog"
    "time"

    "github.com/google/uuid"
    "github.com/thanyakwaonueng/shopgo/lib/database/entity"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)

type Register struct {
    logger                      *slog.Logger
    jwtManager                   jwt.Manager
    domainDb                     *gorm.DB
}

type RequestRegister struct {
    Name            string `json:"name" validate:"required"`
    Email           string `json:"email" validate:"required"`
    Password        string `json:"password" validate:"required"`
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
) *Register {
    return &Register{
        logger:                         logger,
        jwtManager:                     jwtManager,
        domainDb:                       domainDb,
    }
}


func (r *Register) Handle(
	ctx context.Context,
	request RequestRegister,
) (ResultRegister, error) {

	// 1. Check if user already exists
	var existingUser entity.User
	err := r.domainDb.Where("email = ?", request.Email).First(&existingUser).Error

	if err == nil {
		return ResultRegister{}, customerror.NewInternalErr("Email already registered")
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		r.logger.Error("Database error during check", "error", err)
		return ResultRegister{}, customerror.NewInternalErr("Database error")
	}

	// 2. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return ResultRegister{}, customerror.NewInternalErr("Failed to process password")
	}

	var result ResultRegister

	// 3. The Transaction
	err = r.domainDb.Transaction(func(tx *gorm.DB) error {

		// 3a. Create User record
		newUser := entity.User{
			ID:           uuid.New(),
			Email:        request.Email,
			PasswordHash: string(hashedPassword),
			Name:         request.Name,
			Role:         entity.RoleCustomer,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := tx.Create(&newUser).Error; err != nil {
			return err
		}

		// 4. Generate Tokens
		loginToken, err := r.jwtManager.GenerateLoginToken(
			newUser.ID,
            string(newUser.Role),
			false,
		)
		if err != nil {
			return err
		}

		// 5. Build Result to match the nested Postman format
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
		r.logger.Error("Registration transaction failed", "error", err)
		return ResultRegister{}, err
	}

	return result, nil
}
