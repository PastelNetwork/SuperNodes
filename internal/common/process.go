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
	ctx  context.Context
	//wg   sync.WaitGroup
	Cfg  Config
	Log  Logger
}

func (a *Application) Init(name string, configFile string, logFile string){
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
}

func (a *Application) Run(servers []func(a *Application) func() error) {

	a.Log.InfoLog.Printf("")
	a.Log.InfoLog.Printf("=======================================")
	a.Log.InfoLog.Printf("====== %s starting ======", a.name)
	a.Log.InfoLog.Printf("=======================================")

	if servers == nil || len(servers) == 0 {
		panic("runners array can't be empty")
	}

	//a.wg.Add(len(servers))

	var cancel context.CancelFunc
	a.ctx, cancel = context.WithCancel(context.Background())
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		<-signals
		a.Log.InfoLog.Println("program interrupted")
		cancel()
		a.Log.InfoLog.Println("cancel context sent")
	}()

	eg, egCtx := errgroup.WithContext(a.ctx/*context.Background()*/)
	for _, f := range servers {
		eg.Go(f(a))
	}
	go func() {
		<-egCtx.Done()
		cancel()
	}()

	if err := eg.Wait(); err != nil {
		a.Log.ErrorLog.Printf("error in the server goroutines: %s\n", err)
		os.Exit(1)
	}
	a.Log.InfoLog.Println("everything closed successfully")
	a.Log.InfoLog.Println("exiting")
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

			a.Log.InfoLog.Printf("the %s server is closed\n", serverName)
			close(errChan)
			//a.wg.Done()
		}()

		a.Log.InfoLog.Printf("the %s server is starting\n", serverName)

		if err := runServer(a.ctx); err != nil {
			return fmt.Errorf("error running the %s server: %w", serverName, err)
		}

		a.Log.InfoLog.Printf("the %s server is closing\n", serverName)
		err := <-errChan
		//a.wg.Wait()
		return err
	}
}
