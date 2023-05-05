package controller

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/exp/slog"

	"github.com/LuchaComics/cps-backend/adapter/cache/redis"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/jwt"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

type GatewayController interface {
	Login(ctx context.Context, email string, password string) (*user_s.User, string, time.Time, string, time.Time, error)
	GetUserBySessionID(ctx context.Context, sessionID string) (*user_s.User, error)
	RefreshToken(ctx context.Context, value string) (*user_s.User, string, time.Time, string, time.Time, error)
	//TODO: Add more...
}

type GatewayControllerImpl struct {
	Config     *config.Conf
	Logger     *slog.Logger
	UUID       uuid.Provider
	JWT        jwt.Provider
	Password   password.Provider
	Cache      redis.Cacher
	UserStorer user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	jwtp jwt.Provider,
	passwordp password.Provider,
	cache redis.Cacher,
	usr_storer user_s.UserStorer,
) GatewayController {
	s := &GatewayControllerImpl{
		Config:     appCfg,
		Logger:     loggerp,
		UUID:       uuidp,
		JWT:        jwtp,
		Password:   passwordp,
		Cache:      cache,
		UserStorer: usr_storer,
	}

	return s
}

func (impl *GatewayControllerImpl) GetUserBySessionID(ctx context.Context, sessionID string) (*user_s.User, error) {
	userBytes, err := impl.Cache.Get(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if userBytes == nil {
		impl.Logger.Warn("record not found")
		return nil, errors.New("record not found")
	}
	var user user_s.User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		impl.Logger.Error("unmarshalling failed", slog.Any("err", err))
		return nil, err
	}
	return &user, nil
}
