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
	cmdGenerateKey                  *bool         = nil
	cmdService                      *flag.FlagSet = nil
	cmdServiceInstall               *flag.FlagSet = nil
	cmdServiceInstallSvcName        *string       = nil
	cmdServiceInstallSvcDisplay     *string       = nil
	cmdServiceInstallSvcDescription *string       = nil
	cmdServiceInstallSvcApi         *bool         = nil
	cmdServiceStart                 *flag.FlagSet = nil
	cmdServiceStartSvcname          *string       = nil
	cmdServiceStop                  *flag.FlagSet = nil
	cmdServiceStopSvcname           *string       = nil
	cmdServiceUninstall             *flag.FlagSet = nil
	cmdServiceUninstallSvcname      *string       = nil
	cmdMigration                    *flag.FlagSet = nil
	cmdMigrationMigrate             *bool         = nil
	cmdMigrationRollback            *bool         = nil
	cmdMigrationFresh               *bool         = nil
	cmdMigrationRefresh             *bool         = nil
	cmdMigrationCreate              *flag.FlagSet = nil
	cmdMigrationCreateModel         *bool         = nil
	cmdMigrationDirModel            *string       = nil
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
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" service install [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" service start [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" service stop [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" service uninstall [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString("Migration Tool").SetTextColor(stringers.YELLOW)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" migration [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, stringers.NewString(" migration create [flags]").SetTextColor(stringers.GREEN)); err != nil {
			return
		}
	}

	cmdService.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "service [args]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, " Args:"); err != nil {
			return
		}
		if cmdServiceInstall != nil {
			if _, err := fmt.Fprintln(os.Stderr, "  install [flag]"); err != nil {
				return
			}
		}
		if cmdServiceStop != nil {
			if _, err := fmt.Fprintln(os.Stderr, "  stop [flag]"); err != nil {
				return
			}
		}
		if cmdServiceUninstall != nil {
			if _, err := fmt.Fprintln(os.Stderr, "  uninstall [flag]"); err != nil {
				return
			}
		}
	}

	cmdServiceInstall.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "service install [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Flags:"); err != nil {
			return
		}
		cmdServiceInstall.PrintDefaults()
	}

	cmdServiceStop.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "service start [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Flags:"); err != nil {
			return
		}
		cmdServiceStart.PrintDefaults()
	}

	cmdServiceStop.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "service stop [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Flags:"); err != nil {
			return
		}
		cmdServiceStop.PrintDefaults()
	}

	cmdServiceUninstall.Usage = func() {
		if _, err := fmt.Fprintln(os.Stderr, "service uninstall [flags]"); err != nil {
			return
		}
		if _, err := fmt.Fprintln(os.Stderr, "Flags:"); err != nil {
			return
		}
		cmdServiceUninstall.PrintDefaults()
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
	cmdService = flag.NewFlagSet("service", flag.ExitOnError)
	cmdServiceInstall = flag.NewFlagSet("install", flag.ExitOnError)
	cmdServiceInstallSvcName = cmdServiceInstall.String("svcname", "", "service name")
	cmdServiceInstallSvcDisplay = cmdServiceInstall.String("svcdisplay", "", "display name")
	cmdServiceInstallSvcDescription = cmdServiceInstall.String("svcdescription", "", "description")
	cmdServiceInstallSvcApi = cmdServiceInstall.Bool("api", true, "with separation api")
	cmdServiceStart = flag.NewFlagSet("start", flag.ExitOnError)
	cmdServiceStartSvcname = cmdServiceStart.String("svcname", "", "service name")
	cmdServiceStop = flag.NewFlagSet("stop", flag.ExitOnError)
	cmdServiceStopSvcname = cmdServiceStop.String("svcname", "", "service name")
	cmdServiceUninstall = flag.NewFlagSet("uninstall", flag.ExitOnError)
	cmdServiceUninstallSvcname = cmdServiceUninstall.String("svcname", "", "service name")
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
