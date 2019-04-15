package core

import (
	"github.com/gernest/sydent-go/config"
	"github.com/gernest/sydent-go/logger"
	"github.com/gernest/sydent-go/store"
	"go.uber.org/zap"
)

type Ctx struct {
	Config            *config.Matrix
	Log               logger.Logger
	Email             config.Mail
	Store             store.Store
	ReplicationClient config.HTTPClient
}

// Namespace returns a new Ctx with the logger namespaced to ns.
func (ctx *Ctx) Namespace(ns string) *Ctx {
	return &Ctx{
		Config:            ctx.Config,
		Log:               ctx.Log.With(zap.Namespace(ns)),
		Email:             ctx.Email,
		Store:             ctx.Store,
		ReplicationClient: ctx.ReplicationClient,
	}
}
