package handler

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/auth"
	"github.com/shappy0/mantritva/pkg/log"
	"github.com/shappy0/mantritva/pkg/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ConsoleHandler struct {
	Service     service.ConsoleService
	FeedService service.FeedService
}

func NewConsoleHandler(services *service.Services) *ConsoleHandler {
	return &ConsoleHandler{
		Service:     services.ConsoleService,
		FeedService: services.FeedService,
	}
}

func (h *ConsoleHandler) GetUsersList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetUsersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get users list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetStatesList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetStatesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get states list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetPartiesList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetPartiesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get parties list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetAlliancesList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetAlliancesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get alliances list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetMembersList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetMembersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get members list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetMemberByID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	id := c.Param("id")
	idInt, err := utils.StrToInt(id)
	if err != nil {
		utils.Error(c, 500, "Invalid member id", err)
		return
	}
	member, err := h.Service.GetByID(uint(idInt))
	if err != nil {
		utils.Error(c, 500, "Failed to fetch member", err)
		return
	}
	utils.Success(c, "", member)
}

func (h *ConsoleHandler) UpdateMemberByID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	updatedBy, err := auth.GetKeyInCtx[string](c, "name")
	if err != nil {
		utils.Error(c, 401, "Unauthorized Access", nil)
		return
	}
	id := c.Param("id")
	if id == "" {
		utils.Error(c, 400, "Invalid memberID", nil)
		return
	}
	idInt, err := utils.StrToInt(id)
	if err != nil {
		utils.Error(c, 500, "Invalid memberID", err)
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		utils.Error(c, 400, "Invalid form data", err)
		return
	}
	req := map[string]interface{}{}
	for k, v := range form.Value {
		if len(v) > 0 {
			req[k] = v[0]
		}
	}
	req["updated_by"] = updatedBy
	req["updated_at"] = time.Now()

	photo, err := c.FormFile("photo")
	if err != nil && err != http.ErrMissingFile {
		utils.Error(c, 400, "Failed to get photo from request", err)
		return
	}

	updatedMember, err := h.Service.UpdateMemberByID(c, uint(idInt), req, photo)
	if err != nil {
		utils.Error(c, 500, "Failed to update member", err)
		return
	}
	utils.Success(c, "Member updated successfully", updatedMember)
}

