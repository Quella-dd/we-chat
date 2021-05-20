package models

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"

	Message "we-chat/message"

	"gopkg.in/ini.v1"
)

type Manager struct {
	*gorm.DB
	redis *redis.Client
	*UserManager
	*GroupManager
	*WebsocketManager
	*DataCenterManager
	*SessionManager
	*RequestManager
	*MomentManager
	*CommonManager
}

type ManagerIni struct {
	Mysql
	Redis
	Port      string
	SecretKey string
}

type Mysql struct {
	UserName string
	Password string
	DBName   string
	Host     string
}

type Redis struct {
	Address  string
	Password string
}

var ManagerEnv *Manager
var ManagerConfig ManagerIni

func InitManage() {
	ManagerEnv = &Manager{
		UserManager:      NewUserManager(),
		GroupManager:     NewGroupManager(),
		WebsocketManager: NewWebSocketManager(),
		//DataCenterManager: NewDataCenterManager(),
		SessionManager: NewSessionManager(),
		RequestManager: NewRequestManager(),
		MomentManager:  NewMomentManager(),
		CommonManager:  NewCommonManager(),
	}
}

func LoadInit() {
	cfg, err := ini.Load("/home/len/go/src/we-chat/config.ini")
	if err != nil {
		panic(err)
	}
	ManagerConfig.Port = cfg.Section("").Key("Port").String()
	ManagerConfig.SecretKey = cfg.Section("").Key("SecretKey").String()

	ManagerConfig.Mysql.UserName = cfg.Section("mysql").Key("username").String()
	ManagerConfig.Mysql.Password = cfg.Section("mysql").Key("password").String()
	ManagerConfig.Mysql.DBName = cfg.Section("mysql").Key("dbname").String()
	ManagerConfig.Mysql.Host = cfg.Section("mysql").Key("host").String()

	ManagerConfig.Redis.Address = cfg.Section("redis").Key("address").String()
	ManagerConfig.Redis.Password = cfg.Section("redis").Key("password").String()
}

func initMysql() {
	sql := ManagerConfig.Mysql
	url := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", sql.UserName, sql.Password, sql.Host, sql.DBName)
	if db, err := gorm.Open("mysql", url); err != nil {
		panic(err)
	} else {
		ManagerEnv.DB = db
	}
}

func initRedis() {
	options := &redis.Options{
		Addr:     ManagerConfig.Redis.Address,
		Password: ManagerConfig.Redis.Password,
		DB:       0,
	}

	client := redis.NewClient(options)
	if err := client.Ping().Err(); err != nil {
		panic(err)
	} else {
		fmt.Printf("client: %T %+v\n", client, client)
		ManagerEnv.redis = client
	}
	ManagerEnv.DataCenterManager = NewDataCenterManager()
}

func initDataTable() {
	ManagerEnv.DB.AutoMigrate(&User{})
	ManagerEnv.DB.AutoMigrate(&Group{})
	ManagerEnv.DB.AutoMigrate(&Message.RequestMessage{})
	ManagerEnv.DB.AutoMigrate(&Session{})
	ManagerEnv.DB.AutoMigrate(&Request{})
	ManagerEnv.DB.AutoMigrate(&Moment{})
	ManagerEnv.DB.AutoMigrate(&Comment{})
}

func init() {
	InitManage()
	LoadInit()
	initMysql()
	initRedis()
	initDataTable()
}
