package common

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Logger struct {
	ErrorLog   *log.Logger
	WarningLog *log.Logger
	InfoLog    *log.Logger
}

type Application struct {
	Name   string
	ctx    context.Context
	config Config
	logger Logger
	wg     sync.WaitGroup
}

func (a *Application) Run(configFile string, logFile string,
	servers []func(a *Application) func() error) {

	if err := a.config.LoadConfig(configFile); err != nil {
		panic(err)
	}

	file2log, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	a.logger = Logger{
		ErrorLog:   log.New(file2log, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLog: log.New(file2log, "WARN\t", log.Ldate|log.Ltime),
		InfoLog:    log.New(file2log, "INFO\t", log.Ldate|log.Ltime),
	}

	a.logger.InfoLog.Printf("")
	a.logger.InfoLog.Printf("=======================================")
	a.logger.InfoLog.Printf("====== %s starting ======", a.Name)
	a.logger.InfoLog.Printf("=======================================")

	if servers == nil || len(servers) == 0 {
		panic("runners array can't be empty")
	}

	a.wg.Add(len(servers))
	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(context.Background())
	eg, egCtx := errgroup.WithContext(context.Background())

	for _, f := range servers {
		eg.Go(f(a))
	}

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-signals
		a.logger.InfoLog.Println("program interrupted")
		cancel()
		a.logger.InfoLog.Println("cancel context sent")
	}()

	if err := eg.Wait(); err != nil {
		a.logger.ErrorLog.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	a.logger.InfoLog.Println("everything closed successfully")
	a.logger.InfoLog.Println("exiting")
}

func (a *Application) CreateServer(serverName string,
	startServer func(ctx context.Context) error,
	runServer func(ctx context.Context) error,
	stopServer func(ctx context.Context) error) func() error {

	return func() error {

		if err := startServer(a.ctx); err != nil {
			return fmt.Errorf("error starting the %s server: %w", serverName, err)
		}

		errChan := make(chan error, 1)

		go func() {
			<-a.ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := stopServer(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the %s server: %w", serverName, err)
			}

			a.logger.InfoLog.Printf("the %s server is closed\n", serverName)
			close(errChan)
			a.wg.Done()
		}()

		a.logger.InfoLog.Printf("the %s server is starting\n", serverName)

		if err := runServer(a.ctx); err != nil {
			return fmt.Errorf("error running the %s server: %w", serverName, err)
		}

		a.logger.InfoLog.Printf("the %s server is closing\n", serverName)
		err := <-errChan
		a.wg.Wait()
		return err
	}
}
