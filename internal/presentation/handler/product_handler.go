package handler

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/example/cadastro-de-usuarios/internal/application/usecase"
	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ProductHandler struct {
	CreateProductUseCase *usecase.CreateProductUseCase
	UpdateProductUseCase *usecase.UpdateProductUseCase
	DeleteProductUseCase *usecase.DeleteProductUseCase
	ListProductsUseCase  *usecase.ListProductsUseCase
}

func NewProductHandler(createProductUC *usecase.CreateProductUseCase, updateProductUC *usecase.UpdateProductUseCase, deleteProductUC *usecase.DeleteProductUseCase, listProductsUC *usecase.ListProductsUseCase) *ProductHandler {
	return &ProductHandler{
		CreateProductUseCase: createProductUC,
		UpdateProductUseCase: updateProductUC,
		DeleteProductUseCase: deleteProductUC,
		ListProductsUseCase:  listProductsUC,
	}
}

func (h *ProductHandler) CreateProduct(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	var req usecase.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload", "code": "INVALID_INPUT"})
	}

	req.SellerID = userID

	resp, err := h.CreateProductUseCase.Execute(req)
	if err != nil {
		switch err {
		case domain.ErrInvalidProductName, domain.ErrInvalidProductDescription,
			domain.ErrInvalidProductPrice, domain.ErrInvalidProductStock,
			usecase.ErrInvalidProductInput, usecase.ErrInvalidUUIDFormat:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "INVALID_INPUT"})
		case domain.ErrCommunityNotFound:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error(), "code": "COMMUNITY_NOT_FOUND"})
		case usecase.ErrUnauthorizedCreateProduct:
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error", "code": "INTERNAL_ERROR"})
		}
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *ProductHandler) UpdateProduct(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid product ID"})
	}

	var req usecase.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
	}

	req.ID = productID
	req.SellerID = userID

	resp, err := h.UpdateProductUseCase.Execute(req)
	if err != nil {
		switch err {
		case domain.ErrProductNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		case usecase.ErrUnauthorizedUpdateProduct:
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
		case domain.ErrInvalidProductName, domain.ErrInvalidProductDescription,
			domain.ErrInvalidProductPrice, domain.ErrInvalidProductStock:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *ProductHandler) DeleteProduct(c echo.Context) error {
	userID, ok := c.Get("user_id").(string)
	if !ok || userID == "" {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
	}

	productID := c.Param("id")
	if productID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	err := h.DeleteProductUseCase.Execute(productID, userID)
	if err != nil {
		switch err {
		case domain.ErrProductNotFound:
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Product not found"})
		case usecase.ErrForbiddenDeleteProduct:
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Forbidden: product ownership required"})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
		}
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *ProductHandler) ListProducts(c echo.Context) error {
	page := 1
	limit := 30

	if pageStr := c.QueryParam("page"); pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err != nil || p <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid page parameter"})
		}
		page = p
	}

	if limitStr := c.QueryParam("limit"); limitStr != "" {
		l, err := strconv.Atoi(limitStr)
		if err != nil || l <= 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid limit parameter"})
		}
		limit = l
	}

	resp, err := h.ListProductsUseCase.Execute(usecase.ListProductsRequest{
		Page:  page,
		Limit: limit,
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, resp)
}
