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

const selectProductByIDSQL = `
SELECT id, name, description, price, seller_id, stock, category, image_url, community_id, created_at
FROM products
WHERE id = $1
`

const updateProductSQL = `
UPDATE products SET name = $1, description = $2, price = $3, stock = $4 WHERE id = $5
`

const deleteProductSQL = `
DELETE FROM products WHERE id = $1 AND seller_id = $2
`

const countProductsSQL = `SELECT COUNT(*) FROM products`

const listProductsSQL = `
SELECT id, name, description, price, seller_id, stock, category, image_url, community_id, created_at
FROM products
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProductByID(id string) (*domain.Product, error) {
	product := &domain.Product{}
	var communityID sql.NullString
	err := r.db.QueryRow(selectProductByIDSQL, id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price,
		&product.SellerID, &product.Stock, &product.Category, &product.ImageURL,
		&communityID, &product.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get product by id: %w", err)
	}
	if communityID.Valid {
		product.CommunityID = communityID.String
	}
	return product, nil
}

func (r *ProductRepository) UpdateProduct(product *domain.Product) error {
	result, err := r.db.Exec(updateProductSQL, product.Name, product.Description, product.Price, product.Stock, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) DeleteProduct(id string, sellerID string) error {
	result, err := r.db.Exec(deleteProductSQL, id, sellerID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}
	if affected == 0 {
		return domain.ErrProductNotFound
	}
	return nil
}

func (r *ProductRepository) ListProducts(page int, limit int) ([]*domain.Product, int, error) {
	var total int
	err := r.db.QueryRow(countProductsSQL).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query(listProductsSQL, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		product := &domain.Product{}
		var communityID sql.NullString
		err := rows.Scan(
			&product.ID, &product.Name, &product.Description, &product.Price,
			&product.SellerID, &product.Stock, &product.Category, &product.ImageURL,
			&communityID, &product.CreatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan product: %w", err)
		}
		if communityID.Valid {
			product.CommunityID = communityID.String
		}
		products = append(products, product)
	}

	return products, total, nil
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
