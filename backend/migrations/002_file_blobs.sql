-- 2. File blobs (unique file contents)
CREATE TABLE file_blobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(), -- needs pgcrypto
    hash TEXT,               -- SHA-256
    storage_path TEXT NOT NULL,          -- disk/S3 path
    size BIGINT NOT NULL,
    ref_count INT DEFAULT 1
);
