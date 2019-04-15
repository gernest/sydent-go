package service

import (
	"context"

	"github.com/gernest/sydent-go/config"
	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
)

func SessionWithToken(ctx context.Context, coreContext *core.Ctx, sid int64, clientSecret, token string) error {
	s, err := coreContext.Store.GetTokenSessionByID(ctx, sid)
	if err != nil {
		return err
	}
	if s.ClientSecret != clientSecret {
		return models.ErrIncorrectClientSecret
	}
	x := s.Mtime + config.ThreepidSessionValidationTimeout
	if x < models.Time() {
		return models.ErrSessionExpired
	}
	if s.Token != token {
		return models.ErrIncorrectToken
	}
	return nil
}
