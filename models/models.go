package models

var ManageEnv = struct {
	UserManager         *UserManager
	GroupManager        *GroupManager
	//RelationShipManager *RelationShipManager
	WebsocketManager    *WebsocketManager
	DataCenterManager   *DataCenterManager
	SessionManager *SessionManager
	RequestManager * RequestManager
}{}

func InitModels() {
	ManageEnv.UserManager = NewUserManager()
	ManageEnv.GroupManager = NewGroupManager()
	//ManageEnv.RelationShipManager = NewRelationShipManager()
	ManageEnv.WebsocketManager = NewWebSocketManager()
	ManageEnv.DataCenterManager = NewDataCenterManager()
	ManageEnv.SessionManager = NewSessionManager()
	ManageEnv.RequestManager = NewRequestManager()
}
