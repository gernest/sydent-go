package matrixid

import (
	"net/http"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/signedjson"
	"github.com/labstack/echo"
)

func SignED25519(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	serverName := coreContext.Config.Server.Name
	count := m.CountError("sign_ed25519")
	db := coreContext.Store
	return func(ctx echo.Context) error {
		req := ctx.Request()
		m, merr := models.EnsureParams(req, "private_key", "token", "mxid")
		if merr != nil {
			count.Inc()
			RequestError(coreContext.Log, req, merr)
			return ctx.JSON(http.StatusBadRequest, merr)
		}
		privateKeyBase64 := m["private_key"]
		token := m["token"]
		mxid := m["mxid"]
		sender, err := db.GetSenderForToken(req.Context(), token)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, merr)
			return ctx.JSON(http.StatusNotFound,
				models.NewError(
					models.ErrUnrecognized,
					"Didn't recognize token",
				),
			)
		}
		message := signedjson.Message{
			"mxid":   mxid,
			"sender": sender,
			"token":  token,
		}
		key, err := signedjson.DecodeSigningKeyBase64("ed25519", "0", privateKeyBase64)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, merr)
			return ctx.JSON(http.StatusNotFound,
				models.NewError(
					models.ErrUnknown,
					"",
				),
			)
		}
		err = key.Sign(message, serverName)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, merr)
			return ctx.JSON(http.StatusNotFound,
				models.NewError(
					models.ErrUnknown,
					"",
				),
			)
		}
		return ctx.JSON(http.StatusOK, message)
	}
}
