package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/sistema-agendamento/backend/internal/domain"
)

var (
	ErrInvalidToken = errors.New("token inválido ou expirado")
)

type CustomClaims struct {
	UserID   uuid.UUID   `json:"user_id"`
	TenantID uuid.UUID   `json:"tenant_id"`
	Email    string      `json:"email"`
	Name     string      `json:"name"`
	Role     domain.Role `json:"role"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secretKey     []byte
	expireDuration time.Duration
}

func NewJWTService(secretKey string, expireHours int) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		expireDuration: time.Duration(expireHours) * time.Hour,
	}
}

// GenerateToken gera um token JWT com tenant_id, user_id e role
func (s *JWTService) GenerateToken(user *domain.User) (string, error) {
	now := time.Now()
	claims := CustomClaims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Name:     user.Name,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.expireDuration)),
			Issuer:    "sistema-agendamento-api",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken valida o token JWT e retorna os claims customizados
func (s *JWTService) ValidateToken(tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return s.secretKey, nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*CustomClaims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}
