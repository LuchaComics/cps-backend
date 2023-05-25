package controller

import (
	"context"

	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Verify(ctx context.Context, code string) error {
	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByVerificationCode(ctx, code)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return httperror.NewForBadRequestWithSingleField("code", "does not exist")
	}

	//TODO: Handle expiry dates.

	// Verify the user.
	u.WasEmailVerified = true
	if err := impl.UserStorer.UpdateByID(ctx, u); err != nil {
		impl.Logger.Error("update error", slog.Any("err", err))
		return err
	}

	return nil
}
