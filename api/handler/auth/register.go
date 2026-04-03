package handlerauth

import (
    "log/slog"
    "github.com/thanyakwaonueng/shopgo/lib/util/customerror"
    "github.com/thanyakwaonueng/shopgo/api/service/auth/command"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"
    "github.com/mehdihadeli/go-mediatr"
)

type RequestRegister struct {
    Name        string `json:"name" validate:"required"`
    Email       string `json:"email" validate:"required"`
    Password    string `json:"password" validate:"required"`
}

// Register
//
//  @Summary        User Register
//  @Description    Create new user(if the email not exist in database) then sign the token back then auto login
//  @Tags           Auth
//  @Accept         json
//  @Produce        json
//  @Param          request body RequestRegister true "Registration Details"
//  @Success        200
//  @Failer         400
//  @Router         /auth/register [post]
func Register(logger *slog.Logger, validate *validator.Validate) fiber.Handler {
    return func(c *fiber.Ctx) error {
        var request RequestRegister

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

        // Call service to register new user
        result, err := mediatr.Send[command.RequestRegister, command.ResultRegister](
            c.Context(),
            command.RequestRegister{
                Name: request.Name,
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
