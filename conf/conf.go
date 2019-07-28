package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "config.toml", "default config path.")
}

// Init init config.
func Init() (err error) {
	//Conf = Default()
	_, err = toml.DecodeFile(confPath, &Conf)

	return
}

type Config struct {
	DB  *DB
	Log *logrus.Logger
}

type DB struct {
	Host     string
	UserName string
	Password string
	Schema   string
	Port     int
}
