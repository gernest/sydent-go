package service

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
)

func TestVersion(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)
	err := Version(ctx)
	if err != nil {
		t.Fatal(err)
	}
	expect := "{}"
	got := rec.Body.String()
	got = strings.TrimSpace(got)
	if got != expect {
		t.Errorf("expected %q got %q", expect, got)
	}
}
