package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID              string `gorm:"type:text;primarykey"`
	Username        string `gorm:"type:text"`
	Email           string `gorm:"type:text"`
	Password        string `gorm:"type:text"`
	IsEmailVerified *bool  `gorm:"type:bool;default:false"`
	Role            string `gorm:"type:text"`
	Signature       string `gorm:"type:text"` // https://medium.com/swlh/building-a-user-auth-system-with-jwt-using-golang-30892659cc0#06a3

	CreatedAt time.Time      `gorm:"type:timestamp"`
	UpdatedAt time.Time      `gorm:"type:timestamp"`
	DeletedAt gorm.DeletedAt `gorm:"type:timestamp;index"`

	ManagerID *string
	Manager   *User
	Tasks     []Task
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	if u.Signature == "" {
		u.Signature = uuid.New().String()
	}

	return nil
}
