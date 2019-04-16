package store

import (
	"context"
	"time"

	"github.com/gernest/sydent-go/models"
	"github.com/prometheus/client_golang/prometheus"
)

type Metric struct {
	Enabled bool
	Vec     *prometheus.HistogramVec
}

const label = "name"

func NewMetric(opts prometheus.Opts) Metric {
	var m Metric
	m.Vec = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: opts.Namespace,
		Subsystem: opts.Subsystem,
		Name:      "query",
		Buckets:   prometheus.LinearBuckets(10, 5, 6),
	}, []string{label})
	prometheus.MustRegister(m.Vec)
	m.Enabled = true
	return m
}

func (m Metric) observe(labelStr string, fn func()) {
	if m.Enabled {
		now := time.Now()
		fn()
		d := float64(time.Since(now).Nanoseconds()) / float64(time.Second)
		if d < 0 {
			d = 0
		}
		m.Vec.With(prometheus.Labels{label: labelStr}).Observe(d)
	} else {
		fn()
	}
}

// Identity contains all midentity service database facing routines.
type Identity struct {
	db      models.Query
	metrics Metric
}

func New(db models.Query, m Metric) *Identity {
	return &Identity{db: db, metrics: m}
}

func (id *Identity) StoreToken(ctx context.Context, token models.InviteToken) (err error) {
	id.metrics.observe("store_token", func() {
		err = StoreToken(ctx, id.db, token)
	})
	return
}

func (id *Identity) GetTokens(ctx context.Context, medium, address string) (tokens []models.InviteToken, err error) {
	id.metrics.observe("get_tokens", func() {
		tokens, err = GetTokens(ctx, id.db, medium, address)
	})
	return
}

func (id *Identity) MarkTokensAsSent(ctx context.Context, medium, address string) (err error) {
	id.metrics.observe("mark_tokens_as_sent", func() {
		err = MarkTokensAsSent(ctx, id.db, medium, address)
	})
	return
}

func (id *Identity) StoreEphemeralPublicKey(ctx context.Context, publicKey string) (err error) {
	id.metrics.observe("store_ephemeral_public_key", func() {
		err = StoreEphemeralPublicKey(ctx, id.db, publicKey)
	})
	return
}

func (id *Identity) ValidateEphemeralPublicKey(ctx context.Context, publicKey string) (err error) {
	id.metrics.observe("store_ephemeral_public_key", func() {
		err = ValidateEphemeralPublicKey(ctx, id.db, publicKey)
	})
	return
}

func (id *Identity) GetSenderForToken(ctx context.Context, token string) (tokenInfo string, err error) {
	id.metrics.observe("get_sender_for_token", func() {
		tokenInfo, err = GetSenderForToken(ctx, id.db, token)
	})
	return
}

func (id *Identity) SignedAssociationStringForThreepid(ctx context.Context, medium, address string) (ass string, err error) {
	id.metrics.observe("signed_association_string_for_threepid", func() {
		ass, err = SignedAssociationStringForThreepid(ctx, id.db, medium, address)
	})
	return
}

func (id *Identity) GlobalGetMxid(ctx context.Context, medium, address string) (mxid string, err error) {
	id.metrics.observe("global_get_mxid", func() {
		mxid, err = GlobalGetMxid(ctx, id.db, medium, address)
	})
	return
}

func (id *Identity) GetPeerByName(ctx context.Context, name string) (peer *models.Peer, err error) {
	id.metrics.observe("get_peer_by_name", func() {
		peer, err = GetPeerByName(ctx, id.db, name)
	})
	return
}

func (id *Identity) GetAllPeers(ctx context.Context) (peers []models.Peer, err error) {
	id.metrics.observe("get_all_peers", func() {
		peers, err = GetAllPeers(ctx, id.db)
	})
	return
}

func (id *Identity) SetLastSentVersionAndPokeSucceeded(ctx context.Context, peerName, lastSentVersion, lastPokeSucceeded string) (err error) {
	id.metrics.observe("set_last_sent_version_and_poke_succeeded", func() {
		err = SetLastSentVersionAndPokeSucceeded(ctx, id.db, peerName, lastSentVersion, lastPokeSucceeded)
	})
	return
}

