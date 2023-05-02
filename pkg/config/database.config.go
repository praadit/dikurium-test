package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/gorm/utils"
)

func InitDb() *gorm.DB {
	loggerLevel := logger.Silent
	debugLogEnv := []string{"debug", "local"}
	if utils.Contains(debugLogEnv, Env) {
		loggerLevel = logger.Info
	}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  loggerLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	tls := "disable"
	if Config.DbTLS {
		tls = "require"
	}
	dsn := fmt.
		Sprintf(`
			host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Asia/Jakarta`,
			Config.DbHost,
			Config.DbPort,
			Config.DbUser,
			Config.DbPass,
			Config.DbName,
			tls,
		)
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: newLogger,
	})
	if err != nil {
		log.Fatal(err)
	}

	return db
}
