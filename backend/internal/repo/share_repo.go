package repo

import (
    "database/sql"
    "balkanid-capstone/internal/models"
)

type ShareRepo struct {
    DB *sql.DB
}

func NewShareRepo(db *sql.DB) *ShareRepo {
    return &ShareRepo{DB: db}
}

func (r *ShareRepo) CreateShare(req models.ShareRequest, ownerID string) (string, error) {
    var shareID string
    err := r.DB.QueryRow(`
        INSERT INTO public_files (file_id, folder_id, owner_id, visibility, shared_with, download_count)
        VALUES ($1, $2, $3, $4, $5, 0)
        RETURNING id
    `, req.FileID, req.FolderID, ownerID, req.Visibility, req.SharedWith).Scan(&shareID)
    return shareID, err
}

func (r *ShareRepo) GetShareByID(shareID string) (*models.Share, error) {
    var s models.Share
    err := r.DB.QueryRow(`
        SELECT id, file_id, folder_id, owner_id, visibility, shared_with, download_count, created_at
        FROM shares WHERE id=$1
    `, shareID).Scan(&s.ID, &s.FileID, &s.FolderID, &s.OwnerID, &s.Visibility, &s.SharedWith, &s.DownloadCount, &s.CreatedAt)
    if err != nil {
        return nil, err
    }
    return &s, nil
}

func (r *ShareRepo) IncrementDownload(shareID string) error {
    _, err := r.DB.Exec(`
        UPDATE shares SET download_count = download_count + 1 WHERE id=$1
    `, shareID)
    return err
}

func (r *ShareRepo) GetPublicStats(fileID string) (*models.PublicStats, error) {
    var stats models.PublicStats
    err := r.DB.QueryRow(`
        SELECT file_id, download_count FROM shares 
        WHERE file_id=$1 AND visibility='public'
    `, fileID).Scan(&stats.FileID, &stats.DownloadCount)
    if err != nil {
        return nil, err
    }
    return &stats, nil
}