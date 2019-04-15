package service

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gernest/sydent-go/models"

	"github.com/labstack/echo"
)

func TestPublicKey(t *testing.T) {
	crypto := mainContext.Config.Server.Crypto
	key := crypto.Key()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	ctx.SetParamNames("keyId")
	ctx.SetParamValues(key.KeyID())
	err := GetPublicKey(mainContext)(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ex, err := json.Marshal(models.PublicKey{
		Key: crypto.VerifyKey,
	})
	if err != nil {
		t.Fatal(err)
	}
	expect := string(ex)
	got := strings.TrimSpace(rec.Body.String())
	if got != expect {
		t.Errorf("expected %q got %q", expect, got)
	}

	// missing keyId
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	ctx = e.NewContext(req, rec)
	ctx.SetParamNames("keyId")
	ctx.SetParamValues("bad key id")
	err = GetPublicKey(mainContext)(ctx)
	if err != nil {
		t.Fatal(err)
	}
	ex, err = json.Marshal(ErrPublicKeyNotFound)
	if err != nil {
		t.Fatal(err)
	}
	expect = string(ex)
	got = strings.TrimSpace(rec.Body.String())
	if got != expect {
		t.Errorf("expected %q got %q", expect, got)
	}
}

func TestPeers(t *testing.T) {
	p := []models.Peer{
		{
			Name:     "id1.matrixid.org",
			JSONPort: DefaultReplicationPort,
			PublicKeys: map[string]string{
				SIGNING_KEY_ALGORITHM: "key one",
			},
		},
		{
			Name:     "id2.matrixid.org",
			JSONPort: DefaultReplicationPort,
			PublicKeys: map[string]string{
				SIGNING_KEY_ALGORITHM: "key two",
			},
		},
		{
			Name:     "id3.matrixid.org",
			JSONPort: DefaultReplicationPort,
			PublicKeys: map[string]string{
				SIGNING_KEY_ALGORITHM: "key three",
			},
		},
	}
	b, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("tmp/peers.json", b, 0600)
}
