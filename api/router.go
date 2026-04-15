package api

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/lib/jwt"
    "github.com/thanyakwaonueng/shopgo/lib/middleware"

    handlerauth "github.com/thanyakwaonueng/shopgo/api/handler/auth"
    handlercategories "github.com/thanyakwaonueng/shopgo/api/handler/categories"
    handlerproducts "github.com/thanyakwaonueng/shopgo/api/handler/products"
    handlerorders "github.com/thanyakwaonueng/shopgo/api/handler/orders"
    handlerusers "github.com/thanyakwaonueng/shopgo/api/handler/users"

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
    registerAdminRoutes(api, logger, validate, mid)
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

    groupCategories := api.Group("/categories")
    {
        groupCategories.Get("/", handlercategories.GetCategories(logger))
    }

    groupProducts := api.Group("/products")
    {
        groupProducts.Get("/", handlerproducts.GetProducts(logger, validate))
        groupProducts.Get("/:id", handlerproducts.GetProductByID(logger))
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

    groupOrders := api.Group("/orders")
	{
		groupOrders.Use(mid.Authenticated())

		// Customer can place, list own, view own, and cancel own pending orders
		groupOrders.Post("/", handlerorders.CreateOrder(logger, validate))
		groupOrders.Get("/", handlerorders.GetOrders(logger, validate))
		groupOrders.Get("/:id", handlerorders.GetOrderByID(logger))
		groupOrders.Post("/:id/cancel", handlerorders.CancelOrder(logger))
	}
}

func registerAdminRoutes(
    api fiber.Router,
    logger *slog.Logger,
    validate *validator.Validate,
    mid *middleware.FiberMiddleware,
) {

    groupCategories := api.Group("/categories")
    {
        groupCategories.Use(mid.Authenticated())
        groupCategories.Use(mid.AdminOnly())
        groupCategories.Post("/", handlercategories.CreateCategory(logger, validate))
        groupCategories.Put("/:id", handlercategories.UpdateCategory(logger, validate))
        groupCategories.Delete("/:id", handlercategories.DeleteCategory(logger))
    }

    groupProducts := api.Group("/products")
    {
        groupProducts.Use(mid.Authenticated())
        groupProducts.Use(mid.AdminOnly())

        groupProducts.Post("/", handlerproducts.CreateProduct(logger, validate))
        groupProducts.Put("/:id", handlerproducts.UpdateProduct(logger, validate))
        groupProducts.Delete("/:id", handlerproducts.DeleteProduct(logger))
    }

	groupOrders := api.Group("/orders")
	{
		groupOrders.Use(mid.Authenticated())
		groupOrders.Use(mid.AdminOnly())

		// Only Admin can advance the order status
		groupOrders.Patch("/:id/status", handlerorders.UpdateOrderStatus(logger, validate))
		groupOrders.Post("/:id/cancel", handlerorders.CancelOrder(logger))
	}

    groupUsers := api.Group("/users")
    {
        groupUsers.Use(mid.Authenticated())
        groupUsers.Use(mid.AdminOnly())

        groupUsers.Get("/", handlerusers.GetUsers(logger, validate))
        groupUsers.Get("/:id", handlerusers.GetUserByID(logger, validate))
        groupUsers.Get("/:id/role", handlerusers.UpdateUserRole(logger, validate))
    }
}
