package serviceauth

import (
    "log/slog"
    repocustom "github.com/thanyakwaonueng/shopgo/api/repository/custom"
    repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
    "github.com/thanyakwaonueng/shopgo/api/service/auth/command" 
    "github.com/thanyakwaonueng/shopgo/api/service/auth/query" 
    "github.com/thanyakwaonueng/shopgo/lib/jwt" 

    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/gorm"
)

func Register(
    domainDb *gorm.DB,
    logger *slog.Logger,
    jwtManager jwt.Manager,
    repoUser repogeneric.User,
    repoUserRoleById repocustom.UserRoleById,
) {
    // Register New User Register Handler
    serviceRegister := command.NewRegister(logger, jwtManager, domainDb, repoUser)
    err := mediatr.RegisterRequestHandler(serviceRegister)
    if err != nil {
        panic(err)
    }

    // Register Login Handler
    serviceLogin := command.NewLogin(logger, jwtManager, domainDb, repoUser)
    err = mediatr.RegisterRequestHandler(serviceLogin)
    if err != nil {
        panic(err)
    }

	// Register RefreshToken Handler
	serviceRefreshToken := command.NewRefreshToken(logger, jwtManager, domainDb, repoUserRoleById)
	err = mediatr.RegisterRequestHandler(serviceRefreshToken)
	if err != nil {
		panic(err)
	}
    
    // Register GetMe Handler
    serviceGetMe := query.NewGetMeHandler(logger, domainDb, repoUser) 
    err = mediatr.RegisterRequestHandler(serviceGetMe)
	if err != nil {
		panic(err)
	}
}
