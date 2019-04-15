package service

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gernest/sydent-go/models"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/signedjson"
	"github.com/labstack/echo"
)

// Lookup gets a 3pid bound to a matrix user id.
func Lookup(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	serverName := coreContext.Config.Server.Name
	key := coreContext.Config.Server.Crypto.Key()
	db := coreContext.Store
	count := m.CountError("lookup")
	return func(ctx echo.Context) error {
		medium := ctx.QueryParam("medium")
		address := ctx.QueryParam("address")
		req := ctx.Request()
		rs, err := db.SignedAssociationStringForThreepid(req.Context(), medium, address)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			if err == sql.ErrNoRows {
				return ctx.JSON(http.StatusNotFound, models.NewError(
					models.ErrNotFound,
					"",
				))
			}
			return ctx.JSON(http.StatusInternalServerError, models.NewError(
				models.ErrUnknown,
				"",
			))
		}
		if rs == "" {
			return InternalError(ctx)
		}
		var o map[string]interface{}
		err = json.Unmarshal([]byte(rs), &o)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return InternalError(ctx)
		}
		as, ok := o["signatures"]
		if ok {
			am, ok := as.(map[string]interface{})
			if ok {
				_, ok = am[serverName]
				if ok {
					return ctx.JSON(http.StatusOK, o)
				}
			}
		}
		// # We have not yet worked out what the proper trust model should be.
		// #
		// # Maybe clients implicitly trust a server they talk to (and so we
		// # should sign every assoc we return as ourselves, so they can
		// # verify this).
		// #
		// # Maybe clients really want to know what server did the original
		// # verification, and want to only know exactly who signed the assoc.
		// #
		// # Until we work out what we should do, sign all assocs we return as
		// # ourself. This is vaguely ok because there actually is only one
		// # identity server, but it happens to have two names (matrix.org and
		// # vector.im), and so we're not really lying too much.
		// #
		// # We do this when we return assocs, not when we receive them over
		// # replication, so that we can undo this decision in the future if
		// # we wish, without having destroyed the raw underlying data.
		msg := signedjson.Message(o)
		err = key.Sign(msg, serverName)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return InternalError(ctx)
		}
		return ctx.JSON(http.StatusOK, msg)
	}
}
