package domain

import (
	"context"

	"github.com/google/uuid"
)

// -------------------------------------------------------------
// DTOs (REQUEST / RESPONSE)
// -------------------------------------------------------------

type CreatePlanDTO struct {
	Name             string           `json:"name" binding:"required,min=3,max=150"`
	Description      string           `json:"description"`
	Price            float64          `json:"price" binding:"required,gt=0"`
	BillingCycle     PlanBillingCycle `json:"billing_cycle" binding:"required,oneof=MONTHLY QUARTERLY SEMIANNUALLY YEARLY"`
	MaxProfessionals int              `json:"max_professionals" binding:"gte=0"`
	MaxServices      int              `json:"max_services" binding:"gte=0"`
	Features         string           `json:"features"`
	SortOrder        int              `json:"sort_order"`
	IsActive         *bool            `json:"is_active"`
}

type UpdatePlanDTO struct {
	Name             string           `json:"name" binding:"required,min=3,max=150"`
	Description      string           `json:"description"`
	Price            float64          `json:"price" binding:"required,gt=0"`
	BillingCycle     PlanBillingCycle `json:"billing_cycle" binding:"required,oneof=MONTHLY QUARTERLY SEMIANNUALLY YEARLY"`
	MaxProfessionals int              `json:"max_professionals" binding:"gte=0"`
	MaxServices      int              `json:"max_services" binding:"gte=0"`
	Features         string           `json:"features"`
	SortOrder        int              `json:"sort_order"`
	IsActive         *bool            `json:"is_active"`
}

type TogglePlanStatusDTO struct {
	IsActive bool `json:"is_active"`
}

type RegisterTenantWithPlanDTO struct {
	PlanID     uuid.UUID `json:"plan_id" binding:"required"`
	TenantName string    `json:"tenant_name" binding:"required"`
	Slug       string    `json:"slug" binding:"required"`
	Document   string    `json:"document"`
	Phone      string    `json:"phone" binding:"required"`
	City       string    `json:"city"`
	State      string    `json:"state"`
	AdminName  string    `json:"admin_name" binding:"required"`
	AdminEmail string    `json:"admin_email" binding:"required,email"`
	Password   string    `json:"password" binding:"required,min=6"`
}

type OverrideSubscriptionStatusDTO struct {
	Status SubscriptionStatus `json:"status" binding:"required,oneof=ACTIVE PENDING OVERDUE CANCELLED EXPIRED TRIAL"`
}

// -------------------------------------------------------------
// REPOSITORY INTERFACES (PORTS)
// -------------------------------------------------------------

type PlanRepository interface {
	GetPlanByID(ctx context.Context, id uuid.UUID) (*Plan, error)
	ListPlans(ctx context.Context, search string, onlyActive *bool) ([]Plan, error)
	CreatePlan(ctx context.Context, plan *Plan) error
	UpdatePlan(ctx context.Context, plan *Plan) error
	UpdatePlanStatus(ctx context.Context, planID uuid.UUID, isActive bool) error
	DeletePlan(ctx context.Context, planID uuid.UUID) error
}

type SubscriptionRepository interface {
	GetSubscriptionByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	GetSubscriptionByTenantID(ctx context.Context, tenantID uuid.UUID) (*Subscription, error)
	GetSubscriptionByAsaasID(ctx context.Context, asaasSubID string) (*Subscription, error)
	CreateSubscription(ctx context.Context, sub *Subscription) error
	UpdateSubscription(ctx context.Context, sub *Subscription) error
	UpdateSubscriptionStatus(ctx context.Context, subID uuid.UUID, status SubscriptionStatus) error
	ListAllSubscriptions(ctx context.Context, status string, search string) ([]Subscription, error)
	SaveSubscriptionInvoice(ctx context.Context, invoice *SubscriptionInvoice) error
	GetSubscriptionInvoices(ctx context.Context, subscriptionID uuid.UUID) ([]SubscriptionInvoice, error)
}

// -------------------------------------------------------------
// ASAAS INTEGRATION INTERFACE (PORT)
// -------------------------------------------------------------

type AsaasCustomerRequest struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Phone                string `json:"phone,omitempty"`
	MobilePhone          string `json:"mobilePhone,omitempty"`
	CpfCnpj              string `json:"cpfCnpj,omitempty"`
	ExternalReference    string `json:"externalReference,omitempty"`
	NotificationDisabled bool   `json:"notificationDisabled,omitempty"`
}

type AsaasCustomerResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	CpfCnpj   string `json:"cpfCnpj"`
	Deleted   bool   `json:"deleted"`
}

type AsaasSubscriptionRequest struct {
	Customer          string   `json:"customer"`
	BillingType       string   `json:"billingType"` // UNDEFINED, BOLETO, CREDIT_CARD, PIX
	Value             float64  `json:"value"`
	NextDueDate       string   `json:"nextDueDate"` // YYYY-MM-DD
	Cycle             string   `json:"cycle"`       // MONTHLY, QUARTERLY, SEMIANNUALLY, YEARLY
	Description       string   `json:"description,omitempty"`
	ExternalReference string   `json:"externalReference,omitempty"`
}

type AsaasSubscriptionResponse struct {
	ID          string  `json:"id"`
	Customer    string  `json:"customer"`
	Value       float64 `json:"value"`
	NextDueDate string  `json:"nextDueDate"`
	Cycle       string  `json:"cycle"`
	Status      string  `json:"status"` // ACTIVE, INACTIVE, EXPIRED
	BillingType string  `json:"billingType"`
	Description string  `json:"description"`
	BankSlipURL string  `json:"bankSlipUrl"`
	InvoiceURL  string  `json:"paymentLink"`
}

type AsaasPaymentResponse struct {
	ID             string    `json:"id"`
	Customer       string    `json:"customer"`
	Subscription   string    `json:"subscription"`
	Value          float64   `json:"value"`
	NetValue       float64   `json:"netValue"`
	BillingType    string    `json:"billingType"`
	Status         string    `json:"status"` // PENDING, RECEIVED, CONFIRMED, OVERDUE, REFUNDED
	DueDate        string    `json:"dueDate"`
	PaymentDate    *string   `json:"paymentDate"`
	InvoiceURL     string    `json:"invoiceUrl"`
	BankSlipURL    string    `json:"bankSlipUrl"`
	InvoiceNumber  string    `json:"invoiceNumber"`
}

type AsaasPaymentListResponse struct {
	Data       []AsaasPaymentResponse `json:"data"`
	HasMore    bool                   `json:"hasMore"`
	TotalCount int                    `json:"totalCount"`
}

type AsaasProvider interface {
	CreateCustomer(ctx context.Context, req AsaasCustomerRequest) (*AsaasCustomerResponse, error)
	CreateSubscription(ctx context.Context, req AsaasSubscriptionRequest) (*AsaasSubscriptionResponse, error)
	GetSubscription(ctx context.Context, asaasSubID string) (*AsaasSubscriptionResponse, error)
	CancelSubscription(ctx context.Context, asaasSubID string) error
	GetSubscriptionPayments(ctx context.Context, asaasSubID string) ([]AsaasPaymentResponse, error)
}
