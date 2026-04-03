package api

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/middleware"

    handlerauth "github.com/thanyakwaonueng/shopgo/api/handler/auth"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
)

// Register sets up all API routes with appropriate middleware
func Register(
    app *fiber.App,
    logger *slog.Logger,
    validate *validator.Validate,
    jwtManager jwt.Manager,
    mid *middleware.FiberMiddleware,
) {
    api := app.Group("/api/v1")

    registerPublicRoutes(api, logger, validate)
    registerProtectedRoutes(api, logger, validate, mid)
}

func registerPublicRoutes(
    api fiber.Router,
    logger *slog.Logger,
    validate *validator.Validate,
) {
    //not implement yet
    
    groupAuth := api.Group("/auth")
    {
        groupAuth.Post("/register", handlerauth.Register(logger, validate))
    }
}

func registerProtectedRoutes(
    api fiber.Router,
    logger *slog.Logger,
    validate *validator.Validate,
    mid *middleware.FiberMiddleware,
) {
    //not implement yet
}