func (h *ConsoleHandler) GetRulersList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetRulersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get rulers list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetCurrentRulersList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetCurrentRulersList()
	if err != nil {
		utils.Error(c, 500, "Failed to get current rulers list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetLeadershipByRulerID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	rulerIDStr := c.Param("rulerID")
	isCentralStr := c.Query("isCentral")
	// Convert string to int
	rulerID, err := strconv.Atoi(rulerIDStr)
	if err != nil {
		utils.Error(c, 400, "Invalid rulerID", err)
		return
	}
	resp, err := h.Service.GetLeadershipByRulerID(rulerID, utils.StrToBool(isCentralStr))
	if err != nil {
		log.Errorf("Failed to get leadership for rulerID %d: %v", rulerID, err)
		utils.Error(c, 500, "Failed to get leadership for ruler", err)
		return
	}
	var results []dto.LeadershipRes
	for _, r := range resp {
		res := dto.LeadershipRes{}
		err = utils.MapToStruct(r, &res)
		if err != nil {
			log.Errorf("Failed to parse leadership data for rulerID %d: %v", rulerID, err)
			continue
		}
		utils.TransformPositionLabel(&res.Position.Position)
		results = append(results, res)
	}
	utils.Success(c, "", results)
}

func (h *ConsoleHandler) GetNewsPressList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetNewsPressList()
	if err != nil {
		utils.Error(c, 500, "Failed to get news press list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetNewsPressPostByID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	id := c.Param("id")
	resp, err := h.Service.GetNewsPressPostByID(id)
	if err != nil {
		utils.Error(c, 500, "Failed to get news press post", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) UpdateNewsPressPost(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	updatedBy := authUser.Name

	id := c.Param("id")
	if id == "" {
		utils.Error(c, 400, "Invalid newsPressPost ID", nil)
		return
	}

	var req map[string]interface{}
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req["updated_by"] = updatedBy
	req["updated_at"] = time.Now()
	ok, err := h.Service.UpdateNewsPressPost(id, req)
	if err != nil {
		utils.Error(c, 400, "Failed to update NewsPress post", err)
		return
	}
	if status, ok := req["status"].(string); ok {
		if status == string(utils.Published) {
			postId, err := primitive.ObjectIDFromHex(id)
			if err != nil {
				utils.Error(c, 500, "Failed to create newsPress feed record", err)
				return
			}
			feed := dto.FeedReq{
				Type:   utils.NEWSPRESS,
				PostID: postId,
				Status: utils.PostStatus(status),
			}
			_, err = h.FeedService.UpdateStatus(updatedBy, feed)
			if err != nil {
				utils.Error(c, 500, "Failed to create newsPress feed record", err)
				return
			}
		}
	}
	utils.Success(c, "NewsPress Post Updated!", ok)
}

func (h *ConsoleHandler) GetGrievancesList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetGrievancesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get grievances list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GetGrievanceByID(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	grievanceID := c.Param("grievanceID")
	resp, err := h.Service.GetGrievanceByID(grievanceID)
	if err != nil {
		utils.Error(c, 500, "Failed to get grievance", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) GrievanceUpdateStatus(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}

	grievanceID := c.Param("grievanceID")
	var req dto.GrievanceStatus
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	resp, err := h.Service.GrievanceUpdateStatus(grievanceID, req)
	if err != nil {
		utils.Error(c, 500, "Failed to update grievance status", err)
		return
	}
	if req.Status == utils.OPEN {
		postID, err := primitive.ObjectIDFromHex(req.ID)
		if err != nil {
			utils.Error(c, 500, "Invalid grievance ID", err)
			return
		}
		_, err = h.FeedService.Create(dto.FeedReq{
			Type:   utils.GRIEVANCE,
			PostID: postID,
			Status: utils.Published,
		})
		if err != nil {
			utils.Error(c, 500, "Failed to create grievance feed", err)
			return
		} else {
			fmt.Printf("Grievance feed created successfully for grievance ID: %s\n", req.ID)
		}
	}
	utils.Success(c, "Grievance status updated successfully", resp)
}

func (h *ConsoleHandler) ExitLeadership(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	leadershipIDStr := c.Param("id")
	var req dto.LeadershipExitReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	leadershipID, err := strconv.Atoi(leadershipIDStr)
	if err != nil {
		utils.Error(c, 400, "Invalid leadership ID", err)
		return
	}
	req.LeadershipID = uint(leadershipID)
	resp, err := h.Service.ExitLeadership(req)
	if err != nil {
		utils.Error(c, 500, "Failed to exit leadership", err)
		return
	}
	utils.Success(c, "Leadership exited successfully", resp)
}

func (h *ConsoleHandler) UpdateRole(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	var req dto.UpdateUserReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	req.UpdatedAt = time.Now()
	req.UpdatedBy = authUser.UserName
	ok, err := h.Service.UpdateRole(req)
	if err != nil || !ok {
		utils.Error(c, 500, "Failed to update user role", err)
		return
	}
	if ok {
		utils.Success(c, "Role updated successfully", nil)
	}
}

func (h *ConsoleHandler) AddRole(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	var req dto.AddRoleReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	role, err := h.Service.AddRole(authUser.UserName, req)
	if err != nil {
		utils.Error(c, 500, "Failed to add user role", err)
		return
	}
	if role != nil {
		utils.Success(c, "Role added successfully", role)
	} else {
		utils.Error(c, 500, "Failed to add user role", nil)
	}
}

func (h *ConsoleHandler) GetUserRoles(c *gin.Context) {
	userID := c.Param("userID")
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetUserRoles(userID)
	if err != nil {
		utils.Error(c, 500, "Failed to get user roles", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) ListRoleScopes(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	utils.Success(c, "RoleScopes fetched successfully", utils.RoleRoleScopeMap)
}

func (h *ConsoleHandler) GetEnquiriesList(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsAdminOrSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	resp, err := h.Service.GetEnquiriesList()
	if err != nil {
		utils.Error(c, 500, "Failed to get enquiries list", err)
		return
	}
	utils.Success(c, "", resp)
}

func (h *ConsoleHandler) SetupBulkLeader(c *gin.Context) {
	authUser := auth.GetAuthUser(c)
	if !authUser.Found || !authUser.IsSuperAdmin() {
		utils.Error(c, 401, "Access Denied", nil)
		return
	}
	fh, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "Failed to get uploaded file", err)
		return
	}

	f, err := fh.Open()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "Failed to open uploaded file", err)
		return
	}
	defer f.Close()

	ext := filepath.Ext(fh.Filename)
	if ext != ".json" {
		utils.Error(c, http.StatusBadRequest, "Only .json file accepted", err)
		return
	}

	success, errored, err := h.Service.SetupBulkLeader(authUser.UserName, f)
	if err != nil {
		utils.Error(c, 500, "Failed to upload bulk leaders", err)
		return
	}
	utils.Success(c, "Leaders uploaded successfully", map[string]any{
		"success": success,
		"errored": errored,
	})
}
