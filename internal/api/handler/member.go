package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
)

type MemberHandler struct {
	Service service.MemberService
	Report  service.ReportService
}

func NewMemberHandler(svc *service.Services) *MemberHandler {
	return &MemberHandler{
		Service: svc.MemberService,
		Report:  svc.ReportService,
	}
}

func (h *MemberHandler) CreateMember(c *gin.Context) {
	var req dto.MemberReq
	if err := c.ShouldBind(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	photo, err := c.FormFile("photo")
	if err != nil {
		if err.Error() != "http: no such file" {
			utils.Error(c, 400, "Failed to get photo from request", err)
			return
		} else {
			log.Warnf("No member photo uploaded with the request")
		}
	}
	reps, err := h.Service.Create(c, req, photo)
	if err != nil {
		utils.Error(c, 500, "Failed to create member", err)
		return
	}
	utils.Success(c, "Member created successfully", reps)
}

func (h *MemberHandler) ListMember(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch member", nil)
		return
	}
	utils.Success(c, "", val)
}

func (h *MemberHandler) BulkMemberUpload(c *gin.Context) {
	filePath := c.PostForm("filePath")

	if filePath == "" {
		utils.Error(c, 400, "Failed to get filePath from request", nil)
		return
	}

	success, errored, err := h.Service.BulkMemberUpload(c, filePath)
	if err != nil {
		utils.Error(c, 500, "Failed to upload bulk members", err)
		return
	}
	utils.Success(c, "Members uploaded successfully", map[string]any{
		"success": success,
		"errored": errored,
	})
}

func (h *MemberHandler) GetMemberByName(c *gin.Context) {
	name := c.Param("name")
	if name == "" {
		utils.Error(c, 500, "Invalid memberName", nil)
		return
	}
	member, err := h.Service.GetByName(name)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch member", err)
		return
	}
	utils.Success(c, "", member)
}

func (h *MemberHandler) GetMemberCareer(c *gin.Context) {
	var careersDto []dto.CareerRes
	memberID := c.Param("memberID")
	memberIdInt, err := utils.StrToInt(memberID)
	if err != nil {
		utils.Error(c, 400, "Invalid member ID", nil)
		return
	}
	careers, err := h.Service.GetCareerByMemberID(memberIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch member career", err)
		return
	}
	for _, career := range careers {
		careerDto := dto.CareerRes{}
		if err := utils.MapToStruct(career, &careerDto); err != nil {
			log.Errorf("Failed to parse career data. Error: %v", err)
		} else {
			utils.TransformPositionLabel(&careerDto.Position.Position)
			careersDto = append(careersDto, careerDto)
		}
	}
	utils.Success(c, "", careersDto)
}

func (h *MemberHandler) GetGrievancesByMemberID(c *gin.Context) {
	memberID := c.Param("memberID")
	memberIdInt, err := utils.StrToInt(memberID)
	position := c.Query("position")
	if err != nil {
		utils.Error(c, 400, "Invalid member ID", nil)
		return
	}
	pagination := dto.NewPagination(c)
	grievances, err := h.Service.GetGrievancesByMemberID(pagination, memberIdInt, position)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch grievances", err)
		return
	}
	utils.Success(c, "", grievances)
}

func (h *MemberHandler) ListUnaffiliatedMembers(c *gin.Context) {
	val, err := h.Service.ListUnaffiliatedMembers()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch unaffiliated members", nil)
		return
	}
	utils.Success(c, "Unaffiliated members fetched successfully", val)
}

func (h *MemberHandler) GetPostsByMemberID(c *gin.Context) {
	userID, err := auth.GetKeyInCtx[int](c, "user_id")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	memberID := c.Param("memberID")
	memberIdInt, err := utils.StrToInt(memberID)
	if err != nil {
		utils.Error(c, 400, "Invalid member ID", nil)
		return
	}
	pagination := dto.NewPagination(c)
	posts, err := h.Service.GetPostsByMemberID(userID, pagination, memberIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch news press posts", err)
		return
	}
	utils.Success(c, "News press posts fetched successfully", posts)
}

func (h *MemberHandler) GetDevelopmentWelfareByMemberID(c *gin.Context) {
	memberID := c.Param("memberID")
	memberIdInt, err := utils.StrToInt(memberID)
	if err != nil {
		utils.Error(c, 400, "Invalid member ID", nil)
		return
	}
	pagination := dto.NewPagination(c)
	posts, err := h.Service.GetDevelopmentWelfareByMemberID(pagination, memberIdInt)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch development welfare posts", err)
		return
	}
	utils.Success(c, "", posts)
}

func (h *MemberHandler) GetMemberReportCard(c *gin.Context) {
	memberName := c.Param("name")
	if memberName == "" {
		utils.Error(c, 400, "Invalid member name", nil)
		return
	}
	from := c.Param("from")
	if from == "" {
		utils.Error(c, 400, "Please pass From Year", nil)
		return
	}
	fromx, err := utils.StrToInt(from)
	if err != nil {
		utils.Error(c, 400, "Invalid from year", nil)
		return
	}
	to := c.Param("to")
	if to == "" {
		utils.Error(c, 400, "Please pass To Year", nil)
		return
	}
	tox, err := utils.StrToInt(to)
	if err != nil {
		utils.Error(c, 400, "Invalid to year", nil)
		return
	}
	reportCard := dto.MemberReportCard{From: fromx, To: tox}
	member, err := h.Service.GetByName(memberName)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch member detail", err)
		return
	}
	reportCard.Member = dto.Member{
		Name:     member.Name,
		Photo:    member.Photo,
		Position: member.Position,
		Party:    member.Party,
	}
	reportCard.OverallScore = *member.Score
	reportCard.Reports = make([]dto.MemberReport, 0)

	reports, err := h.Report.GetMemberReports(int(member.ID), fromx, tox)
	if err != nil {
	} else {
		reportCard.Reports = reports
	}
	utils.Success(c, "", reportCard)
}

func (h *MemberHandler) GetTopScorer(c *gin.Context) {
	scorers, err := h.Report.GetTopScorer()
	if err != nil {
		utils.Error(c, 500, "Failed to top scorer members", nil)
		return
	}
	if scorers != nil {
		for _, memberScore := range *scorers {
			member, err := h.Service.GetByID(toUint(memberScore["member_id"]))
			if err != nil {
				log.Errorf("Error fetching top scorer's member data %v", err)
				continue
			}
			memberScore["member"] = dto.Member{
				Name:     member.Name,
				Photo:    member.Photo,
				Position: member.Position,
				Party:    member.Party,
			}
		}
	}
	utils.Success(c, "", scorers)
}

func toUint(v interface{}) uint {
	switch val := v.(type) {
	case int:
		return uint(val)
	case int32:
		return uint(val)
	case int64:
		return uint(val)
	case float32:
		return uint(val)
	case float64:
		return uint(val)
	case uint:
		return val
	case uint32:
		return uint(val)
	case uint64:
		return uint(val)
	default:
		return 0
	}
}
