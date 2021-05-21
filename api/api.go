package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"we-chat/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var GE = &models.ManagerEnv

type router struct {
	method  string
	path    string
	handler gin.HandlerFunc
}

var routers = []*router{
	{method: http.MethodPost, path: "/login", handler: Login},
	{method: http.MethodPost, path: "/register", handler: Register},

	{method: http.MethodGet, path: "/user/:id", handler: GetUser},
	{method: http.MethodGet, path: "/userSearch/:name", handler: SearchUsers},

	{method: http.MethodGet, path: "/friends", handler: GetFriends},
	{method: http.MethodPost, path: "/friend/:id", handler: AddFriend},
	{method: http.MethodDelete, path: "/friend/:id", handler: DeleteFriend},

	{method: http.MethodGet, path: "/requests", handler: ListRequests},
	{method: http.MethodGet, path: "/request/:id", handler: GetRequest},
	{method: http.MethodPost, path: "/request/:id", handler: AckRequest},
	{method: http.MethodDelete, path: "/request/:id", handler: DeleteRequest},

	{method: http.MethodGet, path: "/groups", handler: ListGroups},
	{method: http.MethodPost, path: "/group", handler: CreateGroup},
	{method: http.MethodPost, path: "/group/:id", handler: UpdateGroup},
	{method: http.MethodGet, path: "/group/:id", handler: GetGroup},
	{method: http.MethodDelete, path: "/group/:id", handler: DeleteGroup},

	{method: http.MethodPost, path: "/join/group/:id", handler: JoinGroup},
	{method: http.MethodPost, path: "/leave/group/:id", handler: LeaveGroup},

	// TODO: 可以采取多重索引来加速message的查询
	{method: http.MethodGet, path: "/sessions", handler: ListSessions},
	{method: http.MethodPost, path: "/session", handler: CreateSession},
	{method: http.MethodGet, path: "/session/:id", handler: GetSession},
	{method: http.MethodDelete, path: "/session/:id", handler: DeleteSession},

	{method: http.MethodGet, path: "/event", handler: HandlerEvent},

	{method: http.MethodGet, path: "/messages/:id", handler: GetMessages},
	{method: http.MethodPost, path: "/sendMessage", handler: HandlerMessage},

	{method: http.MethodGet, path: "/moments", handler: ListMoments},
	{method: http.MethodGet, path: "/moment/:id", handler: GetMoment},
	{method: http.MethodPost, path: "/moment", handler: CreateMoment},
	{method: http.MethodDelete, path: "/moment/:id", handler: DeleteMoment},
	{method: http.MethodPut, path: "/moment/:id", handler: UpdateMoment},

	{method: http.MethodPost, path: "/comment", handler: CreateComment},
	{method: http.MethodDelete, path: "/comment/:id", handler: DeleteComment},

	// TODO
	{method: http.MethodPost, path: "/RTCRequest", handler: sendRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/", handler: handlerRTCRequest},
	{method: http.MethodPost, path: "/RTCRequest/:id/hangDown", handler: hangDownRTCRequest},
}

func InitRouter() {
	engine := gin.Default()

	for _, router := range routers {
		engine.Handle(router.method, router.path, validateHandler(router.handler))
	}

	port := fmt.Sprintf(":%s", models.ManagerConfig.Port)
	log.Fatal(engine.Run(port))
}

func validateHandler(f gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !strings.Contains(c.Request.URL.Path, "/login") && !strings.Contains(c.Request.URL.Path, "/register") {
			tokenAuth := c.GetHeader("token")

			if strings.Contains(c.Request.URL.Path, "/event") {
				tokenAuth = c.Query("token")
			}

			if tokenAuth == "" {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": errors.New("http.StatusUnauthorized"),
				})
				return
			} else {
				t, err := jwt.ParseWithClaims(tokenAuth, &models.LoginClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(models.ManagerConfig.SecretKey), nil
				})

				if err != nil {
					c.JSON(http.StatusUnauthorized, gin.H{
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

func jsonResult(c *gin.Context, code int, o interface{}) {
	w := c.Writer

	w.WriteHeader(code)

	if v, ok := o.(error); ok {
		w.Write([]byte(v.Error()))
	} else {
		buf, _ := json.Marshal(o)
		w.Write(buf)
	}
	return
}
