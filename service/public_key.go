package service

import (
	"net/http"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/signedjson"
	"github.com/labstack/echo"
)

// GetPublicKey uses key id to search for a stored public key pinned by this
// server.
func GetPublicKey(coreContext *core.Ctx) echo.HandlerFunc {
	key := coreContext.Config.Server.Crypto.Key()
	keyID := key.KeyID()
	pubBase64 := signedjson.EncodeBase64(key.PublicKey)
	return func(ctx echo.Context) error {
		id := ctx.Param("keyId")
		if keyID != id {
			return ctx.JSON(http.StatusNotFound, ErrPublicKeyNotFound)
		}
		return ctx.JSON(http.StatusOK, models.PublicKey{Key: pubBase64})
	}
}

// ErrPublicKeyNotFound shortcut for the error object returned when there is no
// public key found.
var ErrPublicKeyNotFound = models.NewError(
	models.ErrNotFound,
	"The public key was not found",
)

// PublicKeyIsValid checks if a public key is valid.
func PublicKeyIsValid(coreContext *core.Ctx) echo.HandlerFunc {
	key := coreContext.Config.Server.Crypto.Key()
	pubBase64 := signedjson.EncodeBase64(key.PublicKey)
	return func(ctx echo.Context) error {
		pk := ctx.QueryParam("public_key")
		return ctx.JSON(http.StatusOK, models.ValidPubKey{Valid: pk == pubBase64})
	}
}

// EphemeralIsValid checks if a short term public key is valid.
func EphemeralIsValid(coreContext *core.Ctx) echo.HandlerFunc {
	db := coreContext.Store
	return func(ctx echo.Context) error {
		pk := ctx.QueryParam("public_key")
		return ctx.JSON(http.StatusOK, models.ValidPubKey{
			Valid: db.ValidateEphemeralPublicKey(ctx.Request().Context(), pk) == nil,
		})
	}
}
