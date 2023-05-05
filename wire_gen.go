// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/LuchaComics/cps-backend/adapter/cache/redis"
	"github.com/LuchaComics/cps-backend/adapter/storage/mongodb"
	"github.com/LuchaComics/cps-backend/app/gateway/controller"
	controller2 "github.com/LuchaComics/cps-backend/app/user/controller"
	"github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http"
	"github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/cps-backend/inputport/http/user"
	"github.com/LuchaComics/cps-backend/provider/jwt"
	"github.com/LuchaComics/cps-backend/provider/logger"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/time"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

// Injectors from wire.go:

func InitializeEvent() Application {
	slogLogger := logger.NewProvider()
	conf := config.New()
	provider := uuid.NewProvider()
	timeProvider := time.NewProvider()
	jwtProvider := jwt.NewProvider(conf)
	passwordProvider := password.NewProvider()
	cacher := redis.NewCache(conf, slogLogger)
	client := mongodb.NewStorage(conf, slogLogger)
	userStorer := datastore.NewDatastore(conf, slogLogger, client)
	gatewayController := controller.NewController(conf, slogLogger, provider, jwtProvider, passwordProvider, cacher, userStorer)
	middlewareMiddleware := middleware.NewMiddleware(conf, slogLogger, provider, timeProvider, jwtProvider, gatewayController)
	handler := gateway.NewHandler(gatewayController)
	userController := controller2.NewController(conf, slogLogger, provider, passwordProvider, userStorer)
	userHandler := user.NewHandler(userController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, userHandler)
	application := NewApplication(slogLogger, inputPortServer)
	return application
}