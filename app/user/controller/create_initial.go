package controller

import (
	"context"

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
			UserID:                c.UUID.NewUUID(),
			Name:                  "Root Administrator",
			Email:                 c.Config.AppServer.InitialAdminEmail,
			PasswordHash:          passwordHash,
			PasswordHashAlgorithm: c.Password.AlgorithmName(),
		}
		err = c.UserStorer.Create(ctx, m)
		if err != nil {
			c.Logger.Error("database create error", slog.Any("error", err))
			return err
		}
		c.Logger.Info("Root user created.",
			slog.String("user_id", m.UserID),
			slog.String("name", m.Name),
			slog.String("email", m.Email),
			slog.String("password_hash_algorithm", m.PasswordHashAlgorithm),
			slog.String("password_hash", m.PasswordHash))
	} else {
		c.Logger.Info("Root user already exists, skipping creation.")
	}

	// #1 FOR TESTING PURPOSES ONLY
	user, err := c.UserStorer.GetByEmail(ctx, c.Config.AppServer.InitialAdminEmail)
	if err != nil {
		c.Logger.Error("GetByEmail.",
			slog.Any("error", err))
		return err
	}
	c.Logger.Info("retrieved user by email.",
		slog.String("user_id", user.UserID),
		slog.String("name", user.Name),
		slog.String("email", user.Email),
		slog.String("password_hash_algorithm", user.PasswordHashAlgorithm),
		slog.String("password_hash", user.PasswordHash))

	// #2 FOR TESTING PURPOSES ONLY
	user.Name = "Root"
	if err := c.UserStorer.UpdateByUserID(ctx, user); err != nil {
		c.Logger.Error("Update.",
			slog.Any("error", err))
		return err
	}

	return nil
}
