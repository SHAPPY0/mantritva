package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type RulerRepository interface {
	Create(*model.Ruler) error
	List() ([]map[string]interface{}, error)
	GetCentralRuler() (*model.Ruler, error)
	GetStateRuler(stateId string) (*model.Ruler, error)
	ConsoleGetRulersList() (*[]model.Ruler, error)
	ConsoleGetRulersListByStateIDs([]uint) (*[]model.Ruler, error)
	ConsoleGetRulersListByDistrictIDs([]uint) (*[]model.Ruler, error)
	ConsoleGetRulersListByLocalBodyIDs([]uint) (*[]model.Ruler, error)
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

func (r *rulerRepository) List() ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `SELECT
				r.id,
				r.state_id,
				CASE
					WHEN s.name IS null and r.state_id = 0 and r.is_central = true THEN 'Central'
					ELSE s.name
				END AS name,
				r.is_central,
				r.party_id,
				r.is_alliance,
				r.alliance_id
			FROM rulers r
			LEFT JOIN states s on
				s.id = r.state_id
			WHERE
				r.start_date < now()
				and r.end_date > now()
				and r.is_active = true
				and r.is_deleted = false`
	if err := r.DB.Raw(query).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
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

func (r *rulerRepository) ConsoleGetRulersListByStateIDs(stateIDs []uint) (*[]model.Ruler, error) {
	var rulers []model.Ruler
	if err := r.DB.Where("state_id IN ?", stateIDs).Find(&rulers).Error; err != nil {
		return nil, err
	}
	return &rulers, nil
}

func (r *rulerRepository) ConsoleGetRulersListByDistrictIDs(districtIDs []uint) (*[]model.Ruler, error) {
	var rulers []model.Ruler
	if err := r.DB.Where("state_id IN (select state_id from districts where id IN ?)", districtIDs).Find(&rulers).Error; err != nil {
		return nil, err
	}
	return &rulers, nil
}

func (r *rulerRepository) ConsoleGetRulersListByLocalBodyIDs(localBodyIDs []uint) (*[]model.Ruler, error) {
	var rulers []model.Ruler
	if err := r.DB.Where("local_body_id IN ?", localBodyIDs).Find(&rulers).Error; err != nil {
		return nil, err
	}
	return &rulers, nil
}
