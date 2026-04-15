package serviceusers

import (
	"log/slog"
    repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
	"github.com/thanyakwaonueng/shopgo/api/service/users/query"
	"github.com/thanyakwaonueng/shopgo/api/service/users/command"
	"github.com/mehdihadeli/go-mediatr"
	"gorm.io/gorm"
)

func Register(
	domainDb *gorm.DB,
	logger *slog.Logger,
    repoUser repogeneric.User,
) {
	// Register GetUsers Handler
	serviceGetUsers := query.NewGetUsersHandler(logger, domainDb, repoUser)
	err := mediatr.RegisterRequestHandler(serviceGetUsers)
	if err != nil {
        panic(err)
	}

    // Register GetUserByID Handler (Single User)
	serviceGetUserByIDHandler := query.NewGetUserByIDHandler(logger, domainDb, repoUser)
	err = mediatr.RegisterRequestHandler(serviceGetUserByIDHandler)
    if err != nil {
		panic(err)
	}

    // Register UpdateUser Command
	serviceUpdateUserRoleHandler := command.NewUpdateUserRoleHandler(logger, domainDb, repoUser)
	err = mediatr.RegisterRequestHandler(serviceUpdateUserRoleHandler)
    if err != nil {
		panic(err)
	}

}

