package store

import (
	"context"

	"github.com/gernest/sydent-go/models"
)

var _ Store = (*Matrix)(nil)

type Store interface {
	StoreToken(ctx context.Context, token models.InviteToken) error
	GetTokens(ctx context.Context, medium, address string) ([]models.InviteToken, error)
	MarkTokensAsSent(ctx context.Context, medium, address string) error
	StoreEphemeralPublicKey(ctx context.Context, publicKey string) error
	ValidateEphemeralPublicKey(ctx context.Context, publicKey string) error
	GetSenderForToken(ctx context.Context, token string) (string, error)

	SignedAssociationStringForThreepid(ctx context.Context, medium, address string) (string, error)
	GlobalGetMxid(ctx context.Context, medium, address string) (string, error)
	GlobalGetMxids(ctx context.Context, ids [][]string) ([]models.Association, error)
	GlobalRemoveAssociation(ctx context.Context, medium, address string) error
	GlobalAddAssociation(ctx context.Context, as *models.Association, originServer string, originID int64, rawSgnAssoc string) error
	LocalAddOrUpdateAssociation(ctx context.Context, as *models.Association) error
	LocalRemoveAssociation(ctx context.Context, as *models.Association) error
	GetAssociationsAfterID(ctx context.Context, afterID int64, limit int64) ([]models.Association, error)

	GetPeerByName(ctx context.Context, name string) (*models.Peer, error)
	GetAllPeers(ctx context.Context) ([]models.Peer, error)
	SetLastSentVersionAndPokeSucceeded(ctx context.Context, peerName, lastSentVersion, lastPokeSucceeded string) error

	SetSendAttemptNumber(ctx context.Context, sid int64, attemptNo int64) error
	SetValidated(ctx context.Context, sid string, validated int) error
	SetMtime(ctx context.Context, sid int64, mtime int64) error
	GetSessionByID(ctx context.Context, sid int64) (*models.ValidationSession, error)
	GetTokenSessionByID(ctx context.Context, sid int64) (*models.TokenSession, error)
	GlobalLastIDFromServer(ctx context.Context, originServer string) (int64, error)
	GetOrCreateTokenSession(ctx context.Context, medium, address, clientSecret string) (*models.ValidationSession, error)
	GetValidatedSession(ctx context.Context, sid int64, clientSecret string) (*models.ValidationSession, error)

	DB() models.SQL
	Driver() Driver
	Metric() Metric
}

// Matrix universal database handler for the matrix services.
type Matrix struct {
	*Identity
	db      models.SQL
	driver  Driver
	metrics Metric
}

func NewStore(db models.SQL, driver Driver, m Metric) *Matrix {
	return &Matrix{
		Identity: New(db, driver, m),
		db:       db,
		driver:   driver,
		metrics:  m,
	}
}

func (m *Matrix) DB() models.SQL {
	return m.db
}

func (m *Matrix) Driver() Driver {
	return m.driver
}

func (m *Matrix) Metric() Metric {
	return m.metrics
}

// Driver defines methods for retrieving sql query strings used by the identity
// service.
type Driver interface {
	Name() string

	// Param returns a string representing positional parameter. Different
	// databases offers different ways to do this.
	//
	// For instance in postgres you can use $1 to mark first argument and ? is used
	// for sqlite
	Param(index int) string

	StoreToken() string
	GetTokens() string
	MarkTokensAsSent() string
	StoreEphemeralPublicKey() string
	ValidateEphemeralPublicKey() string
	GetSenderForToken() string

	SignedAssociationStringForThreepid() string
	GlobalGetMxid() string
	GlobalGetMxids() string
	GlobalRemoveAssociation() string
	GlobalAddAssociation() string

	GetPeerByName() string
	GetAllPeers() string
	SetLastSentVersionAndPokeSucceeded() string
	SetSendAttemptNumber() string
	SetValidated() string
	SetMtime() string
	GetSessionByID() string
	GetTokenSessionByID() string
	GlobalLastIDFromServer() string
	LocalAddOrUpdateAssociation() string
	GetTokenSession() string
	CreateTokenSession() string
	AddValidationSession() string
	GetAssociationsAfterId() string
	GetLocal3pid() string
	LocalRemoveAssociation() string
	CreateTMPMxid() string
}
