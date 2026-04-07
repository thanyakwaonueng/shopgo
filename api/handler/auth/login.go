package handlerauth

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"
    "github.com/thanyakwaonueng/shopgo/api/service/auth/command"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
)

type RequestLogin struct {
    Email       string `json:"email" validate:"required"`
    Password    string `json:"password" validate:"required"`
}


// Login
//
//  @Summary        User Login
//  @Description    Authenticate user with username and password. Use remember_me=true to extend refresh token expiration to 30 days     (instead of 1 day).
//  @Tags           Auth
//  @Accept         json
//  @Produce        json
//  @Param          body    body        RequestLogin    true    "Login credentials (remember_me is optional, defaults to false)"
//  @Success        200     {object}    command.ResultLogin
//  @Failure        400     {object}    customerror.Model
//  @Router         /auth/login [post]
func Login(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var request RequestLogin

        // Parse request body
        if err := c.BodyParser(&request); err != nil {
            customErr := customerror.NewInternalErr("Invalid body params")
            logger.Error(customErr.Message, customerror.LogErrorKey, err)
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
        }

        // Validate request
        if err := validate.Struct(request); err != nil {
            customErr := customerror.NewInternalErr("Invalid request")
            logger.Error(customErr.Message, customerror.LogErrorKey, err)
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
        }

        // Call service to login with remember me flag(I think this one is nice to have)
        result, err := mediatr.Send[command.RequestLogin, command.ResultLogin](
            c.Context(),
            command.RequestLogin{
                Email: request.Email,
                Password: request.Password,
            },
        )


        if err != nil {
            customErr := customerror.UnmarshalError(err)
            logger.Error(customErr.Message)
            return c.Status(fiber.StatusBadRequest).JSON(customErr)
        }

        return c.JSON(result)
    }
}
