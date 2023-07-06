package api

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/todo/v1"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// List is the list service
type List struct {
	config.Config
	context.Context
	UserID string
}

// NewListService creates a new list service
func NewListService(ctx context.Context, cfg config.Config, id string) *List {
	return &List{
		Config:  cfg,
		Context: ctx,
		UserID:  id,
	}
}

// StoredList is the stored list
type StoredList struct {
	UserID string `bson:"userid" json:"userid"`
	Data   string `bson:"data" json:"data"`
	IV     string `bson:"iv" json:"iv"`
}

// GetList gets a list for the user
func (l *List) GetList() (*StoredList, error) {
	conn, err := grpc.DialContext(l.Context, l.Config.Services.Todo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}
	defer conn.Close()
	g := pb.NewTodoServiceClient(conn)
	resp, err := g.Get(l.Context, &pb.TodoGetRequest{
		UserId: l.UserID,
	})
	if err != nil {
		return nil, logs.Errorf("error getting list: %v", err)
	}
	if resp.GetStatus() != "" {
		return nil, logs.Errorf("error getting list status: %v", resp.GetStatus())
	}

	return &StoredList{
		UserID: resp.GetUserId(),
		Data:   resp.GetData(),
		IV:     resp.GetIv(),
	}, nil
}

// UpdateList updates a list for the user
func (l *List) UpdateList(list *StoredList) (*StoredList, error) {
	conn, err := grpc.DialContext(l.Context, l.Config.Services.Todo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}
	defer conn.Close()
	g := pb.NewTodoServiceClient(conn)

	resp, err := g.Update(l.Context, &pb.TodoInjectRequest{
		UserId: l.UserID,
		Data:   list.Data,
		Iv:     list.IV,
	})
	if err != nil {
		return nil, logs.Errorf("error updating list: %v", err)
	}
	if resp.GetStatus() != "" {
		return nil, logs.Errorf("error updating list status: %v", resp.GetStatus())
	}

	return &StoredList{
		UserID: resp.GetUserId(),
		Data:   resp.GetData(),
		IV:     resp.GetIv(),
	}, nil
}

// DeleteList deletes a list for the user
func (l *List) DeleteList(id string) (*StoredList, error) {
	client, err := config.GetMongoClient(l.Context, l.Config)
	if err != nil {
		return nil, logs.Errorf("error getting mongo client: %v", err)
	}
	defer func() {
		if err := client.Disconnect(l.Context); err != nil {
			logs.Errorf("error disconnecting mongo client: %v", err)
		}
	}()

	return &StoredList{
		UserID: id,
	}, nil
}

// CreateList creates a new list for the user
func (l *List) CreateList(list *StoredList) (*StoredList, error) {
	conn, err := grpc.DialContext(l.Context, l.Config.Services.Todo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}
	defer conn.Close()
	g := pb.NewTodoServiceClient(conn)
	resp, err := g.Insert(l.Context, &pb.TodoInjectRequest{
		UserId: l.UserID,
		Data:   list.Data,
		Iv:     list.IV,
	})
	if err != nil {
		return nil, logs.Errorf("error inserting list: %v", err)
	}
	if resp.GetStatus() != "" {
		return nil, logs.Errorf("error inserting list status: %v", resp.GetStatus())
	}

	return &StoredList{
		UserID: resp.GetUserId(),
		Data:   resp.GetData(),
		IV:     resp.GetIv(),
	}, nil
}
