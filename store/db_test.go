package store

import (
	"context"
	"database/sql"
	"os"
	"reflect"
	"testing"

	"github.com/gernest/sydent-go/embed"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/sydent-go/store/query"
	"github.com/gernest/sydent-go/store/schema"
	_ "github.com/lib/pq"
)

type TestContext struct {
	Ctx    context.Context
	Driver Driver
	Query  models.Query
}

func TestGenValues(t *testing.T) {
	ids := [][]string{
		[]string{"a", "b"},
		[]string{"c", "d"},
		[]string{"e", "f"},
	}
	x, a := genValues(ids)
	ex := "($1,$2),($3,$4),($5,$6)"
	if x != ex {
		t.Errorf("expected %s got %s", ex, x)
	}
	ea := []interface{}{"a", "b", "c", "d", "e", "f"}
	if !reflect.DeepEqual(a, ea) {
		t.Errorf("expected %#v got %#v", ea, a)
	}
}

func TestIdentity(t *testing.T) {
	driverName := "postgres"
	conn := os.Getenv("MATRIXID_DB_CONN")
	db, err := sql.Open(driverName, conn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	drv, err := NewDriver(driverName)
	if err != nil {
		t.Fatal(err)
	}
	tctx := TestContext{
		Ctx:    context.Background(),
		Driver: drv,
		Query:  query.New(db),
	}
	fs := embed.New()

	err = schema.IdentityDown(tctx.Ctx, fs, tctx.Query)
	if err != nil {
		t.Fatal(err)
	}
	err = schema.IdentityUp(tctx.Ctx, fs, tctx.Query)
	if err != nil {
		t.Fatal(err)
	}

	testInvites(t, tctx)
}
