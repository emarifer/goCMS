-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS images (
    uuid VARCHAR(36) DEFAULT(UUID()) PRIMARY KEY,
    name TEXT NOT NULL,
    alt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS images;
-- +goose StatementEnd
