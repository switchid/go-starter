package migration

import (
	"GoStarter/pkg/config"
	"GoStarter/pkg/utils/helpers"
	"bytes"
	"fmt"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
	"unicode"
)

type Migration struct {
	ID        uint      `gorm:"primaryKey"`
	Name      string    `gorm:"unique"`
	Batch     int       `gorm:"default:0"`
	AppliedAt time.Time `gorm:"autoCreateTime"`
}

type Migrations struct {
	MigrationName string
	Up            func(db *gorm.DB) error
	Down          func(db *gorm.DB) error
}

var MigrationsDB []Migrations

func RunMigrate(db *gorm.DB) {
	if !db.Migrator().HasTable(&Migration{}) {
		err := db.AutoMigrate(&Migration{})
		if err != nil {
			log.Fatalf("failed to create migration table: %v", err)
		}
	}
	var batch, countMigrate int
	countMigrate = 0
	err := db.Table("migrations").Select("IFNULL(MAX(Batch),0)+1 AS Batch").Scan(&batch).Error
	if err != nil {
		fmt.Println("Error to get batch:", err)
	} else {
		for _, migration := range MigrationsDB {
			var allMigrations []Migration
			db.Find(&allMigrations, "Name = ?", migration.MigrationName)
			if len(allMigrations) < 1 {
				fmt.Printf("Migration %s : \n", migration.MigrationName)
				fmt.Printf(helpers.Strings.SetTextRed(" Drop migration...\n"))
				err = migration.Down(db)
				if err != nil {
					fmt.Printf(helpers.Strings.SetTextRed("Error applying migration %s: %v \n"), migration.MigrationName, err)
				}

				fmt.Printf(helpers.Strings.SetTextBlue(" Applying migration...\n"))
				err = migration.Up(db)
				if err != nil {
					fmt.Printf(helpers.Strings.SetTextRed("Error applying migration %s: %v \n"), migration.MigrationName, err)
				}

				err = db.Create(&Migration{Name: migration.MigrationName, Batch: batch}).Error
				if err != nil {
					fmt.Printf(helpers.Strings.SetTextRed("Error to insert migrations %s: %v \n"), migration.MigrationName, err)
				}

				countMigrate += 1
				fmt.Printf(helpers.Strings.SetTextGreen(" Migration applied successfully... \n"))
			}
		}
	}
	if countMigrate > 0 {
		fmt.Printf("Migrations completed successfully. \n")
	} else {
		fmt.Printf(helpers.Strings.SetTextRed("Nothing to migrate. \n"))
	}
}

