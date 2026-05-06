package repository

import (
	"devhelper/internal/models"

	"gorm.io/gorm"
)

type SchemaRepo struct {
	db *gorm.DB
}

func NewSchemaRepo(db *gorm.DB) *SchemaRepo {
	return &SchemaRepo{db: db}
}

func (r *SchemaRepo) Create(s *models.JsonSchema) error {
	return r.db.Create(s).Error
}

func (r *SchemaRepo) List(userID uint) ([]models.JsonSchema, error) {
	var schemas []models.JsonSchema
	err := r.db.Where("user_id = ? OR is_public = ?", userID, true).
		Order("created_at DESC").Find(&schemas).Error
	return schemas, err
}

func (r *SchemaRepo) FindByID(id, userID uint) (*models.JsonSchema, error) {
	var s models.JsonSchema
	err := r.db.Where("id = ? AND (user_id = ? OR is_public = ?)", id, userID, true).First(&s).Error
	return &s, err
}

func (r *SchemaRepo) Update(s *models.JsonSchema) error {
	return r.db.Save(s).Error
}

func (r *SchemaRepo) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&models.JsonSchema{}).Error
}
