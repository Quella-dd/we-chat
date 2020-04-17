package api

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	Method  string
	Path    string
	Handler gin.HandlerFunc
}

var Routers = []*Router{
	{Method: "POST", Path: "/login", Handler: Login},
	{Method: "POST", Path: "/registry", Handler: CreateUser},
	{Method: "GET", Path: "/users", Handler: ListUsers},
	{Method: "GET", Path: "/users/:id", Handler: GetUser},

	{Method: "POST", Path: "/user/:id/room/:name", Handler: JoinRoom},
	{Method: "DELETE", Path: "/user/:id/room/:name", Handler: LeaveRoom},
	{Method: "POST", Path: "/user/:id/friend/:name", Handler: AddFriend},
	{Method: "DELETE", Path: "/user/:id/friend/:name", Handler: RemovFriend},

	{Method: "POST", Path: "/room", Handler: CreateRoom},
	{Method: "GET", Path: "/rooms", Handler: ListRooms},
	{Method: "GET", Path: "/room/:id", Handler: GetRoom},
	{Method: "POST", Path: "/room/:name", Handler: UpdateRoom},
	{Method: "DELETE", Path: "/room/:name", Handler: DeleteRoom},

	{Method: "GET", Path: "/friends", Handler: GetFriends},

	{Method: "POST", Path: "/sendMessage", Handler: HandlerMessage},
	{Method: "GET", Path: "/event", Handler: HandlerEvent},
}

func InitRouter() {
	engineer := gin.Default()
	for _, router := range Routers {
		engineer.Handle(router.Method, router.Path, router.Handler)
	}
	engineer.Run(":9999") // :8080
}
