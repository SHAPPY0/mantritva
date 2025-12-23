package repository

import (
	"strings"

	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type MemberRepository interface {
	Create(*model.Member) error
	CreateCriminalRecord(*model.CriminalRecord) error
	CreateCareer(*model.Career) error
	List() ([]model.Member, error)
	GetByName(name string) (*model.Member, error)
	GetByID(id uint) (*model.Member, error)
	ConsoleGetMembersList() (*[]model.Member, error)
	ListUnaffiliatedMembers() ([]model.Member, error)
	GetMemberAndPartyDetails(int, int) (map[string]interface{}, error)
	// GetCareerByMemberID(memberID uint) (*[]model.CareerRes, error)
	GetMemberById(memberID uint) (map[string]interface{}, error)
}

type memberRepository struct {
	DB *gorm.DB
}

func NewMemberRepository(db *gorm.DB) MemberRepository {
	return &memberRepository{
		DB: db,
	}
}

func (r *memberRepository) Create(member *model.Member) error {
	if err := r.DB.Create(member).Error; err != nil {
		return err
	}
	return nil
}

func (r *memberRepository) CreateCriminalRecord(record *model.CriminalRecord) error {
	if err := r.DB.Create(record).Error; err != nil {
		return err
	}
	return nil
}

func (r *memberRepository) CreateCareer(career *model.Career) error {
	if err := r.DB.Create(career).Error; err != nil {
		return err
	}
	return nil
}

func (r *memberRepository) List() ([]model.Member, error) {
	var members []model.Member
	if err := r.DB.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *memberRepository) GetByName(name string) (*model.Member, error) {
	var member model.Member
	query := "SELECT * FROM members m WHERE LOWER(m.name) = ? LIMIT 1"
	if err := r.DB.Raw(query, strings.ToLower(name)).Scan(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) GetByID(id uint) (*model.Member, error) {
	var member model.Member
	if err := r.DB.Where("id = ?", id).First(&member).Error; err != nil {
		return nil, err
	}
	return &member, nil
}

func (r *memberRepository) ConsoleGetMembersList() (*[]model.Member, error) {
	var members []model.Member
	if err := r.DB.Order("name asc").Find(&members).Error; err != nil {
		return nil, err
	}
	return &members, nil
}

func (r *memberRepository) ListUnaffiliatedMembers() ([]model.Member, error) {
	var members []model.Member
	query := `
		SELECT m.*
		FROM members m
		LEFT JOIN party_members pm ON m.id = pm.member_id AND pm.exit_date IS NULL
		WHERE pm.member_id IS NULL
		ORDER BY m.name ASC
	`
	if err := r.DB.Raw(query).Scan(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *memberRepository) GetMemberAndPartyDetails(memberID, partyID int) (map[string]interface{}, error) {
	var results map[string]interface{}
	query := `
		SELECT
			json_build_object(
				'name', m.name,
				'photo', m.photo,
				'position', CASE WHEN l.role IS NOT NULL THEN l.role ELSE l.position || CASE WHEN s.name IS NOT NULL THEN ', ' || s.name || '' ELSE '' END END
			) AS member,
			json_build_object(
				'name', p.name,
				'abbrv', p.abbrv,
				'symbol', p.symbol
			) AS party
		FROM
			members m
		LEFT JOIN leaderships l ON
			l.member_id = m.id
		LEFT JOIN states s ON
			s.id = l.state_id
		LEFT JOIN parties p ON p.id = ?
		WHERE m.id = ?
	`
	if err := r.DB.Raw(query, partyID, memberID).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *memberRepository) GetMemberById(memberID uint) (map[string]interface{}, error) {
	var member map[string]interface{}
	query := `
		SELECT
			m.name,
			m.photo,
			pm."position",
			json_build_object(
				'name', p.name,
				'abbrv', p.abbrv,
				'symbol', p.symbol,
				'position', case when l.role is not null then l.role else l.position || case when s.name is not null then ', ' || s.name || '' else '' end end
			) AS party
		FROM members m
		INNER JOIN party_members pm ON
			pm.member_id = m.id
		INNER JOIN parties p ON
			p.id = pm.party_id
		INNER JOIN leaderships l ON
			l.member_id = m.id and l.is_active = true and l.term_start < now() and l.term_end is null
		LEFT JOIN states s ON
			s.id = l.state_id
		WHERE m.id = ? AND pm.exit_date IS NULL
	`
	if err := r.DB.Raw(query, memberID).Scan(&member).Error; err != nil {
		return nil, err
	}
	return member, nil
}
