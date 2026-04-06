package handler

import (
	"errors"
	"strconv"

	"github.com/gofiber/fiber/v2"

	"no6/backend/internal/config"
	"no6/backend/internal/models"
	"no6/backend/internal/modules/barcode/dto"
	"no6/backend/internal/modules/barcode/service"
	"no6/backend/internal/modules/barcode/utils"
	productUtils "no6/backend/internal/modules/barcode/utils"
)

type BarcodeHandler struct {
	service *service.BarcodeService
	cfg     config.Config
}

func NewBarcodeHandler(service *service.BarcodeService, cfg config.Config) *BarcodeHandler {
	return &BarcodeHandler{
		service: service,
		cfg:     cfg,
	}
}

func (h *BarcodeHandler) List(c *fiber.Ctx) error {
	products, err := h.service.List(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": utils.ErrList.Error(),
		})
	}

	output := make([]dto.BarcodeResponse, 0, len(products))
	for _, product := range products {
		output = append(output, toBarcodeResponse(product, h.cfg))
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": output,
	})
}

func (h *BarcodeHandler) Create(c *fiber.Ctx) error {
	var req dto.CreateBarcodeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.ErrInvalidRequest.Error(),
		})
	}

	product, err := h.service.Create(c.Context(), productUtils.ResolveCode(req))
	if err != nil {
		switch {
		case errors.Is(err, utils.ErrInvalidCode):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": utils.ErrInvalidCode.Error(),
			})
		case errors.Is(err, utils.ErrDuplicate):
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"error": utils.ErrDuplicate.Error(),
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": utils.ErrCreate.Error(),
			})
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": toBarcodeResponse(product, h.cfg),
	})
}

func (h *BarcodeHandler) Delete(c *fiber.Ctx) error {
	idUint64, err := strconv.ParseUint(c.Params("id"), 10, 32)
	if err != nil || idUint64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": utils.ErrInvalidId.Error(),
		})
	}

	if err := h.service.Delete(c.Context(), uint(idUint64)); err != nil {
		if errors.Is(err, utils.ErrNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": utils.ErrNotFound.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": utils.ErrDelete.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "barcode deleted",
	})
}

func toBarcodeResponse(product models.Product, cfg config.Config) dto.BarcodeResponse {
	return dto.BarcodeResponse{
		ID:          product.ID,
		ProductCode: product.ProductCode,
		BarcodePath: product.BarcodePath,
		BarcodeURL:  service.BarcodeURL(cfg, product.BarcodePath),
		CreatedAt:   product.CreatedAt,
	}
}
