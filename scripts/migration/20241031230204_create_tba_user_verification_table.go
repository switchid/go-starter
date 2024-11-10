package migration

import (
	internalMigration "GoStarter/internal/cli/commands/migration"
	models "GoStarter/internal/models/authentication"
	"fmt"
	"gorm.io/gorm"
)

func init() {

	var mTbaUserVerification models.TbaUserVerification

	internalMigration.MigrationsDB = append(internalMigration.MigrationsDB, internalMigration.Migrations{
		MigrationName: "20241031230204_create_tba_user_verification_table",
		Up: func(db *gorm.DB) error {
			errVerifyType := db.Exec("CREATE TYPE verify AS ENUM ('administrator','email','sms');").Error
			if errVerifyType != nil {
				return errVerifyType
			}
			if !db.Migrator().HasTable(&mTbaUserVerification) {
				if err := db.Migrator().CreateTable(&mTbaUserVerification); err != nil {
					_ = fmt.Errorf("failed to create users table: %w", err)
					return err
				}
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			errDropVerifyType := db.Exec("DROP TYPE verify CASCADE;").Error
			if errDropVerifyType != nil {
				return errDropVerifyType
			}
			if db.Migrator().HasTable(&mTbaUserVerification) {
				if err := db.Migrator().DropTable(&mTbaUserVerification); err != nil {
					_ = fmt.Errorf("failed to drop users table: %w", err)
					return err
				}
			}
			return nil
		},
	})
}
