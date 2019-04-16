package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/gernest/sydent-go/models"
	"github.com/gernest/sydent-go/store/query"
)

func StoreToken(ctx context.Context, db models.Query, token models.InviteToken) error {
	_, err := db.ExecContext(ctx, query.StoreToken,
		token.Medium, token.Address, token.RoomID, token.Sender, token.Token,
		time.Now().UTC().Unix(),
	)
	return err
}

func GetTokens(ctx context.Context, db models.Query, medium, address string) ([]models.InviteToken, error) {
	rows, err := db.QueryContext(ctx, query.GetTokens, medium, address)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var result []models.InviteToken
	for rows.Next() {
		var token models.InviteToken
		var sent sql.NullInt64
		var received sql.NullInt64
		err := rows.Scan(
			&token.Medium,
			&token.Address,
			&token.RoomID,
			&token.Sender,
			&token.Token,
			&received,
			&sent,
		)
		if err != nil {
			return nil, err
		}
		if received.Valid {
			token.ReceivedAt = models.FromMS(received.Int64)
		}
		if sent.Valid {
			token.SentAt = models.FromMS(sent.Int64)
		}
		result = append(result, token)
	}
	return result, nil
}

func MarkTokensAsSent(ctx context.Context, db models.Query, medium, address string) error {
	now := time.Now()
	ts := models.MS(&now)
	_, err := db.ExecContext(ctx, query.MarkTokensAsSent, ts, medium, address)
	return err
}

func StoreEphemeralPublicKey(ctx context.Context, db models.Query, publicKey string) error {
	now := time.Now()
	ts := models.MS(&now)
	_, err := db.ExecContext(ctx, query.StoreEphemeralPublicKey, publicKey, ts)
	return err
}

func ValidateEphemeralPublicKey(ctx context.Context, db models.Query, publicKey string) error {
	_, err := db.ExecContext(ctx, query.ValidateEphemeralPublicKey, publicKey)
	return err
}

func GetSenderForToken(ctx context.Context, db models.Query, token string) (string, error) {
	var sender string
	err := db.QueryRowContext(ctx, query.GetSenderForToken, token).Scan(&sender)
	if err != nil {
		return "", err
	}
	return sender, nil
}
