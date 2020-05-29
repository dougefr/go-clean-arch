package api

import (
	"github.com/dougefr/go-clean-code/controller"
	"github.com/gofiber/fiber"
)

// APIServer ...
type APIServer interface {
	Start(port int) error
}

type apiServer struct {
	userController controller.User
}

// NewAPIServer ...
func NewAPIServer(userController controller.User) APIServer {
	return apiServer{
		userController: userController,
	}
}

// Start ...
func (a apiServer) Start(port int) error {
	app := fiber.New()

	app.Post("/user", do(a.userController.CreateUser))

	if err := app.Listen(port); err != nil {
		return err
	}

	return nil
}

func do(fn func(controller.RestRequest) controller.RestResponse) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var req controller.RestRequest
		req.Body = ctx.Fasthttp.PostBody()
		resp := fn(req)
		ctx.Status(resp.StatusCode).SendBytes(resp.Body)
	}
}
