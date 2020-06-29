package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type router struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

var routers = []*router{
	{method: http.MethodPost, path: "/login", handler: Login},
	{method: http.MethodPost, path: "/registry", handler: CreateUser},
	{method: http.MethodGet, path: "/users", handler: ListUsers},
	{method: http.MethodGet, path: "/users/:id", handler: GetUser},

	{method: http.MethodPost, path: "/user/:id/room/:name", handler: JoinRoom},
	{method: http.MethodDelete, path: "/user/:id/room/:name", handler: LeaveRoom},

	{method: http.MethodPost, path: "/room/:id/addUser/:name", handler: AddUserToRoom},
	{method: http.MethodPost, path: "/room/:id/deleteUser/:name", handler: RemoveFromRoom},

	//{method: http.MethodPost, path: "/room", handler: CreateRoom},
	//{method: http.MethodGet, path: "/rooms", handler: ListRooms},
	//{method: http.MethodGet, path: "/room/:id", handler: GetRoom},
	//{method: http.MethodPost, path: "/room/:name", handler: UpdateRoom},
	//{method: http.MethodDelete, path: "/room/:name", handler: DeleteRoom},

	{method: http.MethodGet, path: "/search/users/:name", handler: SearchUsers},

	{method: http.MethodGet, path: "/friends/:id", handler: GetFriends},
	{method: http.MethodPost, path: "/friends/:id", handler: AddFriend},
	{method: http.MethodDelete, path: "/friends/:id", handler: DeleteFriend},

	{method: http.MethodPost, path: "/sendMessage", handler: HandlerMessage},
	{method: http.MethodGet, path: "/event", handler: HandlerEvent},

	{method: http.MethodGet, path: "/messages/:id", handler: GetMessage},
}

func InitRouter() {
	engineer := gin.Default()
	for _, router := range routers {
		engineer.Handle(router.method, router.path, router.handler)
	}
	err := engineer.Run(":9999") // :8080
	if err != nil {
		panic(err)
	}
}

/*
	dataCenter-> (数据转发, 数据存储, 离线数据)
*/
