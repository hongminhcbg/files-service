package store

import (
	"context"
	"file-service/models"

	"gorm.io/gorm"
)

type IStore interface {
	GetTemplateById(ctx context.Context, id int64) (*models.Template, error)
}

type _storeImpl struct {
	*gorm.DB
}

func NewStore(db *gorm.DB) IStore {
	return &_storeImpl{db}
}

func (s *_storeImpl) GetTemplateById(ctx context.Context, id int64) (*models.Template, error) {
	result := new(models.Template)
	db := s.DB.WithContext(ctx)
	err := db.Where("id=?", id).First(result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}
