/*
 dumped from github.com/matrix-org/synapse commit ea00f18135ce30e8415526ce68585ea90da5b856
*/


CREATE TABLE IF NOT EXISTS access_tokens (
    id bigint primary key,
    user_id text NOT NULL,
    device_id text,
    token text NOT NULL,
    last_used bigint UNSIGNED,
    UNIQUE(token)
);
CREATE INDEX IF NOT EXISTS access_tokens_device_id ON access_tokens USING btree (user_id, device_id);


CREATE TABLE IF NOT EXISTS account_data (
    user_id text NOT NULL,
    account_data_type text NOT NULL,
    stream_id bigint NOT NULL,
    content text NOT NULL
    CONSTRAINT account_data_uniqueness UNIQUE (user_id, account_data_type)
);
CREATE INDEX IF NOT EXISTS account_data_stream_id ON account_data USING btree (user_id, stream_id);

CREATE TABLE IF NOT EXISTS account_data_max_stream_id (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    stream_id bigint NOT NULL,
    CONSTRAINT private_user_data_max_stream_id_lock_check CHECK ((lock = 'X'::bpchar))
);

CREATE TABLE IF NOT EXISTS application_services (
    id BIGINT PRIMARY KEY,
    url text,
    token text,
    hs_token text,
    sender text
    UNIQUE(token)
);


CREATE TABLE IF NOT EXISTS application_services_regex (
    id bigint NOT NULL,
    as_id bigint NOT NULL,
    namespace integer,
    regex text
);

CREATE TABLE IF NOT EXISTS application_services_state (
    as_id text NOT NULL,
    state character varying(5),
    last_txn integer
);


CREATE TABLE IF NOT EXISTS application_services_txns (
    as_id text NOT NULL,
    txn_id integer NOT NULL,
    event_ids text NOT NULL
);
CREATE INDEX IF NOT EXISTS application_services_txns_id ON application_services_txns USING btree (as_id);

CREATE TABLE IF NOT EXISTS applied_module_schemas (
    module_name text NOT NULL,
    file text NOT NULL
);

CREATE TABLE IF NOT EXISTS applied_schema_deltas (
    version integer NOT NULL,
    file text NOT NULL
);

CREATE TABLE IF NOT EXISTS appservice_room_list (
    appservice_id text NOT NULL,
    network_id text NOT NULL,
    room_id text NOT NULL
);
CREATE UNIQUE INDEX appservice_room_list_idx ON appservice_room_list USING btree (appservice_id, network_id, room_id);

CREATE TABLE IF NOT EXISTS appservice_stream_position (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    stream_ordering bigint,
    CONSTRAINT appservice_stream_position_lock_check CHECK ((lock = 'X'::bpchar))
);

CREATE TABLE IF NOT EXISTS background_updates (
    update_name text NOT NULL,
    progress_json text NOT NULL,
    depends_on text
);

CREATE TABLE IF NOT EXISTS blocked_rooms (
    room_id text NOT NULL,
    user_id text NOT NULL
);
CREATE UNIQUE INDEX blocked_rooms_idx ON blocked_rooms USING btree (room_id);

CREATE TABLE IF NOT EXISTS cache_invalidation_stream (
    stream_id bigint,
    cache_func text,
    keys text[],
    invalidation_ts bigint
);
CREATE INDEX IF NOT EXISTS cache_invalidation_stream_id ON cache_invalidation_stream USING btree (stream_id);

CREATE TABLE IF NOT EXISTS current_state_delta_stream (
    stream_id bigint NOT NULL,
    room_id text NOT NULL,
    type text NOT NULL,
    state_key text NOT NULL,
    event_id text,
    prev_event_id text
);
CREATE INDEX IF NOT EXISTS current_state_delta_stream_idx ON current_state_delta_stream USING btree (stream_id);

CREATE TABLE IF NOT EXISTS current_state_events (
    event_id text NOT NULL,
    room_id text NOT NULL,
    type text NOT NULL,
    state_key text NOT NULL
);
CREATE INDEX IF NOT EXISTS current_state_events_member_index ON current_state_events USING btree (state_key) WHERE (type = 'm.room.member'::text);

