package controller

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"golang.org/x/exp/slog"

	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
)

func (impl *GatewayControllerImpl) Login(ctx context.Context, email, password string) (*user_s.User, string, time.Time, string, time.Time, error) {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	email = strings.ToLower(email)
	password = strings.ReplaceAll(password, " ", "")

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByEmail(ctx, email)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Verify the inputted password and hashed password match.
	passwordMatch, _ := impl.Password.ComparePasswordAndHash(password, u.PasswordHash)
	if passwordMatch == false {
		impl.Logger.Warn("password check validation error")
		return nil, "", time.Now(), "", time.Now(), errors.New("password do not match with record")
	}

	uBin, err := json.Marshal(u)
	if err != nil {
		impl.Logger.Error("marshalling error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Start our session using an access and refresh token.
	sessionUUID := impl.UUID.NewUUID()

	err = impl.Cache.Set(ctx, sessionUUID, uBin)
	if err != nil {
		impl.Logger.Error("in-memory set error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := impl.JWT.GenerateJWTTokenPair(sessionUUID, 24*time.Hour, 14*24*time.Hour)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, "", time.Now(), "", time.Now(), err
	}

	// Return our auth keys.
	return u, accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, nil
}
