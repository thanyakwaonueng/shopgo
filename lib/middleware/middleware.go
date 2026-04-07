package middleware

import (
//	"encoding/base64"
//	"fmt"
	"log/slog"
	"os"
//	"runtime/debug"
	"strconv"
	"strings"
	"sync"

    "github.com/thanyakwaonueng/shopgo/lib/environment"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"


	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type FiberMiddleware struct {
	corsSetUp      corsSetUp
	logger         *slog.Logger
	jwtManager     jwt.Manager
}

type corsSetUp struct {
	AllowOrigins     string
	AllowCredentials bool
}

// a mutex for synchronizing access to the fiberMiddlewareInstance variable
var fiberMiddlewareLock = &sync.Mutex{}

// a singleton instance of the FiberMiddleware struct
var fiberMiddlewareInstance *FiberMiddleware

// return the singleton instance of the FiberMiddleware
func getFiberMiddlewareInstance(
	logger *slog.Logger,
	jwtManager jwt.Manager,
) *FiberMiddleware {
	if fiberMiddlewareInstance == nil {
		fiberMiddlewareLock.Lock()
		defer fiberMiddlewareLock.Unlock()
		if fiberMiddlewareInstance == nil {
			fiberMiddlewareInstance = createFiberMiddlewareInstance(
				logger,
				jwtManager,
			)
		}
	}

	return fiberMiddlewareInstance
}

// new the fiberMiddlewareInstance and return it out
func NewFiberMiddleware(
	logger *slog.Logger,
	jwtManager jwt.Manager,
) *FiberMiddleware {
	return getFiberMiddlewareInstance(logger, jwtManager)
}

// create the fiberMiddlewareInstance and set up it
func createFiberMiddlewareInstance(
	logger *slog.Logger,
	jwtManager jwt.Manager,
) *FiberMiddleware {
	allowCredential, err := strconv.ParseBool(environment.GetString(environment.AllowCredentialKey))
	if err != nil {
		//message := "Failed to set CORS config"
		//logger.Error(message, customerror.LogPanicKey, err.Error())
		os.Exit(1)
	}

	return &FiberMiddleware{
		corsSetUp: corsSetUp{
			AllowOrigins:     environment.GetString(environment.AllowOriginKey),
			AllowCredentials: allowCredential,
		},
		logger:         logger,
		jwtManager:     jwtManager,
	}
}

// allows servers to specify who can access its resources and what resources can access
func (f *FiberMiddleware) CORS() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     f.corsSetUp.AllowOrigins,
		AllowCredentials: f.corsSetUp.AllowCredentials,
	})
}

func (f *FiberMiddleware) Authenticated() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Get the token from the Header (Bearer <token>)
		tokenStr, err := f.jwtManager.GetAccessTokenFromContext(c)
		if err != nil {
			f.logger.Error("Missing or invalid auth header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Authorization header is required",
			})
		}

		// 2. Validate the signature and extract claims
		claims, err := f.jwtManager.ExtractAccessToken(tokenStr)
		if err != nil {
			f.logger.Error("JWT Validation failed", "error", err.Error())
			
			// Check if expired or just invalid
			status := fiber.StatusUnauthorized
			msg := "Invalid token"
			if strings.Contains(err.Error(), "expired") {
				msg = "Token has expired"
			}

			return c.Status(status).JSON(fiber.Map{
				"message": msg,
			})
		}

		// 3. Store user info in Context (Locals)
		// This allows your handlers to do: c.Locals("userId")
		c.Locals("userId", claims.UserId)
        c.Locals("userRole", claims.Role)

		// 4. Everything is good, go to the next handler!
		return c.Next()
	}
}

// This assume Authenticated() is called
func (f *FiberMiddleware) AdminOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Retrieve the role that Authenticated() stored in Locals
		role, ok := c.Locals("userRole").(string)

		// 2. If it's not there or it's not "admin", block it
		// Use StatusForbidden (403) because we know who they are,
		// but they don't have permission.
		if !ok || role != "admin" {
			f.logger.Warn("Unauthorized access attempt",
				"path", c.Path(),
				"role", role,
			)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "Access denied: Admin role required",
			})
		}

		// 3. User is an admin, let them through!
		return c.Next()
	}
}
