-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS history_records
(
    id                uuid        default gen_random_uuid() primary key,
    game              varchar(32),
    user_id           varchar(128),
    currency          varchar(10),
    start_balance     bigint,
    end_balance       bigint,
    wager             bigint,
    base_award        bigint,
    bonus_award       bigint,

    spin              json,
    restoring_indexes json,
    is_shown          bool,

    created_at        timestamptz default now(),
    updated_at        timestamptz default now()
);

CREATE TRIGGER history_records_updated_at
    BEFORE UPDATE
    ON history_records
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_column();

CREATE INDEX history_records_user_game_index ON history_records (user_id, game);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS history_records_user_game_index;
DROP TRIGGER IF EXISTS history_records_updated_at ON history_records;

DROP TABLE IF EXISTS history_records;
-- +goose StatementEnd
