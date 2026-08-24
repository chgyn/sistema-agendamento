package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
	"github.com/sistema-agendamento/backend/internal/repository/postgres"
)

type TimeSlot struct {
	StartTime        string    `json:"start_time"`        // "09:00"
	EndTime          string    `json:"end_time"`          // "09:30"
	StartDateTime    time.Time `json:"start_datetime"`
	EndDateTime      time.Time `json:"end_datetime"`
	IsAvailable      bool      `json:"is_available"`
	ProfessionalID   uuid.UUID `json:"professional_id"`
	ProfessionalName string    `json:"professional_name"`
}

type DayAvailability struct {
	Date             string     `json:"date"` // "YYYY-MM-DD"
	DayOfWeek        int        `json:"day_of_week"`
	HasAvailability  bool       `json:"has_availability"`
	Slots            []TimeSlot `json:"slots"`
}

type AvailabilityService struct {
	repo *postgres.Repository
}

func NewAvailabilityService(repo *postgres.Repository) *AvailabilityService {
	return &AvailabilityService{repo: repo}
}

// GetAvailableSlotsForDate calcula os horários livres para um serviço e profissional em uma data específica
func (s *AvailabilityService) GetAvailableSlotsForDate(ctx context.Context, tenantID, serviceID, professionalID uuid.UUID, dateStr string) (*DayAvailability, error) {
	// 1. Validar e carregar o serviço
	svc, err := s.repo.GetServiceByID(ctx, tenantID, serviceID)
	if err != nil {
		return nil, err
	}

	// 2. Parse da data (YYYY-MM-DD)
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil, fmt.Errorf("formato de data inválido, use YYYY-MM-DD: %w", err)
	}

	dayOfWeek := int(date.Weekday()) // 0 = Domingo, 1 = Segunda, ..., 6 = Sábado

	// 3. Se profissional foi especificado, calcula para ele. Se uuid.Nil ("qualquer profissional"), busca todos habilitados
	var targetPros []domain.Professional
	if professionalID != uuid.Nil {
		pro, err := s.repo.GetProfessionalByID(ctx, tenantID, professionalID)
		if err != nil {
			return nil, err
		}
		targetPros = append(targetPros, *pro)
	} else {
		targetPros, err = s.repo.ListProfessionalsByService(ctx, tenantID, serviceID)
		if err != nil {
			return nil, err
		}
		if len(targetPros) == 0 {
			// Fallback: pega todos os profissionais ativos
			targetPros, _ = s.repo.ListProfessionals(ctx, tenantID, true)
		}
	}

	allSlots := make([]TimeSlot, 0)
	slotMap := make(map[string]TimeSlot) // Agrupa para caso de múltiplos profissionais

	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.Local)
	endOfDay := time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.Local)

	now := time.Now()

	for _, pro := range targetPros {
		// Horário de trabalho do dia
		var workingHour *domain.WorkingHour
		whs, err := s.repo.GetWorkingHoursByProfessional(ctx, tenantID, pro.ID)
		if err == nil {
			for _, wh := range whs {
				if wh.DayOfWeek == dayOfWeek && wh.IsActive {
					workingHour = &wh
					break
				}
			}
		}

		if workingHour == nil {
			continue // Não trabalha neste dia da semana
		}

		// Bloqueios e exceções do dia
		exceptions, _ := s.repo.ListExceptions(ctx, tenantID, pro.ID, dateStr, dateStr)
		isFullDayBlocked := false
		for _, exp := range exceptions {
			if exp.IsFullDay {
				isFullDayBlocked = true
				break
			}
		}
		if isFullDayBlocked {
			continue
		}

		// Agendamentos já marcados para o profissional no dia
		appointments, _ := s.repo.ListAppointments(ctx, tenantID, postgres.AppointmentFilter{
			ProfessionalID: &pro.ID,
			DateFrom:       &startOfDay,
			DateTo:         &endOfDay,
		})

		// Gera slots no intervalo da jornada
		workStartMinutes := parseTimeToMinutes(workingHour.StartTime)
		workEndMinutes := parseTimeToMinutes(workingHour.EndTime)
		breakStartMinutes := parseTimeToMinutes(workingHour.BreakStart)
		breakEndMinutes := parseTimeToMinutes(workingHour.BreakEnd)

		stepMinutes := 30 // Intervalo de grade (30 min)
		if svc.DurationMinutes < stepMinutes && svc.DurationMinutes > 0 {
			stepMinutes = svc.DurationMinutes
		}

		for m := workStartMinutes; m+svc.DurationMinutes <= workEndMinutes; m += stepMinutes {
			slotStartMinutes := m
			slotEndMinutes := m + svc.DurationMinutes

			slotStartTimeStr := minutesToTimeStr(slotStartMinutes)
			slotEndTimeStr := minutesToTimeStr(slotEndMinutes)

			slotStartDT := time.Date(date.Year(), date.Month(), date.Day(), slotStartMinutes/60, slotStartMinutes%60, 0, 0, time.Local)
			slotEndDT := slotStartDT.Add(time.Duration(svc.DurationMinutes) * time.Minute)

			// 1. Verifica se está no passado (com margem de 10 min)
			if slotStartDT.Before(now.Add(10 * time.Minute)) {
				continue
			}

			// 2. Verifica se colide com intervalo de almoço/pausa
			if breakStartMinutes > 0 && breakEndMinutes > 0 {
				if slotStartMinutes < breakEndMinutes && slotEndMinutes > breakStartMinutes {
					continue
				}
			}

			// 3. Verifica bloqueios parciais de exceção
			isBlocked := false
			for _, exp := range exceptions {
				if exp.StartTime != "" && exp.EndTime != "" {
					expStart := parseTimeToMinutes(exp.StartTime)
					expEnd := parseTimeToMinutes(exp.EndTime)
					if slotStartMinutes < expEnd && slotEndMinutes > expStart {
						isBlocked = true
						break
					}
				}
			}
			if isBlocked {
				continue
			}

			// 4. Verifica colisão com agendamentos existentes (PENDING ou CONFIRMED)
			hasConflict := false
			for _, apt := range appointments {
				if apt.Status == domain.StatusCancelled || apt.Status == domain.StatusNoShow {
					continue
				}
				if slotStartDT.Before(apt.EndAt) && slotEndDT.After(apt.StartAt) {
					hasConflict = true
					break
				}
			}
			if hasConflict {
				continue
			}

			// Slot livre!
			slot := TimeSlot{
				StartTime:        slotStartTimeStr,
				EndTime:          slotEndTimeStr,
				StartDateTime:    slotStartDT,
				EndDateTime:      slotEndDT,
				IsAvailable:      true,
				ProfessionalID:   pro.ID,
				ProfessionalName: pro.Name,
			}

			// Se ainda não adicionamos esse horário ou se estivermos consolidando
			if existing, exists := slotMap[slotStartTimeStr]; !exists || (professionalID != uuid.Nil) {
				slotMap[slotStartTimeStr] = slot
			} else {
				_ = existing
			}
		}
	}

	for _, slot := range slotMap {
		allSlots = append(allSlots, slot)
	}

	// Ordena cronologicamente
	sort.Slice(allSlots, func(i, j int) bool {
		return allSlots[i].StartDateTime.Before(allSlots[j].StartDateTime)
	})

	return &DayAvailability{
		Date:            dateStr,
		DayOfWeek:       dayOfWeek,
		HasAvailability: len(allSlots) > 0,
		Slots:           allSlots,
	}, nil
}

func parseTimeToMinutes(timeStr string) int {
	if timeStr == "" {
		return 0
	}
	parts := strings.Split(timeStr, ":")
	if len(parts) != 2 {
		return 0
	}
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h*60 + m
}

func minutesToTimeStr(minutes int) string {
	h := minutes / 60
	m := minutes % 60
	return fmt.Sprintf("%02d:%02d", h, m)
}
