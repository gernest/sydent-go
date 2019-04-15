package matrixid

import (
	"net/http"
	"strconv"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/labstack/echo"
)

func GetValidated3PID(coreContext *core.Ctx, mx Metric) echo.HandlerFunc {
	count := mx.CountError("get_validated_3pids")
	db := coreContext.Store
	return func(ctx echo.Context) error {
		m, merr := models.EnsureParams(ctx.Request(), "sid", "client_secret")
		if merr != nil {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), merr)
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
		clientSecret := m["client_secret"]
		requestContext := ctx.Request().Context()
		noMatch := models.NewError(
			models.ErrNoValidSession,
			"No valid session was found matching that sid and client secret",
		)
		expired := models.NewError(
			models.ErrSessionExpiredCode,
			"This validation session has expired: call requestToken again",
		)
		notValid := models.NewError(
			models.ErrSessionNotValidatedCode,
			"This validation session has expired: call requestToken again",
		)
		sess, err := db.GetValidatedSession(requestContext, sid, clientSecret)
		if err != nil {
			count.Inc()
			//TODO: figure out proper status codes
			RequestError(coreContext.Log, ctx.Request(), err)
			switch err {
			case models.ErrIncorrectClientSecret:
				RequestError(coreContext.Log, ctx.Request(), noMatch)
				return ctx.JSON(http.StatusOK, noMatch)
			case models.ErrSessionExpired:
				RequestError(coreContext.Log, ctx.Request(), expired)
				return ctx.JSON(http.StatusOK, expired)
			case models.ErrSessionNotValidated:
				RequestError(coreContext.Log, ctx.Request(), notValid)
				return ctx.JSON(http.StatusOK, notValid)
			default:
				RequestError(coreContext.Log, ctx.Request(), noMatch)
				return ctx.JSON(http.StatusOK, noMatch)
			}
		}

		return ctx.JSON(http.StatusOK, sess)
	}
}
