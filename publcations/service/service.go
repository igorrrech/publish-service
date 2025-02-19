package service

import (
	"context"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/igorrrech/publish-service/publications/repo"
	createpost "github.com/igorrrech/publish-service/publications/service/handlers/createPost"
	deletepost "github.com/igorrrech/publish-service/publications/service/handlers/deletePost"
	"github.com/igorrrech/publish-service/publications/service/handlers/health"
	readpost "github.com/igorrrech/publish-service/publications/service/handlers/readPost"
	updatepost "github.com/igorrrech/publish-service/publications/service/handlers/updatePost"
	"github.com/igorrrech/publish-service/publications/service/middleware"
)

type Service struct {
	addr   string
	log    *slog.Logger
	secret string
	pr     repo.PostRepository
	ur     repo.UserRepositrory
}

func NewService(
	Addr string,
	Logger *slog.Logger,
	secret string,
	postRepository repo.PostRepository,
	userRepository repo.UserRepositrory,
) *Service {
	return &Service{
		addr:   Addr,
		log:    Logger,
		secret: secret,
		pr:     postRepository,
	}
}
func (s *Service) Run(
	parentCtx context.Context,
	sutdownTimeout time.Duration,
) {
	ctx, stop := signal.NotifyContext(parentCtx, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	r := gin.Default()

	r.GET("/health", health.HealthCheck())

	publication := r.Group("publication")
	publication.GET("/", readpost.ReadAll(s.pr, s.ur))
	publication.POST("/", createpost.Create(s.pr, s.ur))
	publication.PATCH("/", updatepost.Update(s.pr))
	publication.DELETE("/", deletepost.Delete(s.pr))
	publication.Use(middleware.Auth(s.secret))

	srv := &http.Server{
		Addr:    s.addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.log.Error("server is not started", "error", err.Error())
		}
	}()
	s.log.Info("server startet at", "addres", s.addr)
	<-ctx.Done()
	stop()
	ctx, cancel := context.WithTimeout(ctx, sutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		s.log.Error("forced to stop", "error", err.Error())
	}
	s.log.Info("server stopped")
}
