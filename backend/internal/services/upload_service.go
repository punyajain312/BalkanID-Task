package services

import (
    "crypto/sha256"
    "encoding/hex"
    "io"
    "os"
    "path/filepath"

    "balkanid-capstone/internal/models"
    "balkanid-capstone/internal/repo"
)

type UploadService struct {
    Repo *repo.UploadRepo
}

func NewUploadService(r *repo.UploadRepo) *UploadService {
    return &UploadService{Repo: r}
}

func (s *UploadService) UploadFile(userID string, file io.Reader, filename, mimeType string, size int64) (models.UploadResult, error) {
    // Read file content
    fileBytes, err := io.ReadAll(file)
    if err != nil {
        return models.UploadResult{}, err
    }

    // Compute SHA-256
    hash := sha256.Sum256(fileBytes)
    hashStr := hex.EncodeToString(hash[:])

    // Check if blob exists
    blobID, exists, err := s.Repo.FindBlobByHash(hashStr)
    if err != nil {
        return models.UploadResult{}, err
    }

    if !exists {
        // Save file to disk
        os.MkdirAll("uploads", os.ModePerm)
        path := filepath.Join("uploads", hashStr)
        if err := os.WriteFile(path, fileBytes, 0644); err != nil {
            return models.UploadResult{}, err
        }

        // Insert new blob
        blobID, err = s.Repo.InsertBlob(hashStr, path, size)
        if err != nil {
            return models.UploadResult{}, err
        }
    } else {
        // Increment ref_count
        if err := s.Repo.IncrementBlobRef(blobID); err != nil {
            return models.UploadResult{}, err
        }
    }

    // Insert into files table
    if err := s.Repo.InsertFile(userID, blobID, filename, mimeType, size); err != nil {
        return models.UploadResult{}, err
    }

    return models.UploadResult{Hash: hashStr, Filename: filename}, nil
}