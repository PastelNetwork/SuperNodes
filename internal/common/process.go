package common

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Logger struct {
	ErrorLog   *log.Logger
	WarningLog *log.Logger
	InfoLog    *log.Logger
}

type Application struct {
	name string
	Cfg  Config
	Log  Logger
}

func NewApplication(name string, configFile string, logFile string) *Application {
	a := new(Application)

	a.name = name

	if err := a.Cfg.LoadConfig(configFile); err != nil {
		panic(err)
	}

	file2log, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	a.Log = Logger{
		ErrorLog:   log.New(file2log, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLog: log.New(file2log, "WARN\t", log.Ldate|log.Ltime),
		InfoLog:    log.New(file2log, "INFO\t", log.Ldate|log.Ltime),
	}

	return a
}

func (a *Application) Run(servers []func(ctx context.Context, a *Application) func() error) {

	a.Log.InfoLog.Printf("")
	a.Log.InfoLog.Printf("=======================================")
	a.Log.InfoLog.Printf("====== %s starting ======", a.name)
	a.Log.InfoLog.Printf("=======================================")

	if servers == nil || len(servers) == 0 {
		panic("runners array can't be empty")
	}

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-signals
		a.Log.InfoLog.Println("program interrupted")
		cancel()
		a.Log.InfoLog.Println("cancel context sent")
	}()

	eg, ctx := errgroup.WithContext(ctx)
	for _, f := range servers {
		eg.Go(f(ctx, a))
	}

	if err := eg.Wait(); err != nil {
		a.Log.ErrorLog.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	a.Log.InfoLog.Println("everything closed successfully")
	a.Log.InfoLog.Println("exiting")
}

func (a *Application) CreateServer(ctx context.Context, serverName string,
	startServer func(ctx context.Context) error,
	runServer func(ctx context.Context) error,
	stopServer func(ctx context.Context) error) func() error {

	return func() error {

		if err := startServer(ctx); err != nil {
			return fmt.Errorf("error starting the %s server: %w", serverName, err)
		}

		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			a.Log.InfoLog.Printf("stopping server %s\n", serverName)
			if err := stopServer(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the %s server: %w", serverName, err)
			}
			a.Log.InfoLog.Printf("the %s server is stopped\n", serverName)
			close(errChan)
		}()

		a.Log.InfoLog.Printf("the %s server is starting\n", serverName)

		if err := runServer(ctx); err != nil {
			return fmt.Errorf("error running the %s server: %w", serverName, err)
		}
		a.Log.InfoLog.Printf("the %s server is closing\n", serverName)

		err := <-errChan
		a.Log.InfoLog.Printf("the %s server is closed\n", serverName)
		return err
	}
}
