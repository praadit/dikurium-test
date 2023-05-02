package repos

import (
	"errors"

	"gorm.io/gorm"
)

type BaseRepository[T any] struct {
	model     T
	tableName string
}

func (repo *BaseRepository[T]) GetAll(db *gorm.DB) *gorm.DB {
	return db.Model(&repo.model)
}

func (repo *BaseRepository[T]) GetById(db *gorm.DB, id any, data any) (err error) {
	return db.Model(&repo.model).Where("id", id).First(data).Error
}

func (repo *BaseRepository[T]) Create(db *gorm.DB, newData any) error {
	return db.Create(newData).Error
}

func (repo *BaseRepository[T]) CreateMultiple(db *gorm.DB, newData []T) error {
	return db.CreateInBatches(newData, len(newData)).Error
}

func (repo *BaseRepository[T]) Update(db *gorm.DB, id any, changes map[string]any) error {
	return db.Model(&repo.model).Where("id", id).Updates(changes).Error
}

func (repo *BaseRepository[T]) Delete(db *gorm.DB, id any) error {
	data := repo.model
	err := repo.GetById(db, id, &data)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}
	return db.Delete(&data).Error
}
