package main

import (
	"fmt"
	"github.com/dougefr/go-clean-code/controller/rest"
	"github.com/dougefr/go-clean-code/infrastructure/api"
	"github.com/dougefr/go-clean-code/infrastructure/database"
	"github.com/dougefr/go-clean-code/infrastructure/database/repository"
	"github.com/dougefr/go-clean-code/usecase"
	"os"
)

func main() {
	db, err := database.NewDatabase()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	userRepo := repository.NewUserRepo(db)
	ucCreateUser := usecase.NewCreateUser(userRepo)
	userController := rest.NewUser(ucCreateUser)
	apiServer := api.NewAPIServer(userController)

	if err = apiServer.Start(8080); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
