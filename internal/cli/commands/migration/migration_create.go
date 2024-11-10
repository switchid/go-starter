package migration

import (
	"GoStarter/pkg/utils/paths"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

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
	dir, _ := paths.GetProjectPath()
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
	Id        uint ` + "`" + `gorm:"primaryKey;autoIncrement:true"` + "`" + `
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

		if !paths.FileExistsWithContains(migrationDir, migrationName) {
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

		if !paths.FileExistsWithContains(migrationDir, migrationName) {
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
	projDir, _ := paths.GetProjectPath()
	err = SelfCompile(filepath.Join(projDir, "cmd/cli"), filepath.Join(projDir, "build"), "tools")
	if err != nil {
		fmt.Printf("Migration file failed to create! self compile: %s\n", err.Error())
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
	dir, _ := paths.GetProjectPath()
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

	if !paths.FileExists(filePathModel) {
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
