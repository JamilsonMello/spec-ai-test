package usecase

import (
	"errors"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

var (
	ErrForbiddenDeleteProduct = errors.New("Forbidden: product ownership required")
)

type DeleteProductUseCase struct {
	ProductRepository domain.ProductRepository
}

func NewDeleteProductUseCase(repo domain.ProductRepository) *DeleteProductUseCase {
	return &DeleteProductUseCase{
		ProductRepository: repo,
	}
}

func (uc *DeleteProductUseCase) Execute(id string, sellerID string) error {
	product, err := uc.ProductRepository.GetProductByID(id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return domain.ErrProductNotFound
		}
		return err
	}

	if product.SellerID != sellerID {
		return ErrForbiddenDeleteProduct
	}

	return uc.ProductRepository.DeleteProduct(id, sellerID)
}
