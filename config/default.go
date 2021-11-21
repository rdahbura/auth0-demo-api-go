package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

var (
	MgmtApiClientId     string
	MgmtApiClientSecret string
	MgmtApiAudience     string
	MgmtApiIssuer       string
	MgmtApiBaseUrl      string
	MgmtApiTokenUrl     string
)

var (
	MongoDb  string
	MongoPwd string
	MongoUri string
	MongoUsr string
)

var (
	TokenAudience string
	TokenIssuer   string
)

var (
	Host string
	Mode string
	Port string
)

func Load() {
	godotenv.Load(".env.local")
	godotenv.Load()

	MgmtApiClientId = os.Getenv("MGMT_API_CLIENT_ID")
	MgmtApiClientSecret = os.Getenv("MGMT_API_CLIENT_SECRET")
	MgmtApiAudience = os.Getenv("MGMT_API_AUDIENCE")
	MgmtApiIssuer = os.Getenv("MGMT_API_ISSUER")
	MgmtApiBaseUrl = fmt.Sprintf("%s/api/v2", MgmtApiIssuer)
	MgmtApiTokenUrl = fmt.Sprintf("%s/oauth/token", MgmtApiIssuer)

	MongoDb = os.Getenv("MONGO_DB")
	MongoPwd = os.Getenv("MONGO_PWD")
	MongoUri = os.Getenv("MONGO_URI")
	MongoUsr = os.Getenv("MONGO_USR")

	TokenAudience = os.Getenv("TOKEN_AUDIENCE")
	TokenIssuer = os.Getenv("TOKEN_ISSUER")

	Host = os.Getenv("HOST")
	Mode = os.Getenv("MODE")
	Port = os.Getenv("PORT")
}
