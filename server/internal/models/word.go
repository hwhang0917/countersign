package models

import "time"

type Word struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Word      string    `json:"word"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Word) TableName() string {
	return "words"
}

func (w *Word) BeforeCreateInterface() (err error) {
	w.CreatedAt = time.Now()
	return nil
}
