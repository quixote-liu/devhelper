package api

import (
	"devhelper/internal/models"
	"devhelper/internal/repository"
	"devhelper/internal/utils"

	"github.com/gin-gonic/gin"
)

type JsonHandler struct {
	historyRepo repository.HistoryRepository
	schemaRepo  repository.SchemaRepository
}

func NewJsonHandler(hr repository.HistoryRepository, sr repository.SchemaRepository) *JsonHandler {
	return &JsonHandler{historyRepo: hr, schemaRepo: sr}
}

type jsonInput struct {
	JSON string `json:"json" binding:"required"`
}

type convertInput struct {
	JSON   string `json:"json" binding:"required"`
	Target string `json:"target" binding:"required"`
}

type parseInput struct {
	Content string `json:"content" binding:"required"`
	Source  string `json:"source" binding:"required"`
}

type schemaValidateInput struct {
	Schema string `json:"schema" binding:"required"`
	Data   string `json:"data" binding:"required"`
}

type diffInput struct {
	A string `json:"a" binding:"required"`
	B string `json:"b" binding:"required"`
}

type queryInput struct {
	JSON string `json:"json" binding:"required"`
	Path string `json:"path" binding:"required"`
}

func (h *JsonHandler) Validate(c *gin.Context) {
	var req jsonInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if err := validateJSON(req.JSON); err != nil {
		utils.OK(c, gin.H{"valid": false, "error": err.Error()})
		return
	}
	utils.OK(c, gin.H{"valid": true})
}

func (h *JsonHandler) Format(c *gin.Context) {
	var req struct {
		JSON   string `json:"json" binding:"required"`
		Indent int    `json:"indent"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	if req.Indent == 0 {
		req.Indent = 2
	}
	result, err := formatJSON(req.JSON, req.Indent)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"result": result})
}

func (h *JsonHandler) Minify(c *gin.Context) {
	var req jsonInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := minifyJSON(req.JSON)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"result": result})
}

func (h *JsonHandler) Convert(c *gin.Context) {
	var req convertInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := convertJSON(req.JSON, req.Target)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"result": result})
}

func (h *JsonHandler) Parse(c *gin.Context) {
	var req parseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := parseJSON(req.Content, req.Source)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"result": result})
}

func (h *JsonHandler) GenerateSchema(c *gin.Context) {
	var req jsonInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := generateSchema(req.JSON)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"schema": result})
}

func (h *JsonHandler) ValidateSchema(c *gin.Context) {
	var req schemaValidateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	errs, err := validateSchema(req.Schema, req.Data)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"valid": len(errs) == 0, "errors": errs})
}

func (h *JsonHandler) Diff(c *gin.Context) {
	var req diffInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := diffJSON(req.A, req.B)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"diff": result})
}

func (h *JsonHandler) Query(c *gin.Context) {
	var req queryInput
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	result, err := queryJSON(req.JSON, req.Path)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	utils.OK(c, gin.H{"result": result})
}

// History handlers

type saveHistoryReq struct {
	SessionID string `json:"session_id" binding:"required"`
	Content   string `json:"content" binding:"required"`
	IsBase    bool   `json:"is_base"`
	Note      string `json:"note"`
}

func (h *JsonHandler) SaveHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	uid := userID.(uint)
	var req saveHistoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	entry := &models.JsonHistory{
		UserID:    uid,
		SessionID: req.SessionID,
		SeqNum:    h.historyRepo.NextSeqNum(uid, req.SessionID),
		IsBase:    req.IsBase || !h.historyRepo.HasBase(uid, req.SessionID),
		Content:   req.Content,
		Note:      req.Note,
	}
	if err := h.historyRepo.Create(entry); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, entry)
}

func (h *JsonHandler) GetHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	sessionID := c.Query("session_id")
	if sessionID == "" {
		utils.BadRequest(c, "session_id required")
		return
	}
	items, err := h.historyRepo.ListBySession(userID.(uint), sessionID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, items)
}

func (h *JsonHandler) DeleteHistory(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	if err := h.historyRepo.Delete(id, userID.(uint)); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, nil)
}

// Schema handlers

type saveSchemaReq struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Schema      string `json:"schema" binding:"required"`
	IsPublic    bool   `json:"is_public"`
}

func (h *JsonHandler) ListSchemas(c *gin.Context) {
	userID, _ := c.Get("user_id")
	schemas, err := h.schemaRepo.List(userID.(uint))
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, schemas)
}

func (h *JsonHandler) SaveSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var req saveSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	s := &models.JsonSchema{
		UserID:      userID.(uint),
		Name:        req.Name,
		Description: req.Description,
		Schema:      req.Schema,
		IsPublic:    req.IsPublic,
	}
	if err := h.schemaRepo.Create(s); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, s)
}

func (h *JsonHandler) GetSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	s, err := h.schemaRepo.FindByID(id, userID.(uint))
	if err != nil {
		utils.NotFound(c, "schema not found")
		return
	}
	utils.OK(c, s)
}

func (h *JsonHandler) UpdateSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	s, err := h.schemaRepo.FindByID(id, userID.(uint))
	if err != nil {
		utils.NotFound(c, "schema not found")
		return
	}
	if s.UserID != userID.(uint) {
		utils.Forbidden(c, "not your schema")
		return
	}
	var req saveSchemaReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	s.Name = req.Name
	s.Description = req.Description
	s.Schema = req.Schema
	s.IsPublic = req.IsPublic
	if err := h.schemaRepo.Update(s); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, s)
}

func (h *JsonHandler) DeleteSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var id uint
	if _, err := parseUintParam(c, "id", &id); err != nil {
		return
	}
	if err := h.schemaRepo.Delete(id, userID.(uint)); err != nil {
		utils.InternalError(c, err.Error())
		return
	}
	utils.OK(c, nil)
}
