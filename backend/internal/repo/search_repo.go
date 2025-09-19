package repo

import (
    "database/sql"
    "strconv"
    "time"

    "balkanid-capstone/internal/models"
)

type SearchRepo struct {
    DB *sql.DB
}

func NewSearchRepo(db *sql.DB) *SearchRepo {
    return &SearchRepo{DB: db}
}

func (r *SearchRepo) SearchFiles(userID, q, mime, minSize, maxSize, from, to string) ([]models.File, error) {
    query := `
        SELECT f.id, f.filename, f.mime_type, f.size, f.created_at, b.hash, b.ref_count
        FROM files f
        JOIN file_blobs b ON f.blob_id = b.id
        WHERE f.user_id = $1
    `
    args := []interface{}{userID}
    argIndex := 2

    if q != "" {
        query += " AND f.filename ILIKE $" + strconv.Itoa(argIndex)
        args = append(args, "%"+q+"%")
        argIndex++
    }
    if mime != "" {
        query += " AND f.mime_type = $" + strconv.Itoa(argIndex)
        args = append(args, mime)
        argIndex++
    }
    if minSize != "" {
        if val, err := strconv.ParseInt(minSize, 10, 64); err == nil {
            query += " AND f.size >= $" + strconv.Itoa(argIndex)
            args = append(args, val)
            argIndex++
        }
    }
    if maxSize != "" {
        if val, err := strconv.ParseInt(maxSize, 10, 64); err == nil {
            query += " AND f.size <= $" + strconv.Itoa(argIndex)
            args = append(args, val)
            argIndex++
        }
    }
    if from != "" {
        if t, err := time.Parse("2006-01-02", from); err == nil {
            query += " AND f.created_at >= $" + strconv.Itoa(argIndex)
            args = append(args, t)
            argIndex++
        }
    }
    if to != "" {
        if t, err := time.Parse("2006-01-02", to); err == nil {
            query += " AND f.created_at <= $" + strconv.Itoa(argIndex)
            args = append(args, t)
            argIndex++
        }
    }

    query += " ORDER BY f.created_at DESC"

    rows, err := r.DB.Query(query, args...)
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