func CreateMigration(migrationName string, withModel bool, dirModel string) {
	var dataMigration struct {
		MigrationName string
		NameTable     string
		NameModel     string
		WithModel     bool
		DirModel      string
	}

	timestampMigration := time.Now().Format("20060102150405")
	filenameMigration := fmt.Sprintf("%s_%s.go", timestampMigration, migrationName)
	dir, _ := helpers.GetProjectPath()
	migrationDir := filepath.Join(dir, "scripts/migration")
	err := os.MkdirAll(migrationDir, 0755)
	if err != nil {
		fmt.Println("Error creating migration directory: %\n", err)
		return
	}
	filePathMigration := filepath.Join(migrationDir, filenameMigration)

	dataMigration.MigrationName = strings.Replace(filenameMigration, ".go", "", -1)
	dataMigration.NameTable = getTableName(migrationName)
	dataMigration.NameModel = getModelName(dataMigration.NameTable)
	dataMigration.WithModel = withModel
	dataMigration.DirModel = dirModel

	if strings.Contains(migrationName, "create_") {
		tplMigrationCreate := `package migration
import (
	internalMigration "GoStarter/internal/cli/commands/migration"
	{{ if .WithModel }}
		{{- if eq (.DirModel) "" -}}
	"GoStarter/internal/models"
		{{- else -}}
	"GoStarter/internal/models/{{.DirModel}}"
		{{- end -}}
	{{ else }}
	"time"
	{{ end }}
	"fmt"
	"gorm.io/gorm"
)
{{ if not (.WithModel) }}
type {{.NameModel}} struct {
	Id        uint64 ` + "`" + `gorm:"primaryKey;autoIncrement:true"` + "`" + `
	CreatedAt time.Time
	UpdatedAt time.Time ` + "`" + `gorm:"autoCreateTime:false"` + "`" + `
}
{{ end }}
func init() {
	{{ if .WithModel }}
	var m{{.NameModel}} models.{{.NameModel}}
	{{ else }}
	var m{{.NameModel}} {{.NameModel}}
	{{ end }}
	internalMigration.MigrationsDB = append(internalMigration.MigrationsDB, internalMigration.Migrations{
		MigrationName: "{{.MigrationName}}",
		Up: func(db *gorm.DB) error {
			if !db.Migrator().HasTable({{printf "&m%s" .NameModel}}) {
				if err := db.Migrator().CreateTable({{printf "&m%s" .NameModel}}); err != nil {
					_ = fmt.Errorf("failed to create users table: %w", err)
					return err
				}
			}

			return nil
		},
		Down: func(db *gorm.DB) error {
			if db.Migrator().HasTable({{printf "&m%s" .NameModel}}) {
				if err := db.Migrator().DropTable({{printf "&m%s" .NameModel}}); err != nil {
					_ = fmt.Errorf("failed to drop users table: %w", err)
					return err
				}
			}
			return nil
		},
	})
}
	
{{- if not (.WithModel) -}}
func ({{.NameModel}}) TableName() string {
	return "{{.NameTable}}"
}
{{- end -}}`

		contentMigration := template.Must(template.New(filenameMigration).Parse(tplMigrationCreate))
		var outputMigration bytes.Buffer
		err = contentMigration.Execute(&outputMigration, dataMigration)
		if err != nil {
			return
		}

		if !helpers.FileExistsWithContains(migrationDir, migrationName) {
			err = os.WriteFile(filePathMigration, outputMigration.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error creating migration file: %v\n", err)
				return
			}

			fmt.Printf("Migration file created: %s\n", filePathMigration)

			if withModel {
				if dirModel != "" {
					createModel(dataMigration.NameModel, dataMigration.NameTable, dirModel)
				} else {
					createModel(dataMigration.NameModel, dataMigration.NameTable, "")
				}
			}
		} else {
			fmt.Printf("Migration file failed to create! file is exists: %s\n", filePathMigration)
		}
	} else {
		tplMigration := `package migration
import (
	internalMigration "GoStarter/internal/cli/commands/migration"
	"gorm.io/gorm"
)

func init(){
	internalMigration.MigrationsDB = append(internalMigration.MigrationsDB, internalMigration.Migrations{
		MigrationName: "{{.MigrationName}}",
		Up: func(db *gorm.DB) error {
			
			return nil
		},
		Down: func(db *gorm.DB) error {
			
			return nil
		},
	})
}`
		contentMigration := template.Must(template.New(filenameMigration).Parse(tplMigration))
		var outputMigration bytes.Buffer
		err = contentMigration.Execute(&outputMigration, dataMigration)
		if err != nil {
			return
		}

		if !helpers.FileExistsWithContains(migrationDir, migrationName) {
			err = os.WriteFile(filePathMigration, outputMigration.Bytes(), 0644)
			if err != nil {
				fmt.Printf("Error creating migration file: %v\n", err)
				return
			}

			fmt.Printf("Migration file created: %s\n", filePathMigration)

			if withModel {
				if dirModel != "" {
					createModel(dataMigration.NameModel, dataMigration.NameTable, dirModel)
				} else {
					createModel(dataMigration.NameModel, dataMigration.NameTable, "")
				}
			}
		} else {
			fmt.Printf("Migration file failed to create! file is exists: %s\n", filePathMigration)
		}
	}
}

