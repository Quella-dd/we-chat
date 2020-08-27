package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	confPath = "/webchat/config.json"
	DB       *gorm.DB
	Redis    *redis.Client
)

func parseConfig() *MysqlConfig {
	var conf Conf

	f, err := os.Open(confPath)
	defer f.Close()

	if err != nil {
		panic("database config not found")
	}
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		panic(fmt.Sprintf("database config parse error: %+v\n", err))
	}
	return &conf.Mysql
}

func getMysqlConnecURL(conf *MysqlConfig) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.UserName, conf.PassWord, conf.Host, conf.DataName)
}

func InitDatabase() {
	var err error

	conf := parseConfig()

	if DB, err = gorm.Open("mysql", getMysqlConnecURL(conf)); err != nil {
		panic(err)
	}

	if Redis, err = NewRedisClient(); err != nil {
		panic(err)
	}
}
