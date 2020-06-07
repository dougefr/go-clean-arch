package main

import (
	"context"
	"fmt"
	"github.com/dougefr/go-clean-arch/infrastructure"
	"github.com/dougefr/go-clean-arch/interface/repository"
	"github.com/dougefr/go-clean-arch/interface/restctrl"
	"github.com/dougefr/go-clean-arch/usecase"
	"github.com/gofiber/fiber"
	"os"
)

func main() {
	db, err := infrastructure.NewSQLite3()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	logger, err := infrastructure.NewLogrus("debug")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	userRepo := repository.NewUserRepo(db, logger)
	ucCreateUser := usecase.NewCreateUser(userRepo)
	ucSearchUser := usecase.NewSearchUser(userRepo)
	userController := restctrl.NewUser(ucCreateUser, ucSearchUser, db, logger)

	app := fiber.New()
	app.Post("/user", do(userController.CreateUser))
	app.Get("/user", do(userController.SearchUser))

	logger.Info(context.Background(), "listening to port 8080...")
	if err = app.Listen(8080); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func do(fn func(restctrl.RestRequest) restctrl.RestResponse) func(ctx *fiber.Ctx) {
	return func(ctx *fiber.Ctx) {
		var req restctrl.RestRequest
		req.Body = ctx.Fasthttp.PostBody()
		req.GetQueryParam = func(key string) string {
			return ctx.Query(key)
		}
		resp := fn(req)
		ctx.Status(resp.StatusCode).SendBytes(resp.Body)
	}
}
