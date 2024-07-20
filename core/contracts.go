package core

import (
	"context"
)

type IBaseRepository[IDType, Create, Update, Result any] interface {
	Create(ctx context.Context, entity Create) (Result, *Exception)
	GetByID(ctx context.Context, id IDType) (Result, *Exception)
	GetAll(ctx context.Context, limit, offset int) ([]Result, int, *Exception)
	Update(ctx context.Context, id IDType, changes Update) (Result, *Exception)
	Delete(ctx context.Context, id IDType) *Exception
}

type IBaseService[IDType, Create, Update, Result any] interface {
	Create(ctx context.Context, entity Create) (*Result, *Exception)
	GetByID(ctx context.Context, id IDType) (*Result, *Exception)
	GetAll(ctx context.Context, limit, offset int) (*[]Result, int, *Exception)
	Update(ctx context.Context, id IDType, changes Update) (*Result, *Exception)
	Delete(ctx context.Context, id IDType) *Exception
}
