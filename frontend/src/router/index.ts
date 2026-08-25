import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore, type UserRole } from '../stores/auth'

import HomeView from '../views/HomeView.vue'
import LoginView from '../views/auth/LoginView.vue'
import RegisterTenantView from '../views/auth/RegisterTenantView.vue'
import PublicBookingView from '../views/public/PublicBookingView.vue'
import AdminLayout from '../components/layout/AdminLayout.vue'
import DashboardView from '../views/admin/DashboardView.vue'
import AppointmentsCalendarView from '../views/admin/AppointmentsCalendarView.vue'
import ServicesView from '../views/admin/ServicesView.vue'
import ProfessionalsView from '../views/admin/ProfessionalsView.vue'
import CustomersView from '../views/admin/CustomersView.vue'
import SettingsView from '../views/admin/SettingsView.vue'
import TenantsManagementView from '../views/admin/TenantsManagementView.vue'
import UsersManagementView from '../views/admin/UsersManagementView.vue'
import GlobalAdminsView from '../views/admin/GlobalAdminsView.vue'
import PlansManagementView from '../views/admin/PlansManagementView.vue'
import SubscriptionBillingView from '../views/admin/SubscriptionBillingView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
    },
    {
      path: '/register',
      name: 'register',
      component: RegisterTenantView,
    },
    {
      path: '/agendamento/:slug',
      name: 'public-booking',
      component: PublicBookingView,
    },
    {
      path: '/admin',
      component: AdminLayout,
      meta: { requiresAuth: true },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: DashboardView,
        },
        {
          path: 'estabelecimentos',
          name: 'admin-tenants',
          component: TenantsManagementView,
          meta: { roles: ['ADMIN_GLOBAL'] },
        },
        {
          path: 'planos',
          name: 'admin-plans',
          component: PlansManagementView,
          meta: { roles: ['ADMIN_GLOBAL'] },
        },
        {
          path: 'minha-assinatura',
          name: 'admin-subscription-billing',
          component: SubscriptionBillingView,
          meta: { roles: ['ADMIN_GLOBAL', 'ADMIN_TENANT', 'ADMIN'] },
        },
        {
          path: 'administradores-globais',
          name: 'admin-global-admins',
          component: GlobalAdminsView,
          meta: { roles: ['ADMIN_GLOBAL'] },
        },
        {
          path: 'usuarios',
          name: 'admin-users',
          component: UsersManagementView,
          meta: { roles: ['ADMIN_GLOBAL', 'ADMIN_TENANT', 'ADMIN'] },
        },
        {
          path: 'agenda',
          name: 'admin-agenda',
          component: AppointmentsCalendarView,
        },
        {
          path: 'servicos',
          name: 'admin-services',
          component: ServicesView,
          meta: { roles: ['ADMIN_GLOBAL', 'ADMIN_TENANT', 'ADMIN'] },
        },
        {
          path: 'profissionais',
          name: 'admin-professionals',
          component: ProfessionalsView,
          meta: { roles: ['ADMIN_GLOBAL', 'ADMIN_TENANT', 'ADMIN'] },
        },
        {
          path: 'clientes',
          name: 'admin-customers',
          component: CustomersView,
        },
        {
          path: 'configuracoes',
          name: 'admin-settings',
          component: SettingsView,
          meta: { roles: ['ADMIN_GLOBAL', 'ADMIN_TENANT', 'ADMIN'] },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

// Navigation Guard de Autenticação e Controle de Acesso por Perfil
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.matched.some(record => record.meta.requiresAuth) && !authStore.isAuthenticated) {
    next({ name: 'login' })
    return
  }

  if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    next({ name: 'admin-dashboard' })
    return
  }

  // Verificação de permissões por Role
  const allowedRoles = to.meta.roles as UserRole[] | undefined
  if (allowedRoles && allowedRoles.length > 0 && authStore.user) {
    let currentRole = authStore.user.role
    if (currentRole === 'ADMIN') {
      currentRole = 'ADMIN_TENANT'
    }

    const hasPermission = allowedRoles.some(r => {
      if (r === 'ADMIN') return currentRole === 'ADMIN_TENANT'
      return r === currentRole
    })

    if (!hasPermission) {
      console.warn(`Acesso negado à rota ${to.path} para o perfil ${currentRole}`)
      next({ name: 'admin-dashboard' })
      return
    }
  }

  next()
})

export default router
