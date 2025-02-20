package config

import (
	"database/sql"
	"flag"
	"log"
	"os"
)

type ServerSettings struct {
	AdresRun      string
	AdresBase     string
	AdresAccrual  string
	SecretCookies string
	DB            *sql.DB
}

func InitServerSettings() (*ServerSettings, error) {
	serverSettings := &ServerSettings{}
	ar := flag.String("a", "localhost:8081", "address running server")
	dbStorage := flag.String("d", "", "address connect to DB")
	asa := flag.String("r", "", "address accrual system")
	sec := flag.String("s", "super_secret", "secret for cookies")
	flag.Parse()

	val, ok := os.LookupEnv("RUN_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresRun = *ar
		log.Println("AdresRun from flag:", serverSettings.AdresRun)
	} else {
		serverSettings.AdresRun = val
		log.Println("AdresRun from env:", serverSettings.AdresRun)
	}

	val, ok = os.LookupEnv("DATABASE_URI")
	if !ok || val == "" {
		serverSettings.AdresBase = *dbStorage
		log.Println("DbStorage from flag:", serverSettings.AdresBase)
	} else {
		serverSettings.AdresBase = val
		log.Println("DbStorage from env:", serverSettings.AdresBase)
	}

	val, ok = os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS")
	if !ok || val == "" {
		serverSettings.AdresAccrual = *asa
		log.Println("AdresAccrual from flag:", serverSettings.AdresAccrual)
	} else {
		serverSettings.AdresAccrual = val
		log.Println("AdresAccrual from env:", serverSettings.AdresAccrual)
	}

	val, ok = os.LookupEnv("SECRET_FOR_COOKIES")
	if !ok || val == "" {
		serverSettings.SecretCookies = *sec
		log.Println("SecretCookies from flag:", serverSettings.SecretCookies)
	} else {
		serverSettings.SecretCookies = val
		log.Println("SecretCookies from env:", serverSettings.SecretCookies)
	}
	return serverSettings, nil
}

func (s *ServerSettings) GetServerSettings() *ServerSettings {
	return s
}
