package service

import (
	"context"
	"fmt"
	"net"
	"net/url"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
)

type TokenRequest struct {
	Email        string
	ClientSecret string
	SendAttempt  int64
	NextLink     string
	IP           net.IP
}

func RequestEmailToken(ctx context.Context, coreContext *core.Ctx, req *TokenRequest) (int64, error) {
	db := coreContext.Store
	session, err := db.GetOrCreateTokenSession(ctx,
		"email", req.Email, req.ClientSecret,
	)
	if err != nil {
		return 0, err
	}

	err = db.SetMtime(ctx, session.ID, models.Time())
	if err != nil {
		return 0, err
	}
	if session.SendAttemptNumber >= req.SendAttempt {
		return session.ID, nil
	}
	clientHTTPBase := coreContext.Config.Server.ClientHTTPBase
	link, err := makeValidateLink(clientHTTPBase, session, req.ClientSecret, req.NextLink)
	if err != nil {
		return 0, err
	}
	data := map[string]string{
		"ipaddress": req.IP.String(),
		"link":      link,
		"token":     session.Token,
	}
	err = coreContext.Email.SendMail(ctx, verifyTpl,
		coreContext.Config.Email.Verification.From,
		[]string{req.Email}, data,
	)
	if err != nil {
		return 0, err
	}
	err = db.SetSendAttemptNumber(ctx, session.ID, req.SendAttempt)
	if err != nil {
		return 0, err
	}
	return session.ID, nil
}

func makeValidateLink(clientHTTPBase string, session *models.ValidationSession, clientSecret, nextLink string) (string, error) {
	link := fmt.Sprintf("%s/_matrix/identity/api/v1/validate/email/submitToken", clientHTTPBase)
	q := make(url.Values)
	q.Set("token", session.Token)
	q.Set("client_secret", clientSecret)
	sid := fmt.Sprint(session.ID)
	q.Set("sid", sid)
	if nextLink != "" {
		u, err := url.Parse(nextLink)
		if err != nil {
			return "", err
		}
		nq := u.Query()
		nq.Set("sid", sid)
		u.RawQuery = nq.Encode()
		q.Set("nextLink", u.String())
	}
	return link + "?" + q.Encode(), nil
}
