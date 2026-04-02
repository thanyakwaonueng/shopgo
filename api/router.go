package api

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/middleware"

    "github.com/gofiber/fiber/v2"
)

// Register sets up all API routes with appropriate middleware
func Register(
    app *fiber.App,
    logger *slog.Logger,
    jwtManager jwt.Manager,
    mid *middleware.FiberMiddleware,
) {
    api := app.Group("/api/v1")
    registerPublicRoutes(api, logger)
    registerProtectedRoutes(api, logger, mid)
}

func registerPublicRoutes(
    api fiber.Router,
    logger *slog.Logger,
) {
    //not implement yet
}

func registerProtectedRoutes(
    api fiber.Router,
    logger *slog.Logger,
    mid *middleware.FiberMiddleware,
) {
    //not implement yet
}
