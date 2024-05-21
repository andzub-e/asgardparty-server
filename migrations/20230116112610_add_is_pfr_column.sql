-- +goose Up
-- +goose StatementBegin
ALTER TABLE history_records ADD COLUMN is_pfr bool;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE history_records DROP COLUMN is_pfr;
-- +goose StatementEnd
