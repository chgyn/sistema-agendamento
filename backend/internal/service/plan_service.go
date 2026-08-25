package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type PlanService struct {
	repo *postgres.Repository
}

func NewPlanService(repo *postgres.Repository) *PlanService {
	return &PlanService{repo: repo}
}

func (s *PlanService) Create(ctx context.Context, dto domain.CreatePlanDTO) (*domain.Plan, error) {
	isActive := true
	if dto.IsActive != nil {
		isActive = *dto.IsActive
	}

	plan := domain.Plan{
		ID:               uuid.New(),
		Name:             strings.TrimSpace(dto.Name),
		Description:      strings.TrimSpace(dto.Description),
		Price:            dto.Price,
		BillingCycle:     dto.BillingCycle,
		MaxProfessionals: dto.MaxProfessionals,
		MaxServices:      dto.MaxServices,
		Features:         dto.Features,
		SortOrder:        dto.SortOrder,
		IsActive:         isActive,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := s.repo.CreatePlan(ctx, &plan); err != nil {
		return nil, err
	}

	return &plan, nil
}

func (s *PlanService) Update(ctx context.Context, id uuid.UUID, dto domain.UpdatePlanDTO) (*domain.Plan, error) {
	plan, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return nil, err
	}

	plan.Name = strings.TrimSpace(dto.Name)
	plan.Description = strings.TrimSpace(dto.Description)
	plan.Price = dto.Price
	plan.BillingCycle = dto.BillingCycle
	plan.MaxProfessionals = dto.MaxProfessionals
	plan.MaxServices = dto.MaxServices
	plan.Features = dto.Features
	plan.SortOrder = dto.SortOrder
	if dto.IsActive != nil {
		plan.IsActive = *dto.IsActive
	}
	plan.UpdatedAt = time.Now()

	if err := s.repo.UpdatePlan(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

func (s *PlanService) ToggleStatus(ctx context.Context, id uuid.UUID, isActive bool) error {
	_, err := s.repo.GetPlanByID(ctx, id)
	if err != nil {
		return err
	}
	return s.repo.UpdatePlanStatus(ctx, id, isActive)
}

func (s *PlanService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeletePlan(ctx, id)
}

func (s *PlanService) GetByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	return s.repo.GetPlanByID(ctx, id)
}

func (s *PlanService) List(ctx context.Context, search string, onlyActive *bool) ([]domain.Plan, error) {
	return s.repo.ListPlans(ctx, search, onlyActive)
}

func (s *PlanService) ListActivePublic(ctx context.Context) ([]domain.Plan, error) {
	active := true
	return s.repo.ListPlans(ctx, "", &active)
}
