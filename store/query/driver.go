package query

import (
	"fmt"
)

const StoreToken = `INSERT INTO
invite_tokens (
	medium,
	address,
	room_id,
	sender,
	token,
	received_ts
)
VALUES
($1, $2, $3, $4, $5, $6);`

const GetTokens = `SELECT
    medium,
    address,
    room_id,
    sender,
    token,
    received_ts,
    sent_ts
FROM
    invite_tokens
WHERE
    medium = $1
    AND address = $2;`

const MarkTokensAsSent = `UPDATE
    invite_tokens
SET
    sent_ts = $1
WHERE
    medium = $2
    AND address = $3;`

const StoreEphemeralPublicKey = `INSERT INTO
    ephemeral_public_keys (public_key, persistence_ts)
VALUES
    ($1, $2);`

const ValidateEphemeralPublicKey = `UPDATE
    ephemeral_public_keys
SET
    verify_count = verify_count + 1
WHERE
    public_key = $1;`

const GetSenderForToken = `SELECT
    sender
FROM
    invite_tokens
WHERE
    token = $1;`

const SignedAssociationStringForThreepid = `SELECT
    sgAssoc
FROM
    global_threepid_associations
WHERE
    medium = $1
    and lower(address) = lower($2)
    and notBefore < $3
    and notAfter > $4
ORDER by
    ts desc
LIMIT
    1;`

const GlobalGetMxid = `SELECT
    mxid
FROM
    global_threepid_associations
WHERE
    medium = $1
    and lower(address) = lower($2)
    and notBefore < $3
    and notAfter > $4
ORDER by
    ts desc
LIMIT
    1;`

const CreateTMPMxid = `CREATE TEMPORARY TABLE tmp_getmxids (medium VARCHAR(16), address VARCHAR(256)) ON COMMIT DROP;CREATE INDEX tmp_getmxids_medium_lower_address ON tmp_getmxids (medium, lower(address));`

const GlobalGetMxids = `SELECT
    gte.medium,
    gte.address,
    gte.ts,
    gte.mxid
FROM
    global_threepid_associations gte
    JOIN tmp_getmxids ON gte.medium = tmp_getmxids.medium
    AND lower(gte.address) = lower(tmp_getmxids.address)
WHERE
    gte.notBefore < $1
    AND gte.notAfter > $2
ORDER BY
    gte.medium,
    gte.address,
	gte.ts DESC;`

const GlobalLastIDFromServer = `SELECT
    max(originId),
    count(originId)
FROM
    global_threepid_associations
WHERE
    originServer = $1;`

const GlobalRemoveAssociation = `DELETE FROM
    global_threepid_associations
WHERE
    medium = $1
    AND address = $2`

const GlobalAddAssociation = `insert into
    global_threepid_associations (
        medium,
        address,
        mxid,
        ts,
        notBefore,
        notAfter,
        originServer,
        originId,
        sgAssoc
    ) ($1, $2, $3, $4, $5, $6, $7, $8, $9) ON CONFLICT DO NOTHING;`

const GetPeerByName = `SELECT
    p.name,
    p.port,
    p.lastSentVersion,
    pk.alg,
    pk.key
FROM
    peers p,
    peer_pubkeys pk
WHERE
    p.name = $1
    and pk.peername = p.name
    and p.active = 1;`

const GetAllPeers = `SELECT
    p.name,
    p.port,
    p.lastSentVersion,
    pk.alg,
    pk.key
FROM
    peers p,
    peer_pubkeys pk
WHERE
    pk.peername = p.name
    and p.active = 1;`

const SetLastSentVersionAndPokeSucceeded = `update
    peers
set
    lastSentVersion = $1,
    lastPokeSucceededAt = $2
WHERE
    name = $3;`

const LocalAddOrUpdateAssociation = `INSERT INTO
    local_threepid_associations (
        'medium',
        'address',
        'mxid',
        'ts',
        'notBefore',
        'notAfter'
    )
VALUES
    ($1, $2, $3, $4, $5, $6) ON CONFLICT (medium, address) DO
UPDATE
SET
    mxid = EXCLUDED.mxid,
    ts = EXCLUDED.ts,
    notBefore = EXCLUDED.notBefore,
    notAfter = EXCLUDED.notAfter;`

const GetAssociationsAfterId = `SELECT
    id,
    medium,
    address,
    mxid,
    ts,
    notBefore,
    notAfter
FROM
    local_threepid_associations
WHERE
    id > $1
ORDER by
    id asc
LIMIT
    $2;`

const GetLocal3pid = `SELECT
    COUNT(*)
FROM
    local_threepid_associations
WHERE
    medium = $1
    AND address = $2
    AND mxid = $3;`

const LocalRemoveAssociation = `INSERT INTO
    local_threepid_associations (
        'medium',
        'address',
        'mxid',
        'ts',
        'notBefore',
        'notAfter'
    )
VALUES
    ($1, $2, NULL, $3, NULL, NULL) ON CONFLICT (medium, address) DO
UPDATE
SET
    mxid = EXCLUDED.mxid,
    ts = EXCLUDED.ts,
    notBefore = EXCLUDED.notBefore,
    notAfter = EXCLUDED.notAfter;`

const GetTokenSession = `SELECT
    s.id,
    s.medium,
    s.address,
    s.clientSecret,
    s.validated,
    s.mtime,
    t.token,
    t.sendAttemptNumber
FROM
    threepid_validation_sessions s,
    threepid_token_auths t
WHERE
    s.medium = $1
    and s.address = $2
    and s.clientSecret = $3
	and t.validationSession = s.id;`

const CreateTokenSession = `insert into
    threepid_token_auths (validationSession, token, sendAttemptNumber)
values
    ($1, $2, $3)`

const AddValidationSession = `insert into
    threepid_validation_sessions (medium, address, clientSecret, mtime)
values
    ($1, $2, $3, $4) RETURNING id;`

const SetSendAttemptNumber = `update
    threepid_token_auths
set
    sendAttemptNumber = $1
WHERE
    id = $2`

const SetValidated = `update
    threepid_validation_sessions
set
    validated = $1
WHERE
    id = $2;`

const SetMtime = `update
    threepid_validation_sessions
set
    mtime = $1
WHERE
    id = $2;`

const GetSessionByID = `SELECT
    id,
    medium,
    address,
    clientSecret,
    validated,
    mtime
FROM
    threepid_validation_sessions
WHERE
    id = $1`

const GetTokenSessionByID = `SELECT
    s.id,
    s.medium,
    s.address,
    s.clientSecret,
    s.validated,
    s.mtime,
    t.token,
    t.sendAttemptNumber
FROM
    threepid_validation_sessions s,
    threepid_token_auths t
WHERE
    s.id = $1
    and t.validationSession = s.id;`

func Param(idx int) string {
	return fmt.Sprintf("$%d", idx)
}
