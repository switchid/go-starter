package database

import (
	"GoStarter/internal/pkg/config"
	"GoStarter/pkg/utils/loggers"
	"GoStarter/pkg/utils/stringers"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func Connect() (*gorm.DB, error) {
	logger, err := loggers.NewLogger()
	if err != nil {
		fmt.Printf(stringers.NewString("Failed to set up loggers: %v").SetTextColor(stringers.RED), err)
		return nil, err
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf(stringers.NewString("Failed to load configuration: %v").SetTextColor(stringers.RED), err)
		logger.LogError("Failed to load configuration: %v", err)
		return nil, err
	}

	var db *gorm.DB
	if cfg.GetDBDriver() == "mysql" {
		db, err = gorm.Open(mysql.Open(cfg.GetDBUser()+":"+cfg.GetDBPassword()+"@tcp("+cfg.GetDBHost()+":"+cfg.GetDBPort()+")/"+cfg.GetDBName()+"?parseTime=true"), &gorm.Config{})
		if err != nil {
			fmt.Printf(stringers.NewString("Failed to connect to database mysql: %v").SetTextColor(stringers.RED), err)
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
			fmt.Printf(stringers.NewString("Failed to connect to database postgres: %v").SetTextColor(stringers.RED), err)
			logger.LogError("Failed to connect to database postgres:", err)
			return nil, err
		}
	}

	return db, nil
}
