package main

import (
	"GoStarter/internal/pkg/config"
	"GoStarter/pkg/utils/loggers"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/sys/windows/svc"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type fiberService struct {
	api    *fiber.App
	logger loggers.LoggerPrint
}

func (s *fiberService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown

	changes <- svc.Status{State: svc.StartPending}

	// Start your Fiber app in a separate goroutine
	go s.runFiberApp()

	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}

loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.logger.LogInfo("Received stop command")
				changes <- svc.Status{State: svc.StopPending}
				if err := s.api.Shutdown(); err != nil {
					s.logger.LogError("Error shutting down Fiber app for api: %v", err)
				}
				break loop
			default:
				s.logger.LogError("Unexpected control request #%d", c)
			}
		}
	}

	changes <- svc.Status{State: svc.Stopped}
	return
}

func (s *fiberService) runFiberApp() {
	cfg, err := config.Load()
	if err != nil {
		s.logger.LogError("Error loading config: %v", err)
		os.Exit(1)
	}
	s.api = fiber.New(fiber.Config{})

	s.api.Get("/", func(c *fiber.Ctx) error {
		type User struct {
			Name string
		}

		dataUser := User{
			Name: "test",
		}

		return c.JSON(fiber.Map{
			"status": 200,
			"data":   dataUser,
		})
	})

	s.logger.LogInfo("Starting Fiber app for api on port " + cfg.GetServerApiPort())
	sPort := ":" + cfg.GetServerApiPort()
	if err = s.api.Listen(sPort); err != nil {
		s.logger.LogError("Error starting Fiber app for api: %v", err)
	}
}

func main() {
	var svcname = flag.String("svcname", "", "service name")
	var apiRun = false

	flag.Parse()

	logger, errLog := loggers.NewLogger()
	if errLog != nil {
		fmt.Printf("Failed to set up loggers: %v\n", errLog)
		os.Exit(1)
	}

	service := &fiberService{logger: logger}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	go func() {
		<-quit

		fmt.Println("Shutting down application api...")
		service.logger.LogInfo("Shutting down application api...")

		shutdownTimeout := 1 * time.Second
		if errShutdown := service.api.ShutdownWithTimeout(shutdownTimeout); errShutdown != nil {
			service.logger.LogError("Error shutting down application api: %v\n", errShutdown)
		}
	}()

	if *svcname != "" {
		isService, errIsService := svc.IsWindowsService()
		if errIsService != nil {
			service.logger.LogError("Failed to determine if we are running as a service: %v\n", errIsService)
			os.Exit(1)
		}

		if isService {
			svcName := *svcname
			service.logger.LogInfo("Starting as a Windows service")
			errRun := svc.Run(svcName, service)
			if errRun != nil {
				logger.Printf("Service failed: %v\n", errRun)
				os.Exit(1)
			}
			apiRun = true
		} else {
			service.logger.LogInfo("Starting as a console application api")
			service.runFiberApp()
			apiRun = true
		}
	} else {
		flag.Usage()
		apiRun = false
	}

	if apiRun {
		fmt.Println("Application Api gracefully shutdown.")
		service.logger.LogInfo("Application Api gracefully shutdown.")
	}

}
