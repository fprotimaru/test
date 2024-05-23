-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "rates" (
    "id" UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    "timestamp" int,
    "ask" JSONB,
    "bid" JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "rates";
-- +goose StatementEnd
