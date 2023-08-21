package validate

import (
	"context"

	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/id_checker/v1"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type IdCheckerServiceClientCreator interface {
	NewIdCheckerServiceClient() pb.IdCheckerServiceClient
}

type Validate struct {
	Config *config.Config
	CTX    context.Context
	Client pb.IdCheckerServiceClient
}

type Checker interface {
	ValidateUser(accessToken, userId string) (bool, error)
}

func NewValidate(config *config.Config, ctx context.Context) *Validate {
	return &Validate{
		Config: config,
		CTX:    ctx,
	}
}

func (v *Validate) GetClient() (*Validate, error) {
	conn, err := grpc.DialContext(v.CTX, v.Config.Services.Identity, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}

	v.Client = pb.NewIdCheckerServiceClient(conn)
	return v, nil
}

func (v *Validate) ValidateUser(accessToken, userId string) (bool, error) {
	if v.Config != nil && v.Config.Config.Local.Development {
		return true, nil
	}

	resp, err := v.Client.CheckId(v.CTX, &pb.CheckIdRequest{
		Id:          userId,
		AccessToken: accessToken,
	})
	if err != nil {
		_ = logs.Errorf("error checking id: %v, %s", err, userId)
		return false, nil
	}
	if !resp.GetIsValid() {
		return false, nil
	}

	return true, nil
}
