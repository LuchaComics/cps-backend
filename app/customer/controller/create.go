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

func (impl *CustomerControllerImpl) Create(ctx context.Context, m *user_s.User) (*user_s.User, error) {
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

	// Modify the customer based on role.
	orgID, ok := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	if !ok {
		impl.Logger.Error("wrong format error")
		return nil, fmt.Errorf("%v", "wrong format in user organization id")
	}

	// Add defaults.
	m.Email = strings.ToLower(m.Email)
	m.OrganizationID = orgID
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()
	m.Role = user_s.RetailerCustomerRole
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