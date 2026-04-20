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
    "github.com/thanyakwaonueng/shopgo/lib/util"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"


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
		message := "Failed to set CORS config"
		logger.Error(message, customerror.LogPanicKey, err.Error())
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
            customErr := customerror.New(1, 3, "Insufficient role(Missing or invalid auth header)")
            f.logger.Error(customErr.Message)
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 2. Validate the signature and extract claims
		claims, err := f.jwtManager.ExtractAccessToken(tokenStr)
		if err != nil {
			f.logger.Error("JWT Validation failed", "error", err.Error())
			
			// Check if expired or just invalid
			if strings.Contains(err.Error(), "expired") {
                customErr := customerror.New(1, 2, "Token expired")
                return c.Status(fiber.StatusBadRequest).JSON(customErr)
			}

			//msg := "Missing or invalid token"
            customErr := customerror.New(1, 1, "Missing or invalid token")
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 3. Store user info in Context (Locals)
        userData := util.UserDataCtx{
            UserId: claims.UserId,
            Role: claims.Role,
        }
        util.SetUserDataLocal(c, userData)

		// 4. Everything is good, go to the next handler!
		return c.Next()
	}
}

// This assume Authenticated() is called
func (f *FiberMiddleware) AdminOnly() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// 1. Retrieve the role that Authenticated() stored in Locals
        userData := util.GetUserDataLocal(c)
		role := userData.Role

		// 2. If it's not there or it's not "admin", block it
		if role != "admin" {
            customErr := customerror.New(1, 3, "Insufficient role(Not admin!)")
            f.logger.Error(customErr.Message)
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
		}

		// 3. User is an admin, let them through!
		return c.Next()
	}
}
