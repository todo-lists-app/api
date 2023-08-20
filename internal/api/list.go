package api

import (
	"context"
	"github.com/bugfixes/go-bugfixes/logs"
	pb "github.com/todo-lists-app/protobufs/generated/todo/v1"
	"github.com/todo-lists-app/todo-lists-api/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type TodoServiceClientCreator interface {
	NewTodoServiceClient() pb.TodoServiceClient
}

// List is the list service
type List struct {
	config.Config
	context.Context
	UserID string
	Client pb.TodoServiceClient // Add this line
}

type TodoList interface {
	GetList() (*StoredList, error)
	UpdateList(list *StoredList) (*StoredList, error)
	DeleteList(id string) (*StoredList, error)
	CreateList(list *StoredList) (*StoredList, error)
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

func (l *List) GetClient() (*List, error) {
	conn, err := grpc.DialContext(l.Context, l.Config.Services.Todo, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, logs.Errorf("error dialing grpc: %v", err)
	}

	l.Client = pb.NewTodoServiceClient(conn)
	return l, nil
}

// GetList gets a list for the user
func (l *List) GetList() (*StoredList, error) {
	resp, err := l.Client.Get(l.Context, &pb.TodoGetRequest{
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
	resp, err := l.Client.Update(l.Context, &pb.TodoInjectRequest{
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
		UserID: list.UserID,
		Data:   list.Data,
		IV:     list.IV,
	}, nil
}

// DeleteList deletes a list for the user
func (l *List) DeleteList(id string) (*StoredList, error) {
	resp, err := l.Client.Delete(l.Context, &pb.TodoDeleteRequest{
		UserId: l.UserID,
	})
	if err != nil {
		return nil, logs.Errorf("error deleting list: %v", err)
	}
	if resp.GetStatus() != "" {
		return nil, logs.Errorf("error deleting list status: %v", resp.GetStatus())
	}

	return &StoredList{
		UserID: id,
	}, nil
}

// CreateList creates a new list for the user
func (l *List) CreateList(list *StoredList) (*StoredList, error) {
	resp, err := l.Client.Insert(l.Context, &pb.TodoInjectRequest{
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
