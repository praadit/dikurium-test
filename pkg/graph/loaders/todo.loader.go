package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/praadit/dikurium-test/pkg/models"
)

func (l *Loader) batchTodoByUser(ctx context.Context, keys []string) []*dataloader.Result[[]*models.Todo] {
	results := make([]*dataloader.Result[[]*models.Todo], len(keys))

	for i, r := range results {
		if r == nil {
			results[i] = &dataloader.Result[[]*models.Todo]{}
		}
	}

	todos := []*models.Todo{}
	if err := l.db.Model(&models.Todo{}).Where("user_id in ?", keys).Find(&todos).Error; err != nil {
		fmt.Println("batchTodoByUser error:", err)
		return results
	}

	for i, ID := range keys {
		for _, todo := range todos {
			if todo.UserID == ID {
				results[i].Data = append(results[i].Data, todo)
				break
			}
		}
	}

	return results
}
