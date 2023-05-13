package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	gateway_s "github.com/LuchaComics/cps-backend/app/gateway/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Register(ctx context.Context, req *gateway_s.RegisterRequestIDO) (*gateway_s.RegisterResponseIDO, error) {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	req.Email = strings.ToLower(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByEmail(ctx, req.Email)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u != nil {
		impl.Logger.Warn("user already exists validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
	}

	passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	//TODO: Handle s3.

	u = &user_s.User{
		ID:                        primitive.NewObjectID(),
		FirstName:                 req.FirstName,
		LastName:                  req.LastName,
		Name:                      fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:               fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                     req.Email,
		PasswordHash:              passwordHash,
		PasswordHashAlgorithm:     impl.Password.AlgorithmName(),
		Role:                      user_s.RetailerRole,
		CompanyName:               req.CompanyName,
		Phone:                     req.Phone,
		Country:                   req.Country,
		Region:                    req.Region,
		City:                      req.City,
		PostalCode:                req.PostalCode,
		AddressLine1:              req.AddressLine1,
		HowDidYouHearAboutUs:      req.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther: req.HowDidYouHearAboutUsOther,
		AgreeTOS:                  req.AgreeTOS,
		AgreePromotionsEmail:      req.AgreePromotionsEmail,
		CreatedTime:               time.Now(),
		ModifiedTime:              time.Now(),
	}
	err = impl.UserStorer.Create(ctx, u)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	impl.Logger.Info("User created.",
		slog.Any("_id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	uBin, err := json.Marshal(u)
	if err != nil {
		impl.Logger.Error("marshalling error", slog.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	sessionUUID := impl.UUID.NewUUID()

	err = impl.Cache.SetWithExpiry(ctx, sessionUUID, uBin, rtExpiry)
	if err != nil {
		impl.Logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := impl.JWT.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, err
	}

	// For security.
	u.PasswordHash = ""
	u.PasswordHashAlgorithm = ""

	// Return our auth keys.
	return &gateway_s.RegisterResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}
