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
	// GetCareerByMemberID(memberID uint) (*[]model.CareerRes, error)
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
