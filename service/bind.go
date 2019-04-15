package service

import (
	"context"
	"database/sql"

	"github.com/gernest/signedjson"
	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
)

const associationLifetime = 100 * 365 * 24 * 60 * 60 * 1000

func AddBinding(coreContext *core.Ctx) func(context.Context, string, string, string) (signedjson.Message, error) {
	key := coreContext.Config.Server.Crypto.Key()
	serverName := coreContext.Config.Server.Name
	sign := Signer(coreContext)
	return func(ctx context.Context, medium, address, mxid string) (signedjson.Message, error) {
		createdAt := models.Time()
		expiresAt := createdAt + associationLifetime
		as := &models.Association{
			Medium:    medium,
			Address:   address,
			MatrixID:  mxid,
			Timestamp: createdAt,
			NotBefore: createdAt,
			NotAfter:  expiresAt,
		}
		err := coreContext.Store.LocalAddOrUpdateAssociation(ctx, as)
		if err != nil {
			return nil, err
		}
		tokens, err := coreContext.Store.GetTokens(ctx, medium, address)
		if err != nil {
			return nil, err
		}
		var invites []map[string]interface{}
		for _, token := range tokens {
			m := token.ToMap()
			m["mxid"] = mxid
			signed := signedjson.Message{
				"mxid":  mxid,
				"token": m["token"],
			}
			err = key.Sign(signed, serverName)
			if err != nil {
				return nil, err
			}
			m["signed"] = signed
			invites = append(invites, m)
		}
		if len(invites) > 0 {
			as.ExtraFields = map[string]interface{}{
				"invites": invites,
			}
			err := coreContext.Store.MarkTokensAsSent(ctx, medium, address)
			if err != nil {
				return nil, err
			}
		}
		a, err := sign(as)
		if err != nil {
			return nil, err
		}
		return a, nil
	}
}

func RemoveBinding(coreContext *core.Ctx) func(context.Context, *models.Association) error {
	push := LocalPusher(coreContext)
	return func(ctx context.Context, as *models.Association) error {
		err := coreContext.Store.LocalRemoveAssociation(ctx, as)
		if err != nil {
			return err
		}
		return push(ctx)
	}
}

func LocalPusher(coreContext *core.Ctx) func(context.Context) error {
	push := PushLocal(coreContext)
	name := coreContext.Config.Server.Name
	db := coreContext.Store
	sign := Signer(coreContext)
	return func(ctx context.Context) error {
		lastID, err := db.GlobalLastIDFromServer(ctx, name)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			lastID = -1
		}
		as, err := db.GetAssociationsAfterID(ctx, lastID, 0)
		if err != nil {
			return err
		}
		out := make([]Association, len(as))
		for k, v := range as {
			a, err := sign(&v)
			if err != nil {
				return err
			}
			out[k] = Association{
				OriginID:          v.ID,
				SignedAssociation: a,
			}
		}
		return push(ctx, out)
	}
}

func Signer(coreContext *core.Ctx) func(*models.Association) (signedjson.Message, error) {
	key := coreContext.Config.Server.Crypto.Key()
	name := coreContext.Config.Server.Name
	return func(as *models.Association) (signedjson.Message, error) {
		m := as.ToMap()
		err := key.Sign(m, name)
		if err != nil {
			return nil, err
		}
		return m, nil
	}
}
