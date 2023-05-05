package config

import (
	"log"

	"github.com/joeshaw/envdecode"
)

type Conf struct {
	AppServer serverConf
	DB        dbConfig
	Cache     cacheConfig
	AWS       awsConfig
}

type serverConf struct {
	Port                 string `env:"CPS_BACKEND_PORT,required"`
	IP                   string `env:"CPS_BACKEND_IP,required"`
	HMACSecret           []byte `env:"CPS_BACKEND_HMAC_SECRET,required"`
	HasDebugging         bool   `env:"CPS_BACKEND_HAS_DEBUGGING,default=true"`
	InitialAdminEmail    string `env:"CPS_BACKEND_INITIAL_ADMIN_EMAIL,required"`
	InitialAdminPassword string `env:"CPS_BACKEND_INITIAL_ADMIN_PASSWORD,required"`
}

type dbConfig struct {
	URI  string `env:"CPS_BACKEND_DB_URI,required"`
	Name string `env:"CPS_BACKEND_DB_NAME,required"`
}

type cacheConfig struct {
	Host     string `env:"CPS_BACKEND_CACHE_HOST,required"`
	Port     string `env:"CPS_BACKEND_CACHE_PORT,required"`
	Password string `env:"CPS_BACKEND_CACHE_PASSWORD,required"`
}

type awsConfig struct {
	AccessKey  string `env:"CPS_BACKEND_AWS_ACCESS_KEY,required"`
	SecretKey  string `env:"CPS_BACKEND_AWS_SECRET_KEY,required"`
	Endpoint   string `env:"CPS_BACKEND_AWS_ENDPOINT,required"`
	Region     string `env:"CPS_BACKEND_AWS_REGION,required"`
	BucketName string `env:"CPS_BACKEND_AWS_BUCKET_NAME,required"`
}

func New() *Conf {
	var c Conf
	if err := envdecode.StrictDecode(&c); err != nil {
		log.Fatalf("Failed to decode: %s", err)
	}
	return &c
}
