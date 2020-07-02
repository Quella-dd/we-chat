package database

type Conf struct {
	Mysql MysqlConfig `json:"mysql"`
	Redis ReidsConfig `json:"redis"`
}

type MysqlConfig struct {
	Host     string `json:"host"`
	UserName string `json:"username"`
	PassWord string `json:"password"`
	DataName string `json:"dbname"`
}

type ReidsConfig struct {
	Address  string `json: "addredss"`
	Password string `json: "password"`
}
