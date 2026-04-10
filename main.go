package main

import (
    "log"
    "log/slog"
    "os"
    "path/filepath"
    
    repocustom "github.com/thanyakwaonueng/shopgo/api/repository/custom"
    repogeneric "github.com/thanyakwaonueng/shopgo/api/repository/generic"
    serviceauth "github.com/thanyakwaonueng/shopgo/api/service/auth"
    servicecategories "github.com/thanyakwaonueng/shopgo/api/service/categories"
    serviceproducts "github.com/thanyakwaonueng/shopgo/api/service/products"
    serviceorders "github.com/thanyakwaonueng/shopgo/api/service/orders"
    "github.com/thanyakwaonueng/shopgo/api"
    "github.com/thanyakwaonueng/shopgo/lib/environment"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/database"
    "github.com/thanyakwaonueng/shopgo/lib/middleware"
    "github.com/thanyakwaonueng/shopgo/lib/logging"
    
    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/yokeTH/gofiber-scalar/scalar/v2"
)

func main(){

    environment.New(0)
    logger := logging.New()
    logger.Slogger.Info("Starting ShopGo API...")
    
    // Initialize domain database connection
    domainDsn := environment.GetString(environment.DsnDomainKey)
    domainDb := database.New(domainDsn)

    jwtManager := jwt.New(logger.Slogger)

    // Initialize validator
    validate := validator.New(validator.WithRequiredStructEnabled())

    // Initialize repositories
    repo := initRepo(logger.Slogger) 

    // Register services
    {
        serviceauth.Register(domainDb, logger.Slogger, jwtManager, repo.user, repo.userRoleById)
        servicecategories.Register(domainDb, logger.Slogger, repo.category)
        serviceproducts.Register(domainDb, logger.Slogger)
        serviceorders.Register(domainDb, logger.Slogger)
    }

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

    api.Register(
        app,
        logger.Slogger,
        validate,
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

type repo struct {
    // Generic repo
    user                    repogeneric.User
    category                repogeneric.Category

    // Custom repo
    userRoleById            repocustom.UserRoleById
}

func initRepo(logger *slog.Logger) repo {
    // Init generic repo
    var r repo
    {
        r.user = repogeneric.NewUser(logger)
        r.category = repogeneric.NewCategory(logger)
    }
    // Init custom repo
    {
        r.userRoleById = repocustom.NewUserRoleById(logger)
    }
    return r
}
