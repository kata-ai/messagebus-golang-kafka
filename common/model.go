package common

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/satori/uuid"
)

//BaseModel : struct for common entity model
type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:char(36); primary_key"`
	CreatedAt time.Time  `gorm:"column:created_at;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:updated_at;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

//BeforeCreate : Gorm callback to generate ID before create model
func (b *BaseModel) BeforeCreate(scope *gorm.Scope) error {
	if b.ID == uuid.Nil {
		return scope.SetColumn("ID", uuid.V4)
	}

	return nil
}
