<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- ======================================================== -->
    <!-- 1. DASHBOARD GLOBAL DA PLATAFORMA (ADMIN_GLOBAL) -->
    <!-- ======================================================== -->
    <template v-if="authStore.isGlobalAdmin">
      <!-- Header Global -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight flex items-center gap-2.5 font-display">
            <ShieldCheck class="w-8 h-8 text-purple-400" />
            <span>Dashboard Global da Plataforma</span>
          </h1>
          <p class="text-xs sm:text-sm text-zinc-400 mt-1">
            Visão consolidada de todos os estabelecimentos, agendamentos e faturamento global da plataforma.
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadGlobalKPIs"
            class="flex items-center gap-1.5 px-3.5 py-2.5 rounded-xl bg-[#181922] border border-zinc-700/60 hover:border-zinc-600 text-zinc-300 hover:text-white text-xs font-bold transition shadow-sm cursor-pointer"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/estabelecimentos"
            class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-xs font-bold transition shadow-lg shadow-purple-600/30"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Estabelecimento</span>
          </RouterLink>
        </div>
      </div>

      <!-- Cards de Métricas Globais -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Estabelecimentos -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Estabelecimentos</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ globalStats.total_tenants || 0 }}</h3>
            <p class="text-[11px] text-emerald-400 font-semibold mt-1">
              {{ globalStats.active_tenants || 0 }} ativos / {{ globalStats.inactive_tenants || 0 }} inativos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-purple-500/10 text-purple-400 flex items-center justify-center border border-purple-500/20 shadow-sm">
            <Building2 class="w-6 h-6" />
          </div>
        </div>

        <!-- Faturamento Global Estimado -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Faturamento Global</p>
            <h3 class="text-2xl sm:text-3xl font-black text-emerald-400 font-display">R$ {{ (globalStats.total_revenue || 0).toFixed(2).replace('.', ',') }}</h3>
            <p class="text-[11px] text-zinc-400 font-medium mt-1">Total acumulado na plataforma</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20 shadow-sm">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Agendamentos -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Agendamentos Totais</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ globalStats.total_appointments || 0 }}</h3>
            <p class="text-[11px] text-blue-400 font-semibold mt-1">
              {{ globalStats.completed_appointments || 0 }} concluídos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-blue-500/10 text-blue-400 flex items-center justify-center border border-blue-500/20 shadow-sm">
            <Calendar class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Clientes na Plataforma -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Clientes na Base</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ globalStats.total_customers || 0 }}</h3>
            <p class="text-[11px] text-purple-400 font-semibold mt-1">
              {{ globalStats.total_users || 0 }} usuários do sistema
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-teal-500/10 text-teal-400 flex items-center justify-center border border-teal-500/20 shadow-sm">
            <Users class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Grade de Informações Recentes -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Estabelecimentos Recentes -->
        <div class="glass-panel rounded-2xl border border-zinc-800 overflow-hidden">
          <div class="p-5 border-b border-zinc-800 flex items-center justify-between">
            <div class="flex items-center space-x-2.5">
              <Store class="w-5 h-5 text-purple-400" />
              <h2 class="font-bold text-base text-white font-display">Estabelecimentos Recentes</h2>
            </div>
            <RouterLink to="/admin/estabelecimentos" class="text-xs text-purple-400 hover:text-purple-300 font-bold">
              Ver Todos →
            </RouterLink>
          </div>

          <div v-if="!globalStats.recent_tenants || globalStats.recent_tenants.length === 0" class="p-8 text-center text-zinc-500 text-xs">
            Nenhum estabelecimento cadastrado.
          </div>
          <div v-else class="divide-y divide-zinc-800/60">
            <div v-for="t in globalStats.recent_tenants" :key="t.id" class="p-4 flex items-center justify-between hover:bg-[#16171e]/80 transition">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center font-bold text-sm text-purple-400 shrink-0 font-display">
                  {{ t.name.charAt(0) }}
                </div>
                <div>
                  <h4 class="font-bold text-sm text-white font-display">{{ t.name }}</h4>
                  <p class="text-xs text-zinc-400">/{{ t.slug }} • {{ t.city || 'Brasil' }}</p>
                </div>
              </div>
              <span
                class="px-2.5 py-1 rounded-full text-[10px] font-bold border"
                :class="t.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-rose-950/80 text-rose-300 border-rose-700/60'"
              >
                {{ t.is_active ? '● Ativo' : '○ Inativo' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Últimos Agendamentos na Plataforma -->
        <div class="glass-panel rounded-2xl border border-zinc-800 overflow-hidden">
          <div class="p-5 border-b border-zinc-800 flex items-center justify-between">
            <div class="flex items-center space-x-2.5">
              <Clock class="w-5 h-5 text-orange-400" />
              <h2 class="font-bold text-base text-white font-display">Atividade Recente na Plataforma</h2>
            </div>
          </div>

          <div v-if="!globalStats.recent_appointments || globalStats.recent_appointments.length === 0" class="p-8 text-center text-zinc-500 text-xs">
            Nenhuma atividade de agendamento recente.
          </div>
          <div v-else class="divide-y divide-zinc-800/60">
            <div v-for="apt in globalStats.recent_appointments" :key="apt.id" class="p-4 flex items-center justify-between hover:bg-[#16171e]/80 transition">
              <div>
                <h4 class="font-bold text-sm text-white font-display">{{ apt.customer?.name || 'Cliente' }}</h4>
                <p class="text-xs text-zinc-400 mt-0.5">
                  <span class="text-orange-400 font-semibold">{{ apt.service?.name }}</span> com {{ apt.professional?.name }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-xs font-black text-white font-display">R$ {{ apt.total_price?.toFixed(2).replace('.', ',') }}</p>
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
          <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-display">Dashboard & Indicadores</h1>
          <p class="text-xs sm:text-sm text-zinc-400 mt-1">
            Visão geral do movimento e agendamentos de hoje: <strong class="text-orange-400 font-bold capitalize">{{ todayFormatted }}</strong>
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadKPIs"
            class="flex items-center gap-1.5 px-3.5 py-2.5 rounded-xl bg-[#181922] border border-zinc-700/60 hover:border-zinc-600 text-zinc-300 hover:text-white text-xs font-bold transition shadow-sm cursor-pointer"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/agenda"
            class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-zinc-950 text-xs font-bold transition shadow-glow-sm cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Agendamento</span>
          </RouterLink>
        </div>
      </div>

      <!-- Cards de Métricas Principais (KPIs) Estilo EstiloMarca -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Total Agendamentos Hoje -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Agendamentos Hoje</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ kpis.today_appointments_count || 0 }}</h3>
            <p class="text-[11px] text-orange-400 font-bold mt-1">
              {{ kpis.confirmed_count || 0 }} confirmados / {{ kpis.completed_count || 0 }} concluídos
            </p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-400 flex items-center justify-center border border-orange-500/20 shadow-sm">
            <Calendar class="w-6 h-6" />
          </div>
        </div>

        <!-- Faturamento Estimado Hoje -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Faturamento Hoje</p>
            <h3 class="text-2xl sm:text-3xl font-black text-emerald-400 font-display">R$ {{ (kpis.today_revenue || 0).toFixed(2).replace('.', ',') }}</h3>
            <p class="text-[11px] text-zinc-400 font-medium mt-1">Atendimentos previstos</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 text-emerald-400 flex items-center justify-center border border-emerald-500/20 shadow-sm">
            <DollarSign class="w-6 h-6" />
          </div>
        </div>

        <!-- Total de Clientes Cadastrados -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Total de Clientes</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ kpis.total_customers_count || 0 }}</h3>
            <p class="text-[11px] text-zinc-400 font-medium mt-1">Base fidelizada</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-amber-500/10 text-amber-400 flex items-center justify-center border border-amber-500/20 shadow-sm">
            <Users class="w-6 h-6" />
          </div>
        </div>

        <!-- Profissionais Ativos -->
        <div class="glass-panel p-5 rounded-2xl border border-zinc-800 flex items-center justify-between">
          <div>
            <p class="text-[11px] uppercase tracking-wider text-zinc-400 font-bold mb-1">Profissionais</p>
            <h3 class="text-2xl sm:text-3xl font-black text-white font-display">{{ kpis.active_professionals_count || 0 }}</h3>
            <p class="text-[11px] text-zinc-400 font-medium mt-1">Equipe disponível</p>
          </div>
          <div class="w-12 h-12 rounded-2xl bg-orange-500/10 text-orange-400 flex items-center justify-center border border-orange-500/20 shadow-sm">
            <UserCheck class="w-6 h-6" />
          </div>
        </div>
      </div>

      <!-- Tabela / Timeline de Atendimentos de Hoje -->
      <div class="glass-panel rounded-2xl border border-zinc-800 overflow-hidden">
        <div class="p-5 border-b border-zinc-800 flex items-center justify-between">
          <div class="flex items-center space-x-2.5">
            <Clock class="w-5 h-5 text-orange-400" />
            <h2 class="font-bold text-base text-white font-display">Cronograma de Atendimentos de Hoje</h2>
          </div>
          <RouterLink to="/admin/agenda" class="text-xs text-orange-400 hover:text-orange-300 font-bold">
            Ver Agenda Completa →
          </RouterLink>
        </div>

        <div v-if="isLoading" class="text-center py-12">
          <Loader2 class="w-8 h-8 animate-spin text-orange-500 mx-auto mb-2" />
          <p class="text-xs text-zinc-400">Carregando cronograma...</p>
        </div>

        <div v-else-if="!kpis.today_appointments || kpis.today_appointments.length === 0" class="text-center py-12 px-4">
          <CalendarX class="w-10 h-10 text-zinc-600 mx-auto mb-2" />
          <p class="text-sm font-bold text-zinc-300">Nenhum atendimento agendado para hoje</p>
          <p class="text-xs text-zinc-500 mt-1">Novos agendamentos aparecerão aqui em tempo real.</p>
        </div>

        <div v-else class="divide-y divide-zinc-800/80">
          <div
            v-for="apt in kpis.today_appointments"
            :key="apt.id"
            class="p-4 sm:px-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:bg-[#16171e]/80 transition"
          >
            <div class="flex items-center space-x-4">
              <div class="w-14 text-center py-1.5 px-2 rounded-xl bg-[#14151c] border border-zinc-750">
                <span class="text-xs font-black text-orange-400 block font-display">{{ formatTime(apt.start_at) }}</span>
                <span class="text-[10px] text-zinc-500">{{ apt.duration_minutes }}m</span>
              </div>

              <div>
                <h4 class="font-bold text-sm text-white font-display">{{ apt.customer?.name || 'Cliente' }}</h4>
                <p class="text-xs text-zinc-400 flex items-center gap-2 mt-0.5">
                  <span class="text-orange-400 font-semibold">{{ apt.service?.name }}</span>
                  <span>•</span>
                  <span>Profissional: <strong class="text-zinc-300">{{ apt.professional?.name }}</strong></span>
                </p>
                <p v-if="apt.notes" class="text-[11px] text-amber-400/90 mt-1 italic">
                  Obs: {{ apt.notes }}
                </p>
              </div>
            </div>

            <div class="flex items-center gap-3 sm:self-center self-end">
              <span class="text-sm font-black text-white mr-2 font-display">
                R$ {{ apt.total_price?.toFixed(2).replace('.', ',') }}
              </span>

              <span :class="getStatusBadgeClass(apt.status)">
                {{ formatStatus(apt.status) }}
              </span>

              <div v-if="apt.status === 'CONFIRMED'" class="flex items-center gap-1.5 ml-2">
                <button
                  @click="completeAppointment(apt.id)"
                  class="px-2.5 py-1 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500 hover:text-zinc-950 text-xs font-bold transition cursor-pointer"
                  title="Concluir Atendimento"
                >
                  Concluir ✓
                </button>
                <button
                  @click="cancelAppointment(apt.id)"
                  class="px-2.5 py-1 rounded-lg bg-rose-500/15 text-rose-400 hover:bg-rose-500 hover:text-white text-xs font-bold transition cursor-pointer"
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
      return 'px-2.5 py-1 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-400 border border-emerald-500/30'
    case 'COMPLETED':
      return 'px-2.5 py-1 rounded-full text-xs font-bold bg-blue-500/10 text-blue-400 border border-blue-500/30'
    case 'CANCELLED':
      return 'px-2.5 py-1 rounded-full text-xs font-bold bg-rose-500/10 text-rose-400 border border-rose-500/30'
    default:
      return 'px-2.5 py-1 rounded-full text-xs font-bold bg-amber-500/10 text-amber-400 border border-amber-500/30'
  }
}
</script>
