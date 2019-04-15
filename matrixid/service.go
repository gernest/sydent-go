package matrixid

import (
	"net/http"

	"github.com/gernest/sydent-go/clients"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/logger"
	"github.com/gernest/sydent-go/models"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

//go:generate go run generate_store_driver.go

// @title MatrixID
// @version 1.0
// @description This is implementation of matrix matrix.org identity service

// @license.name MIT

func Service(opts *core.Ctx, m Metric) *echo.Echo {
	e := echo.New()
	matrix := e.Group("/_matrix")
	identityService := matrix.Group("/identity/api")
	matrix.Use(middleware.CORSWithConfig(CORS()))
	identityService.GET("/v1", Version)
	identityService.OPTIONS("/v1", options)
	identityService.GET("/v1/pubkey/ephemeral/isvalid", EphemeralIsValid(opts))
	identityService.GET("/v1/pubkey/:keyId", GetPublicKey(opts))
	identityService.GET("/v1/pubkey/isvalid", PublicKeyIsValid(opts))
	identityService.OPTIONS("/v1/lookup", options)
	identityService.GET("/v1/lookup", Lookup(opts, m))
	identityService.POST("/v1/bulk_lookup", BulkLookup(opts, m))
	identityService.OPTIONS("/v1/bulk_lookup", options)
	identityService.POST("/v1/validate/email/requestToken", EmailRequestCode(opts, m))
	identityService.OPTIONS("/v1/validate/email/requestToken", options)
	identityService.POST("/v1/validate/email/submitToken", PostEmailValidatedCode(opts, m))
	identityService.GET("/v1/validate/email/submitToken", GetEmailValidatedCode(opts, m))
	identityService.POST("/v1/validate/msisdn/requestToken", todo)
	identityService.POST("/v1/validate/msisdn/submitToken", todo)
	identityService.GET("/v1/validate/msisdn/submitToken", todo)
	identityService.GET("/v1/3pid/getValidated3pid", GetValidated3PID(opts, m))
	identityService.OPTIONS("/v1/bind", options)
	identityService.POST("/v1/bind", Bind(opts, m))
	identityService.POST("/v1/unbind", Unbind(opts, clients.Fed))
	identityService.POST("/v1/store-invite", StoreInvite(opts, m))
	identityService.POST("/v1/sign-ed25519", SignED25519(opts, m))
	identityService.OPTIONS("/v1/sign-ed25519", options)
	matrix.POST("/identity/replicate/v1/push", Replicate(opts, m))
	e.GET("/metrics", echo.WrapHandler(promhttp.Handler()))
	return e
}

func todo(ctx echo.Context) error {
	return ctx.JSON(http.StatusNotImplemented, models.NewError(
		models.ErrUnknown,
		"Not implemented",
	))
}

func options(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, map[string]interface{}{})
}

// RequestError logs an error that occurred during request processing. Fields of
// interest are url Path and Method.
func RequestError(lg logger.Logger, req *http.Request, err error) {
	lg.Error(err.Error(),
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
	)
}

// CORS configures echo middleware for identity service cors.
func CORS() middleware.CORSConfig {
	return middleware.CORSConfig{
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderXRequestedWith,
			echo.HeaderContentType,
			echo.HeaderAccept,
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
			echo.OPTIONS,
		},
	}
}

// InternalError renders json for internal server errors. We don't divulge
// reasons for internal errors
func InternalError(ctx echo.Context) error {
	return ctx.JSON(http.StatusInternalServerError, models.NewError(
		models.ErrUnknown,
		http.StatusText(http.StatusInternalServerError),
	))
}
