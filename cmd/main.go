package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/emilien-puget/invoice_microservice/configuration"
	"github.com/emilien-puget/invoice_microservice/invoice"
	"github.com/emilien-puget/invoice_microservice/user"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var ErrStopSignalReceived = errors.New("stop signal received")

// Version set via ldflags
var Version = "local"

const service = "invoice_microservice"

func main() {
	eCfg := configuration.Api{}
	if err := env.Parse(&eCfg, (env.Options{RequiredIfNoDef: true})); err != nil {
		log.Printf("%+v\n", err)
		os.Exit(-1)
	}

	ctx, cl := Init()
	defer cl(nil)

	validate := validator.New()

	e := echo.New()
	defer e.Shutdown(context.Background())

	db, err := initDb(&eCfg.Postgres)
	if err != nil {
		cl(fmt.Errorf("init db:%w", err))
		return
	}
	defer db.Close()

	userRepository := user.NewUserRepository(db)
	invoiceRepository := invoice.NewInvoiceRepository(db)
	usersHandler := user.NewGetAllHandler(userRepository)
	transactionHandler := invoice.NewDoTransactionHandler(invoiceRepository, userRepository, validate)
	invoiceHandler := invoice.NewCreateInvoiceHandler(validate, invoiceRepository, userRepository)
	e.Use(middleware.Logger())
	e.Use(echoprometheus.NewMiddleware(service))
	e.GET("/metrics", echoprometheus.NewHandler())
	e.GET("/users", usersHandler.Handle)
	e.POST("/invoice", invoiceHandler.Handle)
	e.POST("/transaction", transactionHandler.Handle)

	go func() {
		err := e.Start(fmt.Sprintf(":%s", eCfg.Port))
		if err != nil {
			cl(fmt.Errorf("exposed server:%w", err))
		}
	}()

	srv := initInternalSrv(eCfg.InternalPort)
	defer srv.Shutdown(context.Background())
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			cl(fmt.Errorf("internal server:%w", err))
		}
	}()
	log.Printf("starting %s", Version)
	<-ctx.Done()
}

func initDb(c *configuration.Postgres) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Database, c.Sslmode))
	if err != nil {
		return nil, fmt.Errorf("sql.open: %w", err)
	}

	timeout, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	err = db.PingContext(timeout)
	if err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}

func initInternalSrv(internalPort string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	})
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", internalPort),
		Handler: mux,
	}
	return srv
}

func Init() (ctx context.Context, cl func(err error)) {
	ctx = context.Background()
	stopSignal := notifyStopSignal()
	ctx, cancelFunc := context.WithCancelCause(ctx)
	go func() {
		<-stopSignal
		cancelFunc(ErrStopSignalReceived)
	}()

	var once sync.Once
	return ctx, func(err error) {
		once.Do(func() {
			cancelFunc(err)
			loggingExit(ctx)
		})
	}
}

func notifyStopSignal() <-chan os.Signal {
	gracefulStop := make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	return gracefulStop
}

func loggingExit(ctx context.Context) {
	err := context.Cause(ctx)
	if err != nil {
		if errors.Is(err, ErrStopSignalReceived) {
			log.Print("stop signal received")
			return
		}
		if errors.Is(err, context.Canceled) {
			log.Print("context cancel without cause")
			return
		}
		log.Print(err)
		return
	}
}
