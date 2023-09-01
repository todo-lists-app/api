package api

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/user/v1"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountServiceClientCreator interface {
	NewAccountServiceClient() pb.UserServiceClient
}

type Account struct {
	config.Config
	context.Context
	UserID      string
	AccessToken string
	Client      pb.UserServiceClient
}

func NewAccountService(ctx context.Context, cfg config.Config, id, accessToken string) *Account {
	return &Account{
		Config:      cfg,
		Context:     ctx,
		UserID:      id,
		AccessToken: accessToken,
	}
}

func (a *Account) GetClient() (*Account, error) {
	conn, err := grpc.DialContext(a.Context, a.Config.Services.User, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}

	a.Client = pb.NewUserServiceClient(conn)
	return a, nil
}

func (a *Account) DeleteAccount() error {
	resp, err := a.Client.Delete(a.Context, &pb.UserDeleteRequest{
		UserId:      a.UserID,
		AccessToken: a.AccessToken,
	})
	if err != nil {
		return logs.Errorf("error deleting account: %v", err)
	}

	if resp.GetStatus() != "ok" {
		return logs.Errorf("error deleting account: %v", resp.GetStatus())
	}

	return nil
}
