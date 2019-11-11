package conf

import (
	"flag"
	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     *Config
)

func init() {
	flag.StringVar(&confPath, "conf", "config.toml", "default config path.")
}

// Init config.
func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

type Config struct {
	DB              *DB
	Log             *Log
	Memcached       *CacheHost
	Job             *Job
	WebCacheDir     string
	ClassFilePrefix string
	Soc             *Soc
}

type DB struct {
	Host     string
	UserName string
	Password string
	Schema   string
	Port     int
}

type CacheHost struct {
	Hosts []string
}
type Log struct {
	Level string
}
type Job struct {
	EndClassOccurrClassOccurrencesCron string
	DeductUserPlanAfterClass           string
}
type Soc struct {
	Url       string
	SecretID  string
	SecretKey string
}
