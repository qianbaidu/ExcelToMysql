package conf

import (
	"gopkg.in/gcfg.v1"
	"fmt"
	"io/ioutil"
)

const ConfigPath string =  "./conf/config.ini"

type Conf struct {
	Mysql struct {
		      Host     string
		      Port     string
		      Username string
		      Password string
		      Database string
	      }
	Other struct {
		      Savedir string
	      }
}

func ReadConfig() (Config Conf) {
	err := gcfg.ReadFileInto(&Config, ConfigPath)
	if err != nil {
		fmt.Println("Failed to parse config file: %s", err)
	}
	return Config
}

func InitConfig() (Config Conf) {
	var confString = `[Mysql]
host = localhost
port = 3306
username = root
password =
database = Excel
[Other]
savedir = Excel`

	var d1 = []byte(confString)
	ioutil.WriteFile(ConfigPath, d1, 0666)
	Config = ReadConfig()
	return Config
}

