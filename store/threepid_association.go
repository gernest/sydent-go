package store

import (
	"context"
	"fmt"
	"time"

	"github.com/gernest/sydent-go/models"
)

func LocalAddOrUpdateAssociation(ctx context.Context, driver Driver, db models.Query, as *models.Association) error {
	_, err := db.ExecContext(ctx, driver.LocalAddOrUpdateAssociation(),
		as.Medium, as.Address, as.MatrixID, as.Timestamp, as.NotBefore, as.NotAfter)
	return err
}

func GetAssociationsAfterID(ctx context.Context, driver Driver, db models.Query, afterID int64, limit int64) ([]models.Association, error) {
	rows, err := db.QueryContext(ctx, driver.GetAssociationsAfterId(), afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var ass []models.Association
	for rows.Next() {
		var as models.Association
		err := rows.Scan(
			&as.ID,
			&as.Medium,
			&as.Address,
			&as.MatrixID,
			&as.Timestamp,
			&as.NotBefore,
			&as.NotAfter,
		)
		if err != nil {
			return nil, err
		}
		ass = append(ass, as)
	}
	return ass, nil
}

func LocalRemoveAssociation(ctx context.Context, driver Driver, db models.Query, as *models.Association) error {
	var count int64
	err := db.QueryRowContext(ctx,
		driver.GetLocal3pid(), as.Medium, as.Address, as.MatrixID).Scan(&count)
	if err != nil {
		return err
	}
	if count > 0 {
		now := time.Now()
		ts := models.MS(&now)
		_, err = db.ExecContext(ctx, driver.LocalRemoveAssociation(), as.Medium, as.Address, ts)
		return err
	}
	return nil
}

func SignedAssociationStringForThreepid(ctx context.Context, driver Driver, db models.Query, medium, address string) (string, error) {
	now := time.Now()
	ts := models.MS(&now)
	var signed string
	err := db.QueryRowContext(ctx, driver.SignedAssociationStringForThreepid(), medium, address, ts, ts).Scan(&signed)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func GlobalGetMxid(ctx context.Context, driver Driver, db models.Query, medium, address string) (string, error) {
	now := time.Now()
	ts := models.MS(&now)
	var signed string
	err := db.QueryRowContext(ctx, driver.GlobalGetMxid(), medium, address, ts, ts).Scan(&signed)
	if err != nil {
		return "", err
	}
	return signed, nil
}

func GlobalGetMxids(ctx context.Context, driver Driver, db models.Query, ids [][]string) ([]models.Association, error) {
	tx, err := db.(models.SQL).BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, driver.CreateTMPMxid())
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	var results []models.Association
	icap := 0
	for icap < len(ids) {
		var x string
		var a []interface{}
		if len(ids) > icap+500 {
			x, a = mxidsInsertSQL(ids[icap : icap+500])
		} else {
			x, a = mxidsInsertSQL(ids[icap:])
		}
		_, err = tx.ExecContext(ctx, x, a...)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		icap += 500
	}
	now := time.Now()
	ts := models.MS(&now)
	rows, err := tx.QueryContext(ctx, driver.GlobalGetMxids(), ts, ts)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var as models.Association
		err = rows.Scan(
			&as.Medium,
			&as.Address,
			&as.Timestamp,
			&as.MatrixID,
		)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
		results = append(results, as)
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func mxidsInsertSQL(ids [][]string) (string, []interface{}) {
	s := `
INSERT INTO tmp_getmxids (medium, address) VALUES `
	v, args := genValues(ids)
	return s + v + ";", args
}

func genValues(args [][]string) (string, []interface{}) {
	s := ""
	var r []interface{}
	p := 0
	for i, x := range args {
		if i == 0 {
			s = fmt.Sprintf("($%d,$%d)", p+1, p+2)
		} else {
			s += fmt.Sprintf(",($%d,$%d)", p+1, p+2)
		}
		for _, v := range x {
			r = append(r, v)
		}
		p += 2
	}
	return s, r
}

func GlobalAddAssociation(ctx context.Context, driver Driver, db models.Query, as *models.Association, originServer string, originID int64, rawSgnAssoc string) error {
	_, err := db.ExecContext(ctx, driver.GlobalAddAssociation(), as.Medium, as.Address, as.MatrixID,
		as.Timestamp, as.NotBefore, as.NotAfter, originServer, originID, rawSgnAssoc)
	return err
}

func GlobalLastIDFromServer(ctx context.Context, driver Driver, db models.Query, originServer string) (int64, error) {
	var originID int64
	err := db.QueryRowContext(ctx, driver.GlobalLastIDFromServer(), originServer).Scan(&originID)
	if err != nil {
		return 0, err
	}
	return originID, nil
}

func GlobalRemoveAssociation(ctx context.Context, driver Driver, db models.Query, medium, address string) error {
	_, err := db.ExecContext(ctx, driver.GlobalRemoveAssociation(), medium, address)
	return err
}
