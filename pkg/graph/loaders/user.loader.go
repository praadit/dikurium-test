package loader

import (
	"context"
	"fmt"

	"github.com/graph-gophers/dataloader/v7"
	"github.com/praadit/dikurium-test/pkg/models"
)

func (l *Loader) batchUserByTodo(ctx context.Context, keys []string) []*dataloader.Result[*models.User] {
	results := make([]*dataloader.Result[*models.User], len(keys))

	for i, r := range results {
		if r == nil {
			results[i] = &dataloader.Result[*models.User]{}
		}
	}

	users := []*models.User{}
	if err := l.db.Model(&models.User{}).Where("id in ?", keys).Find(&users).Error; err != nil {
		fmt.Println("batchUserByTodo error:", err)
		return results
	}

	for i, ID := range keys {
		for _, user := range users {
			if user.ID == ID {
				results[i].Data = user
				break
			}
		}
	}

	return results
}
