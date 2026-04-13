package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) SaveProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) GetProductByID(id string) (*domain.Product, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Product), args.Error(1)
}

func (m *MockProductRepository) UpdateProduct(product *domain.Product) error {
	args := m.Called(product)
	return args.Error(0)
}

func (m *MockProductRepository) DeleteProduct(id string, sellerID string) error {
	args := m.Called(id, sellerID)
	return args.Error(0)
}

func (m *MockProductRepository) ListProducts(page int, limit int) ([]*domain.Product, int, error) {
	args := m.Called(page, limit)
	if args.Get(0) == nil {
		return nil, args.Int(1), args.Error(2)
	}
	return args.Get(0).([]*domain.Product), args.Int(1), args.Error(2)
}

func TestUpdateProduct_Success(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewUpdateProductUseCase(mockRepo)

	existingProduct := &domain.Product{
		ID:       "product-1",
		Name:     "Old Name",
		Price:    1000,
		Stock:    10,
		SellerID: "seller-1",
		Category: "electronics",
		ImageURL: "http://example.com/image.png",
	}

	mockRepo.On("GetProductByID", "product-1").Return(existingProduct, nil)
	mockRepo.On("UpdateProduct", mock.AnythingOfType("*domain.Product")).Return(nil)

	req := UpdateProductRequest{
		ID:          "product-1",
		Name:        "New Name",
		Description: "New Description",
		Price:       2000,
		Stock:       20,
		SellerID:    "seller-1",
	}

	resp, err := uc.Execute(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "product-1", resp.ID)
	assert.Equal(t, "New Name", resp.Name)
	assert.Equal(t, "New Description", resp.Description)
	assert.Equal(t, int64(2000), resp.Price)
	assert.Equal(t, 20, resp.Stock)
	assert.Equal(t, "seller-1", resp.SellerID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_Unauthorized(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewUpdateProductUseCase(mockRepo)

	existingProduct := &domain.Product{
		ID:       "product-1",
		Name:     "Old Name",
		Price:    1000,
		Stock:    10,
		SellerID: "seller-1",
	}

	mockRepo.On("GetProductByID", "product-1").Return(existingProduct, nil)

	req := UpdateProductRequest{
		ID:          "product-1",
		Name:        "New Name",
		Description: "New Description",
		Price:       2000,
		Stock:       20,
		SellerID:    "other-seller",
	}

	resp, err := uc.Execute(req)

	assert.Nil(t, resp)
	assert.Equal(t, ErrUnauthorizedUpdateProduct, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_NotFound(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewUpdateProductUseCase(mockRepo)

	mockRepo.On("GetProductByID", "nonexistent").Return(nil, domain.ErrProductNotFound)

	req := UpdateProductRequest{
		ID:          "nonexistent",
		Name:        "New Name",
		Description: "New Description",
		Price:       2000,
		Stock:       20,
		SellerID:    "seller-1",
	}

	resp, err := uc.Execute(req)

	assert.Nil(t, resp)
	assert.Equal(t, domain.ErrProductNotFound, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_InvalidPrice(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewUpdateProductUseCase(mockRepo)

	existingProduct := &domain.Product{
		ID:       "product-1",
		Name:     "Old Name",
		Price:    1000,
		Stock:    10,
		SellerID: "seller-1",
		Category: "electronics",
		ImageURL: "http://example.com/image.png",
	}

	mockRepo.On("GetProductByID", "product-1").Return(existingProduct, nil)

	req := UpdateProductRequest{
		ID:          "product-1",
		Name:        "New Name",
		Description: "New Description",
		Price:       -100,
		Stock:       20,
		SellerID:    "seller-1",
	}

	resp, err := uc.Execute(req)

	assert.Nil(t, resp)
	assert.Equal(t, domain.ErrInvalidProductPrice, err)
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct_InvalidName(t *testing.T) {
	mockRepo := new(MockProductRepository)
	uc := NewUpdateProductUseCase(mockRepo)

	existingProduct := &domain.Product{
		ID:       "product-1",
		Name:     "Old Name",
		Price:    1000,
		Stock:    10,
		SellerID: "seller-1",
		Category: "electronics",
		ImageURL: "http://example.com/image.png",
	}

	mockRepo.On("GetProductByID", "product-1").Return(existingProduct, nil)

	req := UpdateProductRequest{
		ID:          "product-1",
		Name:        "",
		Description: "New Description",
		Price:       2000,
		Stock:       20,
		SellerID:    "seller-1",
	}

	resp, err := uc.Execute(req)

	assert.Nil(t, resp)
	assert.Equal(t, domain.ErrInvalidProductName, err)
	mockRepo.AssertExpectations(t)
}
