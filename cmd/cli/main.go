package main

import (
	"GoStarter/internal/cli/commands/generator"
	"GoStarter/internal/cli/commands/migration"
	"GoStarter/internal/cli/commands/services"
	"GoStarter/internal/database"
	"GoStarter/internal/pkg/config"
	"GoStarter/pkg/utils/helpers"
	"GoStarter/pkg/utils/paths"
	_ "GoStarter/scripts/migration"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	var svcApp services.Services
	var svcApi services.Services

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	flag.Parse()

	svcApp.Name = cfg.GetAppName()
	svcApp.Display = cfg.GetAppName()
	svcApp.Description = fmt.Sprintf("Service Application of %s", cfg.GetAppName())

	svcApi.Name = cfg.GetAppName() + "Api"
	svcApi.Display = cfg.GetAppName() + " (API)"
	svcApi.Description = fmt.Sprintf("Service Application of %s", cfg.GetAppName()+" (API)")

	if *cmdGenerateKey {
		generator.GenerateApplicationKey()
		os.Exit(1)
	} else if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "install":
		errInstall := cmdInstall.Parse(os.Args[2:])
		if errInstall != nil {
			return
		}
	case "stop":
		errStop := cmdStop.Parse(os.Args[2:])
		if errStop != nil {
			return
		}
	case "uninstall":
		errorUninstall := cmdUninstall.Parse(os.Args[2:])
		if errorUninstall != nil {
			return
		}
	case "migration":
		errMigragion := cmdMigration.Parse(os.Args[2:])
		if errMigragion != nil {
			return
		}
		if cmdMigration.NArg() > 0 {
			switch os.Args[2] {
			case "create":
				if len(os.Args) <= 3 {
					cmdMigrationCreate.Usage()
					os.Exit(1)
				}
				args := os.Args[3:]
				flagIndex := helpers.GetFlagIndex(args)
				flagArgs := args[flagIndex:]
				errMigrationCreate := cmdMigrationCreate.Parse(flagArgs)
				if errMigrationCreate != nil {
					return
				}
			default:
				fmt.Printf("%q is not valid command.\n", os.Args[2])
				os.Exit(1)
			}
		} else {
			if len(os.Args) < 3 {
				fmt.Printf("Expected Flags:\n")
				cmdMigration.PrintDefaults()
				if cmdMigrationCreate != nil {
					fmt.Printf("Expected Subcommand :\n")
					fmt.Printf(" %s \n", cmdMigrationCreate.Name())
					fmt.Printf("  Flags:\n")
					cmdMigrationCreate.PrintDefaults()
				}
				os.Exit(1)
			}
		}
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(1)
	}

	if cmdInstall.Parsed() {
		if *cmdInstallSvcName != "" {
			svcApp.Name = *cmdInstallSvcName
			svcApp.Display = *cmdInstallSvcName
			svcApp.Description = fmt.Sprintf("Service Application of %s", *cmdInstallSvcName)

			if *cmdInstallSvcApi {
				svcApi.Name = *cmdInstallSvcName + "Api"
				svcApi.Display = *cmdInstallSvcName + " (API)"
				svcApi.Description = fmt.Sprintf("Service Application of %s", *cmdInstallSvcName+" (API)")
			}
		}
		if *cmdInstallSvcDisplay != "" {
			svcApp.Display = *cmdInstallSvcDisplay
			svcApi.Display = *cmdInstallSvcDisplay
		}
		if *cmdInstallSvcDescription != "" {
			svcApp.Description = *cmdInstallSvcDescription
			svcApi.Description = *cmdInstallSvcDescription
		}

		var flags helpers.Flags
		var appFlags string
		var apiFlags string
		exePath, exeErr := paths.GetCurrentExecutableDir()
		if exeErr != nil {
			log.Fatalf("Error getting executable path: %v", exeErr)
		}

		appFlags = flags.SetFlag("svcname", svcApp.Name)
		appPath := filepath.Join(exePath, "app.exe")
		appErr := svcApp.InstallService(appPath, appFlags)
		if appErr != nil {
			log.Printf("Failed to install app service: %v", appErr)
		} else {
			log.Println("App service installed successfully")
		}
		if *cmdInstallSvcApi {
			apiFlags = flags.SetFlag("svcname", svcApi.Name)
			apiPath := filepath.Join(exePath, "api.exe")
			apiErr := svcApi.InstallService(apiPath, apiFlags)
			if apiErr != nil {
				log.Printf("Failed to install api service: %v", apiErr)
			} else {
				log.Printf("Api service installed successfully")
			}
		}
		if cfg != nil {
			cfg.SetServiceAppName(svcApp.Name)
			cfg.SetServiceAppDisplay(svcApp.Display)
			cfg.SetServiceAppDescription(svcApp.Description)
			cfg.SetServiceApiName(svcApi.Name)
			cfg.SetServiceApiDisplay(svcApi.Display)
			cfg.SetServiceApiDescription(svcApi.Description)
			errCfgSave := cfg.SaveConfig()
			if errCfgSave != nil {
				log.Fatalf("Error save service config: %v", exeErr)
			}
		}
	}

	if cmdStop.Parsed() {
		var svcSingle services.Services
		if *cmdStopSvcname != "" {
			svcSingle.Name = *cmdStopSvcname
			stopSingleErr := svcSingle.StopService()
			if stopSingleErr != nil {
				log.Printf("Failed to stop app service: %v", stopSingleErr)
			} else {
				log.Println("App service stopped successfully")
			}
		} else {
			svcApp.Name = *cmdStopSvcname
			svcApi.Name = *cmdStopSvcname

			stopAppErr := svcApp.StopService()
			if stopAppErr != nil {
				log.Printf("Failed to stop app service: %v", stopAppErr)
			} else {
				log.Println("App service stopped successfully")
			}
		}

	}

	if cmdUninstall.Parsed() {
		var svcSingle services.Services
		if *cmdUninstallSvcname != "" {
			svcSingle.Name = *cmdUninstallSvcname
			stopSingleErr := svcSingle.UninstallService()
			if stopSingleErr != nil {
				log.Printf("Failed to stop app service: %v", stopSingleErr)
			} else {
				log.Println("App service stopped successfully")
			}
		} else {
			if *cmdUninstallSvcname != "" {
				svcApp.Name = *cmdUninstallSvcname
				svcApi.Name = *cmdUninstallSvcname + "Api"
			}
			uninstallAppErr := svcApp.UninstallService()
			if uninstallAppErr != nil {
				log.Printf("Failed to uninstall app service: %v", uninstallAppErr)
			} else {
				log.Println("App service uninstalled successfully")
			}
			uninstallApiErr := svcApi.UninstallService()
			if uninstallApiErr != nil {
				log.Printf("Failed to uninstall api service: %v", uninstallApiErr)
			} else {
				log.Println("Api service uninstalled successfully")
			}
		}
	}

	if cmdMigration.Parsed() {
		if *cmdMigrationMigrate {
			db, _ := database.Connect()
			migration.RunMigrate(db)
		}

		if cmdMigrationRollback != nil && *cmdMigrationRollback {
			db, _ := database.Connect()
			migration.RollbackMigrate(db)
		}

		if *cmdMigrationFresh {
			db, _ := database.Connect()
			migration.FreshMigrate(db)
		}

		if *cmdMigrationRefresh {
			db, _ := database.Connect()
			migration.RefreshMigrate(db)
		}
	}

	if cmdMigrationCreate != nil {
		if cmdMigrationCreate.Parsed() {
			migrationName := os.Args[3]

			if *cmdMigrationCreateModel {
				if *cmdMigrationDirModel != "" {
					migration.CreateMigration(migrationName, true, *cmdMigrationDirModel)
				} else {
					migration.CreateMigration(migrationName, true, "")
				}
			} else {
				migration.CreateMigration(migrationName, false, "")
			}
		}
	}

}
