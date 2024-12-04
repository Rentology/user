package server

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
	"user-service/internal/broker"
	"user-service/internal/config"
	"user-service/internal/user/repository"
	service2 "user-service/internal/user/service"
	"user-service/lib/sl"
)

const (
	ctxTimeout = 5
	PprofPort  = "5555"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
	db   *sqlx.DB
	log  *slog.Logger
}

func NewServer(cfg *config.Config, db *sqlx.DB, log *slog.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, db: db, log: log}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(s.cfg.Server.Port),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 5,
	}

	go func() {
		sl.Infof(s.log, "Server is listening on PORT: %d", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.log.Error("Error starting Server: ", err)
			os.Exit(1)
		}
	}()

	go func() {
		sl.Infof(s.log, "Starting Debug Server on PORT: %s", PprofPort)
		if err := http.ListenAndServe(":"+PprofPort, http.DefaultServeMux); err != nil {
			s.log.Error("Error PPROF ListenAndServe: %s", err)
			os.Exit(1)
		}
	}()

	go func() {
		if err := s.runConsumer(); err != nil {
			s.log.Error("Error running Consumer: %s", err)
			os.Exit(1)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	s.log.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)

}

func (s *Server) runConsumer() error {
	brokerConn, err := broker.NewBroker(s.cfg.RabbitMQ.Url)
	if err != nil {
		s.log.Error("Ошибка подключения к RabbitMQ", slog.String("error", err.Error()))
		return err
	}
	// defer brokerConn.Close() перемещен в shutdown секцию
	s.log.Info("RabbitMQ connection established.")
	repo := repository.NewUserRepository(s.db)
	service := service2.NewUserService(s.cfg, repo, s.log)
	userConsumer, err := broker.NewUserConsumer(service, brokerConn, s.cfg.RabbitMQ.QueueName, s.log)
	if err != nil {
		s.log.Error("Ошибка создания consumer", slog.String("error", err.Error()))
		return err
	}

	ctx := context.Background()
	go func() {
		if err := userConsumer.Run(ctx); err != nil {
			s.log.Error("Ошибка запуска consumer", slog.String("error", err.Error()))
		}
	}()

	return nil
}
