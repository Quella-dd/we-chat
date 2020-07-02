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

	// 根据id/name 搜索 users， 是否需要模糊匹配还是准确检查
	// 群聊是私有的，只有内部成员进行邀请，而不能通过群聊名称进行查询并加入
	{method: http.MethodGet, path: "/search/users/:search", handler: SearchUsers},

	{method: http.MethodPost, path: "/room", handler: CreateRoom},
	{method: http.MethodGet, path: "/rooms", handler: ListRooms},
	{method: http.MethodGet, path: "/room/:id", handler: GetRoom},

	// header: {'userID': 'xxx'}, 如果userID 不是room的管理员， 则没有权限更改room信息
	{method: http.MethodPost, path: "/room/:id", handler: UpdateRoom},
	{method: http.MethodDelete, path: "/room/:id", handler: DeleteRoom},

	// user visite other people to join this group
	{method: http.MethodPost, path: "/addUser/:name/room/:id", handler: AddUserToRoom},
	{method: http.MethodPost, path: "/deleteUser/:name/room/:id/", handler: RemoveFromRoom},

	// header: {'userID': 'xxx'}
	{method: http.MethodGet, path: "/friends", handler: GetFriends},
	{method: http.MethodPost, path: "/friends/:id", handler: AddFriend},
	{method: http.MethodDelete, path: "/friends/:id", handler: DeleteFriend},

	// websocket connect
	{method: http.MethodGet, path: "/event", handler: HandlerEvent},

	// dataCenterManger
	{method: http.MethodPost, path: "/sendMessage", handler: HandlerMessage},
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
