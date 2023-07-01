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
	controller6 "github.com/LuchaComics/cps-backend/app/attachment/controller"
	datastore4 "github.com/LuchaComics/cps-backend/app/attachment/datastore"
	controller4 "github.com/LuchaComics/cps-backend/app/comicsub/controller"
	datastore3 "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	controller5 "github.com/LuchaComics/cps-backend/app/customer/controller"
	"github.com/LuchaComics/cps-backend/app/gateway/controller"
	controller3 "github.com/LuchaComics/cps-backend/app/organization/controller"
	datastore2 "github.com/LuchaComics/cps-backend/app/organization/datastore"
	controller2 "github.com/LuchaComics/cps-backend/app/user/controller"
	"github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http"
	"github.com/LuchaComics/cps-backend/inputport/http/attachment"
	"github.com/LuchaComics/cps-backend/inputport/http/comicsub"
	"github.com/LuchaComics/cps-backend/inputport/http/customer"
	"github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/cps-backend/inputport/http/organization"
	"github.com/LuchaComics/cps-backend/inputport/http/user"
	"github.com/LuchaComics/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/cps-backend/provider/jwt"
	"github.com/LuchaComics/cps-backend/provider/kmutex"
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
	userController := controller2.NewController(conf, slogLogger, provider, passwordProvider, organizationStorer, userStorer)
	userHandler := user.NewHandler(userController)
	s3Storager := s3.NewStorage(conf, slogLogger, provider)
	comicSubmissionStorer := datastore3.NewDatastore(conf, slogLogger, client)
	organizationController := controller3.NewController(conf, slogLogger, provider, s3Storager, emailer, organizationStorer, userStorer, comicSubmissionStorer)
	organizationHandler := organization.NewHandler(organizationController)
	kmutexProvider := kmutex.NewProvider()
	cpsrnProvider := cpsrn.NewProvider()
	cbffBuilder := pdfbuilder.NewCBFFBuilder(conf, slogLogger, provider)
	pcBuilder := pdfbuilder.NewPCBuilder(conf, slogLogger, provider)
	ccimgBuilder := pdfbuilder.NewCCIMGBuilder(conf, slogLogger, provider)
	ccscBuilder := pdfbuilder.NewCCSCBuilder(conf, slogLogger, provider)
	ccBuilder := pdfbuilder.NewCCBuilder(conf, slogLogger, provider)
	comicSubmissionController := controller4.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, kmutexProvider, cpsrnProvider, cbffBuilder, pcBuilder, ccimgBuilder, ccscBuilder, ccBuilder, emailer, userStorer, comicSubmissionStorer, organizationStorer)
	comicsubHandler := comicsub.NewHandler(comicSubmissionController)
	customerController := controller5.NewController(conf, slogLogger, provider, s3Storager, passwordProvider, cbffBuilder, emailer, userStorer)
	customerHandler := customer.NewHandler(customerController)
	attachmentStorer := datastore4.NewDatastore(conf, slogLogger, client)
	attachmentController := controller6.NewController(conf, slogLogger, provider, s3Storager, emailer, attachmentStorer, userStorer, comicSubmissionStorer)
	attachmentHandler := attachment.NewHandler(attachmentController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, userHandler, organizationHandler, comicsubHandler, customerHandler, attachmentHandler)
	application := NewApplication(slogLogger, inputPortServer)
	return application
}
