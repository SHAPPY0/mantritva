package repository

import (
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type VidhansabhaRepository interface {
	Create(*model.Vidhansabha) error
	List() ([]model.Vidhansabha, error)
	GetByID(id int) (*model.Vidhansabha, error)
	GetByStateID(int) (*[]dto.VidhanSabhaRes, error)
	GetByStateIDAndDistrictIDAndName(int, int, string) (*model.Vidhansabha, error)
}

type vidhansabhaRepository struct {
	DB *gorm.DB
}

func NewVidhansabhaRepository(db *gorm.DB) VidhansabhaRepository {
	return &vidhansabhaRepository{
		DB: db,
	}
}

func (r *vidhansabhaRepository) Create(vidhanSabha *model.Vidhansabha) error {
	if err := r.DB.Create(vidhanSabha).Error; err != nil {
		return err
	}
	return nil
}

func (r *vidhansabhaRepository) List() ([]model.Vidhansabha, error) {
	var vidhanSabhas []model.Vidhansabha
	if err := r.DB.Find(&vidhanSabhas).Error; err != nil {
		return nil, err
	}
	return vidhanSabhas, nil
}

func (r *vidhansabhaRepository) GetByID(id int) (*model.Vidhansabha, error) {
	var vidhanSabha model.Vidhansabha
	if err := r.DB.First(&vidhanSabha, id).Error; err != nil {
		return nil, err
	}
	return &vidhanSabha, nil
}

func (r *vidhansabhaRepository) GetByStateID(id int) (*[]dto.VidhanSabhaRes, error) {
	var vidhanSabha []dto.VidhanSabhaRes
	query := `SELECT
				vs.id,
				vs.name AS name,
				s.name AS state,
				d.name AS district
			FROM
				vidhansabhas vs
			LEFT JOIN states s ON
				s.id = vs.state_id
			LEFT JOIN districts d ON
				d.id = vs.district_id
			WHERE
				vs.state_id = ?
			ORDER BY vs.name ASC`
	if err := r.DB.Raw(query, id).Scan(&vidhanSabha).Error; err != nil {
		return nil, err
	}
	return &vidhanSabha, nil
}

func (r *vidhansabhaRepository) GetByStateIDAndDistrictIDAndName(stateId, districtId int, name string) (*model.Vidhansabha, error) {
	var vidhanSabha model.Vidhansabha
	if err := r.DB.Where("state_id = ? AND district_id = ? AND name = ?", stateId, districtId, name).First(&vidhanSabha).Error; err != nil {
		return nil, err
	}
	if vidhanSabha.ID <= 0 {
		return nil, nil
	}
	return &vidhanSabha, nil
}
