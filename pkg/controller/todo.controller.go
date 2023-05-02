package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/praadit/dikurium-test/pkg/models"
	"go.uber.org/zap"
)

func (c *Controller) Todos(ctx context.Context) ([]*models.Todo, error) {
	user, err := c.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	uuidUser := uuid.MustParse(user.ID)
	todos, err := c.todoRepo.GetByUser(c.db, uuidUser)
	if err != nil {
		c.logger.Info("Failed to get todo with details",
			zap.String("user", fmt.Sprintf("userId : %s", user.ID)),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
	}

	return todos, err
}

func (c *Controller) CreateTodo(ctx context.Context, input *models.CreateTodoInput) (*models.Todo, error) {
	user, err := c.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	newTodo := &models.Todo{
		UserID:      user.ID,
		Title:       input.Title,
		IsCompleted: false,
	}

	err = c.todoRepo.Create(c.db, newTodo)
	if err != nil {
		c.logger.Info("Failed to create new todo with details",
			zap.String("user", fmt.Sprintf("user : %s", user.Email)),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, errors.New("Failed to create new todo")
	}
	return newTodo, nil
}
func (c *Controller) MarkCompleted(ctx context.Context, todoId uuid.UUID) (*models.Todo, error) {
	user, err := c.GetAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}

	todo, err := c.todoRepo.GetById(todoId)
	if err != nil {
		c.logger.Info("Failed to get todo with details",
			zap.String("todo", fmt.Sprintf("todo id : %s", todoId.String())),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, errors.New("Failed to get todo")
	}

	if todo.UserID != user.ID {
		c.logger.Info("Failed to get todo with details",
			zap.String("todo", fmt.Sprintf("todo id unavailable : %s", todoId.String())),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, errors.New("Todo are unavailable")
	}

	todo.IsCompleted = !todo.IsCompleted

	err = c.todoRepo.Update(c.db, *todo)
	if err != nil {
		c.logger.Info("Failed to update todo with details",
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return nil, errors.New("Failed to update todo")
	}

	return todo, nil
}
func (c *Controller) DeleteTodo(ctx context.Context, todoId uuid.UUID) (bool, error) {
	user, err := c.GetAuthenticatedUser(ctx)
	if err != nil {
		return false, err
	}

	todo, err := c.todoRepo.GetById(todoId)
	if err != nil {
		c.logger.Info("Failed to get todo with details",
			zap.String("todo", fmt.Sprintf("todo id : %s", todoId.String())),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return false, errors.New("Failed to get todo")
	}

	if todo.UserID != user.ID {
		c.logger.Info("Failed to get todo with details",
			zap.String("todo", fmt.Sprintf("todo id unavailable : %s", todoId.String())),
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return false, errors.New("Todo are unavailable")
	}

	err = c.todoRepo.Delete(c.db, todoId)
	if err != nil {
		c.logger.Info("Failed to delete todo with details",
			zap.String("message", fmt.Sprintf("err : %s", err.Error())),
		)
		return false, errors.New("Failed to delete todo")
	}

	return true, nil
}
