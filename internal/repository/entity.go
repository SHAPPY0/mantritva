package repository

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type EntityRepository interface {
	Create(*model.Entity) error
	List() ([]dto.EntityRes, error)
	GetEntityByID(int) (bool, *dto.EntityRes, error)
	GetEntityByName(string) (bool, *dto.EntityRes, error)
}

type entityRepository struct {
	DB *gorm.DB
}

func NewEntityRepository(db *gorm.DB) EntityRepository {
	return &entityRepository{
		DB: db,
	}
}

func (r *entityRepository) Create(entity *model.Entity) error {
	if err := r.DB.Create(entity).Error; err != nil {
		return err
	}
	return nil
}

func (r *entityRepository) List() ([]dto.EntityRes, error) {
	var results []dto.EntityRes
	query := `SELECT e.*, s.name as state_name, d.name as district_name
				FROM entities e
				LEFT JOIN states s ON s.id = e.state_id AND e.is_central = false
				LEFT JOIN districts d ON d.id = e.district_id AND e.is_central = false
				ORDER BY e.name ASC`

	if err := r.DB.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *entityRepository) GetEntityByID(id int) (bool, *dto.EntityRes, error) {
	var result dto.EntityRes
	query := `SELECT e.*, s.name as state_name, d.name as district_name
				FROM entities e
				LEFT JOIN states s ON s.id = e.state_id AND e.is_central = false
				LEFT JOIN districts d ON d.id = e.district_id AND e.is_central = false
				WHERE e.id = ?`

	tx := r.DB.Raw(query, id).Scan(&result)
	if tx.Error != nil {
		return false, nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return false, nil, nil
	}
	return true, &result, nil
}

func (r *entityRepository) GetEntityByName(name string) (bool, *dto.EntityRes, error) {
	var result dto.EntityRes
	query := `SELECT e.*, s.name as state_name, d.name as district_name
				FROM entities e
				LEFT JOIN states s ON s.id = e.state_id AND e.is_central = false
				LEFT JOIN districts d ON d.id = e.district_id AND e.is_central = false
				WHERE e.name = ?`

	tx := r.DB.Raw(query, name).Scan(&result)
	if tx.Error != nil {
		return false, nil, tx.Error
	}
	if tx.RowsAffected == 0 {
		return false, nil, nil
	}
	return true, &result, nil
}
