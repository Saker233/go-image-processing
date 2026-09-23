CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE TABLE images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,

    original_filename TEXT NOT NULL,

    original_s3_key TEXT NOT NULL,
    processed_s3_key TEXT,

    mime_type VARCHAR(100) NOT NULL,

    original_size BIGINT NOT NULL,
    processed_size BIGINT,

    width INT NOT NULL,
    height INT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);