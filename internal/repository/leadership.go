package repository

import (
	"time"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/pkg/utils"
	"gorm.io/gorm"
)

type LeadershipRepository interface {
	Create(*model.Leadership) error
	GetByMemberID(uint) (*model.Leadership, error)
	GetByRulerID(uint, bool) ([]map[string]interface{}, error)
	GetByRulerIDAndPosition(uint, string, bool) ([]model.Leadership, error)
	GetByRulerIDAndPositions(uint, string, []string, bool) ([]map[string]interface{}, error)
	GetCareerByMemberID(int) ([]map[string]interface{}, error)
	GetRajyaSabhMembersByStateID(string) ([]map[string]interface{}, error)
	ExitLeadership(dto.LeadershipExitReq) error
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

func (r *leadershipRepository) GetByMemberID(id uint) (*model.Leadership, error) {
	var leadership model.Leadership
	if err := r.DB.Where("member_id = ? AND term_end IS NULL", id).First(&leadership).Error; err != nil {
		return nil, err
	}
	return &leadership, nil
}

func (r *leadershipRepository) GetByRulerID(rulerID uint, isCentral bool) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `SELECT l.id,
				l.ruler_id,
				l.member_id,
				l.party_id,
				l.is_central,
				l.term_start,
				l.term_end,
				l.is_active,
				json_build_object(
					'name', m.name,
					'position', pm.position,
					'photo', m.photo
				) AS member,
				json_build_object(
					'name', p.name,
					'abbrv', p.abbrv,
					'symbol', p.symbol
				) AS party,
				CASE
					WHEN l.state_id != 0 THEN
						json_build_object(
							'name', s.name,
							'abbrv', s.abbrv
						)
					ELSE NULL
				END AS state,
				json_build_object(
					'position', l."position",
					'role', l.role,
					'departments', l.departments
				) AS position,
				v.name AS vidhansabha,
				l2.name AS loksabha
			FROM leaderships l
			INNER JOIN members m ON m.id = l.member_id
			INNER JOIN party_members pm ON pm.member_id = m.id and pm.exit_code IS NULL
			INNER JOIN parties p ON p.id = l.party_id
			LEFT JOIN states s ON s.id = l.state_id
			LEFT JOIN vidhansabhas v ON v.id = l.vidhansabha_id
			LEFT JOIN loksabhas l2 ON l2.id = l.loksabha_id
			WHERE l.ruler_id = ? AND l.term_start < now()`
	if err := r.DB.Raw(query, rulerID).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *leadershipRepository) GetByRulerIDAndPosition(rulerID uint, position string, isCentral bool) ([]model.Leadership, error) {
	var leadership []model.Leadership
	if err := r.DB.Where("ruler_id = ? AND is_central = ? AND term_end IS NULL AND position = ?", rulerID, isCentral, position).Find(&leadership).Error; err != nil {
		return nil, err
	}
	return leadership, nil
}

func (r *leadershipRepository) GetByRulerIDAndPositions(rulerID uint, stateID string, positions []string, isCentral bool) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `SELECT l.*, 
				JSON_BUILD_OBJECT(
					'name', m.name,
					'photo', m.photo
				) AS member,
				JSON_BUILD_OBJECT(
					'name', p.name,
					'abbrv', p.abbrv,
					'symbol', p.symbol
				) AS party,
				JSON_BUILD_OBJECT(
					'position', l."position",
					'role', l.role,
					'departments', l.departments
				) AS "position",
				JSON_BUILD_OBJECT(
					'state', s.name,
					'vidhansabha', v.name
				) AS region
			FROM leaderships l 
				LEFT JOIN members m ON m.id = l.member_id
				LEFT JOIN parties p ON p.id = l.party_id
				LEFT JOIN states s ON s.id = l.state_id
				LEFT JOIN vidhansabhas v ON v.id = l.vidhansabha_id
			WHERE l.term_end IS NULL AND l.ruler_id  = ? AND l.state_id = ? AND l.position IN (?) AND l.is_central = ?`
	if err := r.DB.Raw(query, rulerID, stateID, positions, isCentral).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *leadershipRepository) GetCareerByMemberID(memberID int) ([]map[string]interface{}, error) {
	var leadership []map[string]interface{}
	query := `SELECT l.term_start, l.term_end,
			JSON_BUILD_OBJECT(
				'position', l."position",
				'role', l.role,
				'departments', l.departments
			) AS "position", 
 			JSON_BUILD_OBJECT(
				'name', p.name,
				'abbrv', p.abbrv,
				'symbol', p.symbol
			) AS party,
			JSON_BUILD_OBJECT(
				'state', s.name,
				'loksabha', l2.name,
				'vidhansabha', v.name
			) AS region
			FROM leaderships l 
				LEFT JOIN parties p ON p.id = l.party_id
				LEFT JOIN states s ON s.id = l.state_id 
				LEFT JOIN loksabhas l2 ON l2.id = l.loksabha_id
				LEFT JOIN vidhansabhas v ON v.id = l.vidhansabha_id
			WHERE l.member_id = ? ORDER BY l.term_start DESC, l.term_end DESC;`
	if err := r.DB.Raw(query, memberID).Scan(&leadership).Error; err != nil {
		return nil, err
	}
	return leadership, nil
}

func (r *leadershipRepository) GetRajyaSabhMembersByStateID(stateId string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	query := `SELECT l.*,
			JSON_BUILD_OBJECT(
				'position', l."position",
				'role', l.role,
				'departments', l.departments
			) AS "position", 
 			JSON_BUILD_OBJECT(
				'name', p.name,
				'abbrv', p.abbrv,
				'symbol', p.symbol
			) AS party,
			JSON_BUILD_OBJECT(
				'state', s.name
			) AS region,
			JSON_BUILD_OBJECT(
				'name', m.name,
				'photo', m.photo
			) AS member
			FROM leaderships l 
				LEFT JOIN parties p ON p.id = l.party_id
				LEFT JOIN states s ON s.id = l.state_id 
				LEFT JOIN members m ON m.id = l.member_id
			WHERE l.state_id = ? AND l."position"  = ? AND l.term_end IS NULL 
			ORDER BY l.term_start DESC, l.term_end DESC;	
	`
	if err := r.DB.Raw(query, stateId, utils.MP_RS).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *leadershipRepository) ExitLeadership(req dto.LeadershipExitReq) error {
	result := r.DB.Model(&model.Leadership{}).
		Where("id = ? AND term_end IS NULL", req.LeadershipID).
		Updates(map[string]interface{}{
			"term_end":   req.TermEnd,
			"is_active":  false,
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
