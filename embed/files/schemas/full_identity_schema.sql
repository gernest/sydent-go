CREATE TABLE IF NOT EXISTS invite_tokens (
    id bigserial primary key,
    medium varchar(16) not null,
    address varchar(256) not null,
    room_id varchar(256) not null,
    sender varchar(256) not null,
    token varchar(256) not null,
    received_ts bigint,
    -- When the invite was received by us from the homeserver
    sent_ts bigint -- When the token was sent by us to the user
);
CREATE INDEX IF NOT EXISTS invite_token_medium_address on invite_tokens(medium, address);
CREATE INDEX IF NOT EXISTS invite_token_token on invite_tokens(token);

CREATE TABLE IF NOT EXISTS ephemeral_public_keys(
    id bigserial primary key,
    public_key varchar(256) not null,
    verify_count bigint default 0,
    persistence_ts bigint
);
CREATE UNIQUE INDEX IF NOT EXISTS ephemeral_public_keys_index on ephemeral_public_keys(public_key);

CREATE TABLE IF NOT EXISTS peers (
    id bigserial primary key,
    name varchar(255) not null,
    port integer default null,
    lastSentVersion integer,
    lastPokeSucceededAt integer,
    active integer not null default 0
);
CREATE UNIQUE INDEX IF NOT EXISTS name on peers(name);

CREATE TABLE IF NOT EXISTS peer_pubkeys (
    id bigserial primary key,
    peername varchar(255) not null,
    alg varchar(16) not null,
    key text not null,
    foreign key (peername) references peers (name)
);
CREATE UNIQUE INDEX IF NOT EXISTS peername_alg on peer_pubkeys(peername, alg);

CREATE TABLE IF NOT EXISTS local_threepid_associations (
    id bigserial primary key,
    medium varchar(16) not null,
    address varchar(256) not null,
    mxid varchar(256) not null,
    ts integer not null,
    notBefore bigint not null,
    notAfter bigint not null
);
CREATE UNIQUE INDEX IF NOT EXISTS medium_address on local_threepid_associations(medium, address);

CREATE TABLE IF NOT EXISTS global_threepid_associations (
    id bigserial primary key,
    medium varchar(16) not null,
    address varchar(256) not null,
    mxid varchar(256) not null,
    ts integer not null,
    notBefore bigint not null,
    notAfter bigint not null,
    originServer varchar(255) not null,
    originId integer not null,
    sgAssoc text not null
);
CREATE INDEX IF NOT EXISTS medium_lower_address on global_threepid_associations (medium, lower(address));
CREATE UNIQUE INDEX IF NOT EXISTS originServer_originId on global_threepid_associations (originServer, originId);COMMIT;

CREATE TABLE IF NOT EXISTS threepid_validation_sessions (
    id bigserial primary key,
    medium varchar(16) not null,
    address varchar(256) not null,
    clientSecret varchar(32) not null,
    validated int default 0,
    mtime bigint not null
);
CREATE TABLE IF NOT EXISTS threepid_token_auths (
    id bigserial primary key,
    validationSession integer not null,
    token varchar(32) not null,
    sendAttemptNumber integer not null,
    foreign key (validationSession) references threepid_validation_sessions(id)
);