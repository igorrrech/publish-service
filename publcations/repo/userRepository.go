package repo

import (
	"context"
	"log/slog"
	"time"

	dbgorm "github.com/igorrrech/publish-service/authorization/pkg/dbGorm"
	models "github.com/igorrrech/publish-service/authorization/service/models"
	"gorm.io/gorm"
)

type UserRepository struct {
	db                *gorm.DB
	logger            *slog.Logger
	operationsTimeout time.Duration
	ctx               context.Context
}

func NewUserRepository(
	dsn string,
	logger *slog.Logger,
) *UserRepository {
	ctx := context.Background()
	db, err := dbgorm.ConnectToDbPg(ctx, dsn, 5*time.Second)
	if err != nil {
		logger.Error("failed to connect to user db", "error", err.Error())
	}

	db.AutoMigrate(&models.User{})

	db.Create(&models.User{
		Phone:    "1",
		Password: "1",
		Role:     "admin",
	})

	return &UserRepository{
		db:                db,
		logger:            logger,
		operationsTimeout: 2 * time.Second,
		ctx:               ctx,
	}
}
func (r UserRepository) GetUserByPhone(phone string) (*models.User, error) {
	var result models.User
	ctx, cancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer cancel()
	err := r.db.Model(models.User{Phone: phone}).WithContext(ctx).First(&result).Error
	return &result, err
}
func (r UserRepository) CreateUser(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(r.ctx, r.operationsTimeout)
	defer cancel()
	return r.db.WithContext(ctx).Create(u).Error
}
