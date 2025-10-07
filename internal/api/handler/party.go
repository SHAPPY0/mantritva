package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/shappy0/mantritva/internal/dto"
	"github.com/shappy0/mantritva/internal/service"
	"github.com/shappy0/mantritva/pkg/utils"
)

type PartyHandler struct {
	Service service.PartyService
}

func NewPartyHandler(serivce service.PartyService) *PartyHandler {
	return &PartyHandler{
		Service: serivce,
	}
}

func (h *PartyHandler) CreateParty(c *gin.Context) {
	var req dto.PartyReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	reps, err := h.Service.Create(req)
	if err != nil {
		utils.Error(c, 500, "Failed to create party", err)
		return
	}
	utils.Success(c, "Party created successfully", reps)
}

func (h *PartyHandler) ListParty(c *gin.Context) {
	val, err := h.Service.List()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch party", nil)
		return
	}
	utils.Success(c, "Party fetched successfully", val)
}

func (h *PartyHandler) GetPartyByAbbrv(c *gin.Context) {
	abbrv := c.Param("abbrv")
	if abbrv == "" {
		utils.Error(c, 400, "Abbreviation is required", nil)
		return
	}
	val, err := h.Service.GetByAbbrv(abbrv)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch party by abbreviation", err)
		return
	}
	if val == nil {
		utils.Error(c, 404, "Party not found", nil)
		return
	}
	utils.Success(c, "Party fetched successfully", val)
}

func (h *PartyHandler) UploadPartySymbol(c *gin.Context) {
	partyId := c.Param("partyId")
	if partyId == "" {
		utils.Error(c, 400, "Party ID is required", nil)
		return
	}

	file, err := c.FormFile("symbol")
	if err != nil {
		utils.Error(c, 400, "Failed to get file from request", err)
		return
	}

	ok, err := h.Service.UploadSymbol(c, partyId, file)
	if err != nil {
		utils.Error(c, 500, "Failed to upload party symbol", err)
		return
	}
	if !ok {
		utils.Error(c, 500, "Internal server error occured while uploading party symbol", err)
		return
	}
	utils.Success(c, "Party symbol uploaded successfully", nil)
}

func (h *PartyHandler) AddMember(c *gin.Context) {
	var req dto.AddMemberReq
	if err := c.BindJSON(&req); err != nil {
		utils.Error(c, 400, "Invalid request data", err)
		return
	}
	if req.PartyId == 0 || req.MemberId == 0 {
		utils.Error(c, 400, "Party ID and Member ID are required", nil)
		return
	}
	ok, err := h.Service.AddMember(req)
	if err != nil {
		utils.Error(c, 500, "Failed to add member to party", err)
		return
	}
	if !ok {
		utils.Error(c, 500, "Internal server error occured while adding member to party", err)
		return
	}
	utils.Success(c, "Member added to party successfully", nil)
}

func (h *PartyHandler) GetPartyMembers(c *gin.Context) {
	partyId := c.Param("partyId")
	if partyId == "" {
		utils.Error(c, 400, "Party ID is required", nil)
		return
	}

	members, err := h.Service.GetPartyMembers(partyId)
	if err != nil {
		utils.Error(c, 500, "Failed to fetch party members", err)
		return
	}
	utils.Success(c, "Party members fetched successfully", members)
}

func (h *PartyHandler) GetRandomParties(c *gin.Context) {
	parties, err := h.Service.GetRandomParties()
	if err != nil {
		utils.Error(c, 500, "Failed to fetch random parties", err)
		return
	}
	utils.Success(c, "", parties)
}
