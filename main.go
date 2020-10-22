package main

import (
	"we-chat/api"
	"we-chat/database"
	"we-chat/models"
)

func main() {
	database.InitDatabase()
	models.InitModels()
	api.InitRouter()
}
