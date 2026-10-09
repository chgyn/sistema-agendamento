<template>
  <div class="min-h-screen bg-[#F5F6FA] dark:bg-[#0d0e12] text-[#202224] dark:text-zinc-100 flex flex-col md:flex-row transition-colors duration-200">
    <!-- Topbar Mobile -->
    <header class="md:hidden flex items-center justify-between p-4 bg-white/95 dark:bg-[#121318]/95 border-b border-gray-100 dark:border-zinc-800 sticky top-0 z-40 backdrop-blur-md">
      <div class="flex items-center space-x-3">
        <div class="w-9 h-9 rounded-xl bg-[#4880FF] text-white flex items-center justify-center shadow-md shadow-blue-500/20 shrink-0 font-display font-black text-lg">
          B
        </div>
        <div class="min-w-0">
          <h2 class="font-bold text-sm text-[#202224] dark:text-white truncate leading-tight">
            {{ authStore.isGlobalAdmin ? 'Plataforma Global' : (authStore.tenant?.name || 'Base Dashboard') }}
          </h2>
          <span class="text-[11px] text-[#718096] dark:text-zinc-400 font-mono truncate block">
            {{ authStore.isGlobalAdmin ? 'Super Admin' : (authStore.tenant?.slug ? `/${authStore.tenant.slug}` : 'Painel de Gestão') }}
          </span>
        </div>
      </div>

      <div class="flex items-center gap-2">
        <!-- Toggle Tema Mobile -->
        <button
          @click="toggleTheme"
          class="p-2 rounded-xl bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#4880FF] dark:hover:text-white transition"
          aria-label="Alternar Tema"
        >
          <Sun v-if="isDark" class="w-4 h-4 text-amber-400" />
          <Moon v-else class="w-4 h-4 text-[#718096]" />
        </button>

        <button
          @click="mobileMenuOpen = !mobileMenuOpen"
          class="p-2 rounded-xl bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#4880FF] dark:hover:text-white transition"
          aria-label="Menu"
        >
          <Menu v-if="!mobileMenuOpen" class="w-5 h-5" />
          <X v-else class="w-5 h-5 text-[#4880FF]" />
        </button>
      </div>
    </header>

    <!-- Sidebar Desktop & Drawer Mobile -->
    <aside
      :class="[
        'bg-white dark:bg-[#121318] border-r border-gray-100 dark:border-zinc-800/80 flex flex-col justify-between shrink-0 transition-all duration-300 z-30',
        isCollapsed ? 'md:w-20' : 'md:w-64',
        mobileMenuOpen ? 'fixed inset-y-0 left-0 w-72 shadow-2xl p-4 flex' : 'hidden md:flex p-4'
      ]"
    >
      <div class="space-y-6 flex flex-col h-full overflow-y-auto overflow-x-hidden">
        <!-- Topo da Sidebar: Logo & Botão de Recolher -->
        <div class="flex items-center justify-between px-2 pt-1">
          <div class="flex items-center gap-3 min-w-0">
            <!-- Logo Icon estilo Figma Base -->
            <div class="w-10 h-10 rounded-2xl bg-gradient-to-tr from-[#386FF0] to-[#5D8EFF] flex items-center justify-center text-white shadow-md shadow-blue-500/25 shrink-0">
              <svg class="w-6 h-6 text-white" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
              </svg>
            </div>
            <div v-if="!isCollapsed" class="min-w-0 flex-1 transition-opacity duration-200">
              <h2 class="font-black text-lg text-[#202224] dark:text-white font-display tracking-tight truncate leading-tight">
                {{ authStore.isGlobalAdmin ? 'Base Admin' : (authStore.tenant?.name || 'Base') }}
              </h2>
              <span class="text-[11px] text-[#718096] dark:text-zinc-400 font-medium truncate block">
                {{ authStore.isGlobalAdmin ? 'Gestão da Plataforma' : (authStore.tenant?.slug ? `/${authStore.tenant.slug}` : 'SaaS Agendamento') }}
              </span>
            </div>
          </div>

          <!-- Botão Desktop para Recolher/Expandir Sidebar -->
          <button
            v-if="!mobileMenuOpen"
            @click="toggleCollapse"
            class="hidden md:flex p-1.5 rounded-lg text-gray-400 hover:text-[#4880FF] hover:bg-blue-50 dark:hover:bg-zinc-800 transition"
            :title="isCollapsed ? 'Expandir Menu' : 'Recolher Menu'"
          >
            <ChevronRight v-if="isCollapsed" class="w-4 h-4" />
            <ChevronLeft v-else class="w-4 h-4" />
          </button>
        </div>

        <!-- Links de Navegação -->
        <nav class="space-y-1.5 flex-1 pt-2">
          <!-- 1. MENU ADMIN GLOBAL -->
          <template v-if="authStore.isGlobalAdmin">
            <div v-if="!isCollapsed" class="px-3 pb-1 pt-2 text-[10px] font-bold text-[#718096] dark:text-zinc-500 uppercase tracking-wider">
              Plataforma
            </div>

            <RouterLink
              to="/admin/dashboard"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/dashboard')"
              :title="isCollapsed ? 'Dashboard Global' : ''"
            >
              <LayoutDashboard class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Dashboard Global</span>
            </RouterLink>

            <RouterLink
              to="/admin/estabelecimentos"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/estabelecimentos')"
              :title="isCollapsed ? 'Estabelecimentos' : ''"
            >
              <Building2 class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Estabelecimentos</span>
            </RouterLink>

            <RouterLink
              to="/admin/planos"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/planos')"
              :title="isCollapsed ? 'Planos de Assinatura' : ''"
            >
              <CreditCard class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Planos de Assinatura</span>
            </RouterLink>

            <RouterLink
              to="/admin/usuarios"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/usuarios')"
              :title="isCollapsed ? 'Usuários' : ''"
            >
              <Users class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Usuários</span>
            </RouterLink>

            <RouterLink
              to="/admin/administradores-globais"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/administradores-globais')"
              :title="isCollapsed ? 'Admins Gerais' : ''"
            >
              <ShieldCheck class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Admins Gerais</span>
            </RouterLink>
          </template>

          <!-- 2. MENU TENANT ADMIN OU OPERADOR -->
          <template v-else>
            <div v-if="!isCollapsed" class="px-3 pb-1 pt-2 text-[10px] font-bold text-[#718096] dark:text-zinc-500 uppercase tracking-wider">
              Menu Principal
            </div>

            <RouterLink
              to="/admin/dashboard"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/dashboard')"
              :title="isCollapsed ? 'Dashboard' : ''"
            >
              <LayoutDashboard class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Dashboard</span>
            </RouterLink>

            <RouterLink
              to="/admin/agenda"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/agenda')"
              :title="isCollapsed ? 'Agenda & Atendimentos' : ''"
            >
              <Calendar class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Schedule List</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/servicos"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/servicos')"
              :title="isCollapsed ? 'Serviços' : ''"
            >
              <Sparkles class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Serviços</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/profissionais"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/profissionais')"
              :title="isCollapsed ? 'Profissionais' : ''"
            >
              <Users class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Profissionais</span>
            </RouterLink>

            <RouterLink
              to="/admin/clientes"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/clientes')"
              :title="isCollapsed ? 'Clientes' : ''"
            >
              <UserSquare2 class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Customer List</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/whatsapp"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/whatsapp')"
              :title="isCollapsed ? 'WhatsApp & IA' : ''"
            >
              <Bot class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">WhatsApp & IA</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/minha-assinatura"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/minha-assinatura')"
              :title="isCollapsed ? 'Minha Assinatura' : ''"
            >
              <CreditCard class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Minha Assinatura</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/usuarios"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/usuarios')"
              :title="isCollapsed ? 'Gestão de Usuários' : ''"
            >
              <UsersRound class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Gestão de Equipe</span>
            </RouterLink>

            <RouterLink
              v-if="!authStore.isOperator"
              to="/admin/configuracoes"
              @click="mobileMenuOpen = false"
              :class="navItemClass('/admin/configuracoes')"
              :title="isCollapsed ? 'Configurações' : ''"
            >
              <Settings class="w-5 h-5 shrink-0" />
              <span v-if="!isCollapsed" class="truncate">Configurações</span>
            </RouterLink>
          </template>
        </nav>

        <!-- Widget Inferior Promo / Upgrade Card (Padrão Figma - Pirâmide / Upgrade Now) -->
        <div
          v-if="!isCollapsed && !authStore.isGlobalAdmin"
          class="mx-1 my-2 p-4 rounded-2xl bg-gradient-to-b from-[#EBF2FE] to-[#F5F8FF] dark:from-[#1b2234] dark:to-[#121622] border border-blue-100/80 dark:border-blue-900/40 text-center relative overflow-hidden"
        >
          <!-- Ícone decorativo estilo pirâmide suave -->
          <div class="w-12 h-12 mx-auto mb-2 rounded-2xl bg-[#4880FF]/15 text-[#4880FF] flex items-center justify-center">
            <Sparkles class="w-6 h-6 animate-pulse" />
          </div>
          <h4 class="font-black text-xs text-[#202224] dark:text-white font-display">
            {{ authStore.tenant?.subscription?.plan?.name || 'Base Pro' }}
          </h4>
          <p class="text-[11px] text-[#718096] dark:text-zinc-400 mt-0.5 mb-3 leading-tight">
            {{ authStore.tenant?.subscription?.status === 'ACTIVE' ? 'Plano Ativo e Regular' : 'Acesse para gerenciar recursos' }}
          </p>
          <RouterLink
            to="/admin/minha-assinatura"
            @click="mobileMenuOpen = false"
            class="inline-block w-full py-2 px-3 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25"
          >
            Upgrade Now
          </RouterLink>
        </div>

        <!-- Usuário Logado & Logout (Padrão Figma: Avatar + Nome + Logout rápido) -->
        <div class="pt-3 border-t border-gray-100 dark:border-zinc-800/80">
          <div class="flex items-center justify-between px-1">
            <div class="flex items-center space-x-3 min-w-0 flex-1">
              <div class="w-9 h-9 rounded-full bg-gradient-to-tr from-[#4880FF] to-[#8280FF] text-white font-black text-xs flex items-center justify-center shrink-0 uppercase shadow-sm">
                {{ userInitials }}
              </div>
              <div v-if="!isCollapsed" class="min-w-0 flex-1">
                <p class="text-xs font-bold text-[#202224] dark:text-white truncate font-display">{{ authStore.user?.name || 'Administrador' }}</p>
                <p class="text-[11px] text-[#718096] dark:text-zinc-400 truncate">{{ roleLabel }}</p>
              </div>
            </div>

            <button
              @click="handleLogout"
              class="p-2 rounded-xl text-gray-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-500/10 transition shrink-0"
              title="Sair da Conta"
            >
              <LogOut class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </aside>

    <!-- Overlay escuro para mobile menu -->
    <div
      v-if="mobileMenuOpen"
      @click="mobileMenuOpen = false"
      class="fixed inset-0 bg-black/40 z-20 md:hidden backdrop-blur-sm"
    ></div>

    <!-- Área de Conteúdo Principal -->
    <div class="flex-1 flex flex-col min-w-0 max-h-screen overflow-hidden">
      <!-- Topbar Desktop (Padrão Figma: Título da Página + Filtro de Período + Toggle Tema) -->
      <header class="hidden md:flex items-center justify-between px-8 py-5 bg-white dark:bg-[#121318] border-b border-gray-100 dark:border-zinc-800/80 z-10 shrink-0">
        <div>
          <h1 class="text-xl sm:text-2xl font-black text-[#202224] dark:text-white font-display tracking-tight">
            {{ pageTitle }}
          </h1>
        </div>

        <div class="flex items-center space-x-4">
          <!-- Botão Ver Página Pública (se Tenant) -->
          <a
            v-if="!authStore.isGlobalAdmin && authStore.tenant?.slug"
            :href="`/agendamento/${authStore.tenant.slug}`"
            target="_blank"
            class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl border border-blue-200 dark:border-blue-900/60 bg-blue-50/60 dark:bg-blue-950/40 text-[#4880FF] text-xs font-bold hover:bg-blue-100/80 transition"
          >
            <ExternalLink class="w-3.5 h-3.5" />
            <span>Página Pública</span>
          </a>

          <!-- Pill de Data (Padrão Figma: "10-06-2026 ▾") -->
          <div class="flex items-center gap-2 px-3 py-1.5 rounded-xl bg-gray-50 dark:bg-zinc-800/80 border border-gray-200/80 dark:border-zinc-700/60 text-xs font-semibold text-[#202224] dark:text-zinc-200">
            <Calendar class="w-3.5 h-3.5 text-[#4880FF]" />
            <span>{{ currentDateFormatted }}</span>
          </div>

          <!-- Toggle de Tema (Light / Dark) -->
          <button
            @click="toggleTheme"
            class="p-2 rounded-xl border border-gray-200/80 dark:border-zinc-700/60 bg-gray-50 dark:bg-zinc-800/80 text-[#718096] dark:text-zinc-300 hover:text-[#4880FF] dark:hover:text-amber-400 transition cursor-pointer"
            :title="isDark ? 'Ativar Modo Claro' : 'Ativar Modo Escuro'"
          >
            <Sun v-if="isDark" class="w-4 h-4 text-amber-400" />
            <Moon v-else class="w-4 h-4 text-[#718096]" />
          </button>
        </div>
      </header>

      <!-- Conteúdo da Rota -->
      <main class="flex-1 overflow-y-auto p-4 sm:p-6 lg:p-8">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  LayoutDashboard, Calendar, Sparkles, Users, UsersRound,
  UserSquare2, Settings, ExternalLink, LogOut, Building2, ShieldCheck, CreditCard, Bot,
  Menu, X, ChevronLeft, ChevronRight, Sun, Moon
} from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

