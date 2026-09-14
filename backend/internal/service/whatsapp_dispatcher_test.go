package service

import (
	"testing"
	"time"
)

func TestCalculateTypingDuration(t *testing.T) {
	tests := []struct {
		name        string
		text        string
		speed       int
		minExpected time.Duration
		maxExpected time.Duration
	}{
		{
			name:        "Mensagem curta respeita mínimo de 1.5s",
			text:        "Oi",
			speed:       35,
			minExpected: 1500 * time.Millisecond,
			maxExpected: 1600 * time.Millisecond,
		},
		{
			name:        "Mensagem média",
			text:        "Olá, tudo bem? Temos horários livres às 14h e 16h amanhã.",
			speed:       35,
			minExpected: 1500 * time.Millisecond,
			maxExpected: 3000 * time.Millisecond,
		},
		{
			name:        "Mensagem extremamente longa respeita teto de 7.5s",
			text:        string(make([]rune, 1000)),
			speed:       35,
			minExpected: 7400 * time.Millisecond,
			maxExpected: 7600 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateTypingDuration(tt.text, tt.speed)
			if got < tt.minExpected || got > tt.maxExpected {
				t.Errorf("CalculateTypingDuration() = %v, esperado entre %v e %v", got, tt.minExpected, tt.maxExpected)
			}
		})
	}
}

func TestCalculateInitialDelay(t *testing.T) {
	for i := 0; i < 20; i++ {
		delay := CalculateInitialDelay(2, 5)
		if delay < 2*time.Second || delay > 6*time.Second {
			t.Errorf("CalculateInitialDelay(2, 5) = %v fora do range esperado (2s a 6s)", delay)
		}
	}
}
