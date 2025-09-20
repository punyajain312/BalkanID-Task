package services

import (
    "errors"
    "balkanid-capstone/internal/models"
    "balkanid-capstone/internal/repo"
)

type ShareService struct {
    Repo *repo.ShareRepo
}

func NewShareService(r *repo.ShareRepo) *ShareService {
    return &ShareService{Repo: r}
}

func (s *ShareService) CreateShare(req models.ShareRequest, ownerID string) (string, error) {
    // Validate visibility
    if req.Visibility != models.VisibilityPrivate &&
       req.Visibility != models.VisibilityPublic &&
       req.Visibility != models.VisibilityRestricted {
        return "", errors.New("invalid visibility option")
    }
    return s.Repo.CreateShare(req, ownerID)
}

func (s *ShareService) AccessShare(shareID, requesterID string) (*models.Share, error) {
    share, err := s.Repo.GetShareByID(shareID)
    if err != nil {
        return nil, err
    }

    // Check access rules
    switch share.Visibility {
    case models.VisibilityPrivate:
        if share.OwnerID != requesterID {
            return nil, errors.New("unauthorized")
        }
    case models.VisibilityRestricted:
        if share.SharedWith != requesterID && share.OwnerID != requesterID {
            return nil, errors.New("unauthorized")
        }
    case models.VisibilityPublic:
        // Anyone can access
    }

    // Increment download count if public
    if share.Visibility == models.VisibilityPublic {
        _ = s.Repo.IncrementDownload(share.ID)
        share.DownloadCount++
    }

    return share, nil
}

func (s *ShareService) GetPublicStats(fileID string) (*models.PublicStats, error) {
    return s.Repo.GetPublicStats(fileID)
}