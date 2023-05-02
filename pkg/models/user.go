package models

import "github.com/praadit/dikurium-test/pkg/utils"

type User struct {
	ID       string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email    string `json:"email"`
	Password string `json:"password"`
	BaseTimeModel
}

func (User) TableName() string {
	return "users"
}

func (en User) ToUpdatable() (fields map[string]interface{}, err error) {
	rawMap, err := utils.ToMap(en)
	if err != nil {
		return nil, err
	}
	//remove relation
	delete(rawMap, "Todo")

	fields = utils.KeyToSnakeCase(rawMap)
	// remove un-updatable field
	delete(fields, "id")
	delete(fields, "created_at")

	return fields, nil
}
