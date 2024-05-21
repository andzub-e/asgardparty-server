-- +goose Up
-- +goose StatementBegin
ALTER TABLE history_records RENAME COLUMN operator TO integrator;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE history_records RENAME COLUMN integrator TO operator;
-- +goose StatementEnd