CREATE TABLE IF NOT EXISTS current_state_resets (
    event_stream_ordering bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS deleted_pushers (
    stream_id bigint NOT NULL,
    app_id text NOT NULL,
    pushkey text NOT NULL,
    user_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS deleted_pushers_stream_id ON deleted_pushers USING btree (stream_id);

CREATE TABLE IF NOT EXISTS destinations (
    destination text NOT NULL,
    retry_last_ts bigint,
    retry_interval integer
);

CREATE TABLE IF NOT EXISTS device_federation_inbox (
    origin text NOT NULL,
    message_id text NOT NULL,
    received_ts bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS device_federation_inbox_sender_id ON device_federation_inbox USING btree (origin, message_id);

CREATE TABLE IF NOT EXISTS device_federation_outbox (
    destination text NOT NULL,
    stream_id bigint NOT NULL,
    queued_ts bigint NOT NULL,
    messages_json text NOT NULL
);
CREATE INDEX IF NOT EXISTS device_federation_outbox_destination_id ON device_federation_outbox USING btree (destination, stream_id);
CREATE INDEX IF NOT EXISTS device_federation_outbox_id ON device_federation_outbox USING btree (stream_id);

CREATE TABLE IF NOT EXISTS device_inbox (
    user_id text NOT NULL,
    device_id text NOT NULL,
    stream_id bigint NOT NULL,
    message_json text NOT NULL
);
CREATE INDEX IF NOT EXISTS device_inbox_stream_id ON device_inbox USING btree (stream_id);
CREATE INDEX IF NOT EXISTS device_inbox_stream_id_user_id ON device_inbox USING btree (stream_id, user_id);
CREATE INDEX IF NOT EXISTS device_inbox_user_stream_id ON device_inbox USING btree (user_id, device_id, stream_id);

CREATE TABLE IF NOT EXISTS device_lists_outbound_last_success (
    destination text NOT NULL,
    user_id text NOT NULL,
    stream_id bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS device_lists_outbound_last_success_idx ON device_lists_outbound_last_success USING btree (destination, user_id, stream_id);

CREATE TABLE IF NOT EXISTS device_lists_outbound_pokes (
    destination text NOT NULL,
    stream_id bigint NOT NULL,
    user_id text NOT NULL,
    device_id text NOT NULL,
    sent boolean NOT NULL,
    ts bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS device_lists_outbound_pokes_id ON device_lists_outbound_pokes USING btree (destination, stream_id);
CREATE INDEX IF NOT EXISTS device_lists_outbound_pokes_stream ON device_lists_outbound_pokes USING btree (stream_id);
CREATE INDEX IF NOT EXISTS device_lists_outbound_pokes_user ON device_lists_outbound_pokes USING btree (destination, user_id);

CREATE TABLE IF NOT EXISTS device_lists_remote_cache (
    user_id text NOT NULL,
    device_id text NOT NULL,
    content text NOT NULL
);

CREATE TABLE IF NOT EXISTS device_lists_remote_extremeties (
    user_id text NOT NULL,
    stream_id text NOT NULL
);

CREATE TABLE IF NOT EXISTS device_lists_stream (
    stream_id bigint NOT NULL,
    user_id text NOT NULL,
    device_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS device_lists_stream_id ON device_lists_stream USING btree (stream_id, user_id);
CREATE INDEX IF NOT EXISTS device_lists_stream_user_id ON device_lists_stream USING btree (user_id, device_id);

CREATE TABLE IF NOT EXISTS device_max_stream_id (
    stream_id bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS devices (
    user_id text NOT NULL,
    device_id text NOT NULL,
    display_name text
);

CREATE TABLE IF NOT EXISTS e2e_device_keys_json (
    user_id text NOT NULL,
    device_id text NOT NULL,
    ts_added_ms bigint NOT NULL,
    key_json text NOT NULL
);

CREATE TABLE IF NOT EXISTS e2e_one_time_keys_json (
    user_id text NOT NULL,
    device_id text NOT NULL,
    algorithm text NOT NULL,
    key_id text NOT NULL,
    ts_added_ms bigint NOT NULL,
    key_json text NOT NULL
);

CREATE TABLE IF NOT EXISTS e2e_room_keys (
    user_id text NOT NULL,
    room_id text NOT NULL,
    session_id text NOT NULL,
    version bigint NOT NULL,
    first_message_index integer,
    forwarded_count integer,
    is_verified boolean,
    session_data text NOT NULL
);
CREATE UNIQUE INDEX e2e_room_keys_idx ON e2e_room_keys USING btree (user_id, room_id, session_id);

CREATE TABLE IF NOT EXISTS e2e_room_keys_versions (
    user_id text NOT NULL,
    version bigint NOT NULL,
    algorithm text NOT NULL,
    auth_data text NOT NULL,
    deleted smallint DEFAULT 0 NOT NULL
);
CREATE UNIQUE INDEX e2e_room_keys_versions_idx ON e2e_room_keys_versions USING btree (user_id, version);

CREATE TABLE IF NOT EXISTS erased_users (
    user_id text NOT NULL
);
CREATE UNIQUE INDEX erased_users_user ON erased_users USING btree (user_id);

CREATE TABLE IF NOT EXISTS event_auth (
    event_id text NOT NULL,
    auth_id text NOT NULL,
    room_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS evauth_edges_id ON event_auth USING btree (event_id);

CREATE TABLE IF NOT EXISTS event_backward_extremities (
    event_id text NOT NULL,
    room_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS ev_b_extrem_id ON event_backward_extremities USING btree (event_id);
CREATE INDEX IF NOT EXISTS ev_b_extrem_room ON event_backward_extremities USING btree (room_id);

CREATE TABLE IF NOT EXISTS event_content_hashes (
    event_id text,
    algorithm text,
    hash bytea
);

CREATE TABLE IF NOT EXISTS event_destinations (
    event_id text NOT NULL,
    destination text NOT NULL,
    delivered_ts bigint DEFAULT 0
);


CREATE TABLE IF NOT EXISTS event_edge_hashes (
    event_id text,
    prev_event_id text,
    algorithm text,
    hash bytea
);


CREATE TABLE IF NOT EXISTS event_edges (
    event_id text NOT NULL,
    prev_event_id text NOT NULL,
    room_id text NOT NULL,
    is_state boolean NOT NULL
);
CREATE INDEX IF NOT EXISTS ev_edges_id ON event_edges USING btree (event_id);
CREATE INDEX IF NOT EXISTS ev_edges_prev_id ON event_edges USING btree (prev_event_id);

CREATE TABLE IF NOT EXISTS event_forward_extremities (
    event_id text NOT NULL,
    room_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS ev_extrem_id ON event_forward_extremities USING btree (event_id);
CREATE INDEX IF NOT EXISTS ev_extrem_room ON event_forward_extremities USING btree (room_id);

CREATE TABLE IF NOT EXISTS event_json (
    event_id text NOT NULL,
    room_id text NOT NULL,
    internal_metadata text NOT NULL,
    json text NOT NULL
);
CREATE INDEX IF NOT EXISTS event_json_room_id ON event_json USING btree (room_id);

CREATE TABLE IF NOT EXISTS event_push_actions (
    room_id text NOT NULL,
    event_id text NOT NULL,
    user_id text NOT NULL,
    profile_tag character varying(32),
    actions text NOT NULL,
    topological_ordering bigint,
    stream_ordering bigint,
    notif smallint,
    highlight smallint
);
CREATE INDEX IF NOT EXISTS event_push_actions_highlights_index ON event_push_actions USING btree (user_id, room_id, topological_ordering, stream_ordering) WHERE (highlight = 1);
CREATE INDEX IF NOT EXISTS event_push_actions_rm_tokens ON event_push_actions USING btree (user_id, room_id, topological_ordering, stream_ordering);
CREATE INDEX IF NOT EXISTS event_push_actions_room_id_user_id ON event_push_actions USING btree (room_id, user_id);
CREATE INDEX IF NOT EXISTS event_push_actions_stream_ordering ON event_push_actions USING btree (stream_ordering, user_id);
CREATE INDEX IF NOT EXISTS event_push_actions_u_highlight ON event_push_actions USING btree (user_id, stream_ordering);

CREATE TABLE IF NOT EXISTS event_push_actions_staging (
    event_id text NOT NULL,
    user_id text NOT NULL,
    actions text NOT NULL,
    notif smallint NOT NULL,
    highlight smallint NOT NULL
);
CREATE INDEX IF NOT EXISTS event_push_actions_staging_id ON event_push_actions_staging USING btree (event_id);

CREATE TABLE IF NOT EXISTS event_push_summary (
    user_id text NOT NULL,
    room_id text NOT NULL,
    notif_count bigint NOT NULL,
    stream_ordering bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS event_push_summary_user_rm ON event_push_summary USING btree (user_id, room_id);

CREATE TABLE IF NOT EXISTS event_push_summary_stream_ordering (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    stream_ordering bigint NOT NULL,
    CONSTRAINT event_push_summary_stream_ordering_lock_check CHECK ((lock = 'X'::bpchar))
);

CREATE TABLE IF NOT EXISTS event_reference_hashes (
    event_id text,
    algorithm text,
    hash bytea
);
CREATE INDEX IF NOT EXISTS event_reference_hashes_id ON event_reference_hashes USING btree (event_id);

CREATE TABLE IF NOT EXISTS event_reports (
    id bigint NOT NULL,
    received_ts bigint NOT NULL,
    room_id text NOT NULL,
    event_id text NOT NULL,
    user_id text NOT NULL,
    reason text,
    content text
);

CREATE TABLE IF NOT EXISTS event_search (
    event_id text,
    room_id text,
    sender text,
    key text,
    vector tsvector,
    origin_server_ts bigint,
    stream_ordering bigint
);
CREATE INDEX IF NOT EXISTS event_search_ev_ridx ON event_search USING btree (room_id);
CREATE UNIQUE INDEX event_search_event_id_idx ON event_search USING btree (event_id);
CREATE INDEX IF NOT EXISTS event_search_fts_idx ON event_search USING gin (vector);

CREATE TABLE IF NOT EXISTS event_signatures (
    event_id text,
    signature_name text,
    key_id text,
    signature bytea
);

CREATE TABLE IF NOT EXISTS event_to_state_groups (
    event_id text NOT NULL,
    state_group bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS events (
    stream_ordering integer NOT NULL,
    topological_ordering bigint NOT NULL,
    event_id text NOT NULL,
    type text NOT NULL,
    room_id text NOT NULL,
    content text,
    unrecognized_keys text,
    processed boolean NOT NULL,
    outlier boolean NOT NULL,
    depth bigint DEFAULT 0 NOT NULL,
    origin_server_ts bigint,
    received_ts bigint,
    sender text,
    contains_url boolean
);
CREATE INDEX IF NOT EXISTS events_order_room ON events USING btree (room_id, topological_ordering, stream_ordering);
CREATE INDEX IF NOT EXISTS events_room_stream ON events USING btree (room_id, stream_ordering);
CREATE INDEX IF NOT EXISTS events_ts ON events USING btree (origin_server_ts, stream_ordering);
CREATE INDEX IF NOT EXISTS event_contains_url_index ON events USING btree (room_id, topological_ordering, stream_ordering) WHERE ((contains_url = true) AND (outlier = false));

CREATE TABLE IF NOT EXISTS ex_outlier_stream (
    event_stream_ordering bigint NOT NULL,
    event_id text NOT NULL,
    state_group bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS federation_stream_position (
    type text NOT NULL,
    stream_id integer NOT NULL
);

CREATE TABLE IF NOT EXISTS feedback (
    event_id text NOT NULL,
    feedback_type text,
    target_event_id text,
    sender text,
    room_id text
);

CREATE TABLE IF NOT EXISTS group_attestations_remote (
    group_id text NOT NULL,
    user_id text NOT NULL,
    valid_until_ms bigint NOT NULL,
    attestation_json text NOT NULL
);
CREATE INDEX IF NOT EXISTS group_attestations_remote_g_idx ON group_attestations_remote USING btree (group_id, user_id);
CREATE INDEX IF NOT EXISTS group_attestations_remote_u_idx ON group_attestations_remote USING btree (user_id);
CREATE INDEX IF NOT EXISTS group_attestations_remote_v_idx ON group_attestations_remote USING btree (valid_until_ms);

CREATE TABLE IF NOT EXISTS group_attestations_renewals (
    group_id text NOT NULL,
    user_id text NOT NULL,
    valid_until_ms bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS group_attestations_renewals_g_idx ON group_attestations_renewals USING btree (group_id, user_id);
CREATE INDEX IF NOT EXISTS group_attestations_renewals_u_idx ON group_attestations_renewals USING btree (user_id);
CREATE INDEX IF NOT EXISTS group_attestations_renewals_v_idx ON group_attestations_renewals USING btree (valid_until_ms);

CREATE TABLE IF NOT EXISTS group_invites (
    group_id text NOT NULL,
    user_id text NOT NULL
);
CREATE UNIQUE INDEX group_invites_g_idx ON group_invites USING btree (group_id, user_id);
CREATE INDEX IF NOT EXISTS group_invites_u_idx ON group_invites USING btree (user_id);


CREATE TABLE IF NOT EXISTS group_roles (
    group_id text NOT NULL,
    role_id text NOT NULL,
    profile text NOT NULL,
    is_public boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS group_room_categories (
    group_id text NOT NULL,
    category_id text NOT NULL,
    profile text NOT NULL,
    is_public boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS group_rooms (
    group_id text NOT NULL,
    room_id text NOT NULL,
    is_public boolean NOT NULL
);
CREATE UNIQUE INDEX group_rooms_g_idx ON group_rooms USING btree (group_id, room_id);
CREATE INDEX IF NOT EXISTS group_rooms_r_idx ON group_rooms USING btree (room_id);

CREATE TABLE IF NOT EXISTS group_summary_roles (
    group_id text NOT NULL,
    role_id text NOT NULL,
    role_order bigint NOT NULL,
    CONSTRAINT group_summary_roles_role_order_check CHECK ((role_order > 0))
);


CREATE TABLE IF NOT EXISTS group_summary_room_categories (
    group_id text NOT NULL,
    category_id text NOT NULL,
    cat_order bigint NOT NULL,
    CONSTRAINT group_summary_room_categories_cat_order_check CHECK ((cat_order > 0))
);

CREATE TABLE IF NOT EXISTS group_summary_rooms (
    group_id text NOT NULL,
    room_id text NOT NULL,
    category_id text NOT NULL,
    room_order bigint NOT NULL,
    is_public boolean NOT NULL,
    CONSTRAINT group_summary_rooms_room_order_check CHECK ((room_order > 0))
);
CREATE UNIQUE INDEX group_summary_rooms_g_idx ON group_summary_rooms USING btree (group_id, room_id, category_id);

CREATE TABLE IF NOT EXISTS group_summary_users (
    group_id text NOT NULL,
    user_id text NOT NULL,
    role_id text NOT NULL,
    user_order bigint NOT NULL,
    is_public boolean NOT NULL
);
CREATE INDEX IF NOT EXISTS group_summary_users_g_idx ON group_summary_users USING btree (group_id);


CREATE TABLE IF NOT EXISTS group_users (
    group_id text NOT NULL,
    user_id text NOT NULL,
    is_admin boolean NOT NULL,
    is_public boolean NOT NULL
);
CREATE UNIQUE INDEX group_users_g_idx ON group_users USING btree (group_id, user_id);
CREATE INDEX IF NOT EXISTS group_users_u_idx ON group_users USING btree (user_id);

CREATE TABLE IF NOT EXISTS groups (
    group_id text NOT NULL,
    name text,
    avatar_url text,
    short_description text,
    long_description text,
    is_public boolean NOT NULL,
    join_policy text DEFAULT 'invite'::text NOT NULL
);
CREATE UNIQUE INDEX groups_idx ON groups USING btree (group_id);


CREATE TABLE IF NOT EXISTS guest_access (
    event_id text NOT NULL,
    room_id text NOT NULL,
    guest_access text NOT NULL
);

CREATE TABLE IF NOT EXISTS history_visibility (
    event_id text NOT NULL,
    room_id text NOT NULL,
    history_visibility text NOT NULL
);

CREATE TABLE IF NOT EXISTS local_group_membership (
    group_id text NOT NULL,
    user_id text NOT NULL,
    is_admin boolean NOT NULL,
    membership text NOT NULL,
    is_publicised boolean NOT NULL,
    content text NOT NULL
);
CREATE INDEX IF NOT EXISTS local_group_membership_g_idx ON local_group_membership USING btree (group_id);
CREATE INDEX IF NOT EXISTS local_group_membership_u_idx ON local_group_membership USING btree (user_id, group_id);

CREATE TABLE IF NOT EXISTS local_group_updates (
    stream_id bigint NOT NULL,
    group_id text NOT NULL,
    user_id text NOT NULL,
    type text NOT NULL,
    content text NOT NULL
);

CREATE TABLE IF NOT EXISTS local_invites (
    stream_id bigint NOT NULL,
    inviter text NOT NULL,
    invitee text NOT NULL,
    event_id text NOT NULL,
    room_id text NOT NULL,
    locally_rejected text,
    replaced_by text
);
CREATE INDEX IF NOT EXISTS local_invites_for_user_idx ON local_invites USING btree (invitee, locally_rejected, replaced_by, room_id);
CREATE INDEX IF NOT EXISTS local_invites_id ON local_invites USING btree (stream_id);

CREATE TABLE IF NOT EXISTS local_media_repository (
    media_id text,
    media_type text,
    media_length integer,
    created_ts bigint,
    upload_name text,
    user_id text,
    quarantined_by text,
    url_cache text,
    last_access_ts bigint
);
CREATE INDEX IF NOT EXISTS local_media_repository_url_idx ON local_media_repository USING btree (created_ts) WHERE (url_cache IS NOT NULL);

CREATE TABLE IF NOT EXISTS local_media_repository_thumbnails (
    media_id text,
    thumbnail_width integer,
    thumbnail_height integer,
    thumbnail_type text,
    thumbnail_method text,
    thumbnail_length integer
);
CREATE INDEX IF NOT EXISTS local_media_repository_thumbnails_media_id ON local_media_repository_thumbnails USING btree (media_id);

CREATE TABLE IF NOT EXISTS local_media_repository_url_cache (
    url text,
    response_code integer,
    etag text,
    expires_ts bigint,
    og text,
    media_id text,
    download_ts bigint
);
CREATE INDEX IF NOT EXISTS local_media_repository_url_cache_by_url_download_ts ON local_media_repository_url_cache USING btree (url, download_ts);
CREATE INDEX IF NOT EXISTS local_media_repository_url_cache_expires_idx ON local_media_repository_url_cache USING btree (expires_ts);
CREATE INDEX IF NOT EXISTS local_media_repository_url_cache_media_idx ON local_media_repository_url_cache USING btree (media_id);

CREATE TABLE IF NOT EXISTS monthly_active_users (
    user_id text NOT NULL,
    "timestamp" bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS monthly_active_users_time_stamp ON monthly_active_users USING btree ("timestamp");
CREATE UNIQUE INDEX monthly_active_users_users ON monthly_active_users USING btree (user_id);

CREATE TABLE IF NOT EXISTS open_id_tokens (
    token text NOT NULL,
    ts_valid_until_ms bigint NOT NULL,
    user_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS open_id_tokens_ts_valid_until_ms ON open_id_tokens USING btree (ts_valid_until_ms);

CREATE TABLE IF NOT EXISTS presence (
    user_id text NOT NULL,
    state character varying(20),
    status_msg text,
    mtime bigint
);

CREATE TABLE IF NOT EXISTS presence_allow_inbound (
    observed_user_id text NOT NULL,
    observer_user_id text NOT NULL
);

CREATE TABLE IF NOT EXISTS presence_list (
    user_id text NOT NULL,
    observed_user_id text NOT NULL,
    accepted boolean NOT NULL
);
CREATE INDEX IF NOT EXISTS presence_list_user_id ON presence_list USING btree (user_id);


CREATE TABLE IF NOT EXISTS presence_stream (
    stream_id bigint,
    user_id text,
    state text,
    last_active_ts bigint,
    last_federation_update_ts bigint,
    last_user_sync_ts bigint,
    status_msg text,
    currently_active boolean
);
CREATE INDEX IF NOT EXISTS presence_stream_id ON presence_stream USING btree (stream_id, user_id);
CREATE INDEX IF NOT EXISTS presence_stream_user_id ON presence_stream USING btree (user_id);

CREATE TABLE IF NOT EXISTS profiles (
    user_id text NOT NULL,
    displayname text,
    avatar_url text
);

CREATE TABLE IF NOT EXISTS public_room_list_stream (
    stream_id bigint NOT NULL,
    room_id text NOT NULL,
    visibility boolean NOT NULL,
    appservice_id text,
    network_id text
);
CREATE INDEX IF NOT EXISTS public_room_list_stream_idx ON public_room_list_stream USING btree (stream_id);
CREATE INDEX IF NOT EXISTS public_room_list_stream_rm_idx ON public_room_list_stream USING btree (room_id, stream_id);

CREATE TABLE IF NOT EXISTS push_rules (
    id bigint NOT NULL,
    user_name text NOT NULL,
    rule_id text NOT NULL,
    priority_class smallint NOT NULL,
    priority integer DEFAULT 0 NOT NULL,
    conditions text NOT NULL,
    actions text NOT NULL
);
CREATE INDEX IF NOT EXISTS push_rules_user_name ON push_rules USING btree (user_name);

CREATE TABLE IF NOT EXISTS push_rules_enable (
    id bigint NOT NULL,
    user_name text NOT NULL,
    rule_id text NOT NULL,
    enabled smallint
);
CREATE INDEX IF NOT EXISTS push_rules_enable_user_name ON push_rules_enable USING btree (user_name);

CREATE TABLE IF NOT EXISTS push_rules_stream (
    stream_id bigint NOT NULL,
    event_stream_ordering bigint NOT NULL,
    user_id text NOT NULL,
    rule_id text NOT NULL,
    op text NOT NULL,
    priority_class smallint,
    priority integer,
    conditions text,
    actions text
);
CREATE INDEX IF NOT EXISTS push_rules_stream_id ON push_rules_stream USING btree (stream_id);
CREATE INDEX IF NOT EXISTS push_rules_stream_user_stream_id ON push_rules_stream USING btree (user_id, stream_id);


CREATE TABLE IF NOT EXISTS pusher_throttle (
    pusher bigint NOT NULL,
    room_id text NOT NULL,
    last_sent_ts bigint,
    throttle_ms bigint
);

CREATE TABLE IF NOT EXISTS pushers (
    id bigint NOT NULL,
    user_name text NOT NULL,
    access_token bigint,
    profile_tag text NOT NULL,
    kind text NOT NULL,
    app_id text NOT NULL,
    app_display_name text NOT NULL,
    device_display_name text NOT NULL,
    pushkey text NOT NULL,
    ts bigint NOT NULL,
    lang text,
    data text,
    last_stream_ordering integer,
    last_success bigint,
    failing_since bigint
);

CREATE TABLE IF NOT EXISTS ratelimit_override (
    user_id text NOT NULL,
    messages_per_second bigint,
    burst_count bigint
);
CREATE UNIQUE INDEX ratelimit_override_idx ON ratelimit_override USING btree (user_id);


CREATE TABLE IF NOT EXISTS receipts_graph (
    room_id text NOT NULL,
    receipt_type text NOT NULL,
    user_id text NOT NULL,
    event_ids text NOT NULL,
    data text NOT NULL
);


CREATE TABLE IF NOT EXISTS receipts_linearized (
    stream_id bigint NOT NULL,
    room_id text NOT NULL,
    receipt_type text NOT NULL,
    user_id text NOT NULL,
    event_id text NOT NULL,
    data text NOT NULL
);
CREATE INDEX IF NOT EXISTS receipts_linearized_id ON receipts_linearized USING btree (stream_id);
CREATE INDEX IF NOT EXISTS receipts_linearized_room_stream ON receipts_linearized USING btree (room_id, stream_id);
CREATE INDEX IF NOT EXISTS receipts_linearized_user ON receipts_linearized USING btree (user_id);


CREATE TABLE IF NOT EXISTS received_transactions (
    transaction_id text,
    origin text,
    ts bigint,
    response_code integer,
    response_json bytea,
    has_been_referenced smallint DEFAULT 0
);
CREATE INDEX IF NOT EXISTS received_transactions_ts ON received_transactions USING btree (ts);

CREATE TABLE IF NOT EXISTS redactions (
    event_id text NOT NULL,
    redacts text NOT NULL
);
CREATE INDEX IF NOT EXISTS redactions_redacts ON redactions USING btree (redacts);

CREATE TABLE IF NOT EXISTS rejections (
    event_id text NOT NULL,
    reason text NOT NULL,
    last_check text NOT NULL
);

CREATE TABLE IF NOT EXISTS remote_media_cache (
    media_origin text,
    media_id text,
    media_type text,
    created_ts bigint,
    upload_name text,
    media_length integer,
    filesystem_id text,
    last_access_ts bigint,
    quarantined_by text
);


CREATE TABLE IF NOT EXISTS remote_media_cache_thumbnails (
    media_origin text,
    media_id text,
    thumbnail_width integer,
    thumbnail_height integer,
    thumbnail_method text,
    thumbnail_type text,
    thumbnail_length integer,
    filesystem_id text
);


CREATE TABLE IF NOT EXISTS remote_profile_cache (
    user_id text NOT NULL,
    displayname text,
    avatar_url text,
    last_check bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS remote_profile_cache_time ON remote_profile_cache USING btree (last_check);
CREATE UNIQUE INDEX remote_profile_cache_user_id ON remote_profile_cache USING btree (user_id);

CREATE TABLE IF NOT EXISTS room_account_data (
    user_id text NOT NULL,
    room_id text NOT NULL,
    account_data_type text NOT NULL,
    stream_id bigint NOT NULL,
    content text NOT NULL
);
CREATE INDEX IF NOT EXISTS room_account_data_stream_id ON room_account_data USING btree (user_id, stream_id);

CREATE TABLE IF NOT EXISTS room_alias_servers (
    room_alias text NOT NULL,
    server text NOT NULL
);
CREATE INDEX IF NOT EXISTS room_alias_servers_alias ON room_alias_servers USING btree (room_alias);

CREATE TABLE IF NOT EXISTS room_aliases (
    room_alias text NOT NULL,
    room_id text NOT NULL,
    creator text
);
CREATE INDEX IF NOT EXISTS room_aliases_id ON room_aliases USING btree (room_id);

CREATE TABLE IF NOT EXISTS room_depth (
    room_id text NOT NULL,
    min_depth integer NOT NULL
);
CREATE INDEX IF NOT EXISTS room_depth_room ON room_depth USING btree (room_id);

CREATE TABLE IF NOT EXISTS room_hosts (
    room_id text NOT NULL,
    host text NOT NULL
);


CREATE TABLE IF NOT EXISTS room_memberships (
    event_id text NOT NULL,
    user_id text NOT NULL,
    sender text NOT NULL,
    room_id text NOT NULL,
    membership text NOT NULL,
    forgotten integer DEFAULT 0,
    display_name text,
    avatar_url text
);
CREATE INDEX IF NOT EXISTS room_memberships_room_id ON room_memberships USING btree (room_id);
CREATE INDEX IF NOT EXISTS room_memberships_user_id ON room_memberships USING btree (user_id);

CREATE TABLE IF NOT EXISTS room_names (
    event_id text NOT NULL,
    room_id text NOT NULL,
    name text NOT NULL
);
CREATE INDEX IF NOT EXISTS room_names_room_id ON room_names USING btree (room_id);

CREATE TABLE IF NOT EXISTS room_tags (
    user_id text NOT NULL,
    room_id text NOT NULL,
    tag text NOT NULL,
    content text NOT NULL
);

CREATE TABLE IF NOT EXISTS room_tags_revisions (
    user_id text NOT NULL,
    room_id text NOT NULL,
    stream_id bigint NOT NULL
);

CREATE TABLE IF NOT EXISTS rooms (
    room_id text NOT NULL,
    is_public boolean,
    creator text
);
CREATE INDEX IF NOT EXISTS public_room_index ON rooms USING btree (is_public);


CREATE TABLE IF NOT EXISTS schema_version (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    version integer NOT NULL,
    upgraded boolean NOT NULL,
    CONSTRAINT schema_version_lock_check CHECK ((lock = 'X'::bpchar))
);


CREATE TABLE IF NOT EXISTS server_keys_json (
    server_name text NOT NULL,
    key_id text NOT NULL,
    from_server text NOT NULL,
    ts_added_ms bigint NOT NULL,
    ts_valid_until_ms bigint NOT NULL,
    key_json bytea NOT NULL
);

CREATE TABLE IF NOT EXISTS server_signature_keys (
    server_name text,
    key_id text,
    from_server text,
    ts_added_ms bigint,
    verify_key bytea
);

CREATE TABLE IF NOT EXISTS server_tls_certificates (
    server_name text,
    fingerprint text,
    from_server text,
    ts_added_ms bigint,
    tls_certificate bytea
);

CREATE TABLE IF NOT EXISTS state_events (
    event_id text NOT NULL,
    room_id text NOT NULL,
    type text NOT NULL,
    state_key text NOT NULL,
    prev_state text
);

CREATE TABLE IF NOT EXISTS state_forward_extremities (
    event_id text NOT NULL,
    room_id text NOT NULL,
    type text NOT NULL,
    state_key text NOT NULL
);
CREATE INDEX IF NOT EXISTS st_extrem_keys ON state_forward_extremities USING btree (room_id, type, state_key);

CREATE TABLE IF NOT EXISTS state_group_edges (
    state_group bigint NOT NULL,
    prev_state_group bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS state_group_edges_idx ON state_group_edges USING btree (state_group);
CREATE INDEX IF NOT EXISTS state_group_edges_prev_idx ON state_group_edges USING btree (prev_state_group);

CREATE SEQUENCE state_group_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;



CREATE TABLE IF NOT EXISTS state_groups (
    id bigint NOT NULL,
    room_id text NOT NULL,
    event_id text NOT NULL
);

CREATE TABLE IF NOT EXISTS state_groups_state (
    state_group bigint NOT NULL,
    room_id text NOT NULL,
    type text NOT NULL,
    state_key text NOT NULL,
    event_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS state_groups_state_id ON state_groups_state USING btree (state_group);

CREATE TABLE IF NOT EXISTS stats_reporting (
    reported_stream_token integer,
    reported_time bigint
);

CREATE TABLE IF NOT EXISTS stream_ordering_to_exterm (
    stream_ordering bigint NOT NULL,
    room_id text NOT NULL,
    event_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS stream_ordering_to_exterm_idx ON stream_ordering_to_exterm USING btree (stream_ordering);
CREATE INDEX IF NOT EXISTS stream_ordering_to_exterm_rm_idx ON stream_ordering_to_exterm USING btree (room_id, stream_ordering);

CREATE TABLE IF NOT EXISTS threepid_guest_access_tokens (
    medium text,
    address text,
    guest_access_token text,
    first_inviter text
);
CREATE UNIQUE INDEX threepid_guest_access_tokens_index ON threepid_guest_access_tokens USING btree (medium, address);

CREATE TABLE IF NOT EXISTS topics (
    event_id text NOT NULL,
    room_id text NOT NULL,
    topic text NOT NULL
);
CREATE INDEX IF NOT EXISTS topics_room_id ON topics USING btree (room_id);

CREATE TABLE IF NOT EXISTS transaction_id_to_pdu (
    transaction_id integer,
    destination text,
    pdu_id text,
    pdu_origin text
);
CREATE INDEX IF NOT EXISTS transaction_id_to_pdu_dest ON transaction_id_to_pdu USING btree (destination);


CREATE TABLE IF NOT EXISTS user_daily_visits (
    user_id text NOT NULL,
    device_id text,
    "timestamp" bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS user_daily_visits_ts_idx ON user_daily_visits USING btree ("timestamp");
CREATE INDEX IF NOT EXISTS user_daily_visits_uts_idx ON user_daily_visits USING btree (user_id, "timestamp");

CREATE TABLE IF NOT EXISTS user_directory (
    user_id text NOT NULL,
    room_id text,
    display_name text,
    avatar_url text
);
CREATE INDEX IF NOT EXISTS user_directory_room_idx ON user_directory USING btree (room_id);
CREATE UNIQUE INDEX user_directory_user_idx ON user_directory USING btree (user_id);

CREATE TABLE IF NOT EXISTS user_directory_search (
    user_id text NOT NULL,
    vector tsvector
);
CREATE INDEX IF NOT EXISTS user_directory_search_fts_idx ON user_directory_search USING gin (vector);
CREATE UNIQUE INDEX user_directory_search_user_idx ON user_directory_search USING btree (user_id);

CREATE TABLE IF NOT EXISTS user_directory_stream_pos (
    lock character(1) DEFAULT 'X'::bpchar NOT NULL,
    stream_id bigint,
    CONSTRAINT user_directory_stream_pos_lock_check CHECK ((lock = 'X'::bpchar))
);

CREATE TABLE IF NOT EXISTS user_filters (
    user_id text,
    filter_id bigint,
    filter_json bytea
);
CREATE INDEX IF NOT EXISTS user_filters_by_user_id_filter_id ON user_filters USING btree (user_id, filter_id);

CREATE TABLE IF NOT EXISTS user_ips (
    user_id text NOT NULL,
    access_token text NOT NULL,
    device_id text,
    ip text NOT NULL,
    user_agent text NOT NULL,
    last_seen bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS user_ips_device_id ON user_ips USING btree (user_id, device_id, last_seen);
CREATE INDEX IF NOT EXISTS user_ips_user_ip ON user_ips USING btree (user_id, access_token, ip);

CREATE TABLE IF NOT EXISTS user_threepids (
    user_id text NOT NULL,
    medium text NOT NULL,
    address text NOT NULL,
    validated_at bigint NOT NULL,
    added_at bigint NOT NULL
);
CREATE INDEX IF NOT EXISTS user_threepids_medium_address ON user_threepids USING btree (medium, address);
CREATE INDEX IF NOT EXISTS user_threepids_user_id ON user_threepids USING btree (user_id);

CREATE TABLE IF NOT EXISTS users (
    name text,
    password_hash text,
    creation_ts bigint,
    admin smallint DEFAULT 0 NOT NULL,
    upgrade_ts bigint,
    is_guest smallint DEFAULT 0 NOT NULL,
    appservice_id text,
    consent_version text,
    consent_server_notice_sent text,
    user_type text DEFAULT NULL,
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS users_in_public_rooms (
    user_id text NOT NULL,
    room_id text NOT NULL
);
CREATE INDEX IF NOT EXISTS users_in_public_rooms_room_idx ON users_in_public_rooms USING btree (room_id);
CREATE UNIQUE INDEX users_in_public_rooms_user_idx ON users_in_public_rooms USING btree (user_id);

CREATE TABLE IF NOT EXISTS users_pending_deactivation (
    user_id text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_who_share_rooms (
    user_id text NOT NULL,
    other_user_id text NOT NULL,
    room_id text NOT NULL,
    share_private boolean NOT NULL
);
CREATE INDEX IF NOT EXISTS users_who_share_rooms_o_idx ON users_who_share_rooms USING btree (other_user_id);
CREATE INDEX IF NOT EXISTS users_who_share_rooms_r_idx ON users_who_share_rooms USING btree (room_id);
CREATE UNIQUE INDEX users_who_share_rooms_u_idx ON users_who_share_rooms USING btree (user_id, other_user_id);






