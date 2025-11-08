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

func InitEnv() {
	// db
	flag.StringVar(&DBUser, "dbuser", "", "username for database")
	flag.StringVar(&DBPassword, "dbpassword", "", "password for database")
	flag.StringVar(&DBHost, "dbhost", "", "host:port for database")
	flag.StringVar(&DBDatabase, "dbdb", "", "name for database")
	flag.StringVar(&DBCert, "dbcert", "", "cert settings for database")

	//api
	flag.StringVar(&Addr, "addr", "", "address for api server")

	// b2
	flag.StringVar(&B2ApplicationKeyId, "b2_app_key_id", "", "B2_APPLICATION_KEY_ID")
	flag.StringVar(&B2ApplicationKey, "b2_app_key", "", "B2_APPLICATION_KEY")
	flag.StringVar(&WorkersLocale, "b2_upload_mkt", "", "One of \"en-us,zh-cn,ja-jp,es-es,en-ca,en-au,de-de,fr-fr,it-it,en-nz,en-gb\", or keep empty.")
	// en-us,zh-cn,ja-jp,es-es,en-ca,en-au,de-de,fr-fr,it-it,en-nz(row),en-gb

	flag.Parse()

	if DBUser == "" {
		DBUser = os.Getenv("dbuser")
	}
	if DBPassword == "" {
		DBPassword = os.Getenv("dbpassword")
	}
	if DBHost == "" {
		DBHost = os.Getenv("dbhost")
	}
	if DBDatabase == "" {
		DBDatabase = os.Getenv("dbdb")
	}
	if DBCert == "" {
		DBCert = os.Getenv("dbcert")
	}

	if Addr == "" {
		Addr = os.Getenv("addr")
	}

	if B2ApplicationKeyId == "" {
		B2ApplicationKeyId = os.Getenv("b2_app_key_id")
	}
	if B2ApplicationKey == "" {
		B2ApplicationKey = os.Getenv("b2_app_key")
	}
	if WorkersLocale == "" {
		WorkersLocale = os.Getenv("b2_upload_mkt")
	}
}
