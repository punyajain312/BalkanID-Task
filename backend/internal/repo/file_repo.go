package repo

import (
    "database/sql"
    "balkanid-capstone/internal/models"
)

type FileRepo struct {
    DB *sql.DB
}

func NewFileRepo(db *sql.DB) *FileRepo {
    return &FileRepo{DB: db}
}

func (r *FileRepo) ListFiles(userID string) ([]models.File, error) {
    rows, err := r.DB.Query(`
        SELECT f.id, f.filename, f.mime_type, f.size, f.created_at,
               b.hash, b.ref_count
        FROM files f
        JOIN file_blobs b ON f.blob_id = b.id
        WHERE f.user_id = $1
        ORDER BY f.created_at DESC
    `, userID)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var files []models.File
    for rows.Next() {
        var f models.File
        if err := rows.Scan(&f.ID, &f.Filename, &f.MimeType, &f.Size, &f.CreatedAt, &f.Hash, &f.RefCount); err != nil {
            return nil, err
        }
        files = append(files, f)
    }
    return files, nil
}

func (r *FileRepo) GetFileBlob(tx *sql.Tx, fileID, userID string) (string, string, error) {
    var blobID, blobPath string
    err := tx.QueryRow(`
        SELECT f.blob_id, b.storage_path
        FROM files f
        JOIN file_blobs b ON f.blob_id = b.id
        WHERE f.id=$1 AND f.user_id=$2
    `, fileID, userID).Scan(&blobID, &blobPath)
    return blobID, blobPath, err
}

func (r *FileRepo) DeleteFileRecord(tx *sql.Tx, fileID, userID string) error {
    _, err := tx.Exec(`DELETE FROM files WHERE id=$1 AND user_id=$2`, fileID, userID)
    return err
}

func (r *FileRepo) DecrementRefCount(tx *sql.Tx, blobID string) (int, error) {
    var refCount int
    err := tx.QueryRow(`
        UPDATE file_blobs SET ref_count = ref_count - 1 WHERE id=$1 RETURNING ref_count
    `, blobID).Scan(&refCount)
    return refCount, err
}

func (r *FileRepo) DeleteBlob(tx *sql.Tx, blobID string) error {
    _, err := tx.Exec(`DELETE FROM file_blobs WHERE id=$1`, blobID)
    return err
}