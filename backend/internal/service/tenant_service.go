package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type DashboardKPIs struct {
	TodayAppointmentsCount int                  `json:"today_appointments_count"`
	ConfirmedCount         int                  `json:"confirmed_count"`
	CompletedCount         int                  `json:"completed_count"`
	TodayRevenue           float64              `json:"today_revenue"`
	TotalCustomersCount    int64                `json:"total_customers_count"`
	ActiveProfessionals    int64                `json:"active_professionals_count"`
	TodayAppointments      []domain.Appointment `json:"today_appointments"`
}

type TenantService struct {
	repo *postgres.Repository
}

func NewTenantService(repo *postgres.Repository) *TenantService {
	return &TenantService{repo: repo}
}

func (s *TenantService) GetDashboardKPIs(ctx context.Context, tenantID uuid.UUID) (*DashboardKPIs, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	apts, err := s.repo.ListAppointments(ctx, tenantID, postgres.AppointmentFilter{
		DateFrom: &startOfDay,
		DateTo:   &endOfDay,
	})
	if err != nil {
		return nil, err
	}

	var confirmed, completed int
	var revenue float64

	for _, a := range apts {
		if a.Status == domain.StatusConfirmed {
			confirmed++
			revenue += a.TotalPrice
		} else if a.Status == domain.StatusCompleted {
			completed++
			revenue += a.TotalPrice
		}
	}

	var totalCust int64
	s.repo.DB().Model(&domain.Customer{}).Where("tenant_id = ?", tenantID).Count(&totalCust)

	var totalPros int64
	s.repo.DB().Model(&domain.Professional{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&totalPros)

	return &DashboardKPIs{
		TodayAppointmentsCount: len(apts),
		ConfirmedCount:         confirmed,
		CompletedCount:         completed,
		TodayRevenue:           revenue,
		TotalCustomersCount:    totalCust,
		ActiveProfessionals:    totalPros,
		TodayAppointments:      apts,
	}, nil
}
