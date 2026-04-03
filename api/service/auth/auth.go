package serviceauth

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/api/service/auth/command" 
    "github.com/thanyakwaonueng/shopgo/lib/jwt" 

    "github.com/mehdihadeli/go-mediatr"
    "gorm.io/gorm"
)

func Register(
    domainDb *gorm.DB,
    logger *slog.Logger,
    jwtManager jwt.Manager,
) {
    //Register New User Register Handler
    serviceRegister := command.NewRegister(logger, jwtManager, domainDb)
    err := mediatr.RegisterRequestHandler(serviceRegister)
    if err != nil {
        panic(err)
    }
}
