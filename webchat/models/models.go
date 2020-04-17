package models

// ManageEnv ...
var ManageEnv = struct {
	UserManager      *UserManager
	RoomManager      *RoomManager
	DataManager      *DataManager
	WebsocketManager *WebsocketManager
}{}

// InitModels ...
func InitModels() {
	ManageEnv.UserManager = NewUserManager()
	ManageEnv.RoomManager = NewRoomManager()
	ManageEnv.DataManager = NewDataManager()
	ManageEnv.WebsocketManager = NewWebSocketManager()

}
