import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '../services/api'

export interface User {
  id: string
  tenant_id: string
  name: string
  email: string
  role: 'ADMIN' | 'OPERATOR'
  is_active: boolean
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
}

export const useAuthStore = defineStore('auth', () => {
  const token = ref<string | null>(localStorage.getItem('token'))
  const user = ref<User | null>(
    localStorage.getItem('user') ? JSON.parse(localStorage.getItem('user')!) : null
  )
  const tenant = ref<Tenant | null>(
    localStorage.getItem('tenant') ? JSON.parse(localStorage.getItem('tenant')!) : null
  )

  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const isAdmin = computed(() => user.value?.role === 'ADMIN')

  function setAuthData(authToken: string, authUser: User, authTenant: Tenant) {
    token.value = authToken
    user.value = authUser
    tenant.value = authTenant
    localStorage.setItem('token', authToken)
    localStorage.setItem('user', JSON.stringify(authUser))
    localStorage.setItem('tenant', JSON.stringify(authTenant))
  }

  function updateTenant(updatedTenant: Tenant) {
    tenant.value = updatedTenant
    localStorage.setItem('tenant', JSON.stringify(updatedTenant))
  }

  async function login(email: string, password: string) {
    const res = await api.post('/auth/login', { email, password })
    if (res.data.success) {
      const { token: t, user: u, tenant: ten } = res.data.data
      setAuthData(t, u, ten)
      return res.data.data
    }
    throw new Error(res.data.error || 'Erro no login')
  }

  async function registerTenant(data: {
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
      const { token: t, user: u, tenant: ten } = res.data.data
      setAuthData(t, u, ten)
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
        localStorage.setItem('user', JSON.stringify(user.value))
        localStorage.setItem('tenant', JSON.stringify(tenant.value))
      }
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = null
    user.value = null
    tenant.value = null
    localStorage.removeItem('token')
    localStorage.removeItem('user')
    localStorage.removeItem('tenant')
  }

  return {
    token,
    user,
    tenant,
    isAuthenticated,
    isAdmin,
    login,
    registerTenant,
    fetchMe,
    updateTenant,
    logout,
  }
})
