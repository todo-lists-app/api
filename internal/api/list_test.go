package api

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/todo-lists-app/protobufs/generated/todo/v1"
	"google.golang.org/grpc"
	"testing"
)

type MockTodoServiceClient struct {
	mock.Mock
}

func (m *MockTodoServiceClient) NewTodoServiceClient() pb.TodoServiceClient {
	return m
}

func (m *MockTodoServiceClient) Get(ctx context.Context, in *pb.TodoGetRequest, opts ...grpc.CallOption) (*pb.TodoRetrieveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.TodoRetrieveResponse), args.Error(1)
}

func (m *MockTodoServiceClient) Update(ctx context.Context, in *pb.TodoInjectRequest, opts ...grpc.CallOption) (*pb.TodoRetrieveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.TodoRetrieveResponse), args.Error(1)
}

func (m *MockTodoServiceClient) Delete(ctx context.Context, in *pb.TodoDeleteRequest, opts ...grpc.CallOption) (*pb.TodoRetrieveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.TodoRetrieveResponse), args.Error(1)
}

func (m *MockTodoServiceClient) Insert(ctx context.Context, in *pb.TodoInjectRequest, opts ...grpc.CallOption) (*pb.TodoRetrieveResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.TodoRetrieveResponse), args.Error(1)
}

func TestGetList(t *testing.T) {
	// Create a mock gRPC client
	mockClient := new(MockTodoServiceClient)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(&pb.TodoRetrieveResponse{
		UserId: "testUserID",
		Data:   "testData",
		Iv:     "testIV",
	}, nil)

	// Inject the mock client directly into the List struct
	list := &List{
		Context: context.Background(),
		UserID:  "testUserID",
		Client:  mockClient,
	}

	// Call the GetList method
	result, err := list.GetList()

	// Assert the results
	assert.Nil(t, err)
	assert.Equal(t, "testUserID", result.UserID)
	assert.Equal(t, "testData", result.Data)
	assert.Equal(t, "testIV", result.IV)
}
