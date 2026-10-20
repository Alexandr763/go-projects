-- +goose Up

CREATE TABLE refresh_tokens (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id uuid NOT NULL,
    token_hash text NOT NULL,
    expires_at timestamptz NOT NULL,
    revoked_at timestamptz,
    created_at timestamptz NOT NULL DEFAULT now(),
    CONSTRAINT refresh_tokens_user_fkey FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT,
    CONSTRAINT refresh_tokens_hash_key UNIQUE (token_hash),
    CONSTRAINT refresh_tokens_hash_present CHECK (token_hash <> '')
);

CREATE INDEX refresh_tokens_user_idx ON refresh_tokens (user_id);

-- +goose Down

DROP TABLE refresh_tokens;
