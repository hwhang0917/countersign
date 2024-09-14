package repository

import (
	"github.com/hwhang0917/countersign/internal/models"
	"gorm.io/gorm"
)

type WordRepository struct {
	db *gorm.DB
}

func NewWordRepository(db *gorm.DB) *WordRepository {
	return &WordRepository{db}
}

func (r *WordRepository) Create(word *models.Word) error {
	return r.db.Create(word).Error
}

func (r *WordRepository) CountAll() (int64, error) {
	var count int64
	err := r.db.Model(&models.Word{}).Count(&count).Error
	return count, err
}

func (r *WordRepository) FindByID(id uint) (*models.Word, error) {
	var word models.Word
	err := r.db.First(&word, id).Error
	return &word, err
}

func (r *WordRepository) IsWordExist(word string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Word{}).Where("word = ?", word).Count(&count).Error
	return count > 0, err
}
