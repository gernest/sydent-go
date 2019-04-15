package service

import (
	"bytes"
	"net/http"
	"strconv"
	"strings"

	"github.com/gernest/sydent-go/config"
	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/labstack/echo"
)

const (
	verifyTpl     = "verification"
	verifyPageTpl = "verify_response"
)

func EmailRequestCode(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	count := m.CountError("email_request_code")
	return func(ctx echo.Context) error {
		req := ctx.Request()
		m, merr := models.EnsureParams(req, "email", "client_secret", "send_attempt")
		if merr != nil {
			count.Inc()
			RequestError(coreContext.Log, req, merr)
			return ctx.JSON(http.StatusBadRequest, merr)
		}
		email := m["email"]
		clientSecret := m["client_secret"]
		sendAttempt, err := strconv.ParseInt(m["send_attempt"], 10, 64)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrInvalidParam,
				"send_attempt is not a valid integer",
			))
		}
		if !config.IsValidEmail(email) {
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrInvalidEmail,
				"Invalid email address",
			))
		}
		var tr TokenRequest
		tr.Email = email
		tr.ClientSecret = clientSecret
		tr.SendAttempt = sendAttempt
		if n, ok := m["next_link"]; ok {
			if !strings.HasPrefix(n, "file:///") {
				tr.NextLink = n
			}
		}

		sid, err := RequestEmailToken(
			req.Context(), coreContext, &tr,
		)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, req, err)
			return ctx.JSON(http.StatusBadRequest, models.NewError(
				models.ErrEmailSendError,
				"Failed to send email",
			))
		}
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"success": true,
			"sid":     sid,
		})
	}
}

func GetEmailValidatedCode(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	tpl := coreContext.Config.GetTemplate()
	v := coreContext.Config.Email.Verification.ResponsePage
	if v == "" {
		v = verifyPageTpl
	}
	count := m.CountError("get_email_validation_code")
	return func(ctx echo.Context) error {
		req := ctx.Request()
		msg := "Verification successful! Please return to your Matrix client to continue."
		status := http.StatusOK
		err := validateEmailRequest(coreContext, req)
		if err != nil {
			count.Inc()
			msg = "Verification failed: you may need to request another verification email"
		} else {
			nextLink := ctx.Param("nextLink")
			if nextLink != "" && !strings.HasPrefix(nextLink, "file:///") {
				status = http.StatusFound
				req.Header.Set("Location", nextLink)
			}
		}
		var buf bytes.Buffer
		err = tpl.ExecuteTemplate(&buf, v, map[string]interface{}{
			"message": msg,
		})
		if err != nil {
			RequestError(coreContext.Log, req, err)
			count.Inc()
			return ctx.HTML(http.StatusInternalServerError, http.StatusText(
				http.StatusInternalServerError,
			))
		}
		return ctx.HTML(status, buf.String())
	}
}

func PostEmailValidatedCode(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	count := m.CountError("post_email_validated_code")
	return func(ctx echo.Context) error {
		err := validateEmailRequest(coreContext, ctx.Request())
		if err != nil {
			count.Inc()
			return ctx.JSON(http.StatusOK, models.Success{
				Success: false,
			})
		}
		return ctx.JSON(http.StatusOK, models.Success{
			Success: true,
		})
	}
}

func validateEmailRequest(coreContext *core.Ctx, req *http.Request) error {
	m, merr := models.EnsureParams(req, "token", "sid", "client_secret")
	if merr != nil {
		return merr
	}
	sid, err := strconv.ParseInt(m["sid"], 10, 64)
	if err != nil {
		RequestError(coreContext.Log, req, err)
		return models.NewError(
			models.ErrInvalidParam,
			"sid is not an int64",
		)
	}
	token := m["token"]
	clientSecret := m["client_secret"]
	err = SessionWithToken(req.Context(), coreContext, sid, clientSecret, token)
	if err != nil {
		RequestError(coreContext.Log, req, err)
		switch err {
		case models.ErrIncorrectClientSecret:
			return models.NewError(
				models.ErrIncorrectClientSecretCode,
				"Client secret does not match the one given when requesting the token",
			)
		case models.ErrSessionExpired:
			return models.NewError(
				models.ErrSessionExpiredCode,
				"This validation session has expired: call requestToken again",
			)
		default:
			return models.NewError(
				models.ErrUnknown,
				"",
			)
		}
	}
	return nil
}
