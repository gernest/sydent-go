package store

import (
	"context"

	"github.com/gernest/sydent-go/models"
)

func GetPeerByName(ctx context.Context, driver Driver, db models.Query, name string) (*models.Peer, error) {
	var peer models.Peer
	rows, err := db.QueryContext(ctx, driver.GetPeerByName(), name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var key, alg string
		err := rows.Scan(
			&peer.Name,
			&peer.Port,
			&peer.LastSentVersion,
			&alg, &key,
		)
		if err != nil {
			return nil, err
		}
		if peer.PublicKeys == nil {
			peer.PublicKeys = make(map[string]string)
		}
		peer.PublicKeys[alg] = key
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	return &peer, nil
}

func GetAllPeers(ctx context.Context, driver Driver, db models.Query) ([]models.Peer, error) {
	rows, err := db.QueryContext(ctx, driver.GetAllPeers())
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cache models.Peer
	cache.PublicKeys = make(map[string]string)
	var peers []models.Peer
	for rows.Next() {
		var p models.Peer
		var key, alg string
		err = rows.Scan(
			&p.Name,
			&p.Port,
			&p.LastSentVersion,
			&key, &alg,
		)
		if p.Name != cache.Name {
			if len(cache.PublicKeys) > 0 {
				p.PublicKeys = make(map[string]string)
				for k, v := range cache.PublicKeys {
					p.PublicKeys[k] = v
				}
				peers = append(peers, p)
				cache.PublicKeys = make(map[string]string)
			}
			cache.Name = p.Name
			cache.Port = p.Port
			cache.LastSentVersion = p.LastSentVersion
		}
		cache.PublicKeys[alg] = key
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}
	if len(cache.PublicKeys) > 0 {
		peers = append(peers, cache)
	}
	return peers, nil
}

func SetLastSentVersionAndPokeSucceeded(ctx context.Context, driver Driver, db models.Query, peerName, lastSentVersion, lastPokeSucceeded string) error {
	_, err := db.ExecContext(ctx, driver.SetLastSentVersionAndPokeSucceeded(), lastSentVersion, lastPokeSucceeded, peerName)
	return err
}
