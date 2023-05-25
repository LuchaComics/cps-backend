// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/LuchaComics/cps-backend/adapter/cache/redis"
	"github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	"github.com/LuchaComics/cps-backend/adapter/storage/mongodb"
	"github.com/LuchaComics/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/cps-backend/app/gateway/controller"
	controller3 "github.com/LuchaComics/cps-backend/app/organization/controller"
	datastore2 "github.com/LuchaComics/cps-backend/app/organization/datastore"
	controller4 "github.com/LuchaComics/cps-backend/app/submission/controller"
	datastore3 "github.com/LuchaComics/cps-backend/app/submission/datastore"
	controller2 "github.com/LuchaComics/cps-backend/app/user/controller"
	"github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http"
	"github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/cps-backend/inputport/http/organization"
	"github.com/LuchaComics/cps-backend/inputport/http/submission"
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
	emailer := mailgun.NewEmailer(conf, slogLogger, provider)
	client := mongodb.NewStorage(conf, slogLogger)
	userStorer := datastore.NewDatastore(conf, slogLogger, client)
	organizationStorer := datastore2.NewDatastore(conf, slogLogger, client)
	gatewayController := controller.NewController(conf, slogLogger, provider, jwtProvider, passwordProvider, cacher, emailer, userStorer, organizationStorer)
	middlewareMiddleware := middleware.NewMiddleware(conf, slogLogger, provider, timeProvider, jwtProvider, gatewayController)
	handler := gateway.NewHandler(gatewayController)
	userController := controller2.NewController(conf, slogLogger, provider, passwordProvider, userStorer)
	userHandler := user.NewHandler(userController)
	s3Storager := s3.NewStorage(conf, slogLogger, provider)
	organizationController := controller3.NewController(conf, slogLogger, provider, s3Storager, emailer, organizationStorer)
	organizationHandler := organization.NewHandler(organizationController)
	cbffBuilder := pdfbuilder.NewCBFFBuilder(conf, slogLogger, provider)
	submissionStorer := datastore3.NewDatastore(conf, slogLogger, client)
	submissionController := controller4.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, cbffBuilder, emailer, submissionStorer)
	submissionHandler := submission.NewHandler(submissionController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, userHandler, organizationHandler, submissionHandler)
	application := NewApplication(slogLogger, inputPortServer)
	return application
}
