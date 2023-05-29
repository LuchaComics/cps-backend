package datastore

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

const (
	UserActiveState      = 1
	UserArchivedState    = 100
	StaffRole            = 1
	RetailerStaffRole    = 2
	RetailerCustomerRole = 3
)

type User struct {
	ID                        primitive.ObjectID `bson:"_id" json:"_id"`
	OrganizationID            primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	Email                     string             `bson:"email" json:"email"`
	PasswordHashAlgorithm     string             `bson:"password_hash_algorithm" json:"password_hash_algorithm,omitempty"`
	PasswordHash              string             `bson:"password_hash" json:"password_hash,omitempty"`
	Role                      int8               `bson:"role" json:"role"`
	WasEmailVerified          bool               `bson:"was_email_verified" json:"was_email_verified"`
	EmailVerificationCode     string             `bson:"email_verification_code,omitempty" json:"email_verification_code,omitempty"`
	EmailVerificationExpiry   time.Time          `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
	Phone                     string             `bson:"phone" json:"phone,omitempty"`
	Country                   string             `bson:"country" json:"country,omitempty"`
	Region                    string             `bson:"region" json:"region,omitempty"`
	City                      string             `bson:"city" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line_1" json:"address_line_1,omitempty"`
	AddressLine2              string             `bson:"address_line_2" json:"address_line_2,omitempty"`
	StoreLogoS3Key            string             `bson:"store_logo_s3_key" json:"store_logo_s3_key,omitempty"`
	StoreLogoTitle            string             `bson:"store_logo_title" json:"store_logo_title,omitempty"`
	StoreLogoFileURL          string             `bson:"store_logo_file_url" json:"store_logo_file_url,omitempty"`     // (Optional, added by endpoint)
	StoreLogoFileURLExpiry    time.Time          `bson:"store_logo_file_url_expiry" json:"store_logo_file_url_expiry"` // (Optional, added by endpoint)
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `bson:"agree_tos" json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `bson:"agree_promotions_email" json:"agree_promotions_email,omitempty"`
	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	State                     int8               `bson:"state" json:"state"`
}

type UserListFilter struct {
	OrganizationID  primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	Role            int8               `bson:"role" json:"role"`
	SortOrder       string             `json:"sort_order"`
	SortField       string             `json:"sort_field"`
	Offset          uint64             `json:"offset"`
	Limit           uint64             `json:"limit"`
	States          []int8             `json:"states"`
	UUIDs           []string           `json:"uuids"`
	ExcludeArchived bool               `json:"exclude_archived"`
	SearchText      string             `json:"search_text"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Email           string             `json:"email"`
	Phone           string             `json:"phone"`
}

type UserListResult struct {
	Results []*User `json:"results"`
}

// UserStorer Interface for user.
type UserStorer interface {
	Create(ctx context.Context, m *User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *User) error
	ListByFilter(ctx context.Context, f *UserListFilter) (*UserListResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
			{"email", "text"},
			{"lexical_name", "text"},
		},
	}
	name, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		panic(err)
	}
	fmt.Println("Name of Index Created: " + name)

	s := &UserStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
