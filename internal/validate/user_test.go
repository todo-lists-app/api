package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/todo-lists-app/protobufs/generated/id_checker/v1"
	"google.golang.org/grpc"
)

type MockIdCheckerServiceClient struct {
	mock.Mock
}

func (m *MockIdCheckerServiceClient) NewIdCheckerServiceClient() pb.IdCheckerServiceClient {
	return m
}

func (m *MockIdCheckerServiceClient) CheckId(ctx context.Context, in *pb.CheckIdRequest, opts ...grpc.CallOption) (*pb.CheckIdResponse, error) {
	args := m.Called(ctx, in)
	return args.Get(0).(*pb.CheckIdResponse), args.Error(1)
}

func TestValidate_ValidateUser(t *testing.T) {
	t.Run("valid user", func(t *testing.T) {
		mockClient := new(MockIdCheckerServiceClient)
		mockClient.On("CheckId", mock.Anything, mock.Anything).Return(&pb.CheckIdResponse{
			IsValid: true,
		}, nil)

		validate := &Validate{
			CTX:    context.Background(),
			Client: mockClient,
		}

		result, err := validate.ValidateUser("testAccessToken", "testUserID")

		assert.Nil(t, err)
		assert.Equal(t, true, result)
	})

	t.Run("invalid user", func(t *testing.T) {
		mockClient := new(MockIdCheckerServiceClient)
		mockClient.On("CheckId", mock.Anything, mock.Anything).Return(&pb.CheckIdResponse{
			IsValid: false,
		}, nil)

		validate := &Validate{
			CTX:    context.Background(),
			Client: mockClient,
		}

		result, err := validate.ValidateUser("testAccessToken", "testUserID")

		assert.Nil(t, err)
		assert.Equal(t, false, result)
	})
}
