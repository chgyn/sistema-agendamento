import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

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
          path: 'agenda',
          name: 'admin-agenda',
          component: AppointmentsCalendarView,
        },
        {
          path: 'servicos',
          name: 'admin-services',
          component: ServicesView,
        },
        {
          path: 'profissionais',
          name: 'admin-professionals',
          component: ProfessionalsView,
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
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      redirect: '/',
    },
  ],
})

// Navigation Guard de Autenticação
router.beforeEach((to, _from, next) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    next({ name: 'login' })
  } else if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    next({ name: 'admin-dashboard' })
  } else {
    next()
  }
})

export default router
