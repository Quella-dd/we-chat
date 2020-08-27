package api

import (
	"encoding/json"
	"io"
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
	{method: http.MethodGet, path: "/search/users/:search", handler: SearchUsers},

	// 群聊是私有的，只有内部成员进行邀请，而不能通过群聊名称进行查询并加入
	{method: http.MethodPost, path: "/room", handler: CreateGroup},
	{method: http.MethodGet, path: "/rooms", handler: ListGroup},
	{method: http.MethodGet, path: "/room/:id", handler: GetGroup},

	// header: {'userID': 'xxx'}, 如果userID 不是room的管理员， 则没有权限更改room信息
	{method: http.MethodPost, path: "/room/:id", handler: UpdateGroup},
	{method: http.MethodDelete, path: "/room/:id", handler: DeleteGroup},

	// user visit other people join in this group
	{method: http.MethodPost, path: "/addUser/:name/room/:id", handler: AddUserToRoom},
	{method: http.MethodPost, path: "/deleteUser/:name/room/:id/", handler: RemoveFromRoom},

	// header: {'userID': 'xxx'}
	{method: http.MethodGet, path: "/friends", handler: GetFriends},
	{method: http.MethodPost, path: "/friends/:id", handler: AddFriend},
	{method: http.MethodDelete, path: "/friends/:id", handler: DeleteFriend},
	{method: http.MethodGet, path: "/requests", handler: ListFriendRequests},
	{method: http.MethodPost, path: "/requests/:id", handler: AckFriendRequest},

	// create ws connection and receive message or event
	{method: http.MethodGet, path: "/event", handler: HandlerEvent},

	// dataCenterManger
	{method: http.MethodPost, path: "/sendMessage", handler: HandlerMessage},
	{method: http.MethodGet, path: "/messages/:id", handler: GetMessage},

	// webcrt singnal server. SDP, ICE信息交换等
	{method: http.MethodPost, path: "/RTCRequest", handler: sendRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/", handler: handlerRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/hangDown", handler: hangDownRTCRequest},
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

func jsonResult(w io.Writer, v interface{}) {
	buf, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	w.Write(buf)
}
