package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	s_d "github.com/LuchaComics/cps-backend/app/submission/datastore"
	u_d "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config/constants"
)

func (c *SubmissionControllerImpl) Create(ctx context.Context, m *s_d.Submission) error {
	// Modify the submission based on role.
	userRole, ok := ctx.Value(constants.SessionUserRole).(int8)
	if ok {
		switch userRole {
		case u_d.RetailerRole:
			// Override state.
			m.State = s_d.SubmissionPendingState

			// Auto-assign the user-if
			m.UserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
			m.UserFirstName = ctx.Value(constants.SessionUserFirstName).(string)
			m.UserLastName = ctx.Value(constants.SessionUserLastName).(string)
			m.UserCompanyName = ctx.Value(constants.SessionUserCompanyName).(string)
			m.ServiceType = s_d.PreScreeningServiceType
		case u_d.StaffRole:
			m.State = s_d.SubmissionActiveState
		default:
			m.State = s_d.SubmissionErrorState
		}
	}

	// Add defaults.
	m.CreatedTime = time.Now()
	m.ModifiedTime = time.Now()

	// Save to our database.
	err := c.SubmissionStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return err
	}
	return err
}
