package router

import (
	"context"
	"flag"
	"time"

	"github.com/caarlos0/env"
	"github.com/labstack/echo"

	"PasManagerGophKeeper/internal/storage"
)

type Config struct {
	ServerAddress  string `env:"RUN_ADDRESS"`
	BDAddress      string `env:"DATABASE_URI"`
	AccrualAddress string `env:"ACCRUAL_SYSTEM_ADDRESS"`
}

type serverKeeper struct {
	Cfg  Config
	Serv *echo.Echo
	DB   storage.DBI
}

func InitServer() *serverKeeper {
	return &serverKeeper{}
}

func (s serverKeeper) Router() error {
	if err := s.parseFlagCfg(); err != nil {
		return err
	}
	if err := s.connectDB(); err != nil {
		return err
	}

	e := echo.New()

	e.Use(s.mwAuthentication)

	e.POST("/api/user/register", s.postAPIUserRegister)
	e.POST("/api/user/login", s.postAPIUserLogin)

	e.POST("/DB/", s.postWrite)
	e.GET("/DB/", s.getReadALL)

	err := e.Start(s.Cfg.ServerAddress)
	if err != nil {
		return err
	}
	return nil
}
func (s *serverKeeper) connectDB() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	var err error
	if s.DB, err = storage.InitDB(); err != nil {
		return err
	}
	if err = s.DB.Connect(ctx, s.Cfg.BDAddress); err != nil {
		return err
	}
	return nil
}
func (s *serverKeeper) parseFlagCfg() error {
	errConfig := env.Parse(&s.Cfg)
	if errConfig != nil {
		return errConfig
	}
	if s.Cfg.ServerAddress == "" {
		flag.StringVar(&s.Cfg.ServerAddress, "a", "localhost:8080", "New RUN_ADDRESS")
	}
	if s.Cfg.BDAddress == "" {
		flag.StringVar(&s.Cfg.BDAddress, "d", "postgres://postgres:0000@localhost:5432/postgres", "New DATABASE_URI")
	}
	flag.Parse()
	return nil
}
