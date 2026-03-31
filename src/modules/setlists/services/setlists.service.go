package services

import (
	"api/src/modules/setlists/entities"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

type SetlistRepository interface {
	Create(setlist *entities.Setlist) error
	FindAll(userID uint) ([]entities.Setlist, error)
	FindOne(id uint, userID uint) (*entities.Setlist, error)
	FindByPublicID(publicID string) (*entities.Setlist, error)
	Update(setlist *entities.Setlist) error
	Delete(id uint, userID uint) error
	AddSong(song *entities.SetlistItem) error
	RemoveSong(id uint, setlistID uint) error
}

type SetlistService struct {
	repository SetlistRepository
}

func NewSetlistService(repository SetlistRepository) *SetlistService {
	return &SetlistService{repository: repository}
}

func (s *SetlistService) Create(title string, userID uint) (*entities.Setlist, error) {
	publicID, _ := generateRandomHex(8)
	setlist := &entities.Setlist{
		Title:    title,
		UserID:   userID,
		PublicID: publicID,
	}
	if err := s.repository.Create(setlist); err != nil {
		return nil, err
	}
	return setlist, nil
}

func (s *SetlistService) FindAll(userID uint) ([]entities.Setlist, error) {
	return s.repository.FindAll(userID)
}

func (s *SetlistService) FindOne(id uint, userID uint) (*entities.Setlist, error) {
	return s.repository.FindOne(id, userID)
}

func (s *SetlistService) FindShared(publicID string) (*entities.Setlist, error) {
	return s.repository.FindByPublicID(publicID)
}

func (s *SetlistService) Update(id uint, userID uint, title string, isPublic bool) (*entities.Setlist, error) {
	setlist, err := s.repository.FindOne(id, userID)
	if err != nil {
		return nil, err
	}
	setlist.Title = title
	setlist.IsPublic = isPublic
	if err := s.repository.Update(setlist); err != nil {
		return nil, err
	}
	return setlist, nil
}

func (s *SetlistService) Delete(id uint, userID uint) error {
	return s.repository.Delete(id, userID)
}

func (s *SetlistService) AddSong(setlistID uint, userID uint, title, artist, url, key string, order int) (*entities.SetlistItem, error) {
	// Check if setlist belongs to user
	_, err := s.repository.FindOne(setlistID, userID)
	if err != nil {
		return nil, fmt.Errorf("setlist not found or access denied")
	}

	song := &entities.SetlistItem{
		SetlistID: setlistID,
		Title:     title,
		Artist:    artist,
		URL:       url,
		Key:       key,
		Order:     order,
	}
	if err := s.repository.AddSong(song); err != nil {
		return nil, err
	}
	return song, nil
}

func (s *SetlistService) RemoveSong(setlistID uint, userID uint, songID uint) error {
	_, err := s.repository.FindOne(setlistID, userID)
	if err != nil {
		return fmt.Errorf("setlist not found or access denied")
	}
	return s.repository.RemoveSong(songID, setlistID)
}

func generateRandomHex(n int) (string, error) {
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
