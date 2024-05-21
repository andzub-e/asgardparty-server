-- +goose Up
-- +goose StatementBegin
ALTER TABLE history_records ADD COLUMN transaction_id uuid;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE history_records DROP COLUMN transaction_id;
-- +goose StatementEnd
