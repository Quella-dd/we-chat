package models

// ManageEnv...
var ManageEnv = struct {
	UserManager         *UserManager
	GroupManager        *GroupManager
	RelationShipManager *RelationShipManager
	WebsocketManager    *WebsocketManager
	DataCenterManager   *DataCenterManager
}{}

// InitModels ...
func InitModels() {
	ManageEnv.UserManager = NewUserManager()
	ManageEnv.GroupManager = NewGroupManager()
	ManageEnv.RelationShipManager = NewRelationShipManager()
	ManageEnv.WebsocketManager = NewWebSocketManager()
	ManageEnv.DataCenterManager = NewDataCenterManager()
}
