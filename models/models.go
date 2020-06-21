package models

// ManageEnv ...
var ManageEnv = struct {
	UserManager         *UserManager
	RoomManager         *RoomManager
	DataManager         *DataManager
	RelationShipManager *RelationShipManager
	WebsocketManager    *WebsocketManager
}{}

// InitModels ...
func InitModels() {
	ManageEnv.UserManager = NewUserManager()
	ManageEnv.RoomManager = NewRoomManager()
	ManageEnv.DataManager = NewDataManager()
	ManageEnv.RelationShipManager = NewRelationShipManager()
	ManageEnv.WebsocketManager = NewWebSocketManager()
}
