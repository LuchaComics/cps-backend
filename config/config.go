package config

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	AppServer  serverConf
	DB         dbConfig
	Cache      cacheConfig
	AWS        awsConfig
	PDFBuilder pdfBuilderConfig
	Emailer    mailgunConfig
}

type serverConf struct {
	Port                         string
	IP                           string
	HMACSecret                   []byte
	HasDebugging                 bool
	InitialAdminEmail            string
	InitialAdminPassword         string
	InitialAdminOrganizationName string
	DomainName                   string
}

type dbConfig struct {
	URI  string
	Name string
}

type cacheConfig struct {
	URI string
}

type awsConfig struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	BucketName string
}

type pdfBuilderConfig struct {
	CBFFTemplatePath  string
	PCTemplatePath    string
	CCIMGTemplatePath string
	CCSCTemplatePath  string
	CCTemplatePath    string
	DataDirectoryPath string
}

type mailgunConfig struct {
	APIKey      string
	Domain      string
	APIBase     string
	SenderEmail string
}

func New() *Conf {
	var c Conf
	c.AppServer.Port = getEnv("CPS_BACKEND_PORT", true)
	c.AppServer.IP = getEnv("CPS_BACKEND_IP", true)
	c.AppServer.HMACSecret = []byte(getEnv("CPS_BACKEND_HMAC_SECRET", true))
	c.AppServer.HasDebugging = getEnvBool("CPS_BACKEND_HAS_DEBUGGING", true, true)
	c.AppServer.InitialAdminEmail = getEnv("CPS_BACKEND_INITIAL_ADMIN_EMAIL", true)
	c.AppServer.InitialAdminPassword = getEnv("CPS_BACKEND_INITIAL_ADMIN_PASSWORD", true)
	c.AppServer.InitialAdminOrganizationName = getEnv("CPS_BACKEND_INITIAL_ADMIN_ORG_NAME", true)
	c.AppServer.DomainName = getEnv("CPS_BACKEND_DOMAIN_NAME", true)

	c.DB.URI = getEnv("CPS_BACKEND_DB_URI", true)
	c.DB.Name = getEnv("CPS_BACKEND_DB_NAME", true)

	c.Cache.URI = getEnv("CPS_BACKEND_CACHE_URI", true)

	c.AWS.AccessKey = getEnv("CPS_BACKEND_AWS_ACCESS_KEY", true)
	c.AWS.SecretKey = getEnv("CPS_BACKEND_AWS_SECRET_KEY", true)
	c.AWS.Endpoint = getEnv("CPS_BACKEND_AWS_ENDPOINT", true)
	c.AWS.Region = getEnv("CPS_BACKEND_AWS_REGION", true)
	c.AWS.BucketName = getEnv("CPS_BACKEND_AWS_BUCKET_NAME", true)

	c.PDFBuilder.CBFFTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CBFF_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.PCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_PC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCIMGTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CCIMG_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCSCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CCSC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.DataDirectoryPath = getEnv("CPS_BACKEND_PDF_BUILDER_DATA_DIRECTORY_PATH", true)

	c.Emailer.APIKey = getEnv("CPS_BACKEND_MAILGUN_API_KEY", true)
	c.Emailer.Domain = getEnv("CPS_BACKEND_MAILGUN_DOMAIN", true)
	c.Emailer.APIBase = getEnv("CPS_BACKEND_MAILGUN_API_BASE", true)
	c.Emailer.SenderEmail = getEnv("CPS_BACKEND_MAILGUN_SENDER_EMAIL", true)

	return &c
}

func getEnv(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func getEnvBool(key string, required bool, defaultValue bool) bool {
	valueStr := getEnv(key, required)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Invalid boolean value for environment variable %s", key)
	}
	return value
}
