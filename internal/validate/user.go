package validate

import (
	"context"
	pb "github.com/todo-lists-app/protobufs/generated/id_checker/v1"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"google.golang.org/grpc"
)

type Validate struct {
	Config *config.Config
	CTX    context.Context
}

func NewValidate(config *config.Config, ctx context.Context) *Validate {
	return &Validate{
		Config: config,
		CTX:    ctx,
	}
}

func (v *Validate) ValidateUser(userId string) (bool, error) {
	if v.Config.Development {
		return true, nil
	}

	conn, err := grpc.DialContext(v.CTX, v.Config.Services.Identity, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	defer conn.Close()

	g := pb.NewIdCheckerServiceClient(conn)
	resp, err := g.CheckId(v.CTX, &pb.CheckIdRequest{
		Id: userId,
	})
	if err != nil {
		return false, err
	}
	if !resp.GetIsValid() {
		return false, nil
	}

	return true, nil
}
