<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col md:flex-row">
    <!-- Sidebar Desktop -->
    <aside class="w-full md:w-64 bg-slate-900/90 border-r border-slate-800/80 flex flex-col justify-between p-4 shrink-0">
      <div class="space-y-6">
        <!-- Logo & Tenant Name -->
        <div class="flex items-center space-x-3 px-2">
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

        <!-- Botão Visualizar Página Pública -->
        <a
          v-if="authStore.tenant?.slug"
          :href="`/agendamento/${authStore.tenant.slug}`"
          target="_blank"
          class="flex items-center justify-center gap-2 w-full py-2 px-3 rounded-xl bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 text-xs font-semibold transition"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Ver Página de Agendamento</span>
        </a>

        <!-- Menu de Navegação -->
        <nav class="space-y-1">
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
            to="/admin/configuracoes"
            class="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition"
            :class="isActive('/admin/configuracoes') ? 'bg-emerald-500 text-slate-950 font-bold shadow-md' : 'text-slate-300 hover:bg-slate-800 hover:text-white'"
          >
            <Settings class="w-4 h-4 shrink-0" />
            <span>Configurações</span>
          </RouterLink>
        </nav>
      </div>

      <!-- Usuário Logado & Logout -->
      <div class="pt-4 mt-6 border-t border-slate-800 space-y-3">
        <div class="flex items-center justify-between px-2">
          <div class="min-w-0">
            <p class="text-xs font-bold text-white truncate">{{ authStore.user?.name }}</p>
            <p class="text-[11px] text-slate-400 truncate">{{ authStore.user?.email }}</p>
          </div>
          <span class="px-2 py-0.5 rounded text-[10px] font-bold bg-slate-800 text-emerald-400 border border-slate-700">
            {{ authStore.user?.role }}
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
import { useRoute, useRouter } from 'vue-router'
import {
  Scissors, LayoutDashboard, Calendar, Sparkles, Users,
  UserSquare2, Settings, ExternalLink, LogOut
} from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

function isActive(path: string) {
  return route.path === path || route.path.startsWith(path + '/')
}

function handleLogout() {
  authStore.logout()
  router.push('/login')
}
</script>
