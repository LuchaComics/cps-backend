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

func (impl *UserControllerImpl) Create(ctx context.Context, m *user_s.User) (*user_s.User, error) {
	// Extract from our session the following data.
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply filtering based on ownership and role.
	if userRole != user_s.StaffRole {
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

	// Modify the user based on role.
	orgID, _ := ctx.Value(constants.SessionUserOrganizationID).(primitive.ObjectID)
	orgName, _ := ctx.Value(constants.SessionUserOrganizationName).(string)

	// Add defaults.
	m.Email = strings.ToLower(m.Email)
	m.OrganizationID = orgID
	m.OrganizationName = orgName
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()
	// m.Role = user_s.RetailerUserRole
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
