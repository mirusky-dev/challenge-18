package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID      string `gorm:"type:text"`
	Summary string `gorm:"type:text"`
	UserID  string `gorm:"type:text"`

	PerformedAt *time.Time     `gorm:"type:timestamp"`
	CreatedAt   time.Time      `gorm:"type:timestamp"`
	UpdatedAt   time.Time      `gorm:"type:timestamp"`
	DeletedAt   gorm.DeletedAt `gorm:"type:timestamp;index"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	return nil
}
