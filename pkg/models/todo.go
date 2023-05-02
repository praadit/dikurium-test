package models

import "github.com/praadit/dikurium-test/pkg/utils"

type Todo struct {
	ID          string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	UserID      string `json:"user_id" gorm:"type:uuid;"`
	Title       string `json:"title"`
	IsCompleted bool   `json:"is_completed"`
	BaseTimeModel
}

func (Todo) TableName() string {
	return "todos"
}

func (en Todo) ToUpdatable() (fields map[string]interface{}, err error) {
	rawMap, err := utils.ToMap(en)
	if err != nil {
		return nil, err
	}
	//remove relation
	delete(rawMap, "user")

	fields = utils.KeyToSnakeCase(rawMap)
	// remove un-updatable field
	delete(fields, "id")
	delete(fields, "user")
	delete(fields, "user_id")
	delete(fields, "created_at")

	return fields, nil
}
