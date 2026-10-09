import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export type UserRole = 'ADMIN_GLOBAL' | 'ADMIN_TENANT' | 'OPERATOR' | 'ADMIN'
export type PlanBillingCycle = 'MONTHLY' | 'QUARTERLY' | 'SEMIANNUALLY' | 'YEARLY'
export type SubscriptionStatus = 'ACTIVE' | 'PENDING' | 'OVERDUE' | 'CANCELLED' | 'EXPIRED' | 'TRIAL'

export type SubscriptionOrigin = 'ASAAS' | 'MANUAL' | 'FREE_PLAN'

export interface Plan {
  id: string
  name: string
  description?: string
  price: number
  billing_cycle: PlanBillingCycle
  max_professionals: number
  max_services: number
  features?: string
  is_active: boolean
  is_free?: boolean
  sort_order: number
  created_at?: string
}

export interface SubscriptionInvoice {
  id: string
  tenant_id: string
  subscription_id: string
  asaas_payment_id: string
  status: string
  value: number
  net_value?: number
  billing_type?: string
  due_date: string
  payment_date?: string
  invoice_url: string
  bank_slip_url?: string
  pix_qr_code_url?: string
  created_at?: string
}

export interface SubscriptionAuditLog {
  id: string
  subscription_id: string
  tenant_id: string
  plan_id: string
  action: string
  previous_status?: SubscriptionStatus
  new_status: SubscriptionStatus
  performed_by_id?: string
  reason: string
  expires_at?: string
  created_at: string
  performed_by?: User
  plan?: Plan
}

export interface Subscription {
  id: string
  tenant_id: string
  plan_id: string
  asaas_customer_id?: string
  asaas_subscription_id?: string
  status: SubscriptionStatus
  billing_cycle: PlanBillingCycle
  price: number
  origin?: SubscriptionOrigin
  next_due_date?: string
  current_period_end?: string
  payment_method: string
  payment_url?: string
  manual_grant_reason?: string
  granted_by_user_id?: string
  granted_by_user?: User
  created_at?: string
  plan?: Plan
  invoices?: SubscriptionInvoice[]
  audit_logs?: SubscriptionAuditLog[]
}

export interface Tenant {
  id: string
  slug: string
  name: string
  document?: string
  phone?: string
  email?: string
  address?: string
  city?: string
  state?: string
  logo_url?: string
  primary_color: string
  slot_interval_minutes: number
  is_active?: boolean
  created_at?: string
  subscription?: Subscription
}

export interface User {
  id: string
  tenant_id?: string | null
  name: string
  email: string
  role: UserRole
  is_active: boolean
  created_at?: string
  tenant?: Tenant
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(
    localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null
  )
  const tenant = ref<Tenant | null>(
    localStorage.getItem('tenant') ? JSON.parse(localStorage.getItem('tenant')!) : null
  )
  const subscription = ref<Subscription | null>(
    localStorage.getItem('subscription') ? JSON.parse(localStorage.getItem('subscription')!) : null
  )

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isGlobalAdmin = computed(() => user.value?.role === 'ADMIN_GLOBAL')
  const isTenantAdmin = computed(() => user.value?.role === 'ADMIN_TENANT' || user.value?.role === 'ADMIN')
  const isOperator = computed(() => user.value?.role === 'OPERATOR')
  const canManageUsers = computed(() => isGlobalAdmin.value || isTenantAdmin.value)
  const canManageTenants = computed(() => isGlobalAdmin.value)

  // Validação do status da assinatura
  const isSubscriptionActive = computed(() => {
    if (isGlobalAdmin.value) return true
    if (!subscription.value) return false
    return subscription.value.status === 'ACTIVE' || subscription.value.status === 'TRIAL'
  })

  // Compatibilidade legada
  const isAdmin = computed(() => isGlobalAdmin.value || isTenantAdmin.value)

  function setAuthData(authToken: string, authUser: User, authTenant: Tenant | null, authSub?: Subscription | null) {
    token.value = authToken
    user.value = authUser
    tenant.value = authTenant
    subscription.value = authSub || authTenant?.subscription || null

    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(authUser))
    if (authTenant) {
      localStorage.setItem('tenant', JSON.stringify(authTenant))
    } else {
      localStorage.removeItem('tenant')
    }
    if (subscription.value) {
      localStorage.setItem('subscription', JSON.stringify(subscription.value))
    } else {
      localStorage.removeItem('subscription')
    }
  }

  function updateTenant(updatedTenant: Tenant) {
    tenant.value = updatedTenant
    localStorage.setItem('tenant', JSON.stringify(updatedTenant))
  }

  function updateSubscription(updatedSub: Subscription) {
    subscription.value = updatedSub
    localStorage.setItem('subscription', JSON.stringify(updatedSub))
  }

  async function login(email: string, password: string) {
    const res = await api.post('/auth/login', { email, password })
    if (res.data.success) {
      const { token: t, user: u, tenant: ten, subscription: sub } = res.data.data
      setAuthData(t, u, ten, sub)
      return res.data.data
    }
    throw new Error(res.data.error || 'Erro no login')
  }

  async function registerTenant(data: {
    plan_id: string
    tenant_name: string
    slug: string
    document?: string
    phone: string
    city?: string
    state?: string
    admin_name: string
    admin_email: string
    password: string
  }) {
    const res = await api.post('/auth/register', data)
    if (res.data.success) {
      const { token: t, user: u, tenant: ten, subscription: sub } = res.data.data
      setAuthData(t, u, ten, sub)
      return res.data.data
    }
    throw new Error(res.data.error || 'Erro no registro')
  }

  async function fetchMe() {
    try {
      const res = await api.get('/auth/me')
      if (res.data.success) {
        user.value = res.data.data.user
        tenant.value = res.data.data.tenant
        subscription.value = res.data.data.subscription || res.data.data.tenant?.subscription || null

        localStorage.setItem('user', JSON.stringify(user.value))
        if (tenant.value) localStorage.setItem('tenant', JSON.stringify(tenant.value))
        if (subscription.value) localStorage.setItem('subscription', JSON.stringify(subscription.value))
      }
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    tenant.value = null
    subscription.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('tenant')
    localStorage.removeItem('subscription')
  }

  return {
    token,
    user,
    tenant,
    subscription,
    isAuthenticated,
    isAdmin,
    isGlobalAdmin,
    isTenantAdmin,
    isOperator,
    canManageUsers,
    canManageTenants,
    isSubscriptionActive,
    login,
    registerTenant,
    fetchMe,
    updateTenant,
    updateSubscription,
    logout,
  }
})
