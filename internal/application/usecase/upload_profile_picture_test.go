package usecase

import (
	"bytes"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"net/textproto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/example/cadastro-de-usuarios/internal/domain"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) SaveUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByUuid(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateProfilePictureURL(userID string, url string) error {
	args := m.Called(userID, url)
	return args.Error(0)
}

func (m *MockUserRepository) ListUsers(filter domain.UserFilter, page int, limit int) ([]*domain.User, int, error) {
	args := m.Called(filter, page, limit)
	return args.Get(0).([]*domain.User), args.Int(1), args.Error(2)
}

type MockFileStorage struct {
	mock.Mock
}

func (m *MockFileStorage) Save(file io.Reader, filename string) (string, error) {
	args := m.Called(file, filename)
	return args.String(0), args.Error(1)
}

func createJPEGFileHeader(t *testing.T, size int64) (multipart.File, *multipart.FileHeader) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	jpeg.Encode(&buf, img, nil)

	content := buf.Bytes()
	if int64(len(content)) < size {
		padding := make([]byte, size-int64(len(content)))
		content = append(content, padding...)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="foto"; filename="photo.jpg"`)
	h.Set("Content-Type", "image/jpeg")
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(10 << 20)
	fileHeaders := form.File["foto"]
	f, _ := fileHeaders[0].Open()

	return f, fileHeaders[0]
}

func createPNGFileHeader(t *testing.T) (multipart.File, *multipart.FileHeader) {
	var buf bytes.Buffer
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	png.Encode(&buf, img)

	content := buf.Bytes()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="foto"; filename="photo.png"`)
	h.Set("Content-Type", "image/png")
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(10 << 20)
	fileHeaders := form.File["foto"]
	f, _ := fileHeaders[0].Open()

	return f, fileHeaders[0]
}

func createPDFFileHeader(t *testing.T) (multipart.File, *multipart.FileHeader) {
	content := []byte("%PDF-1.4 fake pdf content for testing purposes only")

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="foto"; filename="document.pdf"`)
	h.Set("Content-Type", "application/pdf")
	part, _ := writer.CreatePart(h)
	part.Write(content)
	writer.Close()

	reader := multipart.NewReader(body, writer.Boundary())
	form, _ := reader.ReadForm(10 << 20)
	fileHeaders := form.File["foto"]
	f, _ := fileHeaders[0].Open()

	return f, fileHeaders[0]
}

func TestUploadProfilePicture_Success_JPEG(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createJPEGFileHeader(t, 1024)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)
	mockStorage.On("Save", mock.Anything, mock.MatchedBy(func(name string) bool {
		return len(name) > 4
	})).Return("/uploads/profile/test.jpg", nil)
	mockRepo.On("UpdateProfilePictureURL", userID, "/uploads/profile/test.jpg").Return(nil)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Foto atualizada com sucesso", resp.Message)
	assert.Equal(t, "/uploads/profile/test.jpg", resp.URL)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUploadProfilePicture_Success_PNG(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createPNGFileHeader(t)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)
	mockStorage.On("Save", mock.Anything, mock.MatchedBy(func(name string) bool {
		return len(name) > 4
	})).Return("/uploads/profile/test.png", nil)
	mockRepo.On("UpdateProfilePictureURL", userID, "/uploads/profile/test.png").Return(nil)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "Foto atualizada com sucesso", resp.Message)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUploadProfilePicture_FileTooLarge(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createJPEGFileHeader(t, 3*1024*1024)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrFileTooLarge, err)
	mockStorage.AssertNotCalled(t, "Save")
}

func TestUploadProfilePicture_UnsupportedFormat(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createPDFFileHeader(t)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrUnsupportedFileFormat, err)
	mockStorage.AssertNotCalled(t, "Save")
}

func TestUploadProfilePicture_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "nonexistent-user-id"
	file, header := createJPEGFileHeader(t, 1024)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(nil, domain.ErrUserNotFound)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrUserNotFoundUpload, err)
	mockStorage.AssertNotCalled(t, "Save")
	mockRepo.AssertExpectations(t)
}

func TestUploadProfilePicture_EmptyUserID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	file, header := createJPEGFileHeader(t, 1024)
	defer file.Close()

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: "",
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrUserNotFoundUpload, err)
}

func TestUploadProfilePicture_StorageSaveError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createJPEGFileHeader(t, 1024)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)
	mockStorage.On("Save", mock.Anything, mock.Anything).Return("", assert.AnError)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrSaveFileFailed, err)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}

func TestUploadProfilePicture_UpdateProfilePictureURLError(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockStorage := new(MockFileStorage)
	uc := NewUploadProfilePictureUseCase(mockRepo, mockStorage)

	userID := "123e4567-e89b-12d3-a456-426614174000"
	file, header := createJPEGFileHeader(t, 1024)
	defer file.Close()

	mockRepo.On("FindUserByUuid", userID).Return(&domain.User{ID: userID}, nil)
	mockStorage.On("Save", mock.Anything, mock.Anything).Return("/uploads/profile/test.jpg", nil)
	mockRepo.On("UpdateProfilePictureURL", userID, "/uploads/profile/test.jpg").Return(assert.AnError)

	resp, err := uc.Execute(UploadProfilePictureRequest{
		UserID: userID,
		File:   file,
		Header: header,
	})

	assert.Nil(t, resp)
	assert.Equal(t, ErrSaveFileFailed, err)
	mockRepo.AssertExpectations(t)
	mockStorage.AssertExpectations(t)
}
