package usecase

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type mockPostRepository struct {
	mock.Mock
}

func (m *mockPostRepository) SavePost(post *domain.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func (m *mockPostRepository) GetPostByID(id string) (*domain.Post, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Post), args.Error(1)
}

func (m *mockPostRepository) UpdatePost(post *domain.Post) error {
	args := m.Called(post)
	return args.Error(0)
}

func TestCreatePostUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		request     CreatePostRequest
		mockSetup   func(m *mockPostRepository)
		expectedErr error
		expectResp  bool
	}{
		{
			name: "sucesso na criação do post",
			request: CreatePostRequest{
				Title:    "Título válido",
				Content:  "Conteúdo válido do post",
				AuthorID: "user-123",
			},
			mockSetup: func(m *mockPostRepository) {
				m.On("SavePost", mock.AnythingOfType("*domain.Post")).Return(nil)
			},
			expectedErr: nil,
			expectResp:  true,
		},
		{
			name: "erro título curto",
			request: CreatePostRequest{
				Title:    "abc",
				Content:  "Conteúdo válido",
				AuthorID: "user-123",
			},
			mockSetup:   func(m *mockPostRepository) {},
			expectedErr: ErrInvalidRequest,
			expectResp:  false,
		},
		{
			name: "erro título longo",
			request: CreatePostRequest{
				Title:    strings.Repeat("a", 101),
				Content:  "Conteúdo válido",
				AuthorID: "user-123",
			},
			mockSetup:   func(m *mockPostRepository) {},
			expectedErr: ErrInvalidRequest,
			expectResp:  false,
		},
		{
			name: "erro conteúdo vazio",
			request: CreatePostRequest{
				Title:    "Título válido",
				Content:  "",
				AuthorID: "user-123",
			},
			mockSetup:   func(m *mockPostRepository) {},
			expectedErr: ErrInvalidRequest,
			expectResp:  false,
		},
		{
			name: "erro no repositório",
			request: CreatePostRequest{
				Title:    "Título válido",
				Content:  "Conteúdo válido do post",
				AuthorID: "user-123",
			},
			mockSetup: func(m *mockPostRepository) {
				m.On("SavePost", mock.AnythingOfType("*domain.Post")).Return(errors.New("database connection failed"))
			},
			expectedErr: ErrRepositoryFailure,
			expectResp:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(mockPostRepository)
			tt.mockSetup(mockRepo)

			uc := NewCreatePostUseCase(mockRepo)
			resp, err := uc.Execute(tt.request)

			if tt.expectedErr != nil {
				assert.ErrorIs(t, err, tt.expectedErr)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}

			if tt.expectResp {
				assert.NotEmpty(t, resp.ID)
				assert.Equal(t, tt.request.Content, resp.Content)
				assert.Equal(t, tt.request.AuthorID, resp.AuthorID)
				assert.NotEmpty(t, resp.CreatedAt)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCreatePostUseCase_ValidationSkipsRepository(t *testing.T) {
	mockRepo := new(mockPostRepository)
	uc := NewCreatePostUseCase(mockRepo)

	_, err := uc.Execute(CreatePostRequest{
		Title:    "ab",
		Content:  "Conteúdo",
		AuthorID: "user-123",
	})

	assert.ErrorIs(t, err, ErrInvalidRequest)
	mockRepo.AssertNotCalled(t, "SavePost")
}
