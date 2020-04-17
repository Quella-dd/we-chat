package database

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	configPath = "config.json"
	DB         *gorm.DB
)

type Config struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	DataName string `json:"dbname"`
}

func parseConfig() *Config {
	var conf Config
	fd, err := os.Open(configPath)
	defer fd.Close()
	if err != nil {
		panic("database config not found")
	}

	if err := json.NewDecoder(fd).Decode(&conf); err != nil {
		panic(fmt.Sprintf("database config parse error: %+v\n", err))
	}
	return &conf
}

func InitDatabase() {
	var err error
	conf := parseConfig()
	url := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8&parseTime=True&loc=Local", conf.UserName, conf.PassWord, conf.Host, conf.DataName)

	DB, err = gorm.Open("mysql", url)

	if err != nil {
		panic(err)
	}
}
