package migration

import (
	"GoStarter/internal/pkg/config"
	"GoStarter/pkg/utils/stringers"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os/exec"
	"sort"
	"strings"
	"time"
	"unicode"
)

type Migration struct {
	ID        uint      `gorm:"primaryKey;autoIncrement:true"`
	Name      string    `gorm:"size:255; unique"`
	Batch     int       `gorm:"type:INTEGER;default:0"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

type Migrations struct {
	MigrationName string
	Up            func(db *gorm.DB) error
	Down          func(db *gorm.DB) error
}

var MigrationsDB []Migrations

func (Migration) TableName() string {
	return "migration"
}

func RunMigrate(db *gorm.DB) {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %v", errCfg)
	}

	if !db.Migrator().HasTable(&Migration{}) {
		err := db.AutoMigrate(&Migration{})
		if err != nil {
			log.Fatalf("failed to create migration table: %v", err)
		}
	}

	var batch, countMigrate int
	countMigrate = 0
	var errBatch error
	if cfg.GetDBDriver() == "mysql" {
		errBatch = db.Table("migration").Select("IFNULL(MAX(Batch),0)+1 AS Batch").Scan(&batch).Error
	} else if cfg.GetDBDriver() == "postgres" {
		errBatch = db.Table("migration").Select("COALESCE(NULLIF(MAX(Batch),0),0)+1 AS Batch").Scan(&batch).Error
	} else {
		errBatch = nil
	}
	if errBatch != nil {
		fmt.Println("Error to get batch:", errBatch)
	} else {
		for _, migration := range MigrationsDB {
			var allMigrations []Migration
			db.Find(&allMigrations, "Name = ?", migration.MigrationName)
			if len(allMigrations) < 1 {
				fmt.Printf("Migration %s : \n", migration.MigrationName)
				fmt.Printf(stringers.NewString(" Drop migration...\n").SetTextColor(stringers.RED))
				errDown := migration.Down(db)
				if errDown != nil {
					fmt.Printf(stringers.NewString("Error applying migration %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errDown)
				}

				fmt.Printf(stringers.NewString(" Applying migration...\n").SetTextColor(stringers.BLUE))
				errUp := migration.Up(db)
				if errUp != nil {
					fmt.Printf(stringers.NewString("Error applying migration %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errUp)
				}

				errMigration := db.Create(&Migration{Name: migration.MigrationName, Batch: batch}).Error
				if errMigration != nil {
					fmt.Printf(stringers.NewString("Error to insert migrations %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errMigration)
				}

				countMigrate += 1
				fmt.Printf(stringers.NewString(" Migration applied successfully... \n").SetTextColor(stringers.GREEN))
			}
		}
	}

	if countMigrate > 0 {
		fmt.Printf("Migrations completed successfully. \n")
	} else {
		fmt.Printf(stringers.NewString("Nothing to migrate. \n").SetTextColor(stringers.RED))
	}
}

func RollbackMigrate(db *gorm.DB) {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %v", errCfg)
	}

	var batch, countMigrate int
	countMigrate = 0
	var errBatch error
	if cfg.GetDBDriver() == "mysql" {
		errBatch = db.Table("migrations").Select("IFNULL(MAX(Batch),0) AS Batch").Scan(&batch).Error
	} else if cfg.GetDBDriver() == "postgres" {
		errBatch = db.Table("migration").Select("COALESCE(NULLIF(MAX(Batch),0),0) AS Batch").Scan(&batch).Error
	} else {
		errBatch = nil
	}
	if errBatch != nil {
		fmt.Println("Error to get batch:", errBatch)
	} else {
		var rbBatchMigration []Migration
		db.Find(&rbBatchMigration, "Batch = ?", batch)
		if len(rbBatchMigration) > 0 {
			sort.SliceStable(MigrationsDB, func(i, j int) bool {
				// If i starts with "create_" and j doesn't, return false to put i after j
				if strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") && !strings.HasPrefix(MigrationsDB[j].MigrationName, "create_") {
					return false
				}
				// If j starts with "create_" and i doesn't, return true to keep i before j
				if !strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") && strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") {
					return true
				}
				// If both or neither start with "create_", maintain their original order
				return i < j
			})
			for _, rbm := range rbBatchMigration {
				for _, migration := range MigrationsDB {
					if rbm.Name == migration.MigrationName {
						fmt.Printf("Rollback migration %s : \n", migration.MigrationName)
						fmt.Printf(stringers.NewString(" Drop migration...\n").SetTextColor(stringers.RED))
						errUp := migration.Down(db)
						if errUp != nil {
							fmt.Printf(stringers.NewString("Error applying migration %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errUp)
						}
						countMigrate += 1
					}
				}
			}
		}
	}
	if countMigrate > 0 {
		errMigrationDel := db.Delete(&Migration{}, "Batch = ?", batch).Error
		if errMigrationDel != nil {
			fmt.Println("Error to delete migrations:", errMigrationDel)
		}
		fmt.Printf("Rollback migrations completed successfully. \n")
	} else {
		fmt.Printf(stringers.NewString("Nothing to rollback. \n").SetTextColor(stringers.RED))
	}
}

func FreshMigrate(db *gorm.DB) {
	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %v", errCfg)
	}

	var batch, countMigrate int
	countMigrate = 0
	var errBatch error
	if cfg.GetDBDriver() == "mysql" {
		errBatch = db.Table("migrations").Select("IFNULL(MAX(Batch),0)+1 AS Batch").Scan(&batch).Error
	} else if cfg.GetDBDriver() == "postgres" {
		errBatch = db.Table("migration").Select("COALESCE(NULLIF(MAX(Batch),0),0)+1 AS Batch").Scan(&batch).Error
	} else {
		errBatch = nil
	}
	if errBatch != nil {
		fmt.Println("Error to get batch:", errBatch)
	} else {
		var freshMigration []Migration
		errMigration := db.Table("migrations").Select("Name").Scan(&freshMigration).Error
		if errMigration != nil {
			fmt.Println("Error to get migrations:", errMigration)
		} else {
			if len(freshMigration) > 0 {
				sort.SliceStable(MigrationsDB, func(i, j int) bool {
					// If i starts with "create_" and j doesn't, return false to put i after j
					if strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") && !strings.HasPrefix(MigrationsDB[j].MigrationName, "create_") {
						return false
					}
					// If j starts with "create_" and i doesn't, return true to keep i before j
					if !strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") && strings.HasPrefix(MigrationsDB[i].MigrationName, "create_") {
						return true
					}
					// If both or neither start with "create_", maintain their original order
					return i < j
				})
				for _, frm := range freshMigration {
					for _, migration := range MigrationsDB {
						if frm.Name == migration.MigrationName {
							fmt.Printf("Fresh migration %s : \n", migration.MigrationName)
							fmt.Printf(stringers.NewString(" Drop migration...\n").SetTextColor(stringers.RED))
							errDown := migration.Down(db)
							if errDown != nil {
								fmt.Printf(stringers.NewString("Error applying migration %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errDown)
							}

							fmt.Printf(stringers.NewString(" Applying migration...\n").SetTextColor(stringers.BLUE))
							errUp := migration.Up(db)
							if errUp != nil {
								fmt.Printf(stringers.NewString("Error applying migration %s: %v \n").SetTextColor(stringers.RED), migration.MigrationName, errUp)
							}

							countMigrate += 1
						}
					}
				}
			}
		}
	}
}

func RefreshMigrate(db *gorm.DB) {
	var tables []string
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	if cfg.GetDBDriver() == "mysql" {
		err = db.Raw("SHOW TABLES").Scan(&tables).Error

		if err != nil {
			log.Fatal("Failed to fetch tables:", err)
		}

		db.Exec("SET FOREIGN_KEY_CHECKS = 0")
		for _, table := range tables {
			if table != "migrations" {
				fmt.Printf("Refresh table: %s\n", table)
				if err = db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", table)).Error; err != nil {
					log.Printf("Failed to refresh table %s: %v", table, err)
				}
			}
		}
		db.Exec("SET FOREIGN_KEY_CHECKS = 1")

		fmt.Println("All tables have been refresh.")
	} else if cfg.GetDBDriver() == "postgres" {
		err = db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = 'public'").Scan(&tables).Error
		if err != nil {
			log.Fatal("Failed to fetch tables:", err)
		}
		for _, table := range tables {
			if table != "migrations" {
				fmt.Printf("Refresh table: %s\n", table)
				if err = db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s` CASCADE", table)).Error; err != nil {
					log.Printf("Failed to refresh table %s: %v", table, err)
				}
			}
		}
		fmt.Println("All tables have been refresh.")
	} else {
		fmt.Printf(stringers.NewString("Refresh failed database using %s\n").SetTextColor(stringers.RED), cfg.GetDBDriver())
	}
}

func getTableName(s string) string {
	sPos := strings.Index(s, "create_") + len("create_")
	ePos := strings.Index(s, "_table")
	if sPos < ePos {
		return s[sPos:ePos]
	}
	return ""
}

func getModelName(s string) string {
	spl := strings.Split(s, "_")
	if len(spl[0]) > 0 {
		if len(spl) < 2 {
			spl[0] = strings.ToUpper(spl[0])
			return strings.Join(spl, "")
		} else {
			for i, word := range spl {
				runes := []rune(word)
				runes[0] = unicode.ToUpper(runes[0])
				spl[i] = string(runes)
			}
			return strings.Join(spl, "")
		}
	}
	return ""
}

func SelfCompile(filePath, outPath, exeName string) error {
	cmd := exec.Command("go", "build", "-o", outPath+"/"+exeName+".exe", filePath)
	// Run the compilation command
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to compile %s: %s\n%s", filePath, err, output)
	}

	log.Printf("successfully compiled %s to %s", filePath, exeName+".exe")
	return nil
}
