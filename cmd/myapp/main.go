package main

import (
	"GoStarter/internal/middleware"
	"GoStarter/internal/pkg/config"
	"GoStarter/internal/routes"
	"GoStarter/pkg/utils/loggers"
	"GoStarter/web"
	"flag"
	"fmt"
	"github.com/flosch/pongo2/v6"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/django/v3"
	"golang.org/x/sys/windows/svc"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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

func loadMacros(engine *django.Engine, macrosDir string) error {
	return filepath.Walk(macrosDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".django") {
			macroName := strings.TrimSuffix(filepath.Base(path), ".django")
			macroContent, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			// Create a template from the macro content
			var tmpl = pongo2.Must(pongo2.FromString(string(macroContent)))
			//println(tmpl)
			// Create a function that executes the template
			engine.AddFunc(macroName, func(str string) (string, error) {
				return tmpl.Execute(pongo2.Context{"text": str})
			})

		}
		return nil
	})
}

func (s *fiberService) runFiberApp() {
	cfg, err := config.Load()
	if err != nil {
		s.logger.LogError("Error loading config: %v", err)
		os.Exit(1)
	}

	engine := django.NewPathForwardingFileSystem(http.FS(web.TemplatesFS), "/templates", ".django")
	//currentDir, err := os.Getwd()
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//templateDir := filepath.Join(currentDir, "web/templates/components")
	//errMacros := loadMacros(engine, templateDir)
	//if errMacros != nil {
	//	return
	//}

	s.app = fiber.New(fiber.Config{
		Views: engine,
	})

	s.app.Use(middleware.AppMiddleware)

	r := routes.Load(s.app)
	r.AppRoutes()
	r.AuthRoutes()
	r.NotFoundRoutes()

	s.app.Use(func(ctx *fiber.Ctx) error {
		return ctx.RedirectToRoute("not-found", fiber.Map{"name": "", "queries": map[string]string{"urls": strings.TrimLeft(ctx.OriginalURL(), "/")}})
	})

	s.logger.LogInfo("Starting Fiber app on port " + cfg.GetServerAppPort())
	sPort := ":" + cfg.GetServerAppPort()
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
