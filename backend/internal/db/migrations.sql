-- 1. Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- needs pgcrypto
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- 2. File blobs (unique file contents)
CREATE TABLE file_blobs (
    hash TEXT PRIMARY KEY,               -- SHA-256
    storage_path TEXT NOT NULL,          -- disk/S3 path
    size BIGINT NOT NULL,
    ref_count INT DEFAULT 1
);

-- 3. Files (user uploads)
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_hash TEXT NOT NULL REFERENCES file_blobs(hash) ON DELETE CASCADE,
    filename TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
