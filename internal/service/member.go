package service

import (

	// "path/filepath"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/model"
	"github.com/shappy0/mantritva/internal/repository"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/storage"
	"github.com/shappy0/mantritva/pkg/utils"
)

type MemberService interface {
	List() ([]dto.MemberRes, error)
	Create(c *gin.Context, req dto.MemberReq, photo *multipart.FileHeader) (*model.Member, error)
	BulkMemberUpload(c *gin.Context, filePath string) (bool, error)
	GetByName(name string) (*dto.MemberRes, error)
	GetByID(id uint) (*dto.MemberRes, error)
	ListUnaffiliatedMembers() ([]dto.MemberRes, error)
	GetCareerByMemberID(int) ([]map[string]interface{}, error)
	GetPostsByMemberID(int, dto.Pagination, int) (*[]dto.NewsPressPostRes, error)
	GetGrievancesByMemberID(dto.Pagination, int, string) (*[]dto.GrievanceShort, error)
	GetDevelopmentWelfareByMemberID(dto.Pagination, int) (*[]dto.DevelopmentWelfareShort, error)
	// GetByID(id uint) (*dto.Member, error)
}

type memberService struct {
	Repository                   repository.MemberRepository
	StateRepository              repository.StateRepository
	LeadershipRepository         repository.LeadershipRepository
	PartyRepository              repository.PartyRepository
	NewsPressRepository          repository.NewsPressPostRepository
	GrievanceRepository          repository.GrievanceRepository
	ReportRepository             repository.ReportRepository
	DevelopmentWelfareRepository repository.DevelopmentWelfareRepository
	FeedRepository               repository.FeedRepository
	ReactionRepository           repository.ReactionRepository
}

func NewMemberService(repos *repository.Repos) MemberService {
	return &memberService{
		Repository:                   repos.Member,
		StateRepository:              repos.State,
		LeadershipRepository:         repos.Leadership,
		PartyRepository:              repos.Party,
		NewsPressRepository:          repos.NewsPressPost,
		GrievanceRepository:          repos.Grievance,
		ReportRepository:             repos.Report,
		DevelopmentWelfareRepository: repos.DevelopmentWelfare,
		FeedRepository:               repos.Feed,
		ReactionRepository:           repos.Reaction,
	}
}

func (s *memberService) Create(ctx *gin.Context, req dto.MemberReq, photo *multipart.FileHeader) (*model.Member, error) {
	destinationPath := filepath.Join("members")
	storage, err := storage.New(destinationPath)
	if err != nil {
		return nil, err
	}
	fileMeta, err := storage.Save(ctx, photo)
	if err != nil {
		return nil, err
	}
	memberName := req.Name
	if memberName != "" {
		memberName = strings.ToUpper(string(req.Name[0])) + req.Name[1:]
	}
	member := &model.Member{
		Name:             memberName,
		FullName:         req.FullName,
		DOB:              req.DOB,
		Gender:           req.Gender,
		Qualifications:   req.Qualifications,
		MaritalStatus:    req.MaritalStatus,
		FinancialAssets:  req.FinancialAssets,
		SocialMediaLinks: req.SocialMediaLinks,
		Photo:            fileMeta.DownloadURL,
	}
	if err := s.Repository.Create(member); err != nil {
		storage.Delete(fileMeta.AbsolutePath)
		return nil, err
	}
	if len(req.CriminalRecords) > 0 {
		for _, record := range req.CriminalRecords {
			criminalRecord := &model.CriminalRecord{
				MemberID:       member.ID,
				Offense:        record.Offense,
				Status:         record.Status,
				ConvictionDate: record.ConvictionDate,
				Sentence:       record.Sentence,
			}
			if err := s.Repository.CreateCriminalRecord(criminalRecord); err != nil {
				return nil, err
			}
		}
	}

	return member, nil
}

func (s *memberService) List() ([]dto.MemberRes, error) {
	var result []dto.MemberRes
	members, err := s.Repository.List()
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		dtoMember := dto.MemberRes{Member: member}
		//Get Party
		party, err := s.PartyRepository.GetByMemberID(member.ID)
		if err != nil {
			log.Errorf("Failed to get party details by memberID %d", member.ID)
		} else {
			dtoMember.Party = dto.MemberParty{
				Name:        party.Name,
				Abbrv:       party.Abbrv,
				Symbol:      party.Symbol,
				JoiningDate: party.JoiningDate,
				ExitDate:    party.ExitDate,
				Position:    party.Position,
			}
		}
		result = append(result, dtoMember)
	}
	return result, nil
}

