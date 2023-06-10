package controller

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (impl *CustomerControllerImpl) UpdateByID(ctx context.Context, nu *user_s.User) (*user_s.User, error) {
	// Extract from our session the following data.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName, _ := ctx.Value(constants.SessionUserOrganizationName).(string)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	ou.OrganizationID = orgID
	ou.OrganizationName = orgName
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

	if err := impl.UserStorer.UpdateByID(ctx, ou); err != nil {
		impl.Logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}
	return ou, nil
}

func (c *CustomerControllerImpl) CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error) {
	// Fetch the original customer.
	s, err := c.UserStorer.GetByID(ctx, customerID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if s == nil {
		return nil, nil
	}

	// Create our comment.
	comment := &user_s.UserComment{
		ID:               primitive.NewObjectID(),
		Content:          content,
		OrganizationID:   ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID),
		CreatedByUserID:  ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		CreatedByName:    ctx.Value(constants.SessionUserName).(string),
		CreatedAt:        time.Now(),
		ModifiedByUserID: ctx.Value(constants.SessionUserID).(primitive.ObjectID),
		ModifiedByName:   ctx.Value(constants.SessionUserName).(string),
		ModifiedAt:       time.Now(),
	}

	// Add our comment to the comments.
	s.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	s.ModifiedAt = time.Now()
	s.Comments = append(s.Comments, comment)

	// Save to the database the modified customer.
	if err := c.UserStorer.UpdateByID(ctx, s); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return s, nil
}
