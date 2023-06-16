package controller

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

type UserUpdateRequestIDO struct {
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

func (impl *UserControllerImpl) userFromUpdateRequest(requestData *UserUpdateRequestIDO) (*user_s.User, error) {
	passwordHash, err := impl.Password.GenerateHashFromPassword(requestData.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	return &user_s.User{
		OrganizationID:            requestData.OrganizationID,
		FirstName:                 requestData.FirstName,
		LastName:                  requestData.FirstName,
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

func (impl *UserControllerImpl) UpdateByID(ctx context.Context, requestData *UserUpdateRequestIDO) (*user_s.User, error) {
	nu, err := impl.userFromUpdateRequest(requestData)
	if err != nil {
		return nil, err
	}

	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := ctx.Value(constants.SessionUserName).(string)

	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.StaffRole {
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := impl.UserStorer.GetByID(ctx, nu.ID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	// Lookup the organization in our database, else return a `400 Bad Request` error.
	o, err := impl.OrganizationStorer.GetByID(ctx, nu.OrganizationID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if o == nil {
		impl.Logger.Warn("organization does not exist exists validation error")
		return nil, httperror.NewForBadRequestWithSingleField("organization_id", "organization does not exist")
	}

	ou.OrganizationID = o.ID
	ou.OrganizationName = o.Name
	ou.FirstName = nu.FirstName
	ou.LastName = nu.LastName
	ou.Name = fmt.Sprintf("%s %s", nu.FirstName, nu.LastName)
	ou.LexicalName = fmt.Sprintf("%s, %s", nu.LastName, nu.FirstName)
	ou.Email = nu.Email
	ou.Phone = nu.Phone
	ou.Country = nu.Country
	ou.Region = nu.Region
	ou.City = nu.City
	ou.PostalCode = nu.PostalCode
	ou.AddressLine1 = nu.AddressLine1
	ou.AddressLine2 = nu.AddressLine2
	ou.HowDidYouHearAboutUs = nu.HowDidYouHearAboutUs
	ou.HowDidYouHearAboutUsOther = nu.HowDidYouHearAboutUsOther
	ou.AgreePromotionsEmail = nu.AgreePromotionsEmail
	ou.ModifiedByUserID = userID
	ou.ModifiedByName = userName

	if err := impl.UserStorer.UpdateByID(ctx, ou); err != nil {
		impl.Logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}
	return ou, nil
}