func RollbackMigrate(db *gorm.DB) {
	var batch, countMigrate int
	countMigrate = 0
	err := db.Table("migrations").Select("IFNULL(MAX(Batch),0) AS Batch").Scan(&batch).Error
	if err != nil {
		fmt.Println("Error to get batch:", err)
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
						fmt.Printf(helpers.Strings.SetTextRed(" Drop migration...\n"))
						err = migration.Down(db)
						if err != nil {
							fmt.Printf(helpers.Strings.SetTextRed("Error applying migration %s: %v \n"), migration.MigrationName, err)
						}
						countMigrate += 1
					}
				}
			}
		}
	}
	if countMigrate > 0 {
		err = db.Delete(&Migration{}, "Batch = ?", batch).Error
		if err != nil {
			fmt.Println("Error to delete migrations:", err)
		}
		fmt.Printf("Rollback migrations completed successfully. \n")
	} else {
		fmt.Printf(helpers.Strings.SetTextRed("Nothing to rollback. \n"))
	}
}

func FreshMigrate(db *gorm.DB) {
	var batch, countMigrate int
	countMigrate = 0
	err := db.Table("migrations").Select("IFNULL(MAX(Batch),0)+1 AS Batch").Scan(&batch).Error
	if err != nil {
		fmt.Println("Error to get batch:", err)
	} else {
		var freshMigration []Migration
		err = db.Table("migrations").Select("Name").Scan(&freshMigration).Error
		if err != nil {
			fmt.Println("Error to get migrations:", err)
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
							fmt.Printf(helpers.Strings.SetTextRed(" Drop migration...\n"))
							err = migration.Down(db)
							if err != nil {
								fmt.Printf(helpers.Strings.SetTextRed("Error applying migration %s: %v \n"), migration.MigrationName, err)
							}

							fmt.Printf(helpers.Strings.SetTextBlue(" Applying migration...\n"))
							err = migration.Up(db)
							if err != nil {
								fmt.Printf(helpers.Strings.SetTextRed("Error applying migration %s: %v \n"), migration.MigrationName, err)
							}

							//err = db.Create(&Migration{Name: migration.MigrationName, Batch: batch}).Error
							//if err != nil {
							//	fmt.Printf(helpers.Strings.SetTextRed("Error to insert migrations %s: %v \n"), migration.MigrationName, err)
							//}

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
		fmt.Printf(helpers.Strings.SetTextRed("Refresh failed database using %s\n"), cfg.GetDBDriver())
	}
}

func createModel(modelName string, tableName string, modelDir string) {
	var dataModel struct {
		Name      string
		TableName string
		Dir       string
	}
	dataModel.Name = modelName
	dataModel.Dir = modelDir
	dataModel.TableName = tableName

	fileNameModel := fmt.Sprintf("%s.go", modelName)
	dir, _ := helpers.GetProjectPath()
	modelFileDir := filepath.Join(dir, "internal/models/"+dataModel.Dir)
	err := os.MkdirAll(modelFileDir, 0755)
	if err != nil {
		fmt.Println("Error creating model directory: %\n", err)
		return
	}
	filePathModel := filepath.Join(modelFileDir, fileNameModel)

	tplModel := `package models
import "time"
	
type {{.Name}} struct {
	Id        uint64 ` + "`" + `gorm:"primaryKey;autoIncrement:true"` + "`" + `
	CreatedAt time.Time
	UpdatedAt time.Time ` + "`" + `gorm:"autoCreateTime:false"` + "`" + `
}

func ({{.Name}}) TableName() string {
	return "{{.TableName}}"
}`

	contentModel := template.Must(template.New(fileNameModel).Parse(tplModel))
	var outputModel bytes.Buffer
	err = contentModel.Execute(&outputModel, dataModel)
	if err != nil {
		return
	}

	if !helpers.FileExists(filePathModel) {
		err = os.WriteFile(filePathModel, outputModel.Bytes(), 0644)
		if err != nil {
			fmt.Printf("Error creating migration file: %v\n", err)
			return
		}

		fmt.Printf("Model file created: %s\n", filePathModel)
	} else {
		_ = fmt.Errorf("Model file failed to create! file is exists: %s\n", filePathModel)
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
