package genre

import (
	"context"

	"github.com/literalog/library/pkg/models"
)

type Repository interface {
	Create(ctx context.Context, genre *models.Genre) error
	Update(ctx context.Context, genre *models.Genre) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*models.Genre, error)
	GetByName(ctx context.Context, name string) (*models.Genre, error)
	GetAll(ctx context.Context) ([]models.Genre, error)
}
