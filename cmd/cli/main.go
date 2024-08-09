package main

import (
	"GoStarter/internal/cli/commands/migration"
	"GoStarter/internal/database"
	"GoStarter/pkg/config"
	"GoStarter/pkg/utils/helpers"
	"GoStarter/pkg/utils/loggers"
	_ "GoStarter/scripts/migration"
	"flag"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"os"
	"path/filepath"
)

var (
	cmdInstall               = flag.NewFlagSet("install", flag.ExitOnError)
	cmdInstallSvcName        = cmdInstall.String("svcname", "", "service name")
	cmdInstallSvcDisplay     = cmdInstall.String("svcdisplay", "", "display name")
	cmdInstallSvcDescription = cmdInstall.String("svcdescription", "", "description")
	cmdInstallSvcApi         = cmdInstall.Bool("api", true, "with separation api")
	cmdStop                  = flag.NewFlagSet("stop", flag.ExitOnError)
	cmdStopSvcname           = cmdStop.String("svcname", "", "service name")
	cmdUninstall             = flag.NewFlagSet("uninstall", flag.ExitOnError)
	cmdUninstallSvcname      = cmdUninstall.String("svcname", "", "service name")
	cmdMigration             = flag.NewFlagSet("migration", flag.ExitOnError)
	cmdMigrationMigrate      = cmdMigration.Bool("migrate", false, "migrate database [none]")
	cmdMigrationRollback     = cmdMigration.Bool("rollback", false, "rollback migration [none]")
	cmdMigrationFresh        = cmdMigration.Bool("fresh", false, "fresh migration [none]")
	cmdMigrationRefresh      = cmdMigration.Bool("refresh", false, "refresh migration [none]")
	cmdMigrationCreate       = flag.NewFlagSet("create", flag.ExitOnError)
	cmdMigrationCreateModel  = cmdMigrationCreate.Bool("model", false, "add model migration [bool]")
	cmdMigrationDirModel     = cmdMigrationCreate.String("model-dir", "", "model directory [string]")
)

type serviceApp struct {
	name        string
	display     string
	description string
}

func main() {
	var svcApps serviceApp
	flag.Parse()
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	svcApps.name = cfg.GetAppName()
	svcApps.display = cfg.GetAppName()
	svcApps.description = fmt.Sprintf("Service Application of %s", cfg.GetAppName())

	if len(os.Args) < 2 {
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
				flagIndex := getFlagIndex(args)
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
				fmt.Printf("Expected Subcommand :\n")
				fmt.Printf(" %s \n", cmdMigrationCreate.Name())
				fmt.Printf("  Flags:\n")
				cmdMigrationCreate.PrintDefaults()
				os.Exit(1)
			}
		}
	default:
		fmt.Printf("%q is not valid command.\n", os.Args[1])
		os.Exit(1)
	}

	if cmdInstall.Parsed() {
		if *cmdInstallSvcName != "" {
			svcApps.name = *cmdInstallSvcName
			svcApps.display = *cmdInstallSvcName
			svcApps.description = fmt.Sprintf("Service Application of %s", *cmdInstallSvcName)
		}

		ex, err := os.Executable()
		if err != nil {
			log.Fatalf("Error getting executable path: %v", err)
		}
		exePath := filepath.Dir(ex)

		var helperFlag helpers.Flags
		var appFlags string
		if *cmdInstallSvcName != "" {
			appFlags = helperFlag.SetFlag("svcname", *cmdInstallSvcName)
			svcApps.name = *cmdInstallSvcName
		}
		if *cmdInstallSvcDisplay != "" {
			svcApps.display = *cmdInstallSvcDisplay
		}
		if *cmdInstallSvcDescription != "" {
			svcApps.description = *cmdInstallSvcDescription
		}

		appPath := filepath.Join(exePath, "app.exe")

		err = installService(svcApps.name, svcApps.display, svcApps.description, appPath, appFlags)
		if err != nil {
			log.Printf("Failed to install app service: %v", err)
		} else {
			log.Println("App service installed successfully")
		}

		if *cmdInstallSvcApi {
			apiPath := filepath.Join(exePath, "api.exe")

			apiErr := installService(svcApps.name+"api", svcApps.display+" (Api)", svcApps.description+" (Api)", apiPath, "")
			if apiErr != nil {
				log.Printf("Failed to install api service: %v", apiErr)
			} else {
				log.Printf("Api service installed successfully")
			}
		}
	}

	if cmdStop.Parsed() {
		if *cmdStopSvcname != "" {
			svcApps.name = *cmdStopSvcname
		}
		err := stopService(svcApps.name)
		if err != nil {
			log.Printf("Failed to stop app service: %v", err)
		} else {
			log.Println("App service stopped successfully")
		}
	}

	if cmdUninstall.Parsed() {
		if *cmdUninstallSvcname != "" {
			svcApps.name = *cmdUninstallSvcname
		}
		err := uninstallService(svcApps.name)
		if err != nil {
			log.Printf("Failed to uninstall app service: %v", err)
		} else {
			log.Println("App service uninstalled successfully")
		}
	}

	if cmdMigration.Parsed() {
		if *cmdMigrationMigrate {
			db, _ := database.Connect()
			migration.RunMigrate(db)
		}

		if *cmdMigrationRollback {
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

func init() {
	// Customize usage output
	flag.Usage = func() {
		if _, err := fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0]); err != nil {
			fmt.Println("Error writing to stderr:", err)
			return
		}
		if _, err := fmt.Fprintf(os.Stderr, "  %s [args]\n", os.Args[0]); err != nil {
			fmt.Println("Error writing to stderr:", err)
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Service tool"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Install Service"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  install [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  Flags:"); err != nil {
			return
		}
		cmdInstall.PrintDefaults()
		if _, err := fmt.Fprintln(os.Stderr, " Stop Service"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  stop [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  Flags:"); err != nil {
			return
		}
		cmdStop.PrintDefaults()
		if _, err := fmt.Fprintln(os.Stderr, " Uninstall Service"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  uninstall [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  Flags:"); err != nil {
			return
		}
		cmdUninstall.PrintDefaults()
		if _, err := fmt.Fprintln(os.Stderr, "Migration tool \n migration [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Flags:"); err != nil {
			return
		}
		cmdMigration.PrintDefaults()
		if _, err := fmt.Fprintln(os.Stderr, " migration create [name] [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Flags:"); err != nil {
			return
		}
		cmdMigrationCreate.PrintDefaults()
	}

	cmdMigration.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "migration [args] [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Args:"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "  Create [name] [flag]"); err != nil {
			return
		}
		cmdMigrationCreate.PrintDefaults()
		if _, err := fmt.Fprintln(os.Stderr, " Flags:"); err != nil {
			return
		}
		cmdMigration.PrintDefaults()
	}

	cmdMigrationCreate.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "migration create [name] [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Flags:"); err != nil {
			return
		}
		cmdMigrationCreate.PrintDefaults()
	}
}