func (id *Identity) SetSendAttemptNumber(ctx context.Context, sid int64, attemptNo int64) (err error) {
	id.metrics.observe("set_send_attempt_number", func() {
		err = SetSendAttemptNumber(ctx, id.db, sid, attemptNo)
	})
	return
}

func (id *Identity) SetValidated(ctx context.Context, sid string, validated int) (err error) {
	id.metrics.observe("set_validated", func() {
		err = SetValidated(ctx, id.db, sid, validated)
	})
	return
}

func (id *Identity) SetMtime(ctx context.Context, sid int64, mtime int64) (err error) {
	id.metrics.observe("set_mtime", func() {
		err = SetMtime(ctx, id.db, sid, mtime)
	})
	return
}

func (id *Identity) GetSessionByID(ctx context.Context, sid int64) (session *models.ValidationSession, err error) {
	id.metrics.observe("get_session_by_id", func() {
		session, err = GetSessionByID(ctx, id.db, sid)
	})
	return
}

func (id *Identity) GetValidatedSession(ctx context.Context, sid int64, clientSecret string) (session *models.ValidationSession, err error) {
	id.metrics.observe("get_validated_session", func() {
		session, err = GetValidatedSession(ctx, id.db, sid, clientSecret)
	})
	return
}

func (id *Identity) GetTokenSessionByID(ctx context.Context, sid int64) (tokenSession *models.TokenSession, err error) {
	id.metrics.observe("get_token_session_by_id", func() {
		tokenSession, err = GetTokenSessionByID(ctx, id.db, sid)
	})
	return
}

func (id *Identity) GlobalGetMxids(ctx context.Context, ids [][]string) (mxids []models.Association, err error) {
	id.metrics.observe("global_get_mxids", func() {
		mxids, err = GlobalGetMxids(ctx, id.db, ids)
	})
	return
}

func (id *Identity) GlobalLastIDFromServer(ctx context.Context, originServer string) (lastID int64, err error) {
	id.metrics.observe("global_last_id_from_server", func() {
		lastID, err = GlobalLastIDFromServer(ctx, id.db, originServer)
	})
	return
}

func (id *Identity) GlobalAddAssociation(ctx context.Context, as *models.Association, originServer string, originID int64, rawSgnAssoc string) (err error) {
	id.metrics.observe("global_add_association", func() {
		err = GlobalAddAssociation(ctx, id.db, as, originServer, originID, rawSgnAssoc)
	})
	return
}

func (id *Identity) GlobalRemoveAssociation(ctx context.Context, medium, address string) (err error) {
	id.metrics.observe("global_remove_association", func() {
		err = GlobalRemoveAssociation(ctx, id.db, medium, address)
	})
	return
}

func (id *Identity) LocalAddOrUpdateAssociation(ctx context.Context, as *models.Association) (err error) {
	id.metrics.observe("local_add_or_update_association", func() {
		err = LocalAddOrUpdateAssociation(ctx, id.db, as)
	})
	return
}

func (id *Identity) LocalRemoveAssociation(ctx context.Context, as *models.Association) (err error) {
	id.metrics.observe("local_remove_association", func() {
		err = LocalRemoveAssociation(ctx, id.db, as)
	})
	return
}

func (id *Identity) GetAssociationsAfterID(ctx context.Context, afterID int64, limit int64) (as []models.Association, err error) {
	id.metrics.observe("local_get_association_after_id", func() {
		as, err = GetAssociationsAfterID(ctx, id.db, afterID, limit)
	})
	return
}

func (id *Identity) GetOrCreateTokenSession(ctx context.Context, medium, address, clientSecret string) (session *models.ValidationSession, err error) {
	id.metrics.observe("get_or_create_token_session", func() {
		session, err = GetOrCreateTokenSession(ctx, id.db, medium, address, clientSecret)
	})
	return
}
