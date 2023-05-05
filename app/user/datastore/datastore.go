package datastore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

const (
	UserActiveState   = 1
	UserInactiveState = 2
)

type User struct {
	ID                    primitive.ObjectID `bson:"_id"`
	UserID                string             `bson:"user_id"`
	Name                  string             `bson:"name"`
	Email                 string             `bson:"email"`
	PasswordHashAlgorithm string             `bson:"password_hash_algorithm"`
	PasswordHash          string             `bson:"password_hash"`
}

type UserFilter struct {
	SortOrder string   `json:"sort_order"`
	SortField string   `json:"sort_field"`
	Offset    uint64   `json:"offset"`
	Limit     uint64   `json:"limit"`
	States    []int8   `json:"states"`
	UUIDs     []string `json:"uuids"`
}

// UserStorer Interface for user.
type UserStorer interface {
	Create(ctx context.Context, m *User) error
	GetByUserID(ctx context.Context, userID string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByUserID(ctx context.Context, m *User) error
	// //TODO: Add more...
}

type UserStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) UserStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("users")

	s := &UserStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
