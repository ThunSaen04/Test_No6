package service

import (
	"context"
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code39"

	"no6/backend/internal/config"
	"no6/backend/internal/models"
	"no6/backend/internal/modules/barcode/utils"
)

var code39Pattern = regexp.MustCompile(`^[0-9A-Z\-\.\ \$/\+%]+$`)

type BarcodeService struct {
	repo ProductRepository
	cfg  config.Config
}

type ProductRepository interface {
	List(ctx context.Context) ([]models.Product, error)
	FindByProductCode(ctx context.Context, productCode string) (models.Product, error)
	FindByID(ctx context.Context, id uint) (models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, product *models.Product) error
}

func NewBarcodeService(repo ProductRepository, cfg config.Config) *BarcodeService {
	return &BarcodeService{
		repo: repo,
		cfg:  cfg,
	}
}

func (s *BarcodeService) List(ctx context.Context) ([]models.Product, error) {
	return s.repo.List(ctx)
}

func (s *BarcodeService) Create(ctx context.Context, rawCode string) (models.Product, error) {
	productCode, err := normalizeCode39(rawCode)
	if err != nil {
		return models.Product{}, err
	}

	_, err = s.repo.FindByProductCode(ctx, productCode)
	if err == nil {
		return models.Product{}, utils.ErrDuplicate
	}
	if !errors.Is(err, utils.ErrNotFound) {
		return models.Product{}, err
	}

	fileName := fmt.Sprintf("barcode_%d.png", time.Now().UnixNano())
	absolutePath := filepath.Join(s.cfg.BarcodeDir, fileName)
	if err := createCode39PNG(productCode, absolutePath, s.cfg.BarcodeWidth, s.cfg.BarcodeHeight); err != nil {
		return models.Product{}, fmt.Errorf("generate barcode failed: %w", err)
	}

	product := models.Product{
		ProductCode: productCode,
		BarcodePath: "/images/" + fileName,
	}

	if err := s.repo.Create(ctx, &product); err != nil {
		_ = os.Remove(absolutePath)
		return models.Product{}, err
	}

	return product, nil
}

func (s *BarcodeService) Delete(ctx context.Context, id uint) error {
	product, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return utils.ErrNotFound
		}
		return err
	}

	if err := s.repo.Delete(ctx, &product); err != nil {
		return err
	}

	imageFile := filepath.Base(product.BarcodePath)
	if imageFile != "." && imageFile != "" {
		_ = os.Remove(filepath.Join(s.cfg.BarcodeDir, imageFile))
	}

	return nil
}

func BarcodeURL(cfg config.Config, barcodePath string) string {
	if strings.HasPrefix(barcodePath, "http://") || strings.HasPrefix(barcodePath, "https://") {
		return barcodePath
	}
	return cfg.PublicBaseURL + barcodePath
}

func normalizeCode39(raw string) (string, error) {
	value := strings.ToUpper(strings.TrimSpace(raw))
	if value == "" {
		return "", utils.ErrInvalidCode
	}

	if strings.ContainsRune(value, '*') {
		return "", utils.ErrInvalidCode
	}

	if !code39Pattern.MatchString(value) {
		return "", utils.ErrInvalidCode
	}

	return value, nil
}

func createCode39PNG(value, filePath string, width, height int) error {
	bc, err := code39.Encode(value, false, false)
	if err != nil {
		return err
	}

	scaled, err := barcode.Scale(bc, width, height)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return png.Encode(file, scaled)
}
