package repositories

import (
	"context"
	"errors"

	"github.com/gobp/gobp/core"
	"github.com/gobp/gobp/models/entities"

	"gorm.io/gorm"
)

type IUserRepository interface {
	core.IBaseRepository[string, entities.User, entities.User, entities.User]

	FindByUsernameOrEmail(ctx context.Context, username, email string) (entities.User, *core.Exception)
	ChangePassword(ctx context.Context, userID, hashedPassword, signature string) *core.Exception
}

type gormUserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &gormUserRepository{db}
}

func (r *gormUserRepository) Create(ctx context.Context, entity entities.User) (entities.User, *core.Exception) {
	err := r.db.Create(&entity).Error
	if err != nil {
		return entities.User{}, core.Unexpected(core.WithError(err))
	}

	return entity, nil
}

func (r *gormUserRepository) FindByUsernameOrEmail(ctx context.Context, username, email string) (entities.User, *core.Exception) {

	var user entities.User
	if err := r.db.
		Or(entities.User{Username: username}).
		Or(entities.User{Email: email}).
		Preload("Manager").
		First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, core.NotFound()
		}
		return entities.User{}, core.Unexpected(core.WithError(err))
	}

	return user, nil
}

func (r *gormUserRepository) ChangePassword(ctx context.Context, id, hashedPassword, signature string) *core.Exception {

	var user entities.User
	r.db.Where(&entities.User{ID: id}).Find(&user)

	if user.ID != id {
		return core.NotFound()
	}

	user.Password = hashedPassword
	user.Signature = signature

	r.db.Save(&user)

	return nil
}

func (r *gormUserRepository) GetByID(ctx context.Context, id string) (entities.User, *core.Exception) {
	var user entities.User

	if err := r.db.Preload("Manager").Preload("Tasks").First(&user, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.User{}, core.NotFound()
		}
		return entities.User{}, core.Unexpected(core.WithError(err))
	}

	return user, nil
}

func (r *gormUserRepository) GetAll(ctx context.Context, limit, offset int) ([]entities.User, int, *core.Exception) {

	var users []entities.User
	var total int64
	r.db.Model(&entities.User{}).Preload("Manager").Count(&total).Limit(limit).Offset(offset).Find(&users)

	return users, int(total), nil
}

func (r *gormUserRepository) Update(ctx context.Context, id string, changes entities.User) (entities.User, *core.Exception) {
	var user entities.User

	query := r.db.Model(&entities.User{}).Where("id = ?", id).Preload("Manager")
	if changes.Username != "" {
		query = query.Update("username", changes.Username)
	}

	if changes.Email != "" {
		query = query.Update("email", changes.Email)
	}

	if changes.IsEmailVerified != nil {
		query = query.Update("is_email_verified", *changes.IsEmailVerified)
	}

	query.Find(&user)

	return user, nil
}

func (r *gormUserRepository) Delete(ctx context.Context, id string) *core.Exception {

	r.db.Delete(&entities.User{}, "id = ?", id)

	return nil
}
