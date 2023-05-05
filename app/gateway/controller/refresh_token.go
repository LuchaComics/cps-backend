package controller

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
)

func (impl *GatewayControllerImpl) RefreshToken(ctx context.Context, value string) (*user_s.User, string, time.Time, string, time.Time, error) {
	////
	//// Extract the `sessionID` so we can process it.
	////

	sessionID, err := impl.JWT.ProcessJWTToken(value)
	if err != nil {
		impl.Logger.Warn("process jwt refresh token does not exist", slog.String("value", value))
		err := errors.New("jwt refresh token failed")
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Lookup in our in-memory the user record for the `sessionID` or error.
	////

	uBin, err := impl.Cache.Get(ctx, sessionID)
	if err != nil {
		impl.Logger.Error("in-memory set error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	var u *user_s.User
	err = json.Unmarshal(uBin, &u)
	if err != nil {
		impl.Logger.Error("unmarshal error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	////
	//// Generate new access and refresh tokens and return them.
	////

	// Start our session using an access and refresh token.
	newSessionUUID := impl.UUID.NewUUID()

	err = impl.Cache.Set(ctx, newSessionUUID, uBin)
	if err != nil {
		impl.Logger.Error("in-memory set error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := impl.JWT.GenerateJWTTokenPair(newSessionUUID, 100*time.Minute, 200*time.Minute)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Return our auth keys.
	return u, accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, nil
}
