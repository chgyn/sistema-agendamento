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

// GetDashboardAnalytics consolida dados de KPIs, série temporal (Reports) e distribuição (Analytics) no padrão Figma SAAS
func (s *TenantService) GetDashboardAnalytics(ctx context.Context, tenantID uuid.UUID) (*domain.DashboardAnalyticsResponse, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)

	// Agendamentos de hoje
	todayApts, err := s.repo.ListAppointments(ctx, tenantID, postgres.AppointmentFilter{
		DateFrom: &startOfDay,
		DateTo:   &endOfDay,
	})
	if err != nil {
		return nil, err
	}

	var confirmed, completed, pending int
	var todayRevenue float64

	for _, a := range todayApts {
		switch a.Status {
		case domain.StatusConfirmed:
			confirmed++
			todayRevenue += a.TotalPrice
		case domain.StatusCompleted:
			completed++
			todayRevenue += a.TotalPrice
		case domain.StatusPending:
			pending++
		}
	}

	var totalCust int64
	s.repo.DB().WithContext(ctx).Model(&domain.Customer{}).Where("tenant_id = ?", tenantID).Count(&totalCust)

	var totalPros int64
	s.repo.DB().WithContext(ctx).Model(&domain.Professional{}).Where("tenant_id = ? AND is_active = ?", tenantID, true).Count(&totalPros)

	// Gráfico de curva Reports (Spline) - Faixas de horários
	hourlyLabels := []string{"08:00", "10:00", "12:00", "14:00", "16:00", "18:00", "20:00"}
	reportsChart := make([]domain.TimeSeriesPoint, len(hourlyLabels))
	for i, lbl := range hourlyLabels {
		reportsChart[i] = domain.TimeSeriesPoint{
			Label:        lbl,
			Appointments: 0,
			Revenue:      0,
		}
	}

	for _, a := range todayApts {
		hour := a.StartAt.Hour()
		idx := 0
		switch {
		case hour < 10:
			idx = 0
		case hour < 12:
			idx = 1
		case hour < 14:
			idx = 2
		case hour < 16:
			idx = 3
		case hour < 18:
			idx = 4
		case hour < 20:
			idx = 5
		default:
			idx = 6
		}
		reportsChart[idx].Appointments++
		if a.Status == domain.StatusConfirmed || a.Status == domain.StatusCompleted {
			reportsChart[idx].Revenue += a.TotalPrice
		}
	}

	// Gráfico Donut de Status / Analytics
	var totalAptsCount int64
	s.repo.DB().WithContext(ctx).Model(&domain.Appointment{}).Where("tenant_id = ?", tenantID).Count(&totalAptsCount)

	type StatusCount struct {
		Status string
		Count  int
	}
	var statusCounts []StatusCount
	s.repo.DB().WithContext(ctx).Model(&domain.Appointment{}).
		Select("status, count(*) as count").
		Where("tenant_id = ?", tenantID).
		Group("status").
		Scan(&statusCounts)

	statusMap := make(map[string]int)
	for _, sc := range statusCounts {
		statusMap[sc.Status] = sc.Count
	}

	baseTotal := float64(totalAptsCount)
	if baseTotal == 0 {
		baseTotal = 1
	}

	statusDist := []domain.StatusDistribution{
		{
			Status:     "CONFIRMED",
			Label:      "Confirmados",
			Count:      statusMap[string(domain.StatusConfirmed)],
			Percentage: float64(statusMap[string(domain.StatusConfirmed)]) / baseTotal * 100,
			Color:      "#4880FF",
		},
		{
			Status:     "COMPLETED",
			Label:      "Concluídos",
			Count:      statusMap[string(domain.StatusCompleted)],
			Percentage: float64(statusMap[string(domain.StatusCompleted)]) / baseTotal * 100,
			Color:      "#FEC53D",
		},
		{
			Status:     "CANCELLED",
			Label:      "Cancelados",
			Count:      statusMap[string(domain.StatusCancelled)],
			Percentage: float64(statusMap[string(domain.StatusCancelled)]) / baseTotal * 100,
			Color:      "#FF6647",
		},
		{
			Status:     "PENDING",
			Label:      "Pendentes",
			Count:      statusMap[string(domain.StatusPending)],
			Percentage: float64(statusMap[string(domain.StatusPending)]) / baseTotal * 100,
			Color:      "#8280FF",
		},
	}

	// Serviços Mais Populares (Top Services)
	type SvcCount struct {
		ServiceID uuid.UUID
		Name      string
		Price     float64
		Total     int
	}
	var svcCounts []SvcCount
	s.repo.DB().WithContext(ctx).Table("appointments").
		Select("services.id as service_id, services.name as name, services.price as price, count(appointments.id) as total").
		Joins("JOIN services ON services.id = appointments.service_id").
		Where("appointments.tenant_id = ?", tenantID).
		Group("services.id, services.name, services.price").
		Order("total desc").
		Limit(5).
		Scan(&svcCounts)

	topServices := make([]domain.TopServiceItem, 0, len(svcCounts))
	for _, sc := range svcCounts {
		topServices = append(topServices, domain.TopServiceItem{
			ServiceID:     sc.ServiceID,
			ServiceName:   sc.Name,
			Price:         sc.Price,
			BookingsCount: sc.Total,
			Rating:        5.0,
		})
	}

	// Últimos agendamentos recentes (limit 6)
	var recentApts []domain.Appointment
	s.repo.DB().WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Preload("Customer").
		Preload("Professional").
		Preload("Service").
		Order("start_at desc").
		Limit(6).
		Find(&recentApts)

	resp := &domain.DashboardAnalyticsResponse{
		ReportsChart:       reportsChart,
		StatusDistribution: statusDist,
		TopServices:        topServices,
		RecentAppointments: recentApts,
	}
	resp.KPIs.TodayAppointmentsCount = len(todayApts)
	resp.KPIs.ConfirmedCount = confirmed
	resp.KPIs.CompletedCount = completed
	resp.KPIs.PendingCount = pending
	resp.KPIs.TodayRevenue = todayRevenue
	resp.KPIs.TotalCustomersCount = totalCust
	resp.KPIs.ActiveProfessionals = totalPros

	return resp, nil
}
