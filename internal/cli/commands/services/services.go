package services

import (
	"GoStarter/pkg/utils/loggers"
	"fmt"
	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/eventlog"
	"golang.org/x/sys/windows/svc/mgr"
	"log"
	"os"
)

type Services struct {
	Name        string
	Display     string
	Description string
}

func (sc *Services) InstallService(exePath, appFlags string) error {
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

	s, err := m.OpenService(sc.Name)
	if err == nil {
		errClose := s.Close()
		if errClose != nil {
			return errClose
		}
		return fmt.Errorf("service %s already exists", sc.Name)
	}

	s, err = m.CreateService(sc.Name, exePath, mgr.Config{
		DisplayName:      sc.Display,
		Description:      sc.Description,
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

	err = eventlog.Remove(sc.Name)
	if err != nil && !os.IsNotExist(err) {
		log.Printf("failed to remove existing event log source: %v", err)
		// Continue with installation even if this fails
	}

	err = eventlog.InstallAsEventCreate(sc.Name, eventlog.Error|eventlog.Warning|eventlog.Info)
	if err != nil {
		errDel := s.Delete()
		if errDel != nil {
			return errDel
		}
		return fmt.Errorf("SetupEventLogSource() failed: %s", err)
	}

	errRun := sc.RunService()
	if errRun != nil {
		return errRun
	}

	return nil
}

func (sc *Services) RunService() error {
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
	s, err := m.OpenService(sc.Name)
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

func (sc *Services) StopService() error {
	m, err := mgr.Connect()
	if err != nil {
		return fmt.Errorf("failed to connect to service manager: %v", err)
	}
	defer func(m *mgr.Mgr) {
		errDC := m.Disconnect()
		if errDC != nil {

		}
	}(m)

	s, err := m.OpenService(sc.Name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", sc.Name)
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

func (sc *Services) UninstallService() error {
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

	s, err := m.OpenService(sc.Name)
	if err != nil {
		return fmt.Errorf("service %s is not installed", sc.Name)
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

	err = eventlog.Remove(sc.Name)
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
