package repository

import (
	"strings"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type AllianceRepository interface {
	Create(*model.Alliance) error
	List() ([]model.Alliance, error)
	AddPartyToAlliance(*model.AllianceParty) error
	GetByAbbrv(abbrv string) (*dto.AlliancePartyDetail, error)
	GetByID(id int) (*model.Alliance, error)
	ConsoleGetAllianceList() (*[]model.Alliance, error)
}

type allianceRepository struct {
	DB *gorm.DB
}

func NewAllianceRepository(db *gorm.DB) AllianceRepository {
	return &allianceRepository{
		DB: db,
	}
}

func (r *allianceRepository) Create(alliance *model.Alliance) error {
	if err := r.DB.Create(alliance).Error; err != nil {
		return err
	}
	return nil
}

func (r *allianceRepository) List() ([]model.Alliance, error) {
	var alliances []model.Alliance
	if err := r.DB.Order("name ASC").Find(&alliances).Error; err != nil {
		return nil, err
	}
	return alliances, nil
}

func (r *allianceRepository) AddPartyToAlliance(allianceParty *model.AllianceParty) error {
	if err := r.DB.Create(allianceParty).Error; err != nil {
		return err
	}
	return nil
}

func (r *allianceRepository) GetByAbbrv(abbrv string) (*dto.AlliancePartyDetail, error) {
	var allianceDetail dto.AlliancePartyDetail
	var alliance model.Alliance
	var alliancePartyDto []dto.AllianceParty
	var allianceParty []model.AllianceParty
	if err := r.DB.Where("LOWER(abbrv) = ?", strings.ToLower(abbrv)).First(&alliance).Error; err != nil {
		return nil, err
	}
	if err := r.DB.Where("alliance_id = ?", alliance.ID).Find(&allianceParty).Error; err != nil {
		return nil, err
	}
	for _, ap := range allianceParty {
		var allianceParty dto.AllianceParty
		var party model.Party
		if err := r.DB.First(&party, ap.PartyID).Error; err != nil {
			return nil, err
		}
		allianceParty.AllianceParty = ap
		allianceParty.Party = party
		alliancePartyDto = append(alliancePartyDto, allianceParty)
	}
	allianceDetail.Alliance = alliance
	allianceDetail.AllianceParty = alliancePartyDto
	return &allianceDetail, nil
}

func (r *allianceRepository) GetByID(id int) (*model.Alliance, error) {
	var alliance model.Alliance
	if err := r.DB.Where("id = ? AND is_active = ?", id, true).First(&alliance).Error; err != nil {
		return nil, err
	}
	return &alliance, nil
}

func (r *allianceRepository) ConsoleGetAllianceList() (*[]model.Alliance, error) {
	var alliance []model.Alliance
	if err := r.DB.Order("name asc").Find(&alliance).Error; err != nil {
		return nil, err
	}
	return &alliance, nil
}
