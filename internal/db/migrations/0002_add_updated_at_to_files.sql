-- +goose Up
-- +goose StatementBegin
ALTER TABLE files ADD COLUMN updated_at TIMESTAMP NOT NULL DEFAULT NOW();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE files DROP COLUMN updated_at;
-- +goose StatementEnd
