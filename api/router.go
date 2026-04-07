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
    groupAuth := api.Group("/auth")
    {
        groupAuth.Post("/register", handlerauth.Register(logger, validate))
        groupAuth.Post("/login", handlerauth.Login(logger, validate))
        groupAuth.Post("/refresh", handlerauth.RefreshToken(logger))
    }
}

func registerProtectedRoutes(
    api fiber.Router,
    logger *slog.Logger,
    validate *validator.Validate,
    mid *middleware.FiberMiddleware,
) {
    groupAuth := api.Group("/auth")
    {
        groupAuth.Use(mid.Authenticated())
        groupAuth.Get("/me", handlerauth.GetMe(logger))
    }
}
