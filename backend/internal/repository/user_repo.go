package repository

import (
	"devhelper/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(u *models.User) error
	FindByEmail(email string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Update(u *models.User) error
	Delete(id uint) error
	List(page, pageSize int, search string) ([]models.User, int64, error)
	SetAdminByEmail(email string) error
}

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(u *models.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) FindByEmail(email string) (*models.User, error) {
	var u models.User
	err := r.db.Where("email = ?", email).First(&u).Error
	return &u, err
}

func (r *UserRepo) FindByID(id uint) (*models.User, error) {
	var u models.User
	err := r.db.First(&u, id).Error
	return &u, err
}

func (r *UserRepo) FindByUsername(username string) (*models.User, error) {
	var u models.User
	err := r.db.Where("username = ?", username).First(&u).Error
	return &u, err
}

func (r *UserRepo) Update(u *models.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepo) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *UserRepo) List(page, pageSize int, search string) ([]models.User, int64, error) {
	var users []models.User
	var total int64
	q := r.db.Model(&models.User{})
	if search != "" {
		q = q.Where("username LIKE ? OR email LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (r *UserRepo) SetAdminByEmail(email string) error {
	return r.db.Model(&models.User{}).Where("email = ?", email).Update("role", "admin").Error
}
