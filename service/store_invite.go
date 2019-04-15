package service

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/signedjson"
	"github.com/labstack/echo"
)

const inviteTpl = "invite"

func StoreInvite(coreContext *core.Ctx, m Metric) echo.HandlerFunc {
	serverCfg := coreContext.Config.Server
	mail := coreContext.Config.Email
	emailTplName := mail.Invite.Template
	if emailTplName == "" {
		emailTplName = inviteTpl
	}
	count := m.CountError("store_invite")
	db := coreContext.Store
	send := coreContext.Email.SendMail
	return func(ctx echo.Context) error {
		m, err := models.EnsureParams(ctx.Request(), "medium", "address", "room_id", "sender")
		if err != nil {
			RequestError(coreContext.Log, ctx.Request(), err)
			return ctx.JSON(http.StatusBadRequest, err)
		}
		medium := m["medium"]
		address := m["address"]
		roomID := m["room_id"]
		sender := m["sender"]
		requestContext := ctx.Request().Context()
		mxid, err := db.GlobalGetMxid(requestContext, medium, address)
		if err != nil && err != sql.ErrNoRows {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), err)
			return InternalError(ctx)
		}
		if mxid != "" {
			err := models.NewError(
				models.ErrThreepidInUse,
				fmt.Sprintf("Binding to %s is already known", mxid),
			)
			RequestError(coreContext.Log, ctx.Request(), err)
			return ctx.JSON(http.StatusBadRequest, err)
		}
		if medium != "email" {
			err := models.NewError(
				models.ErrUnrecognized,
				fmt.Sprintf("Didn't understand medium %q", medium),
			)
			RequestError(coreContext.Log, ctx.Request(), err)
			return ctx.JSON(http.StatusBadRequest, err)
		}
		token := models.RandomString(128)
		keys, err := signedjson.New("0")
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), err)
			return InternalError(ctx)
		}
		ephemeralPrivateKey := signedjson.EncodeBase64(keys.PrivateKey)
		ephemeralPublicKey := signedjson.EncodeBase64(keys.PublicKey)
		err = db.StoreEphemeralPublicKey(requestContext, ephemeralPublicKey)
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), err)
			return InternalError(ctx)
		}
		err = db.StoreToken(requestContext, models.InviteToken{
			Medium:  medium,
			Address: address,
			RoomID:  roomID,
			Sender:  sender,
			Token:   token,
		})
		if err != nil {
			count.Inc()
			RequestError(coreContext.Log, ctx.Request(), err)
			return InternalError(ctx)
		}
		substitution := map[string]string{}
		query := ctx.Request().URL.Query()
		for k := range query {
			substitution[k] = query.Get(k)
		}
		substitution["token"] = token
		substitution["ephemeral_private_key"] = ephemeralPrivateKey
		if substitution["room_name"] != "" {
			substitution["bracketed_room_name"] = fmt.Sprintf("(%s)", substitution["room_name"])
		}
		err = send(requestContext, emailTplName, mail.Invite.From, []string{address}, substitution)
		if err != nil {
			RequestError(coreContext.Log, ctx.Request(), err)
			return InternalError(ctx)
		}
		serverKeys := serverCfg.Crypto.Key()
		pubKey := signedjson.EncodeBase64(serverKeys.PublicKey)
		baseURL := fmt.Sprintf("%s/_matrix/identity/api/v1",
			serverCfg.ClientHTTPBase,
		)
		return ctx.JSON(http.StatusOK, map[string]interface{}{
			"token":      token,
			"public_key": pubKey,
			"public_keys": []map[string]string{
				map[string]string{
					"public_key":       pubKey,
					"key_validity_url": baseURL + "/pubkey/isvalid",
				},
				map[string]string{
					"public_key":       ephemeralPublicKey,
					"key_validity_url": baseURL + "/pubkey/ephemeral/isvalid",
				},
			},
			"display_name": redact(address),
		})
	}
}

// remove sensitive information from the email address.
func redact(address string) string {
	var p []string
	for _, v := range strings.Split(address, "@") {
		if len(v) > 5 {
			p = append(p, v[:3]+"...")
		} else if len(v) > 1 {
			p = append(p, string(v[0])+"...")
		} else {
			p = append(p, "...")
		}
	}
	return strings.Join(p, "@")
}
