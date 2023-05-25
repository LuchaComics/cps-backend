//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/LuchaComics/cps-backend/adapter/cache/redis"
	"github.com/LuchaComics/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/cps-backend/adapter/pdfbuilder"
	"github.com/LuchaComics/cps-backend/adapter/storage/mongodb"
	s3_storage "github.com/LuchaComics/cps-backend/adapter/storage/s3"
	gateway_c "github.com/LuchaComics/cps-backend/app/gateway/controller"
	submission_c "github.com/LuchaComics/cps-backend/app/submission/controller"
	submission_s "github.com/LuchaComics/cps-backend/app/submission/datastore"
	user_c "github.com/LuchaComics/cps-backend/app/user/controller"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http"
	gateway_http "github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	submission_http "github.com/LuchaComics/cps-backend/inputport/http/submission"
	user_http "github.com/LuchaComics/cps-backend/inputport/http/user"
	"github.com/LuchaComics/cps-backend/provider/jwt"
	"github.com/LuchaComics/cps-backend/provider/logger"
	"github.com/LuchaComics/cps-backend/provider/password"
	"github.com/LuchaComics/cps-backend/provider/time"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

func InitializeEvent() Application {
	// Our application is dependent on the following Golang packages. We need to
	// provide them to Google wire so it can sort out the dependency injection
	// at compile time.
	wire.Build(
		config.New,
		uuid.NewProvider,
		time.NewProvider,
		logger.NewProvider,
		jwt.NewProvider,
		mailgun.NewEmailer,
		password.NewProvider,
		mongodb.NewStorage,
		s3_storage.NewStorage,
		redis.NewCache,
		pdfbuilder.NewCBFFBuilder,
		user_s.NewDatastore,
		user_c.NewController,
		submission_s.NewDatastore,
		submission_c.NewController,
		gateway_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		submission_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		NewApplication)
	return Application{}
}
