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
	return &serverKeeper{Cfg: Config{ServerAddress: ":8080"}}
}

func (s serverKeeper) StartServer() error {
	if err := s.parseFlagCfg(); err != nil {
		return err
	}
	if err := s.connectDB(); err != nil {
		return err
	}

	e := echo.New()

	s.initRouter(e)

	err := e.Start(s.Cfg.ServerAddress)
	if err != nil {
		return err
	}
	return nil
}

func (s *serverKeeper) initRouter(e *echo.Echo) {
	e.POST("/api/user/register", s.postAPIUserRegister)
	e.POST("/api/user/login", s.postAPIUserLogin)

	e.POST("/write/card", s.postWrite, s.mwAuthentication)
	e.POST("/write/password", s.postWrite, s.mwAuthentication)
	e.POST("/write/text", s.postWrite, s.mwAuthentication)
	e.POST("/write/bin", s.postWrite, s.mwAuthentication)

	e.GET("/read/card", s.getReadALL, s.mwAuthentication)
	e.GET("/read/password", s.getReadALL, s.mwAuthentication)
	e.GET("/read/text", s.getReadALL, s.mwAuthentication)
	e.GET("/read/bin", s.getReadALL, s.mwAuthentication)
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
