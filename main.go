package main

import (
    "log"
    "os"
    "path/filepath"

    "github.com/thanyakwaonueng/shopgo/api"
    "github.com/thanyakwaonueng/shopgo/lib/environment"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/middleware"
    "github.com/thanyakwaonueng/shopgo/lib/logging"

    "github.com/gofiber/fiber/v2"
    "github.com/yokeTH/gofiber-scalar/scalar/v2"
)

func main(){

    environment.New(0)
    logger := logging.New()
    logger.Slogger.Info("Starting ShopGo API...")

    jwtManager := jwt.New(logger.Slogger)

    //Initialize Fiber app
    app := fiber.New(fiber.Config{
        AppName: "ShopGo",
        ServerHeader: "ShopGO",
        BodyLimit:    environment.GetRequestMaxBodySizeLimit(environment.RequestMaxBodySizeMB),
    })

	// swagger handler
	if environment.GetString(environment.EnvKey) == "development" {
		docPath := filepath.Join("swagger_doc", "swagger.json")
		doc, err := os.ReadFile(docPath)
		if err != nil {
			logger.Slogger.Warn(
				"swagger doc not found, skip swagger",
				"path", docPath,
				"err", err,
			)
		} else {
			cfg := scalar.Config{
				FileContentString: string(doc),
				Theme:             scalar.ThemePurple,
			}
			app.Get("/docs/*", scalar.New(cfg))
		}
	}

    //Middleware
    mid := middleware.NewFiberMiddleware(
        logger.Slogger,
        jwtManager,
    )
    
    app.Use(mid.CORS()) //allow frontend access

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("ShopGo API is running!")
    })
    
    api.Register(
        app,
        logger.Slogger,
        jwtManager,
        mid,
    )

	// Start server
	servicePort := environment.GetString(environment.ServicePortKey)
	serverAddr := ":" + servicePort
	if err := app.Listen(serverAddr); err != nil {
		logger.Slogger.Error("Failed to start server", "error", err)
		log.Fatal(err)
	}
    
}