const mobileMenuOpen = ref(false)
const isCollapsed = ref(false)
const isDark = ref(false)

onMounted(() => {
  // Carrega estado de colapso da sidebar
  const savedCollapsed = localStorage.getItem('sidebar_collapsed')
  if (savedCollapsed === 'true') {
    isCollapsed.value = true
  }

  // Carrega tema
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    isDark.value = true
    document.documentElement.classList.add('dark')
  } else {
    isDark.value = false
    document.documentElement.classList.remove('dark')
  }
})

function toggleCollapse() {
  isCollapsed.value = !isCollapsed.value
  localStorage.setItem('sidebar_collapsed', isCollapsed.value ? 'true' : 'false')
}

function toggleTheme() {
  isDark.value = !isDark.value
  if (isDark.value) {
    document.documentElement.classList.add('dark')
    localStorage.setItem('theme', 'dark')
  } else {
    document.documentElement.classList.remove('dark')
    localStorage.setItem('theme', 'light')
  }
}

function isActive(path: string) {
  return route.path === path || route.path.startsWith(path + '/')
}

function navItemClass(path: string) {
  const active = isActive(path)
  const base = 'flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition cursor-pointer'
  if (active) {
    return `${base} bg-[#E9F0FE] text-[#4880FF] font-bold dark:bg-[#4880FF]/15 dark:text-[#4880FF]`
  }
  return `${base} text-[#718096] hover:bg-gray-50 hover:text-[#202224] dark:text-zinc-400 dark:hover:bg-zinc-800/60 dark:hover:text-zinc-100`
}

