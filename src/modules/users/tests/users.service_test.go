package tests

import (
	"errors"
	"testing"

	"api/src/modules/users"
	"api/src/modules/users/dto"
	"api/src/modules/users/entities"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository implements UserRepository to simulate a Database layer without querying a real physical DB
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *entities.User) error {
	args := m.Called(user)
	user.ID = 1 // Simulate Auto Increment by DB
	return args.Error(0)
}

func (m *MockUserRepository) FindAll() ([]entities.User, error) {
	args := m.Called()
	return args.Get(0).([]entities.User), args.Error(1)
}

func (m *MockUserRepository) FindById(id uint) (*entities.User, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) Update(user *entities.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) FindByEmail(email string) (*entities.User, error) {
	args := m.Called(email)
	if args.Get(0) != nil {
		return args.Get(0).(*entities.User), args.Error(1)
	}
	return nil, args.Error(1)
}

// Compile assertion that Mock implements the real repo
var _ users.UserRepository = (*MockUserRepository)(nil)

// Test Unit: Service validates object and passes successfully to DB
func TestUserService_Create_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := users.NewUserService(mockRepo)

	d := dto.CreateUserDto{
		Name:     "Testing Name",
		Email:    "test@domain.com",
		Password: "strongpassword",
	}

	// Tell the mock that when 'Create' is called with ANY user pointer, return 'nil' (no error)
	mockRepo.On("Create", mock.AnythingOfType("*entities.User")).Return(nil)

	// Action
	user, err := service.Create(d)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "Testing Name", user.Name)
	assert.Equal(t, "test@domain.com", user.Email)
	assert.Equal(t, uint(1), user.ID) // ID was injected by our mock repo

	// Verifies if the mocked function was actually called!
	mockRepo.AssertExpectations(t)
}

func TestUserService_FindById_ThrowsError_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := users.NewUserService(mockRepo)

	// Instruct the mock to throw a fake Database error
	mockRepo.On("FindById", uint(99)).Return(nil, errors.New("record not found"))

	user, err := service.FindById(99)

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "record not found", err.Error())

	mockRepo.AssertExpectations(t)
}
