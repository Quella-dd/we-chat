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
	confPath = "database/config.json"
	DB       *gorm.DB
	Redis    *redis.Client
)

func parse() *Conf {
	f, err := os.Open(confPath)
	defer f.Close()

	if err != nil {
		panic("database config not found")
	}

	var conf Conf
	if err := json.NewDecoder(f).Decode(&conf); err != nil {
		panic(fmt.Sprintf("database config parse error: %+v\n", err))
	}
	return &conf
}

func NewMysqlClient(conf MysqlConfig) string {
	return fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.UserName, conf.PassWord, conf.Host, conf.DataName)
}

func NewRedisClient(option *redis.Options) (*redis.Client, error) {
	client := redis.NewClient(option)

	if err := client.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

func InitDatabase() {
	var err error
	conf := parse()
	if DB, err = gorm.Open("mysql", NewMysqlClient(conf.Mysql)); err != nil {
		panic(err)
	}
	if Redis, err = NewRedisClient(&redis.Options{
		Addr:     conf.Redis.Address,
		Password: conf.Redis.Password,
		DB:       0,
	}); err != nil {
		panic(err)
	}
}
