package repository

import (
	"github.com/lcs/devhelper/internal/models"
	"gorm.io/gorm"
)

const maxHistoryPerSession = 50

type HistoryRepo struct {
	db *gorm.DB
}

func NewHistoryRepo(db *gorm.DB) *HistoryRepo {
	return &HistoryRepo{db: db}
}

func (r *HistoryRepo) Create(h *models.JsonHistory) error {
	if err := r.db.Create(h).Error; err != nil {
		return err
	}
	return r.pruneSession(h.UserID, h.SessionID)
}

func (r *HistoryRepo) pruneSession(userID uint, sessionID string) error {
	var count int64
	r.db.Model(&models.JsonHistory{}).Where("user_id = ? AND session_id = ?", userID, sessionID).Count(&count)
	if count <= maxHistoryPerSession {
		return nil
	}
	// Delete oldest entries beyond the limit
	var oldest []models.JsonHistory
	r.db.Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("seq_num ASC").
		Limit(int(count - maxHistoryPerSession)).
		Find(&oldest)
	for _, o := range oldest {
		r.db.Delete(&o)
	}
	return nil
}

func (r *HistoryRepo) ListBySession(userID uint, sessionID string) ([]models.JsonHistory, error) {
	var items []models.JsonHistory
	err := r.db.Where("user_id = ? AND session_id = ?", userID, sessionID).
		Order("seq_num ASC").Find(&items).Error
	return items, err
}

func (r *HistoryRepo) FindByID(id, userID uint) (*models.JsonHistory, error) {
	var h models.JsonHistory
	err := r.db.Where("id = ? AND user_id = ?", id, userID).First(&h).Error
	return &h, err
}

func (r *HistoryRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.JsonHistory{}).Error
}

func (r *HistoryRepo) NextSeqNum(userID uint, sessionID string) int {
	var max int
	r.db.Model(&models.JsonHistory{}).
		Where("user_id = ? AND session_id = ?", userID, sessionID).
		Select("COALESCE(MAX(seq_num), 0)").Scan(&max)
	return max + 1
}

func (r *HistoryRepo) HasBase(userID uint, sessionID string) bool {
	var count int64
	r.db.Model(&models.JsonHistory{}).
		Where("user_id = ? AND session_id = ? AND is_base = ?", userID, sessionID, true).
		Count(&count)
	return count > 0
}