func installService(serviceName, displayNameService, descriptionService, exePath, appFlags string) error {
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		errDC := m.Disconnect()
		if errDC != nil {
			return
		}
	}(m)

	s, err := m.OpenService(serviceName)
	if err == nil {
		errClose := s.Close()
		if errClose != nil {
			return errClose
		}
		return fmt.Errorf("service %s already exists", serviceName)
	}

	s, err = m.CreateService(serviceName, exePath, mgr.Config{
		DisplayName:      displayNameService,
		Description:      descriptionService,
		StartType:        mgr.StartAutomatic,
		ServiceStartName: "LocalSystem",
	}, appFlags)
	if err != nil {
		return err
	}
	defer func(s *mgr.Service) {
		errClose := s.Close()
		if errClose != nil {

		}
	}(s)

	err = eventlog.Remove(serviceName)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("failed to remove existing event log source: %v", err)
		// Continue with installation even if this fails
	}

	err = eventlog.InstallAsEventCreate(serviceName, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		errDel := s.Delete()
		if errDel != nil {
			return errDel
		}
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}

	errRun := runService(serviceName)
	if errRun != nil {
		return errRun
	}

	return nil
}

func runService(name string) error {
	lgs, err := loggers.NewLogger()
	m, err := mgr.Connect()
	if err != nil {
		return err
	}
	defer func(m *mgr.Mgr) {
		errDC := m.Disconnect()
		if errDC != nil {

		}
	}(m)
	s, err := m.OpenService(name)
	if err != nil {
		lgs.LogError("couldn't open service: %v", err)
		return fmt.Errorf("could not access service: %v", err)
	}
	defer func(s *mgr.Service) {
		errClose := s.Close()
		if errClose != nil {
			return
		}
	}(s)
	err = s.Start("is", "auto-started")
	if err != nil {
		lgs.LogError("couldn't start service: %v", err)
		return fmt.Errorf("could not start service: %v", err)
	}
	return nil
}

func stopService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %v", err)
	}
	defer func(m *mgr.Mgr) {
		errDC := m.Disconnect()
		if errDC != nil {

		}
	}(m)

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer func(s *mgr.Service) {
		errClose := s.Close()
		if errClose != nil {
			return
		}
	}(s)

	// Stop the service if it's running
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("failed to query service status: %v", err)
	}

	if status.State != svc.Stopped {
		_, err = s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("failed to stop service: %v", err)
		}
	}

	return nil
}

func uninstallService(name string) error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %v", err)
	}
	defer func(m *mgr.Mgr) {
		errDC := m.Disconnect()
		if errDC != nil {
			return
		}
	}(m)

	s, err := m.OpenService(name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", name)
	}
	defer func(s *mgr.Service) {
		errClose := s.Close()
		if errClose != nil {
			return
		}
	}(s)

	// Stop the service if it's running
	status, err := s.Query()
	if err != nil {
		return fmt.Errorf("failed to query service status: %v", err)
	}

	if status.State != svc.Stopped {
		_, err = s.Control(svc.Stop)
		if err != nil {
			return fmt.Errorf("failed to stop service: %v", err)
		}
	}

	err = eventlog.Remove(name)
	if err != nil {
		log.Printf("failed to remove event log source: %v", err)
		// Continue with uninstallation even if this fails
	}

	// Delete the service
	err = s.Delete()
	if err != nil {
		return fmt.Errorf("failed to delete service: %v", err)
	}

	return nil
}

func getFlagIndex(flags []string) int {
	for i, arg := range flags {
		if arg[0] == '-' {
			return i
		}
	}
	return 0
}
