-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS files (
    "id" UUID PRIMARY KEY NOT NULL,
    "name" varchar(255) NOT NULL,
    "mime_type" varchar(255) NOT NULL,
    "size" bigint NOT NULL,
    "url" varchar(255) NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX idx_files_deleted_at ON files (deleted_at); -- +create index

ALTER TABLE users ADD COLUMN avatar_id UUID REFERENCES files(id) ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN avatar_id;
DROP TABLE IF EXISTS files;
-- +goose StatementEnd
