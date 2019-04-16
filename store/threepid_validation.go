package store

import (
	"context"
	"database/sql"

	"github.com/gernest/sydent-go/config"
	"github.com/gernest/sydent-go/models"
	"github.com/gernest/sydent-go/store/query"
)

func GetOrCreateTokenSession(ctx context.Context, db models.Query, medium, address, clientSecret string) (*models.ValidationSession, error) {
	var v models.ValidationSession
	err := db.QueryRowContext(ctx, query.GetTokenSession, medium, address, clientSecret).Scan(
		&v.ID,
		&v.Medium,
		&v.Address,
		&v.ClientSecret,
		&v.Validated,
		&v.Mtime,
		&v.Token,
		&v.SendAttemptNumber,
	)
	if err == sql.ErrNoRows {
		mtime := models.Time()
		sid, err := AddValidationSession(ctx, db, medium, address, clientSecret, mtime)
		if err != nil {
			return nil, err
		}
		tokenString := models.GenerateToken(medium)
		_, err = db.ExecContext(ctx, query.CreateTokenSession, sid, tokenString, -1)
		if err != nil {
			return nil, err
		}
		v.ID = sid
		v.Medium = medium
		v.Address = address
		v.ClientSecret = clientSecret
		v.Mtime = mtime
		v.Token = tokenString
		v.SendAttemptNumber = -1

	}
	return &v, nil
}

func AddValidationSession(ctx context.Context, db models.Query, medium, address, clientSecret string, mtime int64) (int64, error) {
	var id int64
	err := db.QueryRowContext(ctx, query.AddValidationSession, medium, address, clientSecret, mtime).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func SetSendAttemptNumber(ctx context.Context, db models.Query, sid int64, attemptNo int64) error {
	_, err := db.ExecContext(ctx, query.SetSendAttemptNumber, attemptNo, sid)
	return err
}

func SetValidated(ctx context.Context, db models.Query, sid string, validated int) error {
	_, err := db.ExecContext(ctx, query.SetValidated, validated, sid)
	return err
}

func SetMtime(ctx context.Context, db models.Query, sid int64, mtime int64) error {
	_, err := db.ExecContext(ctx, query.SetMtime, mtime, sid)
	return err
}

func GetSessionByID(ctx context.Context, db models.Query, sid int64) (*models.ValidationSession, error) {
	var v models.ValidationSession
	err := db.QueryRowContext(ctx, query.GetSessionByID, sid).Scan(
		&v.ID,
		&v.Medium,
		&v.Address,
		&v.ClientSecret,
		&v.Validated,
		&v.Mtime,
	)
	if err != nil {
		return nil, err
	}
	return &v, err
}

func GetTokenSessionByID(ctx context.Context, db models.Query, sid int64) (*models.TokenSession, error) {
	var v models.TokenSession
	err := db.QueryRowContext(ctx, query.GetTokenSessionByID, sid).Scan(
		&v.ID,
		&v.Medium,
		&v.Address,
		&v.ClientSecret,
		&v.Validated,
		&v.Mtime,
		&v.Token,
		&v.SendAttemptNumber,
	)
	if err != nil {
		return nil, err
	}
	return &v, err
}

// GetValidatedSession returns a validated session with ma matching clientSecret
func GetValidatedSession(ctx context.Context, db models.Query, sid int64, clientSecret string) (*models.ValidationSession, error) {
	sess, err := GetSessionByID(ctx, db, sid)
	if err != nil {
		return nil, err
	}
	if sess.ClientSecret != clientSecret {
		return nil, models.ErrIncorrectClientSecret
	}
	if (sess.Mtime + config.ThreepidSessionValidationLifetime) < models.Time() {
		return nil, models.ErrSessionExpired
	}
	if sess.Validated == 0 {
		return nil, models.ErrSessionNotValidated
	}
	return sess, nil
}