func (s *memberService) BulkMemberUpload(c *gin.Context, filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Decode JSON into struct
	var memberReq dto.MemberReq
	if err := json.NewDecoder(file).Decode(&memberReq); err != nil {
		panic(err)
	}
	return true, nil
}

func (s *memberService) GetByName(name string) (*dto.MemberRes, error) {
	var data dto.MemberRes
	memberName := strings.Replace(name, "_", " ", -1)
	memberName = strings.ToUpper(string(memberName[0])) + memberName[1:]
	member, err := s.Repository.GetByName(memberName)
	if err != nil {
		log.Errorf("Member not found %s", memberName)
		return nil, err
	}
	data.Member = *member
	//Get Party
	party, err := s.PartyRepository.GetByMemberID(member.ID)
	if err != nil {
		log.Warnf("Party not found for member %s", memberName)
	} else {
		data.Party = *party
	}

	//Get position from leadership table
	leadership, err := s.LeadershipRepository.GetByMemberID(member.ID)
	if err != nil {
		log.Warnf("Leadership not found for member %s", memberName)
	} else {
		data.Position = leadership.Position
		data.Departments = leadership.Departments
	}
	// currentYear := time.Now().Format("2006")
	var totalWeightedScore float64
	var totalConfidence float64
	memberScore := &dto.MemberScore{MaxScore: 100}

	reports, err := s.ReportRepository.GetMemberOverallScore(member.ID)
	if err != nil {
		log.Warnf("No reports found for memberID %d", member.ID)
	}
	for _, report := range reports {
		if report.Confidence <= 0 {
			continue
		}
		totalWeightedScore += report.TotalScore * report.Confidence
		totalConfidence += report.Confidence
	}
	if totalConfidence == 0 {
		memberScore.Score = 0
		memberScore.Confidence = 0
	} else {
		memberScore.Score = utils.RoundToOneDecimal(totalWeightedScore / totalConfidence)
		memberScore.Confidence = utils.RoundToOneDecimal(totalConfidence / float64(len(reports)))
	}
	data.Score = memberScore
	return &data, nil
}

func (s *memberService) GetByID(id uint) (*dto.MemberRes, error) {
	var data dto.MemberRes
	member, err := s.Repository.GetByID(id)
	if err != nil {
		return nil, err
	}
	data.Member = *member
	//Get Party
	party, err := s.PartyRepository.GetByMemberID(member.ID)
	if err != nil {
		return nil, err
	}
	data.Party = *party

	//Get position from leadership table
	leadership, err := s.LeadershipRepository.GetByMemberID(member.ID)
	if err != nil {
		return nil, err
	}
	data.Position = leadership.Position
	data.Departments = leadership.Departments

	var totalWeightedScore float64
	var totalConfidence float64
	memberScore := &dto.MemberScore{MaxScore: 100}
	reports, err := s.ReportRepository.GetMemberOverallScore(member.ID)
	if err != nil {
		log.Warnf("No reports found for memberID %d", member.ID)
	}
	for _, report := range reports {
		if report.Confidence <= 0 {
			continue
		}
		totalWeightedScore += report.TotalScore * report.Confidence
		totalConfidence += report.Confidence
	}
	if totalConfidence == 0 {
		memberScore.Score = 0
		memberScore.Confidence = 0
	} else {
		memberScore.Score = utils.RoundToOneDecimal(totalWeightedScore / totalConfidence)
		memberScore.Confidence = utils.RoundToOneDecimal(totalConfidence / float64(len(reports)))
	}
	data.Score = memberScore
	return &data, nil
}

func (s *memberService) GetCareerByMemberID(memberID int) ([]map[string]interface{}, error) {
	careers, err := s.LeadershipRepository.GetCareerByMemberID(memberID)
	if err != nil {
		return nil, err
	}
	return careers, nil
}

func (s *memberService) ListUnaffiliatedMembers() ([]dto.MemberRes, error) {
	var result []dto.MemberRes
	members, err := s.Repository.ListUnaffiliatedMembers()
	if err != nil {
		return nil, err
	}
	for _, member := range members {
		dtoMember := dto.MemberRes{Member: member}
		result = append(result, dtoMember)
	}
	return result, nil
}

