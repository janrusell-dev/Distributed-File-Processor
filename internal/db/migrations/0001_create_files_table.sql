-- +goose Up
CREATE TABLE files (
    id UUID PRIMARY KEY,
    filename TEXT NOT NULL,
    size BIGINT NOT NULL,
    mime_type TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- +goose Down  
DROP TABLE IF EXISTS files;