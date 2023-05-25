package controller

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	domain "github.com/LuchaComics/cps-backend/app/user/datastore"
)

// CreateInitialRootAdmin function creates the initial root administrator if not previously created.
func (c *UserControllerImpl) CreateInitialRootAdmin(ctx context.Context) error {
	doesExist, err := c.UserStorer.CheckIfExistsByEmail(ctx, c.Config.AppServer.InitialAdminEmail)
	if err != nil {
		c.Logger.Error("database check if exists error", slog.Any("error", err))
		return err
	}
	if doesExist == false {
		c.Logger.Info("No root user detected, proceeding to create now...")
		passwordHash, err := c.Password.GenerateHashFromPassword(c.Config.AppServer.InitialAdminPassword)
		if err != nil {
			c.Logger.Error("hashing error", slog.Any("error", err))
			return err
		}
		m := &domain.User{
			ID:                    primitive.NewObjectID(),
			FirstName:             "Root",
			LastName:              "Administrator",
			Name:                  "Root Administrator",
			LexicalName:           "Administrator, Root",
			Email:                 c.Config.AppServer.InitialAdminEmail,
			PasswordHash:          passwordHash,
			PasswordHashAlgorithm: c.Password.AlgorithmName(),
			Role:                  domain.StaffRole,
			WasEmailVerified:      true,
			CreatedAt:           time.Now(),
			ModifiedAt:          time.Now(),
		}
		err = c.UserStorer.Create(ctx, m)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Info("Root user created.",
			slog.Any("_id", m.ID),
			slog.String("name", m.Name),
			slog.String("email", m.Email),
			slog.String("password_hash_algorithm", m.PasswordHashAlgorithm),
			slog.String("password_hash", m.PasswordHash))
	} else {
		c.Logger.Info("Root user already exists, skipping creation.")
	}
	return nil
}
