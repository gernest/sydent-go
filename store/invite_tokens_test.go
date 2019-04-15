package store

import (
	"testing"

	"github.com/gernest/sydent-go/models"
)

func testInvites(t *testing.T, ctx TestContext) {
	sampleTokens := []models.InviteToken{
		{Medium: "email", Address: "email1"},
	}
	t.Run("StoreToken", func(ts *testing.T) {
		for _, v := range sampleTokens {
			err := StoreToken(ctx.Ctx, ctx.Driver, ctx.Query, v)
			if err != nil {
				ts.Error(err)
			}
		}
	})
	t.Run("GetTokens", func(ts *testing.T) {
		a := sampleTokens[0]
		tkn, err := GetTokens(ctx.Ctx, ctx.Driver, ctx.Query, a.Medium, a.Address)
		if err != nil {
			ts.Fatal(err)
		}
		if len(tkn) != 1 {
			ts.Fatalf("expected one token got %d", len(tkn))
		}
	})
	t.Run("MarkTokensAsSent", func(ts *testing.T) {
		a := sampleTokens[0]
		err := MarkTokensAsSent(ctx.Ctx, ctx.Driver, ctx.Query, a.Medium, a.Address)
		if err != nil {
			ts.Fatal(err)
		}
		tkn, err := GetTokens(ctx.Ctx, ctx.Driver, ctx.Query, a.Medium, a.Address)
		if err != nil {
			ts.Fatal(err)
		}
		o := tkn[0]
		if o.SentAt.IsZero() {
			ts.Error("expected sent_ts to be set")
		}
	})

}
