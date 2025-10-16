package repository

import (
	"strings"

	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"gorm.io/gorm"
)

type PartyRepository interface {
	Create(*model.Party) error
	List() ([]model.Party, error)
	GetByAbbrv(abbrv string) (*model.Party, error)
	GetByID(id int) (*model.Party, error)
	UploadSymbol(partyId string, fileName string) error
	AddMember(member *model.PartyMembers) error
	GetPartyMembers(partyId string) ([]model.PartyMembers, error)
	GetByMemberID(memberId uint) (*dto.MemberParty, error)
	ConsoleGetPartiesList() (*[]model.Party, error)
	GetRandomParties() (*[]model.Party, error)
}

type partyRepository struct {
	DB *gorm.DB
}

func NewPartyRepository(db *gorm.DB) PartyRepository {
	return &partyRepository{
		DB: db,
	}
}

func (r *partyRepository) Create(party *model.Party) error {
	if err := r.DB.Create(party).Error; err != nil {
		return err
	}
	return nil
}

func (r *partyRepository) List() ([]model.Party, error) {
	var parties []model.Party
	if err := r.DB.Order("name ASC").Find(&parties).Error; err != nil {
		return nil, err
	}
	return parties, nil
}

func (r *partyRepository) GetByAbbrv(abbrv string) (*model.Party, error) {
	var party model.Party
	if err := r.DB.Where("LOWER(abbrv) = ?", strings.ToLower(abbrv)).Find(&party).Error; err != nil {
		return nil, err
	}
	return &party, nil
}

func (r *partyRepository) GetByID(id int) (*model.Party, error) {
	var party model.Party
	if err := r.DB.Where("id = ?", id).Find(&party).Error; err != nil {
		return nil, err
	}
	return &party, nil
}

func (r *partyRepository) UploadSymbol(partyId string, fileName string) error {
	if err := r.DB.Model(&model.Party{}).Where("id = ?", partyId).Update("symbol", fileName).Error; err != nil {
		return err
	}
	return nil
}

func (r *partyRepository) AddMember(member *model.PartyMembers) error {
	if err := r.DB.Create(member).Error; err != nil {
		return err
	}
	return nil
}

func (r *partyRepository) GetPartyMembers(partyId string) ([]model.PartyMembers, error) {
	var members []model.PartyMembers
	if err := r.DB.Model(&model.PartyMembers{}).Where("party_id = ? AND exit_date IS NULL", partyId).Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (r *partyRepository) GetByMemberID(memberId uint) (*dto.MemberParty, error) {
	var memberParty dto.MemberParty
	query := `SELECT 
				p.name,
				p.abbrv, 
				p.symbol, 
				pm."position",
				pm.joining_date,
				pm.exit_date
			FROM party_members pm
				INNER JOIN parties p ON p.ID = pm.party_id
			WHERE pm.member_id = ?
				AND pm.exit_date IS NULL
				AND p.is_deleted = ?`
	if err := r.DB.Raw(query, memberId, false).Scan(&memberParty).Error; err != nil {
		return nil, err
	}
	return &memberParty, nil
}

func (r *partyRepository) ConsoleGetPartiesList() (*[]model.Party, error) {
	var parties []model.Party
	if err := r.DB.Order("name asc").Find(&parties).Error; err != nil {
		return nil, err
	}
	return &parties, nil
}

func (r *partyRepository) GetRandomParties() (*[]model.Party, error) {
	var parties []model.Party
	query := `SELECT * 
			FROM parties p
			ORDER BY RANDOM() 
			LIMIT 3`
	if err := r.DB.Raw(query).Scan(&parties).Error; err != nil {
		return nil, err
	}
	return &parties, nil
}
