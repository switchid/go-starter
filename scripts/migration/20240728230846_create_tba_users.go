package migration

import (
	internalMigration "GoStarter/internal/cli/commands/migration"
	models "GoStarter/internal/models/authentication"
	"fmt"
	"gorm.io/gorm"
)

func init() {
	var User models.TbaUser
	internalMigration.MigrationsDB = append(internalMigration.MigrationsDB, internalMigration.Migrations{
		MigrationName: "20240728230846_create_tb_users",
		Up: func(db *gorm.DB) error {
			if !db.Migrator().HasTable(&User) {
				if err := db.Migrator().CreateTable(&User); err != nil {
					_ = fmt.Errorf("failed to create users table: %w", err)
					return err
				}
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			if db.Migrator().HasTable(&User) {
				if err := db.Migrator().DropTable(&User); err != nil {
					_ = fmt.Errorf("failed to drop users table: %w", err)
					return err
				}
			}
			return nil
		},
	})
}
