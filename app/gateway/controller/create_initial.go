package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	org_d "github.com/LuchaComics/cps-backend/app/organization/datastore"
	domain "github.com/LuchaComics/cps-backend/app/user/datastore"
)

// createInitialRootAdmin function creates the initial root administrator if not previously created.
func (c *GatewayControllerImpl) createInitialRootAdmin(ctx context.Context) error {
	doesExist, err := c.UserStorer.CheckIfExistsByEmail(ctx, c.Config.AppServer.InitialAdminEmail)
	if err != nil {
		c.Logger.Error("database check if exists error", slog.Any("error", err))
		return err
	}
	if doesExist == false {
		c.Logger.Info("no root user detected, proceeding to create now...")
		passwordHash, err := c.Password.GenerateHashFromPassword(c.Config.AppServer.InitialAdminPassword)
		if err != nil {
			c.Logger.Error("hashing error", slog.Any("error", err))
			return err
		}
		usr := &domain.User{
			ID:                    primitive.NewObjectID(),
			FirstName:             "Root",
			LastName:              "Administrator",
			Name:                  "Root Administrator",
			LexicalName:           "Administrator, Root",
			Email:                 c.Config.AppServer.InitialAdminEmail,
			Status:                domain.UserStatusActive,
			PasswordHash:          passwordHash,
			PasswordHashAlgorithm: c.Password.AlgorithmName(),
			Role:                  domain.UserRoleRoot,
			WasEmailVerified:      true,
			CreatedAt:             time.Now(),
			ModifiedAt:            time.Now(),
			AgreeTOS:              true,
			AgreePromotionsEmail:  true,
		}
		err = c.UserStorer.Create(ctx, usr)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Info("Root user created.",
			slog.Any("id", usr.ID),
			slog.String("name", usr.Name),
			slog.String("email", usr.Email))

		// Create the organization.
		org := &org_d.Organization{
			ID:                 primitive.NewObjectID(),
			ModifiedAt:         time.Now(),
			ModifiedByUserName: usr.Name,
			ModifiedByUserID:   usr.ID,
			Type:               org_d.RootType,
			Status:             org_d.OrganizationActiveStatus,
			Name:               c.Config.AppServer.InitialAdminOrganizationName,
			CreatedAt:          time.Now(),
			CreatedByUserName:  usr.Name,
			CreatedByUserID:    usr.ID,
			Comments:           []*org_d.OrganizationComment{},
		}
		err = c.OrganizationStorer.Create(ctx, org)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Info("Root organization created.",
			slog.Any("id", org.ID),
			slog.String("name", org.Name))

		// Attach the user.
		usr.OrganizationID = org.ID
		usr.OrganizationName = org.Name
		err = c.UserStorer.UpdateByID(ctx, usr)
		if err != nil {
			c.Logger.Error("database update error", slog.Any("error", err))
			return err
		}
		c.Logger.Info("Root organization updated.",
			slog.Any("id", org.ID),
			slog.String("name", org.Name))

	} else {
		c.Logger.Info("Root user already exists, skipping creation.")
	}
	return nil
}
