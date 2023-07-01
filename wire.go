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
	attachment_c "github.com/LuchaComics/cps-backend/app/attachment/controller"
	attachment_s "github.com/LuchaComics/cps-backend/app/attachment/datastore"
	comicsub_c "github.com/LuchaComics/cps-backend/app/comicsub/controller"
	comicsub_s "github.com/LuchaComics/cps-backend/app/comicsub/datastore"
	customer_c "github.com/LuchaComics/cps-backend/app/customer/controller"
	gateway_c "github.com/LuchaComics/cps-backend/app/gateway/controller"
	organization_c "github.com/LuchaComics/cps-backend/app/organization/controller"
	organization_s "github.com/LuchaComics/cps-backend/app/organization/datastore"
	user_c "github.com/LuchaComics/cps-backend/app/user/controller"
	user_s "github.com/LuchaComics/cps-backend/app/user/datastore"
	"github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/inputport/http"
	attachment_http "github.com/LuchaComics/cps-backend/inputport/http/attachment"
	comicsub_http "github.com/LuchaComics/cps-backend/inputport/http/comicsub"
	customer_http "github.com/LuchaComics/cps-backend/inputport/http/customer"
	gateway_http "github.com/LuchaComics/cps-backend/inputport/http/gateway"
	"github.com/LuchaComics/cps-backend/inputport/http/middleware"
	organization_http "github.com/LuchaComics/cps-backend/inputport/http/organization"
	user_http "github.com/LuchaComics/cps-backend/inputport/http/user"
	"github.com/LuchaComics/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/cps-backend/provider/jwt"
	"github.com/LuchaComics/cps-backend/provider/kmutex"
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
		kmutex.NewProvider,
		mailgun.NewEmailer,
		password.NewProvider,
		cpsrn.NewProvider,
		mongodb.NewStorage,
		s3_storage.NewStorage,
		redis.NewCache,
		pdfbuilder.NewCBFFBuilder,
		pdfbuilder.NewPCBuilder,
		pdfbuilder.NewCCIMGBuilder,
		pdfbuilder.NewCCSCBuilder,
		pdfbuilder.NewCCBuilder,
		pdfbuilder.NewCCUGBuilder,
		user_s.NewDatastore,
		user_c.NewController,
		customer_c.NewController,
		organization_s.NewDatastore,
		organization_c.NewController,
		comicsub_s.NewDatastore,
		comicsub_c.NewController,
		gateway_c.NewController,
		attachment_s.NewDatastore,
		attachment_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		customer_http.NewHandler,
		organization_http.NewHandler,
		comicsub_http.NewHandler,
		attachment_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		NewApplication)
	return Application{}
}
