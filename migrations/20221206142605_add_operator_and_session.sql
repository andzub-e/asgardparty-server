-- +goose Up
-- +goose StatementBegin
ALTER TABLE history_records ADD COLUMN operator varchar(32);
ALTER TABLE history_records ADD COLUMN session_token uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE history_records DROP COLUMN operator;
ALTER TABLE history_records DROP COLUMN session_token;
-- +goose StatementEnd
