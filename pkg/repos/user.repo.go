package repos

import (
	"github.com/google/uuid"
	"github.com/praadit/dikurium-test/pkg/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db       *gorm.DB
	baseRepo BaseRepository[models.User]
}

func NewUserRepo(db *gorm.DB) UserRepo {
	en := models.User{}
	baseRepo := BaseRepository[models.User]{
		model:     en,
		tableName: "users",
	}
	return UserRepo{
		db:       db,
		baseRepo: baseRepo,
	}
}

func (repo *UserRepo) GetAll() (list []*models.User, err error) {
	tables := repo.baseRepo.GetAll(repo.db)
	err = tables.Find(&list).Error
	return
}

func (repo *UserRepo) GetById(id uuid.UUID) (data *models.User, err error) {
	data = &models.User{}
	err = repo.baseRepo.GetById(repo.db, id, data)
	return
}

func (repo *UserRepo) Create(db *gorm.DB, newEntity *models.User) (err error) {
	err = repo.baseRepo.Create(db, newEntity)
	return
}

func (repo *UserRepo) Update(db *gorm.DB, updatedEntity models.User) (err error) {
	updatable, err := updatedEntity.ToUpdatable()
	if err != nil {
		return err
	}
	err = repo.baseRepo.Update(db, updatedEntity.ID, updatable)

	return
}

func (repo *UserRepo) Delete(db *gorm.DB, id uuid.UUID) (err error) {
	err = repo.baseRepo.Delete(db, id)
	return
}

func (repo *UserRepo) GetByEmail(email string) (user *models.User, err error) {
	err = repo.baseRepo.GetAll(repo.db).Where("email", email).First(&user).Error
	return
}
