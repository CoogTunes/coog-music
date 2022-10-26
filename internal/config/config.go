package config

import (
	"log"
	"os"

	"github.com/alexedwards/scs/v2"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	UseCache      bool
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	DBDriver      string
	DBSource      string
	ServerAddress string
	Session       *scs.SessionManager
}

func GetDBEnvVariable(key string) string {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
