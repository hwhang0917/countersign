package services

import (
	"errors"

	"github.com/hwhang0917/countersign/internal/models"
	"github.com/hwhang0917/countersign/internal/repository"
)

type WordService struct {
	wordRepo *repository.WordRepository
}

func NewWordService(wordRepo *repository.WordRepository) *WordService {
	return &WordService{wordRepo}
}

func (s *WordService) CreateWord(word string) error {
	if word == "" {
		return errors.New("Word cannot be empty")
	}
	isWordExist, err := s.wordRepo.IsWordExist(word)
	if err != nil {
		return err
	}
	if isWordExist {
		return errors.New("Word already exists")
	}
	return s.wordRepo.Create(&models.Word{Word: word})
}

func (s *WordService) GetWordByID(id int) (string, error) {
	word, err := s.wordRepo.FindByID(id)
	if err != nil {
		return "", err
	}
	return word.Word, nil
}

func (s *WordService) GetLastID() (int64, error) {
	return s.wordRepo.GetLastID()
}
