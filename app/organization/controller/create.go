package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/organization/datastore"
)

func (c *OrganizationControllerImpl) Create(ctx context.Context, m *s_d.Organization) (*s_d.Organization, error) {
	// // Modify the organization based on role.
	// userRole, ok := ctx.Value(constants.SessionUserRole).(int8)
	// if ok {
	// 	switch userRole {
	// 	case u_d.RetailerRole:
	// 		// Override state.
	// 		m.State = s_d.OrganizationPendingState
	//
	// 		// Auto-assign the user-if
	// 		m.UserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// 		m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
	// 		m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)
	// 		m.UserCompanyName = ctx.Value(constants.SessionUserCompanyName).(string)
	// 		m.ServiceType = s_d.PreScreeningServiceType
	// 	case u_d.StaffRole:
	// 		m.State = s_d.OrganizationActiveState
	// 	default:
	// 		m.State = s_d.OrganizationErrorState
	// 	}
	// }

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	m.ModifiedAt = time.Now()

	// Save to our database.
	err := c.OrganizationStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	return m, nil
}
