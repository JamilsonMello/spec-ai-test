package repository

import (
	"database/sql"
	"fmt"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

const insertProductSQL = `
INSERT INTO products (id, name, description, price, seller_id, stock, category, image_url, community_id, created_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) SaveProduct(product *domain.Product) error {
	var communityID interface{}
	if product.CommunityID != "" {
		communityID = product.CommunityID
	}
	_, err := r.db.Exec(insertProductSQL, product.ID, product.Name, product.Description, product.Price, product.SellerID, product.Stock, product.Category, product.ImageURL, communityID, product.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to save product: %w", err)
	}
	return nil
}
