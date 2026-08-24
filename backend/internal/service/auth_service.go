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

type RegisterTenantDTO struct {
	TenantName  string `json:"tenant_name" binding:"required"`
	Slug        string `json:"slug" binding:"required"`
	Document    string `json:"document"`
	Phone       string `json:"phone" binding:"required"`
	City        string `json:"city"`
	State       string `json:"state"`
	AdminName   string `json:"admin_name" binding:"required"`
	AdminEmail  string `json:"admin_email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=6"`
}

type LoginDTO struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token  string       `json:"token"`
	User   *domain.User `json:"user"`
	Tenant *domain.Tenant `json:"tenant"`
}

type AuthService struct {
	repo       *postgres.Repository
	jwtService *jwt.JWTService
}

func NewAuthService(repo *postgres.Repository, jwtSvc *jwt.JWTService) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtService: jwtSvc,
	}
}

func (s *AuthService) RegisterTenant(ctx context.Context, dto RegisterTenantDTO) (*AuthResponse, error) {
	slug := strings.ToLower(strings.TrimSpace(dto.Slug))

	passHash, err := hash.HashPassword(dto.Password)
	if err != nil {
		return nil, err
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
		TenantID:     tenantID,
		Name:         dto.AdminName,
		Email:        dto.AdminEmail,
		PasswordHash: passHash,
		Role:         domain.RoleAdmin,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateUser(ctx, &user); err != nil {
		return nil, err
	}

	token, err := s.jwtService.GenerateToken(&user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:  token,
		User:   &user,
		Tenant: &tenant,
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

	tenant, err := s.repo.GetTenantByID(ctx, user.TenantID)
	if err != nil {
		return nil, domain.ErrTenantNotFound
	}

	token, err := s.jwtService.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:  token,
		User:   user,
		Tenant: tenant,
	}, nil
}
