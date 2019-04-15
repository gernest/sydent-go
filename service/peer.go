package service

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gernest/sydent-go/config"

	"github.com/gernest/sydent-go/core"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/signedjson"
)

const IdentityReplicationPush = "/_matrix/identity/replicate/v1/push"

// ErrNoSignature is returned when there is no signatures field in a signed
// message object.
var ErrNoSignature = errors.New("no signatures found")

const SIGNING_KEY_ALGORITHM = "ed25519"

const DefaultReplicationPort = 1001

// ErrNoMatchingSignature is returned when there is no signatures for this server
// found in a signed association object.
var ErrNoMatchingSignature = errors.New("no matching signatures found")

// Association defines association tuple that is shared during replication.
type Association struct {
	OriginID          int64
	SignedAssociation signedjson.Message
}

type Payload struct {
	SignedAssociations []Association `json:"sgAssocs,omitempty"`
}

// Peer is an interface for replicating messages across matrix peers.
type Peer interface {
	PushUpdates(context.Context, []Association) error
}

// PushFunc defines a function that implements Peer interface.
type PushFunc func(context.Context, []Association) error

// PushUpdates wrapper for implementing Peer interface on pf.
func (pf PushFunc) PushUpdates(ctx context.Context, as []Association) error {
	return pf(ctx, as)
}

// PushLocal pushes associations within the same host/instance. This copies
// associations from local table to the global association table.
func PushLocal(coreContext *core.Ctx) PushFunc {
	serverName := coreContext.Config.Server.Name
	db := coreContext.Store
	return func(ctx context.Context, as []Association) error {
		lastID, err := db.GlobalLastIDFromServer(ctx, coreContext.Config.Server.Name)
		if err != nil {
			if err != sql.ErrNoRows {
				return err
			}
			lastID = -1
		}
		for _, v := range as {
			if v.OriginID > lastID {
				a, err := AssociationFromMap(v.SignedAssociation)
				if err != nil {
					return err
				}
				if a.MatrixID != "" {
					b, err := json.Marshal(v.SignedAssociation)
					if err != nil {
						return err
					}
					err = db.GlobalAddAssociation(ctx, a,
						serverName, v.OriginID, string(b),
					)
					if err != nil {
						return err
					}
				} else {
					err = db.GlobalRemoveAssociation(ctx, a.Medium, a.Address)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	}
}

func AssociationFromMap(m map[string]interface{}) (*models.Association, error) {
	var a models.Association
	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &a)
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetVerifyKeyFromPeer returns the verification key for peer.
func GetVerifyKeyFromPeer(p *models.Peer) (*signedjson.Key, error) {
	if p.PublicKeys == nil || p.PublicKeys[SIGNING_KEY_ALGORITHM] == "" {
		return nil, errors.New("matrixid: no public keys found for peer")
	}
	k := p.PublicKeys[SIGNING_KEY_ALGORITHM]
	return signedjson.DecodeVerifyKeyBase64(SIGNING_KEY_ALGORITHM, "", k)
}

// GetReplicationURLFromPeer returns a url string for replication on peer.
func GetReplicationURLFromPeer(cfg *config.Matrix, peer *models.Peer) string {
	var r string
	for _, p := range cfg.Peers {
		if p.Name == peer.Name {
			r = p.BaseReplicationURL
			break
		}
	}
	if r == "" {
		var port int64 = 1001
		if peer.Port.Valid {
			port = peer.Port.Int64
		}
		r = fmt.Sprintf("https://%s:%d", peer.Name, port)
	}
	if r[len(r)-1] != '/' {
		r += "/"
	}
	return r + "_matrix/identity/replicate/v1/push"
}

func VerifySignedAssociation(ctx context.Context, key *signedjson.Key, serverName string, msg signedjson.Message) error {
	if _, ok := msg["signatures"]; !ok {
		return ErrNoSignature
	}
	keyIDS := msg.SignatureID(serverName)
	if keyIDS == nil {
		return ErrNoMatchingSignature
	}
	if in(key.KeyID(), keyIDS...) {
		return key.Verify(msg, serverName)
	}
	return ErrNoMatchingSignature
}

func in(key string, v ...string) bool {
	for _, s := range v {
		if key == s {
			return true
		}
	}
	return false
}

// PushToRemotePeer pushes signed associations to a remote identity service peer.
//
// replica is the client used for replication.
func PushToRemotePeer(cfg *config.Matrix, peer *models.Peer, replica config.HTTPClient, as []Association) error {
	m := Payload{SignedAssociations: as}
	b, err := json.Marshal(m)
	if err != nil {
		return err
	}
	replicationURL := GetReplicationURLFromPeer(cfg, peer)
	req, err := http.NewRequest("POST", replicationURL, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", config.AgentName)
	res, err := replica.Do(req)
	if err != nil {
		return err
	}
	return res.Body.Close()
}
