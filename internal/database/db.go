package database

import (
	"GoStarter/pkg/config"
	"GoStarter/pkg/utils/helpers"
	"GoStarter/pkg/utils/loggers"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Connect() (*gorm.DB, error) {
	logger, err := loggers.NewLogger()
	if err != nil {
		fmt.Printf(helpers.Strings.SetTextRed("Failed to set up loggers: %v"), err)
		return nil, err
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf(helpers.Strings.SetTextRed("Failed to load configuration: %v"), err)
		logger.LogError("Failed to load configuration: %v", err)
		return nil, err
	}

	var db *gorm.DB
	if cfg.GetDBDriver() == "mysql" {
		db, err = gorm.Open(mysql.Open(cfg.GetDBUser()+":"+cfg.GetDBPassword()+"@tcp("+cfg.GetDBHost()+":"+cfg.GetDBPort()+")/"+cfg.GetDBName()+"?parseTime=true"), &gorm.Config{})
		if err != nil {
			fmt.Printf(helpers.Strings.SetTextRed("Failed to connect to database mysql: %v"), err)
			logger.LogError("Failed to connect to database mysql:", err)
			return nil, err
		}
	} else if cfg.GetDBDriver() == "postgres" {
		db, err = gorm.Open(postgres.Open("host="+cfg.GetDBHost()+" user="+cfg.GetDBUser()+" password="+cfg.GetDBPassword()+" dbname="+cfg.GetDBName()+" port="+cfg.GetDBPort()+" sslmode=disable TimeZone=Asia/Jakarta"), &gorm.Config{
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true,
			},
		})
		if err != nil {
			fmt.Printf(helpers.Strings.SetTextRed("Failed to connect to database postgres: %v"), err)
			logger.LogError("Failed to connect to database postgres:", err)
			return nil, err
		}
	}

	return db, nil
}
