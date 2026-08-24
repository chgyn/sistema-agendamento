<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header & Filtros -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">Agenda de Atendimentos</h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Gerencie e acompanhe todos os agendamentos da sua barbearia/salão
        </p>
      </div>

      <button
        @click="openNewModal"
        class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition shadow-md self-start sm:self-auto"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Agendamento</span>
      </button>
    </div>

    <!-- Barra de Filtros -->
    <div class="glass-panel p-4 rounded-2xl border border-slate-800 grid grid-cols-1 sm:grid-cols-4 gap-3">
      <div>
        <label class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">Data:</label>
        <input
          v-model="filters.date"
          @change="loadAppointments"
          type="date"
          class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:outline-none focus:ring-2 focus:ring-emerald-500"
        />
      </div>

      <div>
        <label class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">Profissional:</label>
        <select
          v-model="filters.professional_id"
          @change="loadAppointments"
          class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
          <option value="">Todos os Profissionais</option>
          <option v-for="pro in professionals" :key="pro.id" :value="pro.id">{{ pro.name }}</option>
        </select>
      </div>

      <div>
        <label class="block text-[11px] font-semibold uppercase tracking-wider text-slate-400 mb-1">Status:</label>
        <select
          v-model="filters.status"
          @change="loadAppointments"
          class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:outline-none focus:ring-2 focus:ring-emerald-500"
        >
          <option value="">Todos os Status</option>
          <option value="CONFIRMED">Confirmado</option>
          <option value="COMPLETED">Concluído</option>
          <option value="CANCELLED">Cancelado</option>
        </select>
      </div>

      <div class="flex items-end">
        <button
          @click="resetFilters"
          class="w-full py-2 px-3 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
        >
          Limpar Filtros
        </button>
      </div>
    </div>

    <!-- Lista de Agendamentos -->
    <div class="glass-panel rounded-2xl border border-slate-800 overflow-hidden">
      <div v-if="isLoading" class="text-center py-16">
        <Loader2 class="w-8 h-8 animate-spin text-emerald-500 mx-auto mb-2" />
        <p class="text-xs text-slate-400">Carregando agendamentos...</p>
      </div>

      <div v-else-if="appointments.length === 0" class="text-center py-16 px-4">
        <CalendarX class="w-10 h-10 text-slate-600 mx-auto mb-2" />
        <p class="text-sm font-semibold text-slate-300">Nenhum agendamento encontrado para este filtro</p>
        <p class="text-xs text-slate-500 mt-1">Crie um agendamento manual ou altere a data selecionada.</p>
      </div>

      <div v-else class="divide-y divide-slate-800/80">
        <div
          v-for="apt in appointments"
          :key="apt.id"
          class="p-4 sm:px-6 flex flex-col md:flex-row md:items-center justify-between gap-4 hover:bg-slate-900/50 transition"
        >
          <div class="flex items-start sm:items-center space-x-4">
            <!-- Bloco de Horário -->
            <div class="w-16 text-center py-2 px-2.5 rounded-xl bg-slate-900 border border-slate-800 shrink-0">
              <span class="text-xs font-black text-emerald-400 block">{{ formatTime(apt.start_at) }}</span>
              <span class="text-[10px] text-slate-400 block">{{ formatDateShort(apt.start_at) }}</span>
            </div>

            <!-- Dados do Atendimento -->
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-sm text-white">{{ apt.customer?.name || 'Cliente' }}</h3>
                <span class="text-xs text-slate-400">({{ apt.customer?.phone }})</span>
              </div>
              <p class="text-xs text-slate-400 mt-0.5 flex flex-wrap items-center gap-2">
                <span class="text-white font-medium">{{ apt.service?.name }}</span>
                <span>•</span>
                <span>Duração: {{ apt.duration_minutes }} min</span>
                <span>•</span>
                <span>Profissional: <strong class="text-emerald-400">{{ apt.professional?.name }}</strong></span>
              </p>
              <p v-if="apt.notes" class="text-[11px] text-amber-400/90 mt-1">
                Obs: {{ apt.notes }}
              </p>
              <p v-if="apt.cancellation_reason" class="text-[11px] text-red-400/90 mt-1">
                Motivo cancelamento: {{ apt.cancellation_reason }}
              </p>
            </div>
          </div>

          <!-- Preço, Status e Ações -->
          <div class="flex items-center gap-3 sm:self-center self-end">
            <span class="text-sm font-extrabold text-white">
              R$ {{ apt.total_price?.toFixed(2) }}
            </span>

            <span :class="getStatusBadgeClass(apt.status)">
              {{ formatStatus(apt.status) }}
            </span>

            <!-- Ações -->
            <div class="flex items-center gap-1.5 ml-2">
              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="completeAppointment(apt.id)"
                class="px-2.5 py-1.5 rounded-lg bg-emerald-500/20 text-emerald-400 hover:bg-emerald-500 hover:text-slate-950 text-xs font-bold transition"
                title="Concluir Atendimento"
              >
                Concluir ✓
              </button>

              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="openRescheduleModal(apt)"
                class="px-2.5 py-1.5 rounded-lg bg-slate-800 text-slate-300 hover:bg-slate-700 text-xs font-semibold transition"
                title="Reagendar"
              >
                Reagendar
              </button>

              <button
                v-if="apt.status === 'CONFIRMED'"
                @click="cancelAppointment(apt.id)"
                class="px-2.5 py-1.5 rounded-lg bg-red-500/10 text-red-400 hover:bg-red-500 hover:text-white text-xs font-bold transition"
                title="Cancelar"
              >
                Cancelar
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Novo Agendamento Administrativo -->
    <div v-if="showNewModal" class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="glass-panel w-full max-w-lg rounded-2xl border border-slate-800 p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <h3 class="font-bold text-lg text-white">Novo Agendamento Manual</h3>
          <button @click="showNewModal = false" class="text-slate-400 hover:text-white text-lg">✕</button>
        </div>

        <form class="space-y-3" @submit.prevent="submitNewAppointment">
          <div v-if="modalError" class="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-400 text-xs">
            {{ modalError }}
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Serviço *</label>
              <select
                v-model="newAptForm.service_id"
                required
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              >
                <option value="" disabled>Selecione</option>
                <option v-for="s in services" :key="s.id" :value="s.id">{{ s.name }} (R$ {{ s.price.toFixed(2) }})</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Profissional *</label>
              <select
                v-model="newAptForm.professional_id"
                required
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              >
                <option value="" disabled>Selecione</option>
                <option v-for="p in professionals" :key="p.id" :value="p.id">{{ p.name }}</option>
              </select>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Data *</label>
              <input
                v-model="newAptForm.date"
                type="date"
                required
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Horário (HH:MM) *</label>
              <input
                v-model="newAptForm.time"
                type="time"
                required
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              />
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">Nome do Cliente *</label>
              <input
                v-model="newAptForm.customer_name"
                type="text"
                required
                placeholder="Ex: Carlos Silva"
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 mb-1">WhatsApp / Telefone *</label>
              <input
                v-model="newAptForm.customer_phone"
                type="text"
                required
                placeholder="(11) 99999-8888"
                class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Observações</label>
            <input
              v-model="newAptForm.notes"
              type="text"
              placeholder="Opcional"
              class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-800">
            <button
              type="button"
              @click="showNewModal = false"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
            >
              Cancelar
            </button>
            <button
              type="submit"
              :disabled="isSubmittingModal"
              class="px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition disabled:opacity-50"
            >
              {{ isSubmittingModal ? 'Salvando...' : 'Confirmar Agendamento' }}
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Reagendamento -->
    <div v-if="showRescheduleModal" class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="glass-panel w-full max-w-md rounded-2xl border border-slate-800 p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <h3 class="font-bold text-lg text-white">Reagendar Atendimento</h3>
          <button @click="showRescheduleModal = false" class="text-slate-400 hover:text-white text-lg">✕</button>
        </div>

        <form class="space-y-3" @submit.prevent="submitReschedule">
          <div v-if="modalError" class="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-400 text-xs">
            {{ modalError }}
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Nova Data *</label>
            <input
              v-model="rescheduleForm.date"
              type="date"
              required
              class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 mb-1">Novo Horário (HH:MM) *</label>
            <input
              v-model="rescheduleForm.time"
              type="time"
              required
              class="w-full px-3 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-slate-800">
            <button
              type="button"
              @click="showRescheduleModal = false"
              class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition"
            >
              Fechar
            </button>
            <button
              type="submit"
              class="px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition"
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
import { onMounted, reactive, ref } from 'vue'
import { Plus, CalendarX, Loader2 } from 'lucide-vue-next'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '../../services/api'

const appointments = ref<any[]>([])
const professionals = ref<any[]>([])
const services = ref<any[]>([])
const isLoading = ref(true)

const filters = reactive({
  date: '',
  professional_id: '',
  status: '',
})

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
    return format(new Date(dateStr), 'dd/MM', { locale: ptBR })
  } catch {
    return dateStr
  }
}

function formatStatus(status: string) {
  switch (status) {
    case 'CONFIRMED': return 'Confirmado'
    case 'COMPLETED': return 'Concluído'
    case 'CANCELLED': return 'Cancelado'
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
