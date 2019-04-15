package schema

import (
	"context"
	"io/ioutil"

	"github.com/gernest/sydent-go/embed"
	"github.com/gernest/sydent-go/models"
)

const homeServerUpSQL = "/schemas/full_home_server_schema.sql"
const identityUpSQL = "/schemas/full_identity_schema.sql"
const identityDownSQL = "/schemas/full_identity_schema_down.sql"

func IdentityUp(ctx context.Context, fs embed.Embed, db models.Query) error {
	return execFile(ctx, fs, identityUpSQL, db)
}

func execFile(ctx context.Context, fs embed.Embed, name string, db models.Query) error {
	f, err := fs.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}
	_, err = db.ExecContext(ctx, string(b))
	return err
}

func IdentityDown(ctx context.Context, fs embed.Embed, db models.Query) error {
	return execFile(ctx, fs, identityDownSQL, db)
}
