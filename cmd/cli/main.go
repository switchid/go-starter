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

	helpers.EnableANSI()

	cfg, errCfg := config.Load()
	if errCfg != nil {
		log.Fatalf("Failed to load configuration: %v", errCfg)
	} else {
		svcApp.Name = cfg.GetAppName()
		svcApp.Display = cfg.GetAppName()
		svcApp.Description = fmt.Sprintf("Service Application of %s", cfg.GetAppName())

		svcApi.Name = cfg.GetAppName() + "Api"
		svcApi.Display = cfg.GetAppName() + " (API)"
		svcApi.Description = fmt.Sprintf("Service Application of %s", cfg.GetAppName()+" (API)")
	}

	flag.Parse()

	if *cmdGenerateKey {
		generator.GenerateApplicationKey()
		os.Exit(1)
	} else if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "service":
		errService := cmdService.Parse(os.Args[2:])
		if errService != nil {
			return
		}
		if cmdService.NArg() > 0 {
			switch os.Args[2] {
			case "install":
				if len(os.Args) < 3 {
					cmdServiceInstall.Usage()
					os.Exit(1)
				}
				args := os.Args[3:]
				flagIndex := helpers.GetFlagIndex(args)
				flagArgs := args[flagIndex:]
				errServiceInstall := cmdServiceInstall.Parse(flagArgs)
				if errServiceInstall != nil {
					return
				}
			case "start":
				if len(os.Args) < 3 {
					cmdServiceStart.Usage()
					os.Exit(1)
				}
				args := os.Args[3:]
				flagIndex := helpers.GetFlagIndex(args)
				flagArgs := args[flagIndex:]
				errServiceStart := cmdServiceStart.Parse(flagArgs)
				if errServiceStart != nil {
					return
				}
			case "stop":
				if len(os.Args) < 3 {
					cmdServiceStop.Usage()
					os.Exit(1)
				}
				args := os.Args[3:]
				flagIndex := helpers.GetFlagIndex(args)
				flagArgs := args[flagIndex:]
				errServiceStop := cmdServiceStop.Parse(flagArgs)
				if errServiceStop != nil {
					return
				}
			case "uninstall":
				if len(os.Args) < 3 {
					cmdServiceUninstall.Usage()
					os.Exit(1)
				}
				args := os.Args[3:]
				flagIndex := helpers.GetFlagIndex(args)
				flagArgs := args[flagIndex:]
				errMigrationCreate := cmdServiceUninstall.Parse(flagArgs)
				if errMigrationCreate != nil {
					return
				}
			}
		} else {
			cmdService.Usage()
		}
	case "migration":
		errMigration := cmdMigration.Parse(os.Args[2:])
		if errMigration != nil {
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

	if cmdServiceInstall.Parsed() {
		if *cmdServiceInstallSvcName != "" {
			svcApp.Name = *cmdServiceInstallSvcName
			svcApp.Display = *cmdServiceInstallSvcName
			svcApp.Description = fmt.Sprintf("Service Application of %s", *cmdServiceInstallSvcName)

			if *cmdServiceInstallSvcApi {
				svcApi.Name = *cmdServiceInstallSvcName + "Api"
				svcApi.Display = *cmdServiceInstallSvcName + " (API)"
				svcApi.Description = fmt.Sprintf("Service Application of %s", *cmdServiceInstallSvcName+" (API)")
			}
		}
		if *cmdServiceInstallSvcDisplay != "" {
			svcApp.Display = *cmdServiceInstallSvcDisplay
			svcApi.Display = *cmdServiceInstallSvcDisplay
		}
		if *cmdServiceInstallSvcDescription != "" {
			svcApp.Description = *cmdServiceInstallSvcDescription
			svcApi.Description = *cmdServiceInstallSvcDescription
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
		if *cmdServiceInstallSvcApi {
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

	if cmdServiceStart.Parsed() {
		var svcSingle services.Services
		if *cmdServiceStartSvcname != "" {
			svcSingle.Name = *cmdServiceStartSvcname
			startSingleErr := svcSingle.RunService()
			if startSingleErr != nil {
				log.Printf("Failed to start app service: %v", startSingleErr)
			} else {
				log.Println("App service started successfully")
			}
		} else {
			startAppErr := svcApp.RunService()
			if startAppErr != nil {
				log.Printf("Failed to start app service: %v", startAppErr)
			} else {
				log.Printf("\033[31mApp service started successfully")
			}
			startApiErr := svcApi.RunService()
			if startApiErr != nil {
				log.Printf("Failed to start api service: %v", startApiErr)
			} else {
				log.Printf("Api service started successfully")
			}
		}
	}

	if cmdServiceStop.Parsed() {
		var svcSingle services.Services
		if *cmdServiceStopSvcname != "" {
			svcSingle.Name = *cmdServiceStopSvcname
			stopSingleErr := svcSingle.StopService()
			if stopSingleErr != nil {
				log.Printf("Failed to stop app service: %v", stopSingleErr)
			} else {
				log.Println("App service stopped successfully")
			}
		} else {
			stopAppErr := svcApp.StopService()
			if stopAppErr != nil {
				log.Printf("Failed to stop app service: %v", stopAppErr)
			} else {
				log.Println("App service stopped successfully")
			}
			stopApiErr := svcApi.StopService()
			if stopApiErr != nil {
				log.Printf("Failed to stop api service: %v", stopApiErr)
			} else {
				log.Printf("Api service stopped successfully")
			}
		}
	}

	if cmdServiceUninstall.Parsed() {
		var svcSingle services.Services
		if *cmdServiceUninstallSvcname != "" {
			svcSingle.Name = *cmdServiceUninstallSvcname
			stopSingleErr := svcSingle.UninstallService()
			if stopSingleErr != nil {
				log.Printf("Failed to stop app service: %v", stopSingleErr)
			} else {
				log.Println("App service stopped successfully")
			}
		} else {
			if *cmdServiceUninstallSvcname != "" {
				svcApp.Name = *cmdServiceUninstallSvcname
				svcApi.Name = *cmdServiceUninstallSvcname + "Api"
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
