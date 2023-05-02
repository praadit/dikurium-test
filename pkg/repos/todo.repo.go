package repos

import (
	"github.com/google/uuid"
	"github.com/praadit/dikurium-test/pkg/models"
	"gorm.io/gorm"
)

type TodoRepo struct {
	db       *gorm.DB
	baseRepo BaseRepository[models.Todo]
}

func NewTodoRepo(db *gorm.DB) TodoRepo {
	en := models.Todo{}
	baseRepo := BaseRepository[models.Todo]{
		model:     en,
		tableName: "todos",
	}
	return TodoRepo{
		db:       db,
		baseRepo: baseRepo,
	}
}

func (repo *TodoRepo) GetAll() (list []*models.Todo, err error) {
	tables := repo.baseRepo.GetAll(repo.db)
	err = tables.Find(&list).Error
	return
}

func (repo *TodoRepo) GetById(id uuid.UUID) (data *models.Todo, err error) {
	data = &models.Todo{}
	err = repo.baseRepo.GetById(repo.db, id, data)
	return
}

func (repo *TodoRepo) Create(db *gorm.DB, newEntity *models.Todo) (err error) {
	err = repo.baseRepo.Create(db, newEntity)
	return
}

func (repo *TodoRepo) Update(db *gorm.DB, updatedEntity models.Todo) (err error) {
	updatable, err := updatedEntity.ToUpdatable()
	if err != nil {
		return err
	}
	err = repo.baseRepo.Update(db, updatedEntity.ID, updatable)

	return
}

func (repo *TodoRepo) Delete(db *gorm.DB, id uuid.UUID) (err error) {
	err = repo.baseRepo.Delete(db, id)
	return
}

func (repo *TodoRepo) GetByUser(db *gorm.DB, userId uuid.UUID) (todos []*models.Todo, err error) {
	err = repo.baseRepo.GetAll(repo.db).Where("user_id", userId).Find(&todos).Error
	return
}
