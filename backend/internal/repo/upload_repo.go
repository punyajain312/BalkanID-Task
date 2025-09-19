package repo

import (
    "database/sql"
)

type UploadRepo struct {
    DB *sql.DB
}

func NewUploadRepo(db *sql.DB) *UploadRepo {
    return &UploadRepo{DB: db}
}

func (r *UploadRepo) FindBlobByHash(hash string) (string, bool, error) {
    var blobID string
    var exists bool
    err := r.DB.QueryRow(`SELECT id, true FROM file_blobs WHERE hash=$1`, hash).
        Scan(&blobID, &exists)

    if err == sql.ErrNoRows {
        return "", false, nil
    }
    if err != nil {
        return "", false, err
    }
    return blobID, true, nil
}

func (r *UploadRepo) InsertBlob(hash, path string, size int64) (string, error) {
    var blobID string
    err := r.DB.QueryRow(`
        INSERT INTO file_blobs (hash, storage_path, size, ref_count)
        VALUES ($1, $2, $3, 1) RETURNING id
    `, hash, path, size).Scan(&blobID)
    return blobID, err
}

func (r *UploadRepo) IncrementBlobRef(blobID string) error {
    _, err := r.DB.Exec(`UPDATE file_blobs SET ref_count = ref_count + 1 WHERE id=$1`, blobID)
    return err
}

func (r *UploadRepo) InsertFile(userID, blobID, filename, mimeType string, size int64) error {
    _, err := r.DB.Exec(`
        INSERT INTO files (user_id, blob_id, filename, mime_type, size)
        VALUES ($1, $2, $3, $4, $5)
    `, userID, blobID, filename, mimeType, size)
    return err
}