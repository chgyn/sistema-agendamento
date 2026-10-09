<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header da Página -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl sm:text-2xl font-black text-[#202224] dark:text-white font-display tracking-tight">
          Schedule List (Agenda de Atendimentos)
        </h2>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-zinc-400 mt-0.5">
          Visualize a escala diária, filtre por profissional e gerencie agendamentos em tempo real.
        </p>
      </div>

      <button
        @click="openNewModal"
        class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 self-start sm:self-auto cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>+ Add New</span>
      </button>
    </div>

    <!-- Layout em Duas Colunas (Figma 4:1635 - Schedule List) -->
    <div class="grid grid-cols-1 lg:grid-cols-12 gap-6 items-start">
      <!-- ======================================================== -->
      <!-- COLUNA DA ESQUERDA: MINI CALENDÁRIO + FILTRO DE PROFISSIONAIS -->
      <!-- ======================================================== -->
      <div class="lg:col-span-4 space-y-6">
        <!-- Botão Primário "+ Create Schedule" -->
        <button
          @click="openNewModal"
          class="w-full py-3 px-4 rounded-2xl bg-[#4880FF] hover:bg-[#386FF0] text-white font-bold text-sm transition shadow-md shadow-blue-500/25 flex items-center justify-center gap-2 cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>+ Create Schedule</span>
        </button>

        <!-- Mini Calendário Interativo do Mês (Figma Widget) -->
        <div class="saas-card p-5">
          <!-- Cabeçalho do Mês -->
          <div class="flex items-center justify-between mb-4">
            <h4 class="font-bold text-sm text-[#202224] dark:text-white capitalize font-display">
              {{ currentMonthName }} {{ currentYear }}
            </h4>
            <div class="flex items-center gap-1">
              <button
                @click="prevMonth"
                class="p-1 rounded-lg text-gray-400 hover:text-[#4880FF] hover:bg-gray-100 dark:hover:bg-zinc-800 transition"
              >
                <ChevronLeft class="w-4 h-4" />
              </button>
              <button
                @click="nextMonth"
                class="p-1 rounded-lg text-gray-400 hover:text-[#4880FF] hover:bg-gray-100 dark:hover:bg-zinc-800 transition"
              >
                <ChevronRight class="w-4 h-4" />
              </button>
            </div>
          </div>

          <!-- Dias da Semana (D, S, T, Q, Q, S, S) -->
          <div class="grid grid-cols-7 text-center text-[11px] font-bold text-[#718096] dark:text-zinc-500 mb-2">
            <span>D</span>
            <span>S</span>
            <span>T</span>
            <span>Q</span>
            <span>Q</span>
            <span>S</span>
            <span>S</span>
          </div>

          <!-- Grade de Dias -->
          <div class="grid grid-cols-7 gap-1 text-center text-xs">
            <button
              v-for="(day, idx) in calendarDays"
              :key="idx"
              @click="selectCalendarDay(day)"
              :disabled="!day.isCurrentMonth"
              class="w-8 h-8 mx-auto rounded-full flex items-center justify-center font-medium transition cursor-pointer"
              :class="getDayClasses(day)"
            >
              {{ day.dayNumber }}
            </button>
          </div>
        </div>

        <!-- Bloco People / Filtro de Profissionais (Figma 4:1635) -->
        <div class="saas-card p-5 space-y-4">
          <div class="flex items-center justify-between">
            <h4 class="font-bold text-sm text-[#202224] dark:text-white font-display">
              People (Profissionais)
            </h4>
            <button
              v-if="filters.professional_id"
              @click="selectProfessional('')"
              class="text-[11px] text-[#4880FF] hover:underline font-semibold"
            >
              Ver Todos
            </button>
          </div>

          <!-- Campo de Busca de Profissional -->
          <div class="relative">
            <Search class="w-3.5 h-3.5 text-gray-400 absolute left-3 top-3 pointer-events-none" />
            <input
              v-model="searchPeople"
              type="text"
              placeholder="Buscar profissional..."
              class="w-full pl-8 pr-3 py-2 rounded-xl bg-gray-50 dark:bg-zinc-800/80 border border-gray-200/80 dark:border-zinc-700/60 text-xs text-[#202224] dark:text-zinc-200 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF] transition"
            />
          </div>

          <!-- Lista de Profissionais com Avatar -->
          <div class="space-y-2 max-h-64 overflow-y-auto">
            <div
              v-for="pro in filteredProfessionals"
              :key="pro.id"
              @click="selectProfessional(pro.id)"
              class="flex items-center gap-3 p-2.5 rounded-xl cursor-pointer transition"
              :class="filters.professional_id === pro.id ? 'bg-[#E9F0FE] text-[#4880FF] font-bold dark:bg-[#4880FF]/15' : 'hover:bg-gray-50/80 dark:hover:bg-zinc-800/50'"
            >
              <div class="w-8 h-8 rounded-full bg-gradient-to-tr from-[#4880FF] to-[#8280FF] text-white font-bold text-xs flex items-center justify-center shrink-0 uppercase shadow-sm">
                {{ pro.name.charAt(0) }}
              </div>
              <div class="min-w-0 flex-1">
                <p class="text-xs font-bold truncate leading-tight">{{ pro.name }}</p>
                <p class="text-[10px] text-[#718096] dark:text-zinc-400 truncate">{{ pro.phone || 'Profissional Ativo' }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- ======================================================== -->
      <!-- COLUNA DA DIREITA: LISTA DE AGENDAMENTOS (SCHEDULE LIST) -->
      <!-- ======================================================== -->
      <div class="lg:col-span-8 space-y-4">
        <!-- Barra Superior de Filtros de Status (Pílulas Estilo Figma) -->
        <div class="saas-card p-4 flex flex-wrap items-center justify-between gap-3">
          <div class="flex items-center gap-2">
            <button
              @click="setStatusFilter('')"
              class="px-3 py-1.5 rounded-xl text-xs font-bold transition cursor-pointer"
              :class="filters.status === '' ? 'bg-[#4880FF] text-white shadow-sm' : 'bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#202224] dark:hover:text-white'"
            >
              Todos
            </button>
            <button
              @click="setStatusFilter('CONFIRMED')"
              class="px-3 py-1.5 rounded-xl text-xs font-bold transition cursor-pointer"
              :class="filters.status === 'CONFIRMED' ? 'bg-[#4880FF] text-white shadow-sm' : 'bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#202224] dark:hover:text-white'"
            >
              Confirmados
            </button>
            <button
              @click="setStatusFilter('COMPLETED')"
              class="px-3 py-1.5 rounded-xl text-xs font-bold transition cursor-pointer"
              :class="filters.status === 'COMPLETED' ? 'bg-[#4880FF] text-white shadow-sm' : 'bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#202224] dark:hover:text-white'"
            >
              Concluídos
            </button>
            <button
              @click="setStatusFilter('CANCELLED')"
              class="px-3 py-1.5 rounded-xl text-xs font-bold transition cursor-pointer"
              :class="filters.status === 'CANCELLED' ? 'bg-[#4880FF] text-white shadow-sm' : 'bg-gray-100 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#202224] dark:hover:text-white'"
            >
              Cancelados
            </button>
          </div>

          <button
            v-if="filters.date || filters.professional_id || filters.status"
            @click="resetFilters"
            class="text-xs text-[#718096] hover:text-[#4880FF] font-semibold transition cursor-pointer"
          >
            Limpar Filtros ✕
          </button>
        </div>

        <!-- Estado de Carregamento -->
        <div v-if="isLoading" class="saas-card py-16 text-center">
          <Loader2 class="w-8 h-8 animate-spin text-[#4880FF] mx-auto mb-2" />
          <p class="text-xs text-[#718096]">Carregando atendimentos da agenda...</p>
        </div>

        <!-- Estado Vazio -->
        <div v-else-if="appointments.length === 0" class="saas-card py-16 px-4 text-center">
          <CalendarX class="w-10 h-10 text-gray-300 dark:text-zinc-600 mx-auto mb-2" />
          <h4 class="font-bold text-sm text-[#202224] dark:text-white font-display">Nenhum agendamento encontrado</h4>
          <p class="text-xs text-[#718096] dark:text-zinc-400 mt-1">
            Selecione outro dia no calendário à esquerda ou crie um novo atendimento.
          </p>
          <button
            @click="openNewModal"
            class="mt-4 px-4 py-2 rounded-xl bg-[#4880FF] text-white text-xs font-bold hover:bg-[#386FF0] transition"
          >
            + Criar Agendamento
          </button>
        </div>

        <!-- Lista de Cards de Atendimento (Padrão Figma Schedule List) -->
        <div v-else class="space-y-3">
          <div
            v-for="apt in appointments"
            :key="apt.id"
            class="saas-card p-4 sm:p-5 flex flex-col sm:flex-row sm:items-center justify-between gap-4 hover:shadow-card-hover transition"
          >
            <!-- Lado Esquerdo: Data, Hora, Categoria/Serviço e Cliente -->
            <div class="flex items-start sm:items-center gap-4 min-w-0">
              <!-- Data Badge com Ícone de Calendário -->
              <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-gray-50 dark:bg-zinc-800 text-xs font-bold text-[#202224] dark:text-zinc-200 border border-gray-200/80 dark:border-zinc-700/60 shrink-0">
                <Calendar class="w-3.5 h-3.5 text-[#4880FF]" />
                <span>{{ formatDateShort(apt.start_at) }}</span>
              </div>

              <!-- Hora Badge com Ícone de Relógio -->
              <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-gray-50 dark:bg-zinc-800 text-xs font-bold text-[#202224] dark:text-zinc-200 border border-gray-200/80 dark:border-zinc-700/60 shrink-0">
                <Clock class="w-3.5 h-3.5 text-[#FEC53D]" />
                <span>{{ formatTime(apt.start_at) }}</span>
              </div>

              <!-- Pílula de Categoria / Serviço (Figma Pill) -->
              <span class="px-3 py-1.5 rounded-full text-xs font-bold bg-[#E9F0FE] text-[#4880FF] dark:bg-blue-950/60 dark:text-blue-300 shrink-0 truncate max-w-[150px]">
                {{ apt.service?.name }}
              </span>

              <!-- Dados do Cliente & Profissional -->
              <div class="min-w-0 hidden md:block">
                <h4 class="font-bold text-xs text-[#202224] dark:text-white truncate font-display">
                  {{ apt.customer?.name || 'Cliente' }}
                </h4>
                <p class="text-[11px] text-[#718096] dark:text-zinc-400 truncate">
                  com {{ apt.professional?.name }}
                </p>
              </div>
            </div>

            <!-- Lado Direito: Preço, Status e Botões de Ação Circulares (Figma Action Icons) -->
            <div class="flex items-center justify-between sm:justify-end gap-3 shrink-0">
              <span class="text-xs font-black text-[#202224] dark:text-white font-display">
                R$ {{ (apt.total_price || 0).toFixed(2).replace('.', ',') }}
              </span>

              <!-- Botão Amarelo Circular (Editar / Reagendar - Figma Pencil Icon) -->
              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="openRescheduleModal(apt)"
                class="w-8 h-8 rounded-full bg-[#FFF7E6] text-[#FFB800] hover:bg-[#FFEEC2] dark:bg-amber-950/50 dark:text-amber-400 flex items-center justify-center transition cursor-pointer"
                title="Reagendar"
              >
                <Pencil class="w-3.5 h-3.5" />
              </button>

              <!-- Botão Verde Circular (Concluir Atendimento) -->
              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="completeAppointment(apt.id)"
                class="w-8 h-8 rounded-full bg-[#E6FBF2] text-[#10B981] hover:bg-[#C7F7E3] dark:bg-emerald-950/50 dark:text-emerald-400 flex items-center justify-center transition cursor-pointer"
                title="Concluir"
              >
                <Check class="w-3.5 h-3.5" />
              </button>

              <!-- Botão Vermelho Circular (Cancelar - Figma Trash Icon) -->
              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="cancelAppointment(apt.id)"
                class="w-8 h-8 rounded-full bg-[#FFEFE7] text-[#FF6647] hover:bg-[#FFD9CE] dark:bg-rose-950/50 dark:text-rose-400 flex items-center justify-center transition cursor-pointer"
                title="Cancelar Atendimento"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- ======================================================== -->
    <!-- MODAL: Novo Agendamento Administrativo -->
    <!-- ======================================================== -->
    <div v-if="showNewModal" class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="saas-card w-full max-w-lg p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-zinc-800 pb-3">
          <h3 class="font-bold text-base text-[#202224] dark:text-white font-display">Novo Agendamento</h3>
          <button @click="showNewModal = false" class="text-gray-400 hover:text-gray-600 text-lg font-bold">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="submitNewAppointment">
          <div v-if="modalError" class="p-3 rounded-xl bg-rose-50 text-rose-600 text-xs font-semibold">
            {{ modalError }}
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Serviço *</label>
              <select
                v-model="newAptForm.service_id"
                required
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              >
                <option value="" disabled>Selecione</option>
                <option v-for="s in services" :key="s.id" :value="s.id">{{ s.name }} (R$ {{ s.price.toFixed(2).replace('.', ',') }})</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Profissional *</label>
              <select
                v-model="newAptForm.professional_id"
                required
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              >
                <option value="" disabled>Selecione</option>
                <option v-for="p in professionals" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Data *</label>
              <input
                v-model="newAptForm.date"
                type="date"
                required
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Horário (HH:MM) *</label>
              <input
                v-model="newAptForm.time"
                type="time"
                required
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Nome do Cliente *</label>
              <input
                v-model="newAptForm.customer_name"
                type="text"
                required
                placeholder="Ex: Carlos Silva"
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white placeholder-gray-400 text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">WhatsApp / Telefone *</label>
              <input
                v-model="newAptForm.customer_phone"
                type="text"
                required
                placeholder="(11) 99999-8888"
                class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white placeholder-gray-400 text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Observações</label>
            <input
              v-model="newAptForm.notes"
              type="text"
              placeholder="Opcional"
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white placeholder-gray-400 text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-gray-100 dark:border-zinc-800">
            <button
              type="button"
              @click="showNewModal = false"
              class="px-4 py-2.5 rounded-xl bg-gray-100 hover:bg-gray-200 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 text-xs font-bold transition"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="isSubmittingModal"
              class="px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 disabled:opacity-50 cursor-pointer"
            >
              {{ isSubmittingModal ? 'Salvando...' : 'Confirmar Agendamento' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- ======================================================== -->
    <!-- MODAL: Reagendamento -->
    <!-- ======================================================== -->
    <div v-if="showRescheduleModal" class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="saas-card w-full max-w-md p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-zinc-800 pb-3">
          <h3 class="font-bold text-base text-[#202224] dark:text-white font-display">Reagendar Atendimento</h3>
          <button @click="showRescheduleModal = false" class="text-gray-400 hover:text-gray-600 text-lg font-bold">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="submitReschedule">
          <div v-if="modalError" class="p-3 rounded-xl bg-rose-50 text-rose-600 text-xs font-semibold">
            {{ modalError }}
          </div>

          <div>
            <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Nova Data *</label>
            <input
              v-model="rescheduleForm.date"
              type="date"
              required
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#718096] dark:text-zinc-400 uppercase tracking-wider mb-1">Novo Horário (HH:MM) *</label>
            <input
              v-model="rescheduleForm.time"
              type="time"
              required
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-gray-100 dark:border-zinc-800">
            <button
              type="button"
              @click="showRescheduleModal = false"
              class="px-4 py-2.5 rounded-xl bg-gray-100 hover:bg-gray-200 dark:bg-zinc-800 text-[#718096] dark:text-zinc-300 text-xs font-bold transition"
            >
              Fechar
            </button>
            <button
              type="submit"
              class="px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
            >
              Salvar Reagendamento
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import {
  Plus, CalendarX, Loader2, Calendar, Clock,
  Search, ChevronLeft, ChevronRight, Pencil, Check, Trash2
} from 'lucide-vue-next'
import { format, startOfMonth, endOfMonth, eachDayOfInterval, getDay, isSameDay } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '../../services/api'

const appointments = ref<any[]>([])
const professionals = ref<any[]>([])
const services = ref<any[]>([])
const isLoading = ref(true)
const searchPeople = ref('')

const filters = reactive({
  date: format(new Date(), 'yyyy-MM-dd'),
  professional_id: '',
  status: '',
})

// Calendário Interativo
const calendarDate = ref(new Date())

const currentMonthName = computed(() => {
  return format(calendarDate.value, 'MMMM', { locale: ptBR })
})

const currentYear = computed(() => {
  return format(calendarDate.value, 'yyyy')
})

const calendarDays = computed(() => {
  const start = startOfMonth(calendarDate.value)
  const end = endOfMonth(calendarDate.value)
  const daysInMonth = eachDayOfInterval({ start, end })

  // Dias em branco antes do início do mês para alinhar com o dia da semana
  const startDayOfWeek = getDay(start) // 0 = Domingo
  const leadingBlanks = Array.from({ length: startDayOfWeek }, (_, i) => ({
    dayNumber: '',
    date: null,
    isCurrentMonth: false,
    isSelected: false,
  }))

  const monthDays = daysInMonth.map(d => {
    const isSelected = filters.date ? isSameDay(d, new Date(filters.date + 'T12:00:00')) : false
    return {
      dayNumber: format(d, 'd'),
      date: d,
      isCurrentMonth: true,
      isSelected,
    }
  })

  return [...leadingBlanks, ...monthDays]
})

function prevMonth() {
  const d = new Date(calendarDate.value)
  d.setMonth(d.getMonth() - 1)
  calendarDate.value = d
}

function nextMonth() {
  const d = new Date(calendarDate.value)
  d.setMonth(d.getMonth() + 1)
  calendarDate.value = d
}

function selectCalendarDay(day: any) {
  if (!day.date) return
  filters.date = format(day.date, 'yyyy-MM-dd')
  loadAppointments()
}

function getDayClasses(day: any) {
  if (!day.isCurrentMonth) return 'text-transparent cursor-default'
  if (day.isSelected) return 'bg-[#4880FF] text-white font-bold shadow-sm'
  return 'text-[#202224] dark:text-zinc-200 hover:bg-blue-50 dark:hover:bg-zinc-800'
}

// Filtro de Pessoas
const filteredProfessionals = computed(() => {
  if (!searchPeople.value) return professionals.value
  const q = searchPeople.value.toLowerCase()
  return professionals.value.filter(p => p.name.toLowerCase().includes(q))
})

function selectProfessional(id: string) {
  filters.professional_id = id
  loadAppointments()
}

function setStatusFilter(status: string) {
  filters.status = status
  loadAppointments()
}

// Modais
const showNewModal = ref(false)
const showRescheduleModal = ref(false)
const isSubmittingModal = ref(false)
const modalError = ref('')
const selectedAptForReschedule = ref<any>(null)

const newAptForm = reactive({
  service_id: '',
  professional_id: '',
  date: format(new Date(), 'yyyy-MM-dd'),
  time: '10:00',
  customer_name: '',
  customer_phone: '',
  notes: '',
})

const rescheduleForm = reactive({
  date: format(new Date(), 'yyyy-MM-dd'),
  time: '14:00',
})

onMounted(async () => {
  await Promise.all([loadAppointments(), loadAuxData()])
})

async function loadAuxData() {
  try {
    const [prosRes, svcsRes] = await Promise.all([
      api.get('/admin/professionals'),
      api.get('/admin/services'),
    ])
    if (prosRes.data.success) professionals.value = prosRes.data.data
    if (svcsRes.data.success) services.value = svcsRes.data.data
  } catch (err) {
    console.error(err)
  }
}

async function loadAppointments() {
  isLoading.value = true
  try {
    const params: any = {}
    if (filters.date) {
      params.date_from = filters.date
      params.date_to = filters.date
    }
    if (filters.professional_id) params.professional_id = filters.professional_id
    if (filters.status) params.status = filters.status

    const res = await api.get('/admin/appointments', { params })
    if (res.data.success) {
      appointments.value = res.data.data
    }
  } catch (err) {
    console.error('Erro ao listar agendamentos:', err)
  } finally {
    isLoading.value = false
  }
}

function resetFilters() {
  filters.date = ''
  filters.professional_id = ''
  filters.status = ''
  loadAppointments()
}

function openNewModal() {
  modalError.value = ''
  newAptForm.customer_name = ''
  newAptForm.customer_phone = ''
  newAptForm.notes = ''
  showNewModal.value = true
}

async function submitNewAppointment() {
  isSubmittingModal.value = true
  modalError.value = ''

  try {
    const startAt = `${newAptForm.date}T${newAptForm.time}:00`
    await api.post('/admin/appointments', {
      service_id: newAptForm.service_id,
      professional_id: newAptForm.professional_id,
      start_at: startAt,
      customer_name: newAptForm.customer_name,
      customer_phone: newAptForm.customer_phone,
      notes: newAptForm.notes,
    })

    showNewModal.value = false
    await loadAppointments()
  } catch (err: any) {
    modalError.value = err.response?.data?.error || 'Erro ao criar agendamento'
  } finally {
    isSubmittingModal.value = false
  }
}

function openRescheduleModal(apt: any) {
  modalError.value = ''
  selectedAptForReschedule.value = apt
  showRescheduleModal.value = true
}

async function submitReschedule() {
  if (!selectedAptForReschedule.value) return
  modalError.value = ''

  try {
    const startAt = `${rescheduleForm.date}T${rescheduleForm.time}:00`
    await api.patch(`/admin/appointments/${selectedAptForReschedule.value.id}/reschedule`, {
      start_at: startAt,
    })

    showRescheduleModal.value = false
    await loadAppointments()
  } catch (err: any) {
    modalError.value = err.response?.data?.error || 'Erro ao reagendar'
  }
}

async function completeAppointment(id: string) {
  try {
    await api.patch(`/admin/appointments/${id}/complete`)
    await loadAppointments()
  } catch (err) {
    alert('Erro ao concluir agendamento')
  }
}

async function cancelAppointment(id: string) {
  const reason = prompt('Motivo do cancelamento:')
  if (reason === null) return
  try {
    await api.patch(`/admin/appointments/${id}/cancel`, { reason })
    await loadAppointments()
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

function formatDateShort(dateStr: string) {
  if (!dateStr) return ''
  try {
    return format(new Date(dateStr), 'dd MMM', { locale: ptBR })
  } catch {
    return dateStr
  }
}
</script>
