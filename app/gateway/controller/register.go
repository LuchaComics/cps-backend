package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/exp/slog"

	gateway_s "github.com/LuchaComics/cps-backend/app/gateway/datastore"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Register(ctx context.Context, req *gateway_s.RegisterRequestIDO) error {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	req.Email = strings.ToLower(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByEmail(ctx, req.Email)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return err
	}
	if u != nil {
		impl.Logger.Warn("user already exists validation error")
		return httperror.NewForBadRequestWithSingleField("email", "email is not unique")
	}

	passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return err
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
		WasEmailVerified:          false,
		EmailVerificationCode:     impl.UUID.NewUUID(),
		EmailVerificationExpiry:   time.Now().Add(72 * time.Hour),
	}
	err = impl.UserStorer.Create(ctx, u)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return err
	}
	impl.Logger.Info("User created.",
		slog.Any("_id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	// uBin, err := json.Marshal(u)
	// if err != nil {
	// 	impl.Logger.Error("marshalling error", slog.Any("err", err))
	// 	return nil, err
	// }

	if err := impl.SendVerificationEmail(u.Email, u.EmailVerificationCode, u.FirstName); err != nil {
		impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
		return err
	}

	return nil
}
