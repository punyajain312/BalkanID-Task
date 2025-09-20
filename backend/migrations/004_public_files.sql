DROP TABLE IF EXISTS public_files CASCADE;
-- 4. Publicly shared files
CREATE TABLE IF NOT EXISTS public_files (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    file_id UUID NOT NULL REFERENCES files(id) ON DELETE CASCADE,
    folder_id UUID, -- optional if you have folders
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    visibility TEXT DEFAULT 'public',
    download_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);