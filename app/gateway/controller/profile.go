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

func (impl *GatewayControllerImpl) Profile(ctx context.Context) (*user_s.User, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}
	return u, nil
}

func (impl *GatewayControllerImpl) ProfileUpdate(ctx context.Context, nu *user_s.User) error {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if ou == nil {
		impl.Logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	ou.FirstName = nu.FirstName
	ou.LastName = nu.LastName
	ou.Name = fmt.Sprintf("%s %s", nu.FirstName, nu.LastName)
	ou.LexicalName = fmt.Sprintf("%s, %s", nu.LastName, nu.FirstName)
	ou.Email = nu.Email
	ou.CompanyName = nu.CompanyName
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

	if err := impl.UserStorer.UpdateByID(ctx, ou); err != nil {
		impl.Logger.Error("user update by id error", slog.Any("error", err))
		return err
	}
	return nil
}
