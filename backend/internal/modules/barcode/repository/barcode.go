package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"no6/backend/internal/models"
	"no6/backend/internal/modules/barcode/utils"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) List(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	if err := r.db.WithContext(ctx).Order("id desc").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) FindByProductCode(ctx context.Context, productCode string) (models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).Where("product_code = ?", productCode).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Product{}, utils.ErrNotFound
		}
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepository) FindByID(ctx context.Context, id uint) (models.Product, error) {
	var product models.Product
	if err := r.db.WithContext(ctx).First(&product, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.Product{}, utils.ErrNotFound
		}
		return models.Product{}, err
	}
	return product, nil
}

func (r *ProductRepository) Create(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *ProductRepository) Delete(ctx context.Context, product *models.Product) error {
	return r.db.WithContext(ctx).Delete(product).Error
}
