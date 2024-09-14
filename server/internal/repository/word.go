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

func (r *WordRepository) GetLastID() (int64, error) {
	var lastId int64
	err := r.db.Model(&models.Word{}).Select("id").Order("id desc").Limit(1).Scan(&lastId).Error
	return lastId, err
}

func (r *WordRepository) FindByID(id int) (*models.Word, error) {
	var word models.Word
	err := r.db.First(&word, id).Error
	return &word, err
}

func (r *WordRepository) IsWordExist(word string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Word{}).Where("word = ?", word).Count(&count).Error
	return count > 0, err
}
