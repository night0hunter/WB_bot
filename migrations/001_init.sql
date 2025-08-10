-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS bookings (
	id bigserial NOT NULL,
	chat_id int8 NOT NULL,
	from_date date NOT NULL,
	to_date date NOT NULL,
	draft_id uuid NOT NULL,
	protection int4 NOT NULL,
	warehouse int4 NOT NULL,
	coeff_limit int4 NOT NULL,
	supply_type int4 NOT NULL,
	created_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	updated_at timestamptz DEFAULT CURRENT_TIMESTAMP NOT NULL,
	preorder_id int8 NULL,
	is_active bool DEFAULT true NULL,
	CONSTRAINT bookings_pkey PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS state (
	chat_id int8 NOT NULL,
	seq_name int4 NULL,
	com_name int4 NULL,
	mes_id int4 NULL,
	info jsonb NULL,
	keyboard jsonb NULL,
	CONSTRAINT state_pkey PRIMARY KEY (chat_id)
);

CREATE TABLE IF NOT EXISTS supplies (
	chat_id int8 NOT NULL,
	from_date date NOT NULL,
	to_date date NOT NULL,
	warehouse int4 NOT NULL,
	coeff_limit int4 NOT NULL,
	supply_type int4 NOT NULL,
	is_active bool DEFAULT true NOT NULL,
	created_at timestamptz DEFAULT now() NOT NULL,
	updated_at timestamptz DEFAULT now() NOT NULL,
	id bigserial NOT NULL,
	tracking_date timestamptz NULL,
	CONSTRAINT supplies_pkey PRIMARY KEY (id)
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
DROP TABLE IF EXISTS supplies;
DROP TABLE IF EXISTS state;
-- +goose StatementEnd