const userInitials = computed(() => {
  const name = authStore.user?.name || 'Admin'
  const parts = name.split(' ')
  if (parts.length > 1) {
    return (parts[0][0] + parts[1][0]).toUpperCase()
  }
  return name.slice(0, 2).toUpperCase()
})

const roleLabel = computed(() => {
  if (authStore.isGlobalAdmin) return 'Super Admin'
  if (authStore.isTenantAdmin) return 'Admin da Conta'
  if (authStore.isOperator) return 'Operador'
  return authStore.user?.role || 'Usuário'
})

const pageTitle = computed(() => {
  switch (route.name) {
    case 'admin-dashboard':
      return authStore.isGlobalAdmin ? 'Dashboard Global' : 'Dashboard'
    case 'admin-agenda':
      return 'Agenda & Atendimentos'
    case 'admin-servicos':
      return 'Gestão de Serviços'
    case 'admin-profissionais':
      return 'Profissionais & Horários'
    case 'admin-clientes':
      return 'Base de Clientes'
    case 'admin-whatsapp':
      return 'WhatsApp & Automação IA'
    case 'admin-subscription-billing':
      return 'Minha Assinatura'
    case 'admin-users':
      return 'Gestão de Usuários'
    case 'admin-configuracoes':
      return 'Configurações'
    case 'admin-tenants':
      return 'Gestão de Estabelecimentos'
    case 'admin-plans':
      return 'Planos de Assinatura'
    case 'admin-global-admins':
      return 'Administradores Gerais'
    default:
      return 'Painel de Gestão'
  }
})

const currentDateFormatted = computed(() => {
  return new Intl.DateTimeFormat('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric'
  }).format(new Date())
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
