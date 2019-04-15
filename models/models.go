package models

import (
	"database/sql"
	"time"
)

// common http verbs
const (
	Get         = "GET"
	Post        = "POST"
	Put         = "Put"
	Options     = "OPTIONS"
	Delete      = "DELETE"
	ContentType = "Content-Type"
	ContentJSON = "application/json"
)

// ContextKey keys used by all services
type ContextKey uint

// context keys used by all services
const (

	// Logs
	CtxLogSilent ContextKey = iota
	CtxLogBase
	CtxLogIdentity
	CtxLogAccount

	// Databases
	CtxStore
)

type InviteToken struct {
	ID         int64
	Medium     string
	Address    string
	RoomID     string
	Sender     string
	Token      string
	ReceivedAt time.Time
	SentAt     time.Time
}

func (i *InviteToken) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"medium":  i.Medium,
		"address": i.Address,
		"room_id": i.RoomID,
		"token":   i.Token,
		"sender":  i.Sender,
	}
}

type EphemeralPublicKey struct {
	ID          int64
	PublicKey   string
	VerifyCount int64
	UpdatedAt   time.Time
}

type ServerVerifyKey struct {
	Name      string
	From      string
	Timestamp int64
	VerifyKey []byte
}

type Peer struct {
	Name                string        `json:"name"`
	Port                sql.NullInt64 `json:"-"`
	JSONPort            int64         `json:"port,omitempty"`
	LastSentVersion     sql.NullInt64 `json:"-"`
	LastPokeSucceededAt sql.NullInt64 `json:"-"`

	// 1 is for active and 0 is for not active
	Active     int64             `json:"-"`
	PublicKeys map[string]string `json:"public_keys"`
}

// Data is an arbitrary json object
type Data map[string]interface{}

type Association struct {
	ID          int64                  `json:"-"`
	Medium      string                 `json:"medium"`
	Address     string                 `json:"address"`
	MatrixID    string                 `json:"mxid"`
	Timestamp   int64                  `json:"ts"`
	NotBefore   int64                  `json:"notBefore"`
	NotAfter    int64                  `json:"notAfter"`
	ExtraFields map[string]interface{} `json:"-"`
}

func (a *Association) ToMap() map[string]interface{} {
	m := map[string]interface{}{
		"medium":     a.Medium,
		"address":    a.Address,
		"mxid":       a.MatrixID,
		"ts":         a.Timestamp,
		"not_before": a.NotBefore,
		"not_after":  a.NotAfter,
	}
	for k, v := range a.ExtraFields {
		m[k] = v
	}
	return m
}

type ValidationSession struct {
	ID                int64  `json:"-"`
	Medium            string `json:"medium"`
	Address           string `json:"address"`
	ClientSecret      string `json:"-"`
	Validated         int    `json:"-"`
	Mtime             int64  `json:"validated_at"`
	Token             string `json:"-"`
	SendAttemptNumber int64  `json:"-"`
}

type TokenSession struct {
	ValidationSession
	SendAttemptNumber int64
}

type BulkLookupRequest struct {
	Threepids [][]string `json:"threepids"`
}

type Success struct {
	Success bool `json:"success"`
}

// PublicKey represent a publick key response object that is returned by the
// identity server.
type PublicKey struct {
	Key string `json:"public_key" example:"ed25519:0"`
}

// ValidPubKey is a response object for  valid public key.
type ValidPubKey struct {
	Valid bool `json:"valid"`
}

// IDLookupResponse is the response sent by identity server for /lookup request.
type IDLookupResponse struct {
	// The 3pid address of the user being looked up, matching the address requested.
	Address string `json:"address"`
	Medium  string `json:"medium"`
	Mxid    string `json:"mxid"`
	// A unix timestamp before which the association is not known to be valid
	NotBefore int64 `json:"not_before"`
	// A unix timestamp after which the association is not known to be valid.
	NotAfter int64 `json:"not_after"`
	// The unix timestamp at which the association was verified.
	Timestamp int64 `json:"ts"`
	// The signatures of the verifying identity servers which show that the
	// association should be trusted, if you trust the verifying identity servers
	Signature map[string]map[string]string `json:"signature"`
}

type Profile struct {
	ID          string `json:"user_id,omitempty"`
	DisplayName string `json:"displayname,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
	LastChecked int64  `json:"last_check,omitempty"`
}
