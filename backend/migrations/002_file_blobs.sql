DROP TABLE IF EXISTS file_blobs CASCADE;
-- 2. File blobs (unique file contents)
CREATE TABLE file_blobs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    hash TEXT UNIQUE NOT NULL,
    storage_path TEXT NOT NULL,
    size BIGINT NOT NULL,
    ref_count INT DEFAULT 1
);