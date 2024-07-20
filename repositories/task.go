package repositories

import (
	"context"
	"errors"

	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/models/entities"
	"golang.org/x/exp/slices"

	"gorm.io/gorm"
)

type ITaskRepository interface {
	core.IBaseRepository[string, entities.Task, entities.Task, entities.Task]
}

type gormTaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) ITaskRepository {
	return &gormTaskRepository{db}
}

func (r *gormTaskRepository) Create(ctx context.Context, entity entities.Task) (entities.Task, *core.Exception) {
	err := r.db.Create(&entity).Error
	if err != nil {
		return entities.Task{}, core.Unexpected(core.WithError(err))
	}

	return entity, nil
}

func (r *gormTaskRepository) GetByID(ctx context.Context, id string) (entities.Task, *core.Exception) {
	appCtx, _ := core.FromContext(ctx)

	baseQuery := r.db.Model(&entities.Task{})

	if !slices.Contains(appCtx.Roles(), "manager") {
		baseQuery = baseQuery.Where("user_id = ?", appCtx.UserID())
	}

	var task entities.Task
	if err := baseQuery.First(&task, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Task{}, core.NotFound()
		}
		return entities.Task{}, core.Unexpected(core.WithError(err))
	}

	return task, nil
}

func (r *gormTaskRepository) GetAll(ctx context.Context, limit, offset int) ([]entities.Task, int, *core.Exception) {
	appCtx, _ := core.FromContext(ctx)

	var tasks []entities.Task
	var total int64
	baseQuery := r.db.Model(&entities.Task{})

	if !slices.Contains(appCtx.Roles(), "manager") {
		baseQuery = baseQuery.Where("user_id = ?", appCtx.UserID())
	}

	baseQuery.Count(&total).Limit(limit).Offset(offset).Find(&tasks)

	return tasks, int(total), nil
}

func (r *gormTaskRepository) Update(ctx context.Context, id string, changes entities.Task) (entities.Task, *core.Exception) {
	appCtx, _ := core.FromContext(ctx)

	baseQuery := r.db.Model(&entities.Task{})
	if !slices.Contains(appCtx.Roles(), "manager") {
		baseQuery = baseQuery.Where("user_id = ?", appCtx.UserID())
	}

	baseQuery = baseQuery.Where("id = ?", id)
	if changes.Summary != "" {
		baseQuery = baseQuery.Update("summary", changes.Summary)
	}

	// if changes.UserID != "" {
	// 	baseQuery = baseQuery.Update("user_id", changes.UserID)
	// }

	if changes.PerformedAt != nil {
		baseQuery = baseQuery.Update("performed_at", changes.PerformedAt)
	}

	var task entities.Task
	if err := baseQuery.Find(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.Task{}, core.NotFound()
		}
		return entities.Task{}, core.Unexpected(core.WithError(err))
	}

	return task, nil
}

func (r *gormTaskRepository) Delete(ctx context.Context, id string) *core.Exception {

	r.db.Delete(&entities.Task{}, "id = ?", id)

	return nil
}
