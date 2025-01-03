-- +goose Up
-- +goose StatementBegin
CREATE TYPE todo_status AS ENUM ('Done', 'NotStarted');
CREATE TABLE IF NOT EXISTS todos (
    "id" UUID PRIMARY KEY NOT NULL,
    "user_id" UUID NOT NULL,
    "name" varchar(255) NOT NULL,
    "status" todo_status NOT NULL,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE INDEX idx_todos_deleted_at ON todos (deleted_at); -- +create index
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_deleted_at; -- +drop index first
DROP TABLE IF EXISTS todos;
DROP TYPE IF EXISTS todo_status;
-- +goose StatementEnd
