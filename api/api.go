package api

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"we-chat/models"

	"github.com/gin-gonic/gin"
)

type router struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

var routers = []*router{
	{method: http.MethodPost, path: "/login", handler: Login},
	{method: http.MethodPost, path: "/register", handler: Register},

	{method: http.MethodGet, path: "/sessions", handler: ListSessions},
	{method: http.MethodDelete, path: "/session/:id", handler: DeleteSession},

	{method: http.MethodGet, path: "/user/:id", handler: GetUser},
	{method: http.MethodGet, path: "/user/:name", handler: SearchUsers},

	{method: http.MethodGet, path: "/friends", handler: GetFriends},
	{method: http.MethodPost, path: "/friend/:id", handler: AddFriend},
	{method: http.MethodDelete, path: "/friend/:id", handler: DeleteFriend},

	{method: http.MethodGet, path: "/requests", handler: ListRequests},
	{method: http.MethodPost, path: "/request/:id", handler: AckRequest},
	{method: http.MethodDelete, path: "/request/:id", handler: DeleteRequest},

	{method: http.MethodGet, path: "/groups", handler: ListGroups},
	{method: http.MethodPost, path: "/group", handler: CreateGroup},
	{method: http.MethodPost, path: "/group/:id", handler: UpdateGroup},
	{method: http.MethodGet, path: "/group/:id", handler: GetGroup},
	{method: http.MethodDelete, path: "/group/:id", handler: DeleteGroup},

	{method: http.MethodPost, path: "/join/:name/group/:id", handler: JoinGroup},
	{method: http.MethodPost, path: "/leave/:name/group/:id/", handler: LeaveGroup},


	// TODO: Refactor
	// create ws connection and receive message or event
	{method: http.MethodGet, path: "/event", handler: HandlerEvent},
	// dataCenterManger, 当用户A发送消息用户B, 则分别对用户A和用户B创建Session, 如果session存在，则更新session
	{method: http.MethodPost, path: "/sendMessage", handler: HandlerMessage},
	{method: http.MethodGet, path: "/messages/:id", handler: GetMessage},

	// TODO:
	{method: http.MethodPost, path: "/RTCRequest", handler: sendRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/", handler: handlerRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/hangDown", handler: hangDownRTCRequest},
}

func InitRouter() {
	engineer := gin.Default()
	engineer.Use()
	for _, router := range routers {
		engineer.Handle(router.method, router.path, validateHandler(router.handler))
	}
	err := engineer.Run(":9999") // :8080
	if err != nil {
		panic(err)
	}
}

func validateHandler(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "/login") && !strings.Contains(c.Request.URL.Path, "/register") {
			tokenAuth := c.GetHeader("token")
			if tokenAuth == "" {
				c.JSON(http.StatusUnauthorized, gin.H {
					"error": errors.New("http.StatusUnauthorized"),
				})
				return
			} else {
				t, err := jwt.ParseWithClaims(tokenAuth, &models.LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(models.SecretKey), nil
				})

				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H {
						"error": errors.New("http.StatusUnauthorized"),
					})
					return
				}
				if claims, ok := t.Claims.(*models.LoginClaims); ok && t.Valid {
					c.Set("userName", claims.UserName)
					c.Set("userID", claims.UserID)
				}
			}
		}
		f(c)
	}
}

// TODO: 封装error包, 用来wrap error