package main

import (
	"webchat/api"
	"webchat/database"
	"webchat/models"
)

func main() {
	database.InitDatabase()
	models.InitModels()
	api.InitRouter()
}
