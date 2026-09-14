<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col md:flex-row">
    <!-- Sidebar Desktop -->
    <aside class="w-full md:w-64 bg-slate-900/90 border-r border-slate-800/80 flex flex-col justify-between p-4 shrink-0">
      <div class="space-y-6">
        <!-- Header: Administrador Geral vs Tenant -->
        <div v-if="authStore.isGlobalAdmin" class="flex items-center space-x-3 px-2">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-purple-600 to-indigo-500 flex items-center justify-center shadow-lg shadow-purple-500/20 shrink-0">
            <ShieldCheck class="w-5 h-5 text-white font-bold" />
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="font-bold text-sm text-white truncate leading-tight">
              Plataforma
            </h2>
            <div class="flex items-center gap-1.5 mt-0.5">
              <span class="w-2 h-2 rounded-full bg-purple-400 animate-pulse"></span>
              <span class="text-[11px] text-purple-300 font-semibold truncate">Gestão Global Multi-Tenant</span>
            </div>
          </div>
        </div>

        <div v-else class="flex items-center space-x-3 px-2">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-glow shrink-0">
            <Scissors class="w-5 h-5 text-slate-950 font-bold" />
          </div>
          <div class="min-w-0 flex-1">
            <h2 class="font-bold text-sm text-white truncate leading-tight">
              {{ authStore.tenant?.name || 'Painel de Gestão' }}
            </h2>
            <div class="flex items-center gap-1.5 mt-0.5">
              <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
              <span class="text-[11px] text-slate-400 font-mono truncate">/{{ authStore.tenant?.slug }}</span>
            </div>
          </div>
        </div>

        <!-- Botão Visualizar Página Pública (Apenas para Tenants) -->
        <a
          v-if="!authStore.isGlobalAdmin && authStore.tenant?.slug"
          :href="`/agendamento/${authStore.tenant.slug}`"
          target="_blank"
          class="flex items-center justify-center gap-2 w-full py-2 px-3 rounded-xl bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 text-xs font-semibold transition shadow-sm"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Ver Página de Agendamento</span>
        </a>

        <!-- Menu de Navegação Dinâmico por Perfil -->
        <nav class="space-y-1">
          <!-- 1. MENU ADMINISTRADOR GERAL -->
          <template v-if="authStore.isGlobalAdmin">
            <div class="px-3 pb-1 pt-1 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
              Administração Global
            </div>

            <RouterLink
              to="/admin/dashboard"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/dashboard') ? 'bg-purple-600 text-white font-bold shadow-md shadow-purple-600/20' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <LayoutDashboard class="w-4 h-4 shrink-0" />
              <span>Dashboard Global</span>
            </RouterLink>

            <RouterLink
              to="/admin/estabelecimentos"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/estabelecimentos') ? 'bg-purple-600 text-white font-bold shadow-md shadow-purple-600/20' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Building2 class="w-4 h-4 shrink-0" />
              <span>Estabelecimentos</span>
            </RouterLink>

            <RouterLink
              to="/admin/planos"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/planos') ? 'bg-purple-600 text-white font-bold shadow-md shadow-purple-600/20' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <CreditCard class="w-4 h-4 shrink-0" />
              <span>Planos de Assinatura</span>
            </RouterLink>

            <RouterLink
              to="/admin/usuarios"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/usuarios') ? 'bg-purple-600 text-white font-bold shadow-md shadow-purple-600/20' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Users class="w-4 h-4 shrink-0" />
              <span>Usuários da Plataforma</span>
            </RouterLink>

            <RouterLink
              to="/admin/administradores-globais"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/administradores-globais') ? 'bg-purple-600 text-white font-bold shadow-md shadow-purple-600/20' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <ShieldCheck class="w-4 h-4 shrink-0" />
              <span>Administradores Gerais</span>
            </RouterLink>
          </template>

          <!-- 2. MENU ADMINISTRADOR DO ESTABELECIMENTO -->
          <template v-else-if="authStore.isTenantAdmin">
            <div class="px-3 pb-1 pt-1 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
              Gestão do Estabelecimento
            </div>

            <RouterLink
              to="/admin/dashboard"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/dashboard') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <LayoutDashboard class="w-4 h-4 shrink-0" />
              <span>Dashboard</span>
            </RouterLink>

            <RouterLink
              to="/admin/agenda"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/agenda') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Calendar class="w-4 h-4 shrink-0" />
              <span>Agenda & Atendimentos</span>
            </RouterLink>

            <RouterLink
              to="/admin/servicos"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/servicos') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Sparkles class="w-4 h-4 shrink-0" />
              <span>Serviços</span>
            </RouterLink>

            <RouterLink
              to="/admin/profissionais"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/profissionais') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Users class="w-4 h-4 shrink-0" />
              <span>Profissionais & Horários</span>
            </RouterLink>

            <RouterLink
              to="/admin/clientes"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/clientes') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <UserSquare2 class="w-4 h-4 shrink-0" />
              <span>Clientes</span>
            </RouterLink>

            <RouterLink
              to="/admin/minha-assinatura"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/minha-assinatura') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <CreditCard class="w-4 h-4 shrink-0" />
              <span>Minha Assinatura</span>
            </RouterLink>

            <RouterLink
              to="/admin/usuarios"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/usuarios') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <UsersRound class="w-4 h-4 shrink-0" />
              <span>Gestão de Usuários</span>
            </RouterLink>

            <RouterLink
              to="/admin/whatsapp"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/whatsapp') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Bot class="w-4 h-4 shrink-0" />
              <span>WhatsApp & IA</span>
            </RouterLink>

            <RouterLink
              to="/admin/configuracoes"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/configuracoes') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Settings class="w-4 h-4 shrink-0" />
              <span>Configurações</span>
            </RouterLink>
          </template>

          <!-- 3. MENU OPERADOR -->
          <template v-else-if="authStore.isOperator">
            <div class="px-3 pb-1 pt-1 text-[10px] font-bold text-slate-400 uppercase tracking-wider">
              Operações
            </div>

            <RouterLink
              to="/admin/dashboard"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/dashboard') ? 'bg-sky-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <LayoutDashboard class="w-4 h-4 shrink-0" />
              <span>Dashboard</span>
            </RouterLink>

            <RouterLink
              to="/admin/agenda"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/agenda') ? 'bg-sky-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <Calendar class="w-4 h-4 shrink-0" />
              <span>Agenda & Atendimentos</span>
            </RouterLink>

            <RouterLink
              to="/admin/clientes"
              class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
              :class="isActive('/admin/clientes') ? 'bg-sky-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
            >
              <UserSquare2 class="w-4 h-4 shrink-0" />
              <span>Clientes</span>
            </RouterLink>
          </template>
        </nav>
      </div>

      <!-- Usuário Logado & Logout -->
      <div class="pt-4 mt-6 border-t border-slate-800 space-y-3">
        <div class="flex items-center justify-between px-2">
          <div class="min-w-0 flex-1 mr-2">
            <p class="text-xs font-bold text-white truncate">{{ authStore.user?.name }}</p>
            <p class="text-[11px] text-slate-400 truncate">{{ authStore.user?.email }}</p>
          </div>
          <span
            class="px-2 py-0.5 rounded text-[10px] font-bold border shrink-0"
            :class="roleBadgeClasses"
          >
            {{ roleLabel }}
          </span>
        </div>

        <button
          @click="handleLogout"
          class="flex items-center gap-2 w-full py-2 px-3 rounded-xl text-xs font-semibold text-red-400 hover:bg-red-500/10 transition"
        >
          <LogOut class="w-3.5 h-3.5" />
          <span>Sair da Conta</span>
        </button>
      </div>
    </aside>

    <!-- Conteúdo Principal -->
    <main class="flex-1 overflow-y-auto max-h-screen p-4 sm:p-6 lg:p-8">
      <RouterView />
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import {
  Scissors, LayoutDashboard, Calendar, Sparkles, Users, UsersRound,
  UserSquare2, Settings, ExternalLink, LogOut, Building2, ShieldCheck, CreditCard, Bot
} from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

function isActive(path: string) {
  return route.path === path || route.path.startsWith(path + '/')
}

const roleLabel = computed(() => {
  if (authStore.isGlobalAdmin) return 'Admin Geral'
  if (authStore.isTenantAdmin) return 'Admin'
  if (authStore.isOperator) return 'Operador'
  return authStore.user?.role || 'Usuário'
})

const roleBadgeClasses = computed(() => {
  if (authStore.isGlobalAdmin) {
    return 'bg-purple-950/80 text-purple-300 border-purple-700/60 shadow-sm shadow-purple-500/20'
  }
  if (authStore.isTenantAdmin) {
    return 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60'
  }
  return 'bg-sky-950/80 text-sky-300 border-sky-700/60'
})

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
