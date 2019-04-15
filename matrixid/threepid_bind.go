package matrixid

import (
	"net/http"
	"strconv"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/labstack/echo"
)

// Bind returns a handler that binds third party id to matrix id. You can think
// of this in the following terms.
//
// Say, you own an email address foo@example.com . The server running this
// handler is serving requests on bar.com. Now, you want to use bar.com for
// identification with people you want to chat with in the matrix multiverse.
//
// So you tell bar.com that you are the owner of foo@example.com , bar.com will
// send you a verification email to make sure you are you, after verifying the
// email now bar.com knows who you are, so it gives you another id which is
// recognized by other applications running in the multiverse example
// @foo:bar.com.
//
// You will become @foo:bar.com to the rest of the matrix multiverse and
// foo@example.com to bar.com only.
//
// This handler takes care of the binding of foo@example.com => @foo:bar.com.
// Note that, verification of who you are is done by another handler, this just
// make sure this server remembers who you really are and assignment of your
// matrix id.
func Bind(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	bind := AddBinding(coreContext)
	count := m.CountError("bind")
	db := coreContext.Store
	return func(ctx echo.Context) error {
		req := ctx.Request()
		m, merr := models.EnsureParams(req, "sid", "client_secret", "mxid")
		if merr != nil {
			count.Inc()
			return ctx.JSON(http.StatusBadRequest, merr)
		}
		sid, err := strconv.ParseInt(m["sid"], 10, 64)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), merr)
			return ctx.JSON(http.StatusBadRequest,
				models.NewError(
					models.ErrInvalidParam,
					"sid is not valid int64",
				),
			)
		}
		mxid := m["mxid"]
		clientSecret := m["client_secret"]
		requestContext := req.Context()
		noMatch := models.NewError(
			models.ErrNoValidSession,
			"No valid session was found matching that sid and client secret",
		)
		s, err := db.GetValidatedSession(requestContext, sid, clientSecret)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			switch err {
			case models.ErrIncorrectClientSecret:
				return ctx.JSON(http.StatusBadRequest, noMatch)
			case models.ErrSessionExpired:
				return ctx.JSON(http.StatusBadRequest, models.NewError(
					models.ErrSessionExpiredCode,
					"This validation session has expired: call requestToken again",
				))
			case models.ErrSessionNotValidated:
				return ctx.JSON(http.StatusBadRequest, models.NewError(
					models.ErrSessionNotValidatedCode,
					"This validation session has not yet been completed",
				))
			default:
				return ctx.JSON(http.StatusBadRequest, noMatch)
			}
		}
		sgAss, err := bind(requestContext, s.Medium, s.Address, mxid)
		if err != nil {
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusBadRequest, noMatch)
		}
		return ctx.JSON(http.StatusOK, sgAss)
	}
}
