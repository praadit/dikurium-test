package controller

import (
	"github.com/praadit/dikurium-test/pkg/repos"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Controller struct {
	db       *gorm.DB
	logger   *zap.Logger
	userRepo repos.UserRepo
	todoRepo repos.TodoRepo
}

func InitController(db *gorm.DB, logger *zap.Logger) *Controller {
	return &Controller{
		db:       db,
		logger:   logger,
		userRepo: repos.NewUserRepo(db),
		todoRepo: repos.NewTodoRepo(db),
	}
}
