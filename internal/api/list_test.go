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

func TestList_GetList(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	mockClient.On("Get", mock.Anything, mock.Anything).Return(&pb.TodoRetrieveResponse{
		UserId: "testUserID",
		Data:   "testData",
		Iv:     "testIV",
	}, nil)

	list := &List{
		Context: context.Background(),
		UserID:  "testUserID",
		Client:  mockClient,
	}

	result, err := list.GetList()

	assert.Nil(t, err)
	assert.Equal(t, "testUserID", result.UserID)
	assert.Equal(t, "testData", result.Data)
	assert.Equal(t, "testIV", result.IV)
}

func TestList_CreateList(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	mockClient.On("Insert", mock.Anything, mock.Anything).Return(&pb.TodoRetrieveResponse{
		UserId: "testUserID",
		Data:   "testData",
		Iv:     "testIV",
	}, nil)

	list := &List{
		Context: context.Background(),
		UserID:  "testUserID",
		Client:  mockClient,
	}

	storedList := &StoredList{
		UserID: "testUserID",
		Data:   "testData",
		IV:     "testIV",
	}

	result, err := list.CreateList(storedList)

	assert.Nil(t, err)
	assert.Equal(t, "testUserID", result.UserID)
	assert.Equal(t, "testData", result.Data)
	assert.Equal(t, "testIV", result.IV)
}

func TestList_UpdateList(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	mockClient.On("Update", mock.Anything, mock.Anything).Return(&pb.TodoRetrieveResponse{
		UserId: "testUserID",
		Data:   "testData",
		Iv:     "testIV",
	}, nil)

	list := &List{
		Context: context.Background(),
		UserID:  "testUserID",
		Client:  mockClient,
	}

	storedList := &StoredList{
		UserID: "testUserID",
		Data:   "testData",
		IV:     "testIV",
	}

	result, err := list.UpdateList(storedList)

	assert.Nil(t, err)
	assert.Equal(t, "testUserID", result.UserID)
	assert.Equal(t, "testData", result.Data)
	assert.Equal(t, "testIV", result.IV)
}

func TestList_DeleteList(t *testing.T) {
	mockClient := new(MockTodoServiceClient)
	mockClient.On("Delete", mock.Anything, mock.Anything).Return(&pb.TodoRetrieveResponse{
		UserId: "testUserID",
	}, nil)

	list := &List{
		Context: context.Background(),
		UserID:  "testUserID",
		Client:  mockClient,
	}

	result, err := list.DeleteList("testUserID")

	assert.Nil(t, err)
	assert.Equal(t, "testUserID", result.UserID)
	assert.Equal(t, "", result.Data)
	assert.Equal(t, "", result.IV)
}
