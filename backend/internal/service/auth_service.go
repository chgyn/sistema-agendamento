package service

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
	"github.com/sistema-agendamento/backend/pkg/hash"
	"github.com/sistema-agendamento/backend/pkg/jwt"
)

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token        string               `json:"token"`
	User         *domain.User         `json:"user"`
	Tenant       *domain.Tenant       `json:"tenant"`
	Subscription *domain.Subscription `json:"subscription,omitempty"`
}

type AuthService struct {
	repo                *postgres.Repository
	jwtService          *jwt.JWTService
	subscriptionService *SubscriptionService
}

func NewAuthService(repo *postgres.Repository, jwtSvc *jwt.JWTService, subService *SubscriptionService) *AuthService {
	return &AuthService{
		repo:                repo,
		jwtService:          jwtSvc,
		subscriptionService: subService,
	}
}

func (s *AuthService) RegisterTenant(ctx context.Context, dto domain.RegisterTenantWithPlanDTO) (*AuthResponse, error) {
	slug := strings.ToLower(strings.TrimSpace(dto.Slug))

	passHash, err := hash.HashPassword(dto.Password)
	if err != nil {
		return nil, err
	}

	// Valida existência do plano selecionado
	var plan *domain.Plan
	if dto.PlanID != uuid.Nil {
		plan, err = s.repo.GetPlanByID(ctx, dto.PlanID)
		if err != nil {
			return nil, domain.ErrPlanNotFound
		}
	} else {
		// Se não foi fornecido (fallback), busca o primeiro plano ativo
		active := true
		plans, err := s.repo.ListPlans(ctx, "", &active)
		if err != nil || len(plans) == 0 {
			return nil, domain.ErrPlanNotFound
		}
		plan = &plans[0]
	}

	tenantID := uuid.New()
	tenant := domain.Tenant{
		ID:                  tenantID,
		Slug:                slug,
		Name:                dto.TenantName,
		Document:            dto.Document,
		Phone:               dto.Phone,
		City:                dto.City,
		State:               dto.State,
		SlotIntervalMinutes: 30,
		PrimaryColor:        "#10b981",
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := s.repo.CreateTenant(ctx, &tenant); err != nil {
		return nil, err
	}

	user := domain.User{
		ID:           uuid.New(),
		TenantID:     &tenantID,
		Name:         dto.AdminName,
		Email:        dto.AdminEmail,
		PasswordHash: passHash,
		Role:         domain.RoleAdminTenant,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	// Criação da assinatura vinculada ao plano e Asaas
	var subscription *domain.Subscription
	if s.subscriptionService != nil {
		sub, err := s.subscriptionService.CreateTenantSubscription(ctx, &tenant, plan, &user)
		if err != nil {
			// Não bloqueia o cadastro do tenant caso a API Asaas esteja offline em dev, sub persiste em mock
			return nil, err
		}
		subscription = sub
		tenant.Subscription = sub
	}

	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		User:         &user,
		Tenant:       &tenant,
		Subscription: subscription,
	}, nil
}

func (s *AuthService) Login(ctx context.Context, dto LoginDTO) (*AuthResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, dto.Email)
	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	if !user.IsActive {
		return nil, domain.ErrUnauthorized
	}

	if !hash.CheckPassword(dto.Password, user.PasswordHash) {
		return nil, domain.ErrInvalidCredentials
	}

	var tenant *domain.Tenant
	var sub *domain.Subscription
	if user.TenantID != nil && *user.TenantID != uuid.Nil {
		tenant, err = s.repo.GetTenantByID(ctx, *user.TenantID)
		if err != nil {
			return nil, domain.ErrTenantNotFound
		}
		if tenant != nil {
			sub, _ = s.repo.GetSubscriptionByTenantID(ctx, tenant.ID)
		}
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:        token,
		User:         user,
		Tenant:       tenant,
		Subscription: sub,
	}, nil
}
