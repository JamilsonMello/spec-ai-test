package usecase

import (
	"time"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type ListProductsRequest struct {
	Page  int
	Limit int
}

type ProductItem struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	SellerID    string    `json:"seller_id"`
	Stock       int       `json:"stock"`
	Category    string    `json:"category"`
	ImageURL    string    `json:"image_url"`
	CommunityID string    `json:"community_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListProductsResponse struct {
	Products []ProductItem `json:"products"`
	Total    int           `json:"total"`
	Page     int           `json:"page"`
	Limit    int           `json:"limit"`
}

type ListProductsUseCase struct {
	ProductRepository domain.ProductRepository
}

func NewListProductsUseCase(repo domain.ProductRepository) *ListProductsUseCase {
	return &ListProductsUseCase{
		ProductRepository: repo,
	}
}

func (uc *ListProductsUseCase) Execute(req ListProductsRequest) (*ListProductsResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 30
	}

	products, total, err := uc.ProductRepository.ListProducts(req.Page, req.Limit)
	if err != nil {
		return nil, err
	}

	items := make([]ProductItem, 0, len(products))
	for _, p := range products {
		items = append(items, ProductItem{
			ID:          p.ID,
			Name:        p.Name,
			Description: p.Description,
			Price:       p.Price,
			SellerID:    p.SellerID,
			Stock:       p.Stock,
			Category:    p.Category,
			ImageURL:    p.ImageURL,
			CommunityID: p.CommunityID,
			CreatedAt:   p.CreatedAt,
		})
	}

	return &ListProductsResponse{
		Products: items,
		Total:    total,
		Page:     req.Page,
		Limit:    req.Limit,
	}, nil
}
