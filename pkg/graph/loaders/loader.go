package loader

import (
	"context"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/praadit/dikurium-test/pkg/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Loader struct {
	db         *gorm.DB
	logger     *zap.Logger
	UserByTodo *dataloader.Loader[string, *models.User]
	TodoByUser *dataloader.Loader[string, []*models.Todo]
}

func NewLoader(db *gorm.DB, logger *zap.Logger) *Loader {
	loader := &Loader{
		db:     db,
		logger: logger,
	}
	loader.UserByTodo = dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[*models.User] {
		return loader.batchUserByTodo(ctx, ids)
	})
	loader.TodoByUser = dataloader.NewBatchedLoader(func(ctx context.Context, ids []string) []*dataloader.Result[[]*models.Todo] {
		return loader.batchTodoByUser(ctx, ids)
	})
	return loader
}
