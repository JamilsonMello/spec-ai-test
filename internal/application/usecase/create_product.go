package usecase

import (
	"errors"
	"html"
	"time"

	"github.com/google/uuid"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrUnauthorizedCreateProduct = errors.New("Unauthorized")
	ErrInvalidProductInput       = errors.New("invalid product input")
)

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       int64   `json:"price"`
	Stock       int     `json:"stock"`
	Category    string  `json:"category"`
	ImageURL    string  `json:"image_url"`
	CommunityID *string `json:"community_id,omitempty"`
	SellerID    string  `json:"-"`
}

type CreateProductResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	SellerID    string `json:"seller_id"`
	Stock       int    `json:"stock"`
	Category    string `json:"category"`
	ImageURL    string `json:"image_url"`
	CommunityID string `json:"community_id,omitempty"`
	CreatedAt   string `json:"created_at"`
}

type CreateProductUseCase struct {
	ProductRepository   domain.ProductRepository
	CommunityRepository domain.CommunityRepository
}

func NewCreateProductUseCase(productRepo domain.ProductRepository, communityRepo domain.CommunityRepository) *CreateProductUseCase {
	return &CreateProductUseCase{
		ProductRepository:   productRepo,
		CommunityRepository: communityRepo,
	}
}

func (uc *CreateProductUseCase) Execute(req CreateProductRequest) (*CreateProductResponse, error) {
	if req.SellerID == "" {
		return nil, ErrUnauthorizedCreateProduct
	}

	req.Name = html.EscapeString(req.Name)
	req.Description = html.EscapeString(req.Description)

	if req.CommunityID != nil && *req.CommunityID != "" {
		if _, err := uuid.Parse(*req.CommunityID); err != nil {
			return nil, ErrInvalidUUIDFormat
		}

		exists, err := uc.CommunityRepository.ExistsByID(*req.CommunityID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, domain.ErrCommunityNotFound
		}
	}

	product := &domain.Product{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		SellerID:    req.SellerID,
		Stock:       req.Stock,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		CreatedAt:   time.Now(),
	}

	if req.CommunityID != nil && *req.CommunityID != "" {
		product.CommunityID = *req.CommunityID
	}

	if err := product.Validate(); err != nil {
		return nil, err
	}

	if err := uc.ProductRepository.SaveProduct(product); err != nil {
		return nil, err
	}

	return &CreateProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		SellerID:    product.SellerID,
		Stock:       product.Stock,
		Category:    product.Category,
		ImageURL:    product.ImageURL,
		CommunityID: product.CommunityID,
		CreatedAt:   product.CreatedAt.Format(time.RFC3339),
	}, nil
}
