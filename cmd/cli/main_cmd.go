package main

import (
	"GoStarter/pkg/utils/helpers"
	"GoStarter/pkg/utils/paths"
	"GoStarter/pkg/utils/stringers"
	"flag"
	"fmt"
	"os"
)

var (
	cmdGenerateKey           *bool         = nil
	cmdInstall               *flag.FlagSet = nil
	cmdInstallSvcName        *string       = nil
	cmdInstallSvcDisplay     *string       = nil
	cmdInstallSvcDescription *string       = nil
	cmdInstallSvcApi         *bool         = nil
	cmdStop                  *flag.FlagSet = nil
	cmdStopSvcname           *string       = nil
	cmdUninstall             *flag.FlagSet = nil
	cmdUninstallSvcname      *string       = nil
	cmdMigration             *flag.FlagSet = nil
	cmdMigrationMigrate      *bool         = nil
	cmdMigrationRollback     *bool         = nil
	cmdMigrationFresh        *bool         = nil
	cmdMigrationRefresh      *bool         = nil
	cmdMigrationCreate       *flag.FlagSet = nil
	cmdMigrationCreateModel  *bool         = nil
	cmdMigrationDirModel     *string       = nil
)

func init() {
	commandSet()

	flag.Usage = func() {
		exeName, errExeName := paths.GetCurrentExecutableName()
		if errExeName != nil {
			_ = fmt.Errorf("error get executable name %s", errExeName)
			return
		}
		if _, err := fmt.Fprintf(os.Stderr, stringers.NewString("Usage of %s:\n").SetTextColor(stringers.BLUE), exeName); err != nil {
			fmt.Println("Error writing to stderr:", err)
			return
		}
		if _, err := fmt.Fprintf(os.Stderr, stringers.NewString("  %s [args]\n").SetTextColor(stringers.BLUE), exeName); err != nil {
			fmt.Println("Error writing to stderr:", err)
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString("Application tool").SetTextColor(stringers.YELLOW)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" -generate-key").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString("Service tool").SetTextColor(stringers.YELLOW)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" install [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" stop [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" uninstall [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString("Migration Tool").SetTextColor(stringers.YELLOW)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" migration [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
	}

	cmdMigration.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "migration [args] [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Args:"); err != nil {
			return
		}
		if cmdMigrationCreate != nil {
			if _, err := fmt.Fprintln(os.Stderr, "  Create [name] [flag]"); err != nil {
				return
			}
			cmdMigrationCreate.PrintDefaults()
		}
		if _, err := fmt.Fprintln(os.Stderr, " Flags:"); err != nil {
			return
		}
		cmdMigration.PrintDefaults()
	}

	if cmdMigrationCreate != nil {
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
}

func commandSet() {
	cmdGenerateKey = flag.Bool("generate-key", false, "")
	cmdInstall = flag.NewFlagSet("install", flag.ExitOnError)
	cmdInstallSvcName = cmdInstall.String("svcname", "", "service name")
	cmdInstallSvcDisplay = cmdInstall.String("svcdisplay", "", "display name")
	cmdInstallSvcDescription = cmdInstall.String("svcdescription", "", "description")
	cmdInstallSvcApi = cmdInstall.Bool("api", true, "with separation api")
	cmdStop = flag.NewFlagSet("stop", flag.ExitOnError)
	cmdStopSvcname = cmdStop.String("svcname", "", "service name")
	cmdUninstall = flag.NewFlagSet("uninstall", flag.ExitOnError)
	cmdUninstallSvcname = cmdUninstall.String("svcname", "", "service name")
	cmdMigration = flag.NewFlagSet("migration", flag.ExitOnError)
	cmdMigrationMigrate = cmdMigration.Bool("migrate", false, "migrate database [none]")
	cmdMigrationFresh = cmdMigration.Bool("fresh", false, "fresh migration [none]")
	cmdMigrationRefresh = cmdMigration.Bool("refresh", false, "refresh migration [none]")

	if helpers.DevMode() {
		cmdMigrationRollback = cmdMigration.Bool("rollback", false, "rollback migration [none]")
		cmdMigrationCreate = flag.NewFlagSet("create", flag.ExitOnError)
		cmdMigrationCreateModel = cmdMigrationCreate.Bool("model", false, "add model migration [bool]")
		cmdMigrationDirModel = cmdMigrationCreate.String("model-dir", "", "model directory [string]")
	} else {
		cmdMigrationRollback = nil
		cmdMigrationCreate = nil
	}
}
