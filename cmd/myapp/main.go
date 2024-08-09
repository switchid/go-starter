package main

import (
	"GoStarter/pkg/config"
	"GoStarter/pkg/utils/loggers"
	"GoStarter/web"
	"flag"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"golang.org/x/sys/windows/svc"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type fiberService struct {
	app    *fiber.App
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
				if err := s.app.Shutdown(); err != nil {
					s.logger.LogError("Error shutting down Fiber app: %v", err)
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
	engine := django.NewPathForwardingFileSystem(http.FS(web.TemplatesFS), "/templates", ".django")
	s.app = fiber.New(fiber.Config{
		Views: engine,
	})

	s.app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("layouts/main", fiber.Map{
			"Title": "Hello World",
		}, "index")
	})

	s.logger.LogInfo("Starting Fiber app on port " + cfg.GetServerPort())
	sPort := ":" + cfg.GetServerPort()
	if err = s.app.Listen(sPort); err != nil {
		s.logger.LogError("Error starting Fiber app: %v", err)
	}
}

func main() {
	var svcname = flag.String("svcname", "", "service name")
	var appRun = false

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

		fmt.Println("Shutting down application...")
		service.logger.LogInfo("Shutting down application...")

		shutdownTimeout := 1 * time.Second
		if errShutdown := service.app.ShutdownWithTimeout(shutdownTimeout); errShutdown != nil {
			service.logger.LogError("Error shutting down application: %v\n", errShutdown)
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
			appRun = true
		} else {
			service.logger.LogInfo("Starting as a console application")
			service.runFiberApp()
			appRun = true
		}
	} else {
		flag.Usage()
		appRun = false
	}

	if appRun {
		fmt.Println("Application gracefully shutdown.")
		service.logger.LogInfo("Application gracefully shutdown.")
	}
}
