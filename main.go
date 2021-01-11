package main

import (
	"we-chat/api"
	_ "we-chat/models"
)

func main() {
	api.InitRouter()
}
