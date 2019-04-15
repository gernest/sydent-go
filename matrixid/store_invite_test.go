package matrixid

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gernest/sydent-go/config"

	"github.com/labstack/echo"
)

func TestStoreInvite(t *testing.T) {
	type emailResp struct {
		from string
		to   []string
		msg  string
	}
	var received *emailResp
	mail, err := config.New(TestEmailClient{
		host: "localhost",
		send: func(from string, to []string, msg []byte) error {
			received = &emailResp{
				from: from,
				to:   to,
				msg:  string(msg),
			}
			return nil
		},
	}, mainContext.Config.GetTemplate())
	if err != nil {
		t.Fatal(err)
	}
	lg := &TestLogger{}
	sr := `{
		"medium": "email",
		"address": "foo@bar.baz",
		"room_id": "!something:example.tld",
		"sender": "@bob:example.com"
	  }`

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader(sr))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	tctx := *mainContext
	tctx.Email = mail
	tctx.Log = lg
	err = StoreInvite(&tctx, &TestMetric{})(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if received != nil {
	}
	got := rec.Body.String()
	t.Log(got)
	t.Log(lg)
	if received != nil {
		t.Log(received.msg)
	}
}