func (s *memberService) GetPostsByMemberID(userID int, pagination dto.Pagination, memberID int) (*[]dto.NewsPressPostRes, error) {
	var newsPressPostRes []dto.NewsPressPostRes
	posts, err := s.NewsPressRepository.GetPostsByMemberID(pagination, memberID)
	if err != nil {
		return nil, err
	}
	for _, post := range *posts {
		dtoPost := dto.NewsPressPostRes{NewsPressPost: post, Questionnaire: new([]dto.Questionnaire)}
		memberAndParty, err := s.Repository.GetMemberAndPartyDetails(memberID, int(post.PartyId))
		if err != nil {
			return nil, err
		}
		if err := utils.MapToStruct(memberAndParty, &dtoPost); err != nil {
			return nil, err
		}
		//Get Questionnairs
		ques, err := s.ReportRepository.GetQuestionnairesByCategory(string(utils.NEWSPRESS))
		if err != nil {
			log.Errorf("No questionnairs found for NewsPress post %v", post.ID)
		} else {
			for _, q := range *ques {
				*dtoPost.Questionnaire = append(*dtoPost.Questionnaire, dto.Questionnaire{
					ID:       q.ID,
					Title:    q.Title,
					VoteType: &q.VoteType,
				})
			}
		}
		if userID > 0 {
			reaction, err := s.ReactionRepository.GetByUserIDAndPostID(userID, post.ID)
			if err != nil {

			} else if reaction != nil {
				dtoPost.Reaction = utils.Vote(reaction.Reaction)
			}
		}
		newsPressPostRes = append(newsPressPostRes, dtoPost)
	}
	return &newsPressPostRes, nil
}

func (s *memberService) GetGrievancesByMemberID(pagination dto.Pagination, memberID int, position string) (*[]dto.GrievanceShort, error) {
	var grievancesRes []dto.GrievanceShort
	grievances, err := s.GrievanceRepository.GetByMemberID(pagination, memberID, position)
	if err != nil {
		return nil, err
	}
	for _, grievance := range *grievances {
		dtoGrievance := dto.GrievanceShort{
			GrievanceID:    grievance.GrievanceID,
			Title:          grievance.Title,
			Type:           grievance.Type,
			SubType:        grievance.SubType,
			Level:          grievance.Level,
			CenterSubLevel: grievance.CenterSubLevel,
			DateOfIncident: utils.DateOnly{Time: &grievance.DateOfIncident},
			Location:       grievance.Location,
			Impact:         grievance.Impact,
			FinalStatus:    grievance.FinalStatus,
			CreatedAt:      grievance.CreatedAt,
		}
		grievancesRes = append(grievancesRes, dtoGrievance)
	}
	return &grievancesRes, nil
}

func (s *memberService) GetDevelopmentWelfareByMemberID(pagination dto.Pagination, memberID int) (*[]dto.DevelopmentWelfareShort, error) {
	var DevelopmentWelfareShort []dto.DevelopmentWelfareShort
	developments, err := s.DevelopmentWelfareRepository.GetByMemberID(memberID)
	if err != nil {
		return nil, err
	}
	for _, dev := range developments {
		dtoDev := dto.DevelopmentWelfareShort{
			ID:               dev.ID.Hex(),
			Category:         dev.Category,
			Type:             dev.Type,
			Name:             dev.Name,
			Description:      dev.Description,
			AnnouncementDate: dev.AnnouncementDate,
			CreatedAt:        dev.CreatedAt,
			UpvoteCount:      *dev.UpvoteCount,
			DownvoteCount:    *dev.DownvoteCount,
			Status:           utils.DevelopmentStatus(dev.Status),
		}
		DevelopmentWelfareShort = append(DevelopmentWelfareShort, dtoDev)
	}
	return &DevelopmentWelfareShort, nil
}

func (s *memberService) GetMemberInfo(memberID string) (*dto.Member, error) {
	memId, err := utils.StrToInt(memberID)
	if err != nil {
		return nil, err
	}
	member, err := s.Repository.GetByID(uint(memId))
	if err != nil {
		return nil, err
	}
	memberInfo := &dto.Member{Name: member.Name, Photo: member.Photo}
	leader, err := s.LeadershipRepository.GetByMemberID(member.ID)
	if err == nil {
		if leader.Role != "" {
			memberInfo.Position = leader.Role
		} else {
			memberInfo.Position = leader.Position
		}
	}
	if leader != nil {
		if leader.StateID > 0 {
			state, err := s.StateRepository.GetByID(int(leader.StateID))
			if err != nil {
				log.Errorf("Error while getting state name of member %s", memberID)
			} else {
				memberInfo.Position += ", " + state.Name
			}
		}
	}
	memberParty, err := s.PartyRepository.GetByMemberID(member.ID)
	if err != nil {
		log.Errorf("No memberParty found for member %d", member.ID)
	} else if memberParty != nil {
		memberInfo.Position = memberParty.Position
		memberInfo.Party = dto.MemberParty{Name: memberParty.Name, Abbrv: memberParty.Abbrv, Symbol: memberParty.Symbol}
	}

	return memberInfo, nil
}
