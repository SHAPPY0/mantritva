package repository

import (
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type LeadershipRepository interface {
	Create(*model.Leadership) error
	GetByMemberID(id uint) (*model.Leadership, error)
	GetByRulerID(rulerID uint) ([]model.Leadership, error)
	GetByRulerIDAndPosition(rulerID uint, position string) ([]model.Leadership, error)
}

type leadershipRepository struct {
	DB *gorm.DB
}

func NewLeadershipRepository(db *gorm.DB) LeadershipRepository {
	return &leadershipRepository{
		DB: db,
	}
}

func (r *leadershipRepository) Create(leadership *model.Leadership) error {
	if err := r.DB.Create(leadership).Error; err != nil {
		return err
	}
	return nil
}

// func (r *leadershipRepository) GetPositionByMemberID(id uint) (string, error) {
// 	var leadership model.Leadership
// 	var position string
// 	if err := r.DB.Where("id = ? AND term_end IS NULL", id).First(&leadership).Error; err != nil {
// 		return "", err
// 	}
// 	if leadership.Role == "" {
// 		position = leadership.Position
// 	} else {
// 		position = leadership.Role
// 	}

// 	return position, nil
// }

func (r *leadershipRepository) GetByMemberID(id uint) (*model.Leadership, error) {
	var leadership model.Leadership
	if err := r.DB.Where("member_id = ? AND term_end IS NULL", id).First(&leadership).Error; err != nil {
		return nil, err
	}
	return &leadership, nil
}

func (r *leadershipRepository) GetByRulerID(rulerID uint) ([]model.Leadership, error) {
	var leadership []model.Leadership
	if err := r.DB.Where("ruler_id = ? AND is_central = ? AND term_end IS NULL", rulerID, true).Find(&leadership).Error; err != nil {
		return nil, err
	}
	return leadership, nil
}

func (r *leadershipRepository) GetByRulerIDAndPosition(rulerID uint, position string) ([]model.Leadership, error) {
	var leadership []model.Leadership
	if err := r.DB.Where("ruler_id = ? AND is_central = ? AND term_end IS NULL AND position = ?", rulerID, true, position).Find(&leadership).Error; err != nil {
		return nil, err
	}
	return leadership, nil
}
