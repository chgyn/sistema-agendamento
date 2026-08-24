package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/middleware"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/response"
)

type ProfessionalHandler struct {
	repo *postgres.Repository
}

func NewProfessionalHandler(repo *postgres.Repository) *ProfessionalHandler {
	return &ProfessionalHandler{repo: repo}
}

func (h *ProfessionalHandler) List(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	pros, err := h.repo.ListProfessionals(c.Request.Context(), tenantID, false)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar profissionais")
		return
	}
	response.Success(c, pros)
}

func (h *ProfessionalHandler) GetByID(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	pro, err := h.repo.GetProfessionalByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Profissional não encontrado")
		return
	}
	response.Success(c, pro)
}

func (h *ProfessionalHandler) Create(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)

	var req struct {
		Name       string      `json:"name" binding:"required"`
		Email      string      `json:"email"`
		Phone      string      `json:"phone"`
		Title      string      `json:"title"`
		Specialty  string      `json:"specialty"`
		Bio        string      `json:"bio"`
		AvatarURL  string      `json:"avatar_url"`
		IsActive   *bool       `json:"is_active"`
		ServiceIDs []uuid.UUID `json:"service_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	isActive := true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	pro := domain.Professional{
		ID:        uuid.New(),
		TenantID:  tenantID,
		Name:      req.Name,
		Email:     req.Email,
		Phone:     req.Phone,
		Title:     req.Title,
		Specialty: req.Specialty,
		Bio:       req.Bio,
		AvatarURL: req.AvatarURL,
		IsActive:  isActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.repo.CreateProfessional(c.Request.Context(), &pro, req.ServiceIDs); err != nil {
		response.InternalServerError(c, "Erro ao cadastrar profissional")
		return
	}

	// Cria horários de trabalho padrão (Seg-Sex 08:00 - 18:00)
	defaultWhs := make([]domain.WorkingHour, 0)
	for day := 1; day <= 5; day++ {
		defaultWhs = append(defaultWhs, domain.WorkingHour{
			ID:             uuid.New(),
			TenantID:       tenantID,
			ProfessionalID: pro.ID,
			DayOfWeek:      day,
			StartTime:      "08:00",
			EndTime:        "18:00",
			BreakStart:     "12:00",
			BreakEnd:       "13:00",
			IsActive:       true,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}
	_ = h.repo.SaveWorkingHours(c.Request.Context(), tenantID, pro.ID, defaultWhs)

	response.Created(c, pro, "Profissional cadastrado com sucesso")
}

func (h *ProfessionalHandler) Update(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	pro, err := h.repo.GetProfessionalByID(c.Request.Context(), tenantID, id)
	if err != nil {
		response.NotFound(c, "Profissional não encontrado")
		return
	}

	var req struct {
		Name       string      `json:"name"`
		Email      string      `json:"email"`
		Phone      string      `json:"phone"`
		Title      string      `json:"title"`
		Specialty  string      `json:"specialty"`
		Bio        string      `json:"bio"`
		AvatarURL  string      `json:"avatar_url"`
		IsActive   *bool       `json:"is_active"`
		ServiceIDs []uuid.UUID `json:"service_ids"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos")
		return
	}

	if req.Name != "" {
		pro.Name = req.Name
	}
	pro.Email = req.Email
	pro.Phone = req.Phone
	pro.Title = req.Title
	pro.Specialty = req.Specialty
	pro.Bio = req.Bio
	if req.AvatarURL != "" {
		pro.AvatarURL = req.AvatarURL
	}
	if req.IsActive != nil {
		pro.IsActive = *req.IsActive
	}
	pro.UpdatedAt = time.Now()

	if err := h.repo.UpdateProfessional(c.Request.Context(), pro, req.ServiceIDs); err != nil {
		response.InternalServerError(c, "Erro ao atualizar profissional")
		return
	}

	response.Success(c, pro, "Profissional atualizado com sucesso")
}

func (h *ProfessionalHandler) Delete(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.repo.DeleteProfessional(c.Request.Context(), tenantID, id); err != nil {
		response.InternalServerError(c, "Erro ao remover profissional")
		return
	}

	response.Success(c, nil, "Profissional removido com sucesso")
}

// GetWorkingHours busca grade horária do profissional
func (h *ProfessionalHandler) GetWorkingHours(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	whs, err := h.repo.GetWorkingHoursByProfessional(c.Request.Context(), tenantID, id)
	if err != nil {
		response.InternalServerError(c, "Erro ao buscar horários")
		return
	}
	response.Success(c, whs)
}

// SaveWorkingHours atualiza os dias e horários de trabalho
func (h *ProfessionalHandler) SaveWorkingHours(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req []struct {
		DayOfWeek  int    `json:"day_of_week"`
		StartTime  string `json:"start_time"`
		EndTime    string `json:"end_time"`
		BreakStart string `json:"break_start"`
		BreakEnd   string `json:"break_end"`
		IsActive   bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	whs := make([]domain.WorkingHour, 0)
	for _, item := range req {
		whs = append(whs, domain.WorkingHour{
			DayOfWeek:  item.DayOfWeek,
			StartTime:  item.StartTime,
			EndTime:    item.EndTime,
			BreakStart: item.BreakStart,
			BreakEnd:   item.BreakEnd,
			IsActive:   item.IsActive,
		})
	}

	if err := h.repo.SaveWorkingHours(c.Request.Context(), tenantID, id, whs); err != nil {
		response.InternalServerError(c, "Erro ao salvar horários de trabalho")
		return
	}

	response.Success(c, whs, "Horários atualizados com sucesso")
}

// CreateException adiciona um bloqueio de agenda ou folga
func (h *ProfessionalHandler) CreateException(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	proID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	var req struct {
		Date      string `json:"date" binding:"required"` // "YYYY-MM-DD"
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Reason    string `json:"reason"`
		IsFullDay bool   `json:"is_full_day"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "Dados inválidos")
		return
	}

	exp := domain.AvailabilityException{
		ID:             uuid.New(),
		TenantID:       tenantID,
		ProfessionalID: proID,
		Date:           req.Date,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Reason:         req.Reason,
		IsFullDay:      req.IsFullDay,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := h.repo.CreateException(c.Request.Context(), &exp); err != nil {
		response.InternalServerError(c, "Erro ao criar bloqueio de agenda")
		return
	}

	response.Created(c, exp, "Bloqueio registrado com sucesso")
}

// ListExceptions lista bloqueios do profissional
func (h *ProfessionalHandler) ListExceptions(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	proID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")

	exps, err := h.repo.ListExceptions(c.Request.Context(), tenantID, proID, dateFrom, dateTo)
	if err != nil {
		response.InternalServerError(c, "Erro ao listar bloqueios")
		return
	}

	response.Success(c, exps)
}

// DeleteException remove bloqueio
func (h *ProfessionalHandler) DeleteException(c *gin.Context) {
	tenantID, _ := middleware.GetTenantID(c)
	expID, err := uuid.Parse(c.Param("exception_id"))
	if err != nil {
		response.BadRequest(c, "ID inválido")
		return
	}

	if err := h.repo.DeleteException(c.Request.Context(), tenantID, expID); err != nil {
		response.InternalServerError(c, "Erro ao remover bloqueio")
		return
	}

	response.Success(c, nil, "Bloqueio removido com sucesso")
}
