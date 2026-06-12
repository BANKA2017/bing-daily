package dbio

import (
	"flag"
	"os"
)

// sql
var DBUser string
var DBPassword string
var DBHost string
var DBDatabase string
var DBCert string

var Addr string

// b2 bucket
var B2ApplicationKeyId string
var B2ApplicationKey string
var WorkersLocale string

// test
var TestMode bool

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func InitEnv() {
	// db
	flag.StringVar(&DBUser, "dbuser", GetEnv("dbuser", ""), "username for database")
	flag.StringVar(&DBPassword, "dbpassword", GetEnv("dbpassword", ""), "password for database")
	flag.StringVar(&DBHost, "dbhost", GetEnv("dbhost", ""), "host:port for database")
	flag.StringVar(&DBDatabase, "dbdb", GetEnv("dbdb", ""), "name for database")
	flag.StringVar(&DBCert, "dbcert", GetEnv("dbcert", ""), "cert settings for database")

	//api
	flag.StringVar(&Addr, "addr", GetEnv("addr", ""), "address for api server")

	// b2
	flag.StringVar(&B2ApplicationKeyId, "b2_app_key_id", GetEnv("b2_app_key_id", ""), "B2_APPLICATION_KEY_ID")
	flag.StringVar(&B2ApplicationKey, "b2_app_key", GetEnv("b2_app_key", ""), "B2_APPLICATION_KEY")
	flag.StringVar(&WorkersLocale, "b2_upload_mkt", GetEnv("b2_upload_mkt", ""), "One of \"en-us,zh-cn,ja-jp,es-es,en-ca,en-au,de-de,fr-fr,it-it,en-gb,en-in,pt-br\", or keep empty.")
	// en-us,zh-cn,ja-jp,es-es,en-ca,en-au(row)=en-nz(row),de-de,fr-fr,it-it,en-gb

	// test
	flag.BoolVar(&TestMode, "testmode", false, "test mode")

	flag.Parse()
}
