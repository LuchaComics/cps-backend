package datastore

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

const (
	UserActiveState   = 1
	UserInactiveState = 2
	StaffRole         = 1
	RetailerRole      = 2
)

type User struct {
	ID                        primitive.ObjectID `bson:"_id" json:"_id"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	Email                     string             `bson:"email" json:"email"`
	PasswordHashAlgorithm     string             `bson:"password_hash_algorithm,omitempty" json:"password_hash_algorithm,omitempty"`
	PasswordHash              string             `bson:"password_hash,omitempty" json:"password_hash,omitempty"`
	Role                      int8               `bson:"role" json:"role"`
	WasEmailVerified          bool               `bson:"was_email_verified" json:"was_email_verified"`
	EmailVerificationCode     string             `bson:"email_verification_code,omitempty" json:"email_verification_code,omitempty"`
	EmailVerificationExpiry   time.Time          `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
	CompanyName               string             `bson:"company_name,omitempty" json:"company_name,omitempty"`
	Phone                     string             `bson:"phone,omitempty" json:"phone,omitempty"`
	Country                   string             `bson:"country,omitempty" json:"country,omitempty"`
	Region                    string             `bson:"region,omitempty" json:"region,omitempty"`
	City                      string             `bson:"city,omitempty" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code,omitempty" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line_1,omitempty" json:"address_line_1,omitempty"`
	AddressLine2              string             `bson:"address_line_2,omitempty" json:"address_line_2,omitempty"`
	StoreLogoS3Key            string             `bson:"store_logo_s3_key,omitempty" json:"store_logo_s3_key,omitempty"`
	StoreLogoTitle            string             `bson:"store_logo_title,omitempty" json:"store_logo_title,omitempty"`
	StoreLogoFileURL          string             `bson:"store_logo_file_url,omitempty" json:"store_logo_file_url,omitempty"`     // (Optional, added by endpoint)
	StoreLogoFileURLExpiry    time.Time          `bson:"store_logo_file_url_expiry,omitempty" json:"store_logo_file_url_expiry"` // (Optional, added by endpoint)
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us,omitempty" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other,omitempty" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `bson:"agree_tos,omitempty" json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `bson:"agree_promotions_email,omitempty" json:"agree_promotions_email,omitempty"`
	CreatedTime               time.Time          `bson:"created_time,omitempty" json:"created_time,omitempty"`
	ModifiedTime              time.Time          `bson:"modified_time,omitempty" json:"modified_time,omitempty"`
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
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *User) error
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
