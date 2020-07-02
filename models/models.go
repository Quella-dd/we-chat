package models

// ManageEnv...
var ManageEnv = struct {
	UserManager *UserManager
	RoomManager *RoomManager
	RelationShipManager *RelationShipManager
	WebsocketManager    *WebsocketManager
	DataCenterManager   *DataCenterManager
}{}

// InitModels ...
func InitModels() {
	ManageEnv.UserManager = NewUserManager()
	ManageEnv.RoomManager = NewRoomManager()
	ManageEnv.RelationShipManager = NewRelationShipManager()
	ManageEnv.WebsocketManager = NewWebSocketManager()
	ManageEnv.DataCenterManager = NewDataCenterManager()
}
