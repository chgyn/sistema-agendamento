<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- ======================================================== -->
    <!-- 1. DASHBOARD GLOBAL DA PLATAFORMA (ADMIN_GLOBAL) -->
    <!-- ======================================================== -->
    <template v-if="authStore.isGlobalAdmin">
      <!-- Header Global -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
            <ShieldCheck class="w-7 h-7 text-purple-400" />
            <span>Dashboard Global da Plataforma</span>
          </h1>
          <p class="text-xs sm:text-sm text-slate-400 mt-1">
            Visão consolidada de todos os estabelecimentos, agendamentos e faturamento global da plataforma.
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadGlobalKPIs"
            class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/estabelecimentos"
            class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-xs font-bold transition shadow-lg shadow-purple-600/30"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Estabelecimento</span>
          </RouterLink>
        </div>
      </div>

      <!-- Cards de Métricas Globais -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Estabelecimentos -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Estabelecimentos</p>
            <h3 class="text-2xl font-black text-white">{{ globalStats.total_tenants || 0 }}</h3>
            <p class="text-[11px] text-emerald-400 font-medium mt-1">
              {{ globalStats.active_tenants || 0 }} ativos / {{ globalStats.inactive_tenants || 0 }} inativos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-purple-500/10 text-purple-400 flex items-center justify-center border border-purple-500/20">
            <Building2 class="w-6 h-6" />
          </div>
        </div>

        <!-- Faturamento Global Estimado -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Faturamento Global</p>
            <h3 class="text-2xl font-black text-emerald-400">R$ {{ (globalStats.total_revenue || 0).toFixed(2) }}</h3>
            <p class="text-[11px] text-slate-400 font-medium mt-1">Total acumulado na plataforma</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Agendamentos -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Agendamentos Totais</p>
            <h3 class="text-2xl font-black text-white">{{ globalStats.total_appointments || 0 }}</h3>
            <p class="text-[11px] text-blue-400 font-medium mt-1">
              {{ globalStats.completed_appointments || 0 }} atendimentos concluídos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-blue-500/10 text-blue-400 flex items-center justify-center border border-blue-500/20">
            <Calendar class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Clientes na Plataforma -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Clientes na Base</p>
            <h3 class="text-2xl font-black text-white">{{ globalStats.total_customers || 0 }}</h3>
            <p class="text-[11px] text-purple-400 font-medium mt-1">
              {{ globalStats.total_users || 0 }} usuários e operadores
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-teal-500/10 text-teal-400 flex items-center justify-center border border-teal-500/20">
            <Users class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Grade de Informações Recentes: Estabelecimentos & Agendamentos -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Estabelecimentos Recentes -->
        <div class="glass-panel rounded-2xl border border-slate-800 overflow-hidden">
          <div class="p-5 border-b border-slate-800 flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <Store class="w-5 h-5 text-purple-400" />
              <h2 class="font-bold text-base text-white">Estabelecimentos Recentes</h2>
            </div>
            <RouterLink to="/admin/estabelecimentos" class="text-xs text-purple-400 hover:text-purple-300 font-semibold">
              Ver Todos →
            </RouterLink>
          </div>

          <div v-if="!globalStats.recent_tenants || globalStats.recent_tenants.length === 0" class="p-8 text-center text-slate-500 text-xs">
            Nenhum estabelecimento cadastrado.
          </div>
          <div v-else class="divide-y divide-slate-800/60">
            <div v-for="t in globalStats.recent_tenants" :key="t.id" class="p-4 flex items-center justify-between hover:bg-slate-900/50 transition">
              <div class="flex items-center gap-3">
                <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center font-bold text-sm text-purple-400 shrink-0">
                  {{ t.name.charAt(0) }}
                </div>
                <div>
                  <h4 class="font-bold text-sm text-white">{{ t.name }}</h4>
                  <p class="text-xs text-slate-400">/{{ t.slug }} • {{ t.city || 'Brasil' }}</p>
                </div>
              </div>
              <span
                class="px-2.5 py-0.5 rounded text-[10px] font-bold border"
                :class="t.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-rose-950/80 text-rose-300 border-rose-700/60'"
              >
                {{ t.is_active ? 'Ativo' : 'Inativo' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Últimos Agendamentos na Plataforma -->
        <div class="glass-panel rounded-2xl border border-slate-800 overflow-hidden">
          <div class="p-5 border-b border-slate-800 flex items-center justify-between">
            <div class="flex items-center space-x-2">
              <Clock class="w-5 h-5 text-emerald-400" />
              <h2 class="font-bold text-base text-white">Atividade Recente na Plataforma</h2>
            </div>
          </div>

          <div v-if="!globalStats.recent_appointments || globalStats.recent_appointments.length === 0" class="p-8 text-center text-slate-500 text-xs">
            Nenhuma atividade de agendamento recente.
          </div>
          <div v-else class="divide-y divide-slate-800/60">
            <div v-for="apt in globalStats.recent_appointments" :key="apt.id" class="p-4 flex items-center justify-between hover:bg-slate-900/50 transition">
              <div>
                <h4 class="font-bold text-sm text-white">{{ apt.customer?.name || 'Cliente' }}</h4>
                <p class="text-xs text-slate-400">
                  <span class="text-emerald-400 font-medium">{{ apt.service?.name }}</span> com {{ apt.professional?.name }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-xs font-bold text-white">R$ {{ apt.total_price?.toFixed(2) }}</p>
                <span :class="getStatusBadgeClass(apt.status)" class="inline-block mt-1">
                  {{ formatStatus(apt.status) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ======================================================== -->
    <!-- 2. DASHBOARD DO TENANT (ADMIN_TENANT / OPERATOR) -->
    <!-- ======================================================== -->
    <template v-else>
      <!-- Header Tenant -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-black text-white tracking-tight">Dashboard & Indicadores</h1>
          <p class="text-xs sm:text-sm text-slate-400 mt-1">
            Visão geral do movimento e agendamentos de hoje: <strong class="text-emerald-400">{{ todayFormatted }}</strong>
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadKPIs"
            class="flex items-center gap-1.5 px-3 py-2 rounded-xl bg-slate-900 border border-slate-800 hover:border-slate-700 text-slate-300 text-xs font-semibold transition"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/agenda"
            class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition shadow-md"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Agendamento</span>
          </RouterLink>
        </div>
      </div>

      <!-- Cards de Métricas Principais (KPIs) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Agendamentos Hoje -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Agendamentos Hoje</p>
            <h3 class="text-2xl font-black text-white">{{ kpis.today_appointments_count || 0 }}</h3>
            <p class="text-[11px] text-emerald-400 font-medium mt-1">
              {{ kpis.confirmed_count || 0 }} confirmados / {{ kpis.completed_count || 0 }} concluídos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20">
            <Calendar class="w-6 h-6" />
          </div>
        </div>

        <!-- Faturamento Estimado Hoje -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Faturamento Hoje</p>
            <h3 class="text-2xl font-black text-emerald-400">R$ {{ (kpis.today_revenue || 0).toFixed(2) }}</h3>
            <p class="text-[11px] text-slate-400 font-medium mt-1">Atendimentos previstos</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Clientes Cadastrados -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Total de Clientes</p>
            <h3 class="text-2xl font-black text-white">{{ kpis.total_customers_count || 0 }}</h3>
            <p class="text-[11px] text-slate-400 font-medium mt-1">Base fidelizada</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-teal-500/10 text-teal-400 flex items-center justify-center border border-teal-500/20">
            <Users class="w-6 h-6" />
          </div>
        </div>

        <!-- Profissionais Ativos -->
        <div class="glass-panel p-5 rounded-2xl border border-slate-800 flex items-center justify-between">
          <div>
            <p class="text-xs uppercase tracking-wider text-slate-400 font-semibold mb-1">Profissionais</p>
            <h3 class="text-2xl font-black text-white">{{ kpis.active_professionals_count || 0 }}</h3>
            <p class="text-[11px] text-slate-400 font-medium mt-1">Equipe disponível</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-cyan-500/10 text-cyan-400 flex items-center justify-center border border-cyan-500/20">
            <UserCheck class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Tabela / Timeline de Atendimentos de Hoje -->
      <div class="glass-panel rounded-2xl border border-slate-800 overflow-hidden">
        <div class="p-5 border-b border-slate-800 flex items-center justify-between">
          <div class="flex items-center space-x-2">
            <Clock class="w-5 h-5 text-emerald-400" />
            <h2 class="font-bold text-base text-white">Cronograma de Atendimentos de Hoje</h2>
          </div>
          <RouterLink to="/admin/agenda" class="text-xs text-emerald-400 hover:text-emerald-300 font-semibold">
            Ver Agenda Completa →
          </RouterLink>
        </div>

        <div v-if="isLoading" class="text-center py-12">
          <Loader2 class="w-8 h-8 animate-spin text-emerald-500 mx-auto mb-2" />
          <p class="text-xs text-slate-400">Carregando cronograma...</p>
        </div>

        <div v-else-if="!kpis.today_appointments || kpis.today_appointments.length === 0" class="text-center py-12 px-4">
          <CalendarX class="w-10 h-10 text-slate-600 mx-auto mb-2" />
          <p class="text-sm font-semibold text-slate-300">Nenhum atendimento agendado para hoje</p>
          <p class="text-xs text-slate-500 mt-1">Novos agendamentos aparecerão aqui em tempo real.</p>
        </div>

        <div v-else class="divide-y divide-slate-800/80">
          <div
            v-for="apt in kpis.today_appointments"
            :key="apt.id"
            class="p-4 sm:px-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-slate-900/50 transition"
          >
            <div class="flex items-center space-x-4">
              <div class="w-14 text-center py-1 px-2 rounded-xl bg-slate-900 border border-slate-800">
                <span class="text-xs font-bold text-white block">{{ formatTime(apt.start_at) }}</span>
                <span class="text-[10px] text-slate-500">{{ apt.duration_minutes }}m</span>
              </div>

              <div>
                <h4 class="font-bold text-sm text-white">{{ apt.customer?.name || 'Cliente' }}</h4>
                <p class="text-xs text-slate-400 flex items-center gap-2 mt-0.5">
                  <span class="text-emerald-400 font-medium">{{ apt.service?.name }}</span>
                  <span>•</span>
                  <span>Profissional: <strong class="text-slate-300">{{ apt.professional?.name }}</strong></span>
                </p>
                <p v-if="apt.notes" class="text-[11px] text-amber-400/90 mt-1 italic">
                  Obs: {{ apt.notes }}
                </p>
              </div>
            </div>

            <div class="flex items-center gap-3 sm:self-center self-end">
              <span class="text-sm font-extrabold text-white mr-2">
                R$ {{ apt.total_price?.toFixed(2) }}
              </span>

              <span :class="getStatusBadgeClass(apt.status)">
                {{ formatStatus(apt.status) }}
              </span>

              <div v-if="apt.status === 'CONFIRMED'" class="flex items-center gap-1.5 ml-2">
                <button
                  @click="completeAppointment(apt.id)"
                  class="px-2.5 py-1 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500 hover:text-slate-950 text-xs font-bold transition"
                  title="Concluir Atendimento"
                >
                  Concluir ✓
                </button>
                <button
                  @click="cancelAppointment(apt.id)"
                  class="px-2.5 py-1 rounded-lg bg-red-500/10 text-red-400 hover:bg-red-500 hover:text-white text-xs font-bold transition"
                  title="Cancelar"
                >
                  ✕
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import {
  Calendar, Clock, DollarSign, Users, UserCheck, Plus,
  RefreshCw, CalendarX, Loader2, Building2, Store, ShieldCheck
} from 'lucide-vue-next'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '../../services/api'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const isLoading = ref(true)
const kpis = ref<any>({})
const globalStats = ref<any>({})

const todayFormatted = computed(() => {
  return format(new Date(), "EEEE, dd 'de' MMMM 'de' yyyy", { locale: ptBR })
})

onMounted(() => {
  if (authStore.isGlobalAdmin) {
    loadGlobalKPIs()
  } else {
    loadKPIs()
  }
})

async function loadGlobalKPIs() {
  isLoading.value = true
  try {
    const res = await api.get('/admin/global-dashboard')
    if (res.data.success) {
      globalStats.value = res.data.data
    }
  } catch (err) {
    console.error('Erro ao carregar dados do dashboard global:', err)
  } finally {
    isLoading.value = false
  }
}

async function loadKPIs() {
  isLoading.value = true
  try {
    const res = await api.get('/admin/dashboard')
    if (res.data.success) {
      kpis.value = res.data.data
    }
  } catch (err) {
    console.error('Erro ao carregar KPIs do tenant:', err)
  } finally {
    isLoading.value = false
  }
}

async function completeAppointment(id: string) {
  try {
    await api.patch(`/admin/appointments/${id}/complete`)
    await loadKPIs()
  } catch (err) {
    alert('Erro ao concluir agendamento')
  }
}

async function cancelAppointment(id: string) {
  const reason = prompt('Informe o motivo do cancelamento:')
  if (reason === null) return
  try {
    await api.patch(`/admin/appointments/${id}/cancel`, { reason })
    await loadKPIs()
  } catch (err) {
    alert('Erro ao cancelar agendamento')
  }
}

function formatTime(dateStr: string) {
  if (!dateStr) return ''
  try {
    return format(new Date(dateStr), 'HH:mm')
  } catch {
    return dateStr
  }
}

function formatStatus(status: string) {
  switch (status) {
    case 'CONFIRMED': return 'Confirmado'
    case 'COMPLETED': return 'Concluído'
    case 'CANCELLED': return 'Cancelado'
    case 'PENDING': return 'Pendente'
    default: return status
  }
}

function getStatusBadgeClass(status: string) {
  switch (status) {
    case 'CONFIRMED':
      return 'px-2.5 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20'
    case 'COMPLETED':
      return 'px-2.5 py-1 rounded-full text-xs font-semibold bg-blue-500/10 text-blue-400 border border-blue-500/20'
    case 'CANCELLED':
      return 'px-2.5 py-1 rounded-full text-xs font-semibold bg-red-500/10 text-red-400 border border-red-500/20'
    default:
      return 'px-2.5 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/20'
  }
}
</script>

