package setlists

import (
	"api/src/modules/setlists/entities"
	"gorm.io/gorm"
)

type SetlistRepository struct {
	db *gorm.DB
}

func NewSetlistRepository(db *gorm.DB) *SetlistRepository {
	return &SetlistRepository{db: db}
}

func (r *SetlistRepository) Create(setlist *entities.Setlist) error {
	return r.db.Create(setlist).Error
}

func (r *SetlistRepository) FindAll(userID uint) ([]entities.Setlist, error) {
	var setlists []entities.Setlist
	if err := r.db.Where("user_id = ?", userID).Find(&setlists).Error; err != nil {
		return nil, err
	}
	return setlists, nil
}

func (r *SetlistRepository) FindOne(id uint, userID uint) (*entities.Setlist, error) {
	var setlist entities.Setlist
	if err := r.db.Preload("Songs").Where("id = ? AND user_id = ?", id, userID).First(&setlist).Error; err != nil {
		return nil, err
	}
	return &setlist, nil
}

func (r *SetlistRepository) FindByPublicID(publicID string) (*entities.Setlist, error) {
	var setlist entities.Setlist
	if err := r.db.Preload("Songs").Where("public_id = ? AND is_public = ?", publicID, true).First(&setlist).Error; err != nil {
		return nil, err
	}
	return &setlist, nil
}

func (r *SetlistRepository) Update(setlist *entities.Setlist) error {
	return r.db.Save(setlist).Error
}

func (r *SetlistRepository) Delete(id uint, userID uint) error {
	return r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&entities.Setlist{}).Error
}

func (r *SetlistRepository) AddSong(song *entities.SetlistItem) error {
	return r.db.Create(song).Error
}

func (r *SetlistRepository) RemoveSong(id uint, setlistID uint) error {
	return r.db.Where("id = ? AND setlist_id = ?", id, setlistID).Delete(&entities.SetlistItem{}).Error
}

func (r *SetlistRepository) UpdateSong(song *entities.SetlistItem) error {
	return r.db.Save(song).Error
}

