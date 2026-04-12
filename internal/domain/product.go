package domain

import (
	"time"
)

type Product struct {
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

func (p *Product) Validate() error {
	if len(p.Name) < 1 || len(p.Name) > 60 {
		return ErrInvalidProductName
	}
	if len(p.Description) < 1 || len(p.Description) > 600 {
		return ErrInvalidProductDescription
	}
	if p.Price <= 0 {
		return ErrInvalidProductPrice
	}
	if p.Stock < 0 || p.Stock > 99999 {
		return ErrInvalidProductStock
	}
	return nil
}
