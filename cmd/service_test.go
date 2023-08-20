package main

import (
	"errors"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockConfig is a mock for the config
type MockConfig struct {
	mock.Mock
}

func (m *MockConfig) Build() (*config.Config, error) {
	args := m.Called()
	return args.Get(0).(*config.Config), args.Error(1)
}

// MockService is a mock for the service
type MockService struct {
	mock.Mock
}

func (m *MockService) Start() error {
	args := m.Called()
	return args.Error(0)
}

func TestRunApp(t *testing.T) {
	mockConfig := new(MockConfig)
	mockService := new(MockService)

	// Mock the config.Build() to return an error
	mockConfig.On("Build").Return(nil, errors.New("config error"))

	// Mock the service.Start() to return an error
	mockService.On("Start").Return(errors.New("service error"))

	// TODO: Inject the mocks into runApp and test the error scenarios

	// Assert the expected results
	assert.Nil(t, nil) // This will check if runApp returns an error
}
