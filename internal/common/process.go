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
	ErrorLog *log.Logger
	WarningLog  *log.Logger
	InfoLog  *log.Logger
}

func Run(application string,
		configFile string, logFile string,
		servers []func(ctx context.Context, config *Config, logger *Logger, wg *sync.WaitGroup) func() error) {

	var config Config
	if err := config.LoadConfig(configFile); err != nil {
		panic(err)
	}

	file2log, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	logger := &Logger{
		ErrorLog: log.New(file2log, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		WarningLog: log.New(file2log, "WARN\t", log.Ldate|log.Ltime),
		InfoLog: log.New(file2log, "INFO\t", log.Ldate|log.Ltime),
	}

	logger.InfoLog.Printf("")
	logger.InfoLog.Printf("=======================================")
	logger.InfoLog.Printf("====== %s starting ======", application)
	logger.InfoLog.Printf("=======================================")

	if servers == nil || len(servers) == 0 {
		panic("runners array can't be empty")
	}

	var wg sync.WaitGroup
	wg.Add(len(servers))
	ctx, cancel := context.WithCancel(context.Background())
	eg, egCtx := errgroup.WithContext(context.Background())

	for _, f := range servers {
		eg.Go(f(ctx, &config, logger, &wg))
	}

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-signals
		logger.InfoLog.Println("program interrupted")
		cancel()
		logger.InfoLog.Println("cancel context sent")
	}()

	if err := eg.Wait(); err != nil {
		logger.ErrorLog.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	logger.InfoLog.Println("everything closed successfully")
	logger.InfoLog.Println("exiting")
}

func CreateServer(name string, ctx context.Context, config *Config, logger *Logger, wg *sync.WaitGroup,
	startServer func(ctx context.Context) error,
	runServer func(ctx context.Context) error,
	stopServer func(ctx context.Context) error) func() error {

	return func() error {

		if err := startServer(ctx); err != nil {
			return fmt.Errorf("error starting the %s server: %w", name, err)
		}

		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			if err := stopServer(shutCtx); err != nil {
				errChan <- fmt.Errorf("error shutting down the %s server: %w", name, err)
			}

			logger.ErrorLog.Printf("the %s server is closed\n", name)
			close(errChan)
			wg.Done()
		}()

		logger.InfoLog.Printf("the %s server is starting\n", name)

		if err := runServer(ctx); err != nil {
			return fmt.Errorf("error running the %s server: %w", name, err)
		}

		logger.InfoLog.Printf("the %s server is closing\n", name)
		err := <-errChan
		wg.Wait()
		return err
	}
}
