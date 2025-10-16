package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type RulerRepository interface {
	Create(*model.Ruler) error
	List() ([]model.Ruler, error)
	GetCentralRuler() (*model.Ruler, error)
	GetStateRuler(stateId string) (*model.Ruler, error)
	ConsoleGetRulersList() (*[]model.Ruler, error)
}

type rulerRepository struct {
	DB *gorm.DB
}

func NewRulersRepository(db *gorm.DB) RulerRepository {
	return &rulerRepository{
		DB: db,
	}
}

func (r *rulerRepository) Create(ruler *model.Ruler) error {
	if err := r.DB.Create(ruler).Error; err != nil {
		return err
	}
	return nil
}

func (r *rulerRepository) List() ([]model.Ruler, error) {
	var rulers []model.Ruler
	if err := r.DB.Find(&rulers).Error; err != nil {
		return nil, err
	}
	return rulers, nil
}

func (r *rulerRepository) GetCentralRuler() (*model.Ruler, error) {
	var ruler model.Ruler
	currentDate := time.Now().Format("2006-01-02")
	if err := r.DB.Where("is_central = ? and start_date < ? and end_date > ? and is_active = ?", true, currentDate, currentDate, true).First(&ruler).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No central ruler found
		}
		return nil, err
	}
	return &ruler, nil
}

func (r *rulerRepository) GetStateRuler(stateId string) (*model.Ruler, error) {
	var ruler model.Ruler
	currentDate := time.Now().Format("2006-01-02")
	if err := r.DB.Where("state_id = ? AND is_central = false AND start_date < ? AND end_date > ? AND is_active = ?", stateId, currentDate, currentDate, true).First(&ruler).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No state ruler found
		}
		return nil, err
	}
	return &ruler, nil
}

func (r *rulerRepository) ConsoleGetRulersList() (*[]model.Ruler, error) {
	var rulers []model.Ruler
	if err := r.DB.Find(&rulers).Error; err != nil {
		return nil, err
	}
	return &rulers, nil
}
