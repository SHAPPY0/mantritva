package repository

import (
	"strings"

	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type StateRepository interface {
	Create(*model.State) error
	List() ([]model.State, error)
	Get(name string) (*model.StateDetail, error)
	GetByID(id int) (*model.State, error)
	ConsoleGetStateList() (*[]model.State, error)
	ConsoleGetStateListByID([]uint) (*[]model.State, error)
	ConsoleGetStateListByDistrictIDs([]uint) (*[]model.State, error)
}

type stateRepository struct {
	DB *gorm.DB
}

func NewStateRepository(db *gorm.DB) StateRepository {
	return &stateRepository{
		DB: db,
	}
}

func (r *stateRepository) Create(state *model.State) error {
	if err := r.DB.Create(state).Error; err != nil {
		return err
	}
	return nil
}

func (r *stateRepository) List() ([]model.State, error) {
	var states []model.State
	if err := r.DB.Find(&states).Error; err != nil {
		return nil, err
	}
	return states, nil
}

func (r *stateRepository) Get(name string) (*model.StateDetail, error) {
	var stateDetail model.StateDetail
	query := `SELECT s.*, 
					COUNT(distinct d.id) AS district_count, 
					COUNT(distinct l.id) AS lok_sabha_count,
					COUNT(distinct v.id) AS vidhan_sabha_count
				FROM states s
					LEFT JOIN districts d ON d.state_id = s.id
					LEFT JOIN loksabhas l ON l.state_id = s.id
					LEFT JOIN vidhansabhas v ON v.state_id = s.id
				WHERE s.slug = ? AND s.is_deleted = false 
				GROUP BY s.id`
	if err := r.DB.Raw(query, strings.ToLower(name)).Scan(&stateDetail).Error; err != nil {
		return nil, err
	}
	return &stateDetail, nil
}

func (r *stateRepository) GetByID(id int) (*model.State, error) {
	var state model.State
	if err := r.DB.Where("id = ?", id).First(&state).Error; err != nil {
		return nil, err
	}
	return &state, nil
}

func (r *stateRepository) ConsoleGetStateList() (*[]model.State, error) {
	var states []model.State
	if err := r.DB.Order("name asc").Find(&states).Error; err != nil {
		return nil, err
	}
	return &states, nil
}

func (r *stateRepository) ConsoleGetStateListByID(stateIDs []uint) (*[]model.State, error) {
	var states []model.State
	if err := r.DB.Where("id IN ?", stateIDs).Find(&states).Error; err != nil {
		return nil, err
	}
	return &states, nil
}

func (r *stateRepository) ConsoleGetStateListByDistrictIDs(districtIDs []uint) (*[]model.State, error) {
	var states []model.State
	if err := r.DB.Where("id IN (SELECT state_id FROM districts WHERE id IN ?)", districtIDs).Find(&states).Error; err != nil {
		return nil, err
	}
	return &states, nil
}
