package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
)

type Cacher interface {
	Shutdown()
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, val []byte) error
	SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error
}

type cache struct {
	Client *redis.Client
	Logger *slog.Logger
}

func NewCache(cfg *c.Conf, logger *slog.Logger) Cacher {
	logger.Debug("cache initializing...")
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Cache.Host, cfg.Cache.Port),
		Password: cfg.Cache.Password, // no password set
		DB:       0,                  // use default DB
	})

	// Confirm connection made with our application and redis.
	response, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		logger.Error("cache failed pinging redis", slog.Any("err", err))
		panic("")
	}

	logger.Debug("cache initialized with successful redis ping response", slog.Any("resp", response))
	return &cache{
		Client: rdb,
		Logger: logger,
	}
}

func (s *cache) Shutdown() {
	s.Client.Close()
}

func (s *cache) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := s.Client.Get(ctx, key).Result()
	if err != nil {
		s.Logger.Error("cache get failed", slog.Any("error", err))
		return nil, err
	}
	return []byte(val), nil
}

func (s *cache) Set(ctx context.Context, key string, val []byte) error {
	err := s.Client.Set(ctx, key, val, 0).Err()
	if err != nil {
		s.Logger.Error("cache set failed", slog.Any("error", err))
		return err
	}
	return nil
}

func (s *cache) SetWithExpiry(ctx context.Context, key string, val []byte, expiry time.Duration) error {
	err := s.Client.Set(ctx, key, val, expiry).Err()
	if err != nil {
		s.Logger.Error("cache set with expiry failed", slog.Any("error", err))
		return err
	}
	return nil
}
