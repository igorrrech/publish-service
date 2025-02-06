package service

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igorrrech/publish-service/authorization/repo"
	"github.com/igorrrech/publish-service/authorization/service/handlers/health"
	"github.com/igorrrech/publish-service/authorization/service/handlers/login"
	"github.com/igorrrech/publish-service/authorization/service/handlers/refresh"
)

// SSO
type AuthService struct {
	addr       string
	ur         repo.UserRepository  //for login
	tr         repo.TokenRepository //for login and refresh
	accessTTL  time.Duration
	refreshTTL time.Duration
	logger     *slog.Logger
}

func NewAuthService(
	addr string,
	ur repo.UserRepository,
	tr repo.TokenRepository,
	accessTTL time.Duration,
	refreshTTL time.Duration,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		addr:       addr,
		ur:         ur,
		tr:         tr,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
		logger:     logger,
	}
}
func (s *AuthService) Run(
	parentCtx context.Context,
	shutdownTimeout time.Duration,
) {
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	r := gin.Default()
	r.GET("/health", health.Health())
	r.POST("/login", login.Login(s.accessTTL, s.refreshTTL, s.ur, s.tr))
	r.POST("/refresh", refresh.Refresh(s.accessTTL, s.refreshTTL, s.tr))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("server is not started", "error", err.Error())
		}
	}()
	//wait for context done or get signal
	s.logger.Info("server started at", "addres", s.addr)
	<-ctx.Done()
	//clear context
	stop()
	//grace server cancelleration context
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.logger.Error("server forced to stop", "error", err.Error())
	}
	s.logger.Info("server stoped")
}
