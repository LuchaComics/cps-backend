package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

type UserCreateRequestIDO struct {
	OrganizationID            primitive.ObjectID `bson:"organization_id" json:"organization_id,omitempty"`
	FirstName                 string             `json:"first_name"`
	LastName                  string             `json:"last_name"`
	Email                     string             `json:"email"`
	Password                  string             `json:"password"`
	PasswordRepeated          string             `json:"password_repeated"`
	Phone                     string             `json:"phone,omitempty"`
	Country                   string             `json:"country,omitempty"`
	Region                    string             `json:"region,omitempty"`
	City                      string             `json:"city,omitempty"`
	PostalCode                string             `json:"postal_code,omitempty"`
	AddressLine1              string             `json:"address_line_1,omitempty"`
	AddressLine2              string             `json:"address_line_2,omitempty"`
	HowDidYouHearAboutUs      int8               `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTOS                  bool               `json:"agree_tos,omitempty"`
	AgreePromotionsEmail      bool               `json:"agree_promotions_email,omitempty"`
	Status                    int8               `bson:"status" json:"status"`
	Role                      int8               `bson:"role" json:"role"`
}

func (impl *UserControllerImpl) userFromCreateRequest(requestData *UserCreateRequestIDO) (*user_s.User, error) {
	passwordHash, err := impl.Password.GenerateHashFromPassword(requestData.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	return &user_s.User{
		OrganizationID:            requestData.OrganizationID,
		FirstName:                 requestData.FirstName,
		LastName:                  requestData.LastName,
		Email:                     requestData.Email,
		PasswordHash:              passwordHash,
		PasswordHashAlgorithm:     impl.Password.AlgorithmName(),
		Phone:                     requestData.Phone,
		Country:                   requestData.Country,
		Region:                    requestData.Region,
		City:                      requestData.City,
		PostalCode:                requestData.PostalCode,
		AddressLine1:              requestData.AddressLine1,
		AddressLine2:              requestData.AddressLine2,
		HowDidYouHearAboutUs:      requestData.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: requestData.HowDidYouHearAboutUsOther,
		AgreeTOS:                  requestData.AgreeTOS,
		AgreePromotionsEmail:      requestData.AgreePromotionsEmail,
		Status:                    requestData.Status,
		Role:                      requestData.Role,
	}, nil
}

func (impl *UserControllerImpl) Create(ctx context.Context, requestData *UserCreateRequestIDO) (*user_s.User, error) {
	m, err := impl.userFromCreateRequest(requestData)
	if err != nil {
		return nil, err
	}

	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	// Apply filtering based on ownership and role.
	if userRole != user_s.UserRoleRoot {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByEmail(ctx, m.Email)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u != nil {
		impl.Logger.Warn("user already exists validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
	}

	// Lookup the organization in our database, else return a `400 Bad Request` error.
	o, err := impl.OrganizationStorer.GetByID(ctx, m.OrganizationID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if o == nil {
		impl.Logger.Warn("organization does not exist exists validation error")
		return nil, httperror.NewForBadRequestWithSingleField("organization_id", "organization does not exist")
	}

	// Modify the user based on role.

	// Add defaults.
	m.Email = strings.ToLower(m.Email)
	m.OrganizationID = o.ID
	m.OrganizationName = o.Name
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.CreatedByUserID = userID
	m.CreatedByName = userName
	m.ModifiedAt = time.Now()
	m.ModifiedByUserID = userID
	m.ModifiedByName = userName
	m.Name = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
	m.LexicalName = fmt.Sprintf("%s, %s", m.LastName, m.FirstName)
	m.WasEmailVerified = true

	// Save to our database.
	if err := impl.UserStorer.Create(ctx, m); err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	return m, nil
}
