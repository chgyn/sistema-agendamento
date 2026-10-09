<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight font-['Outfit']">Equipe de Profissionais & Horários</h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-surface-400 mt-1">
          Gerencie os profissionais, jornadas semanais de trabalho e bloqueios de agenda
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 self-start sm:self-auto cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Profissional</span>
      </button>
    </div>

    <!-- Lista de Profissionais -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="pro in professionals"
        :key="pro.id"
        class="saas-card p-6 flex flex-col justify-between hover:border-[#4880FF]/40 transition group"
      >
        <div>
          <div class="flex items-center space-x-3.5 mb-3.5">
            <div class="relative shrink-0">
              <img
                v-if="pro.avatar_url"
                :src="pro.avatar_url"
                :alt="pro.name"
                class="w-14 h-14 rounded-2xl object-cover ring-2 ring-[#4880FF]/30"
              />
              <div v-else class="w-14 h-14 rounded-2xl bg-[#4880FF]/10 border border-blue-200 dark:border-blue-900/40 text-[#4880FF] flex items-center justify-center font-bold text-lg font-['Outfit']">
                {{ pro.name.charAt(0) }}
              </div>
              <span class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 rounded-full bg-emerald-500 border-2 border-white dark:border-[#121318]"></span>
            </div>

            <div class="min-w-0 flex-1">
              <h3 class="font-bold text-base text-[#202224] dark:text-white truncate font-['Outfit'] group-hover:text-[#4880FF] transition">{{ pro.name }}</h3>
              <p class="text-xs text-[#4880FF] font-bold truncate">{{ pro.title || pro.specialty || 'Profissional' }}</p>
              <p class="text-[11px] text-gray-400 dark:text-surface-400 truncate">{{ pro.phone || pro.email || 'Sem contato cadastrado' }}</p>
            </div>
          </div>

          <p v-if="pro.specialty" class="text-xs text-gray-600 dark:text-surface-300 mb-2">
            <strong class="text-gray-700 dark:text-surface-200 font-semibold">Especialidade:</strong> {{ pro.specialty }}
          </p>
          <p v-if="pro.bio" class="text-xs text-gray-500 dark:text-surface-400 line-clamp-2 mb-4 leading-relaxed">
            {{ pro.bio }}
          </p>
        </div>

        <div class="space-y-2 pt-3 border-t border-gray-100 dark:border-surface-800">
          <div class="grid grid-cols-2 gap-2">
            <button
              @click="openWorkingHoursModal(pro)"
              class="py-2.5 px-3 rounded-xl bg-blue-50 hover:bg-blue-100 text-[#4880FF] dark:bg-blue-950/40 dark:hover:bg-blue-900/50 dark:text-blue-300 text-xs font-bold transition flex items-center justify-center gap-1.5 border border-blue-200 dark:border-blue-900/40 cursor-pointer"
            >
              <Clock class="w-3.5 h-3.5" />
              <span>Grade Horária</span>
            </button>
            <button
              @click="openExceptionsModal(pro)"
              class="py-2.5 px-3 rounded-xl bg-amber-50 hover:bg-amber-100 text-amber-700 dark:bg-amber-950/40 dark:hover:bg-amber-900/50 dark:text-amber-300 text-xs font-bold transition flex items-center justify-center gap-1.5 border border-amber-200 dark:border-amber-900/40 cursor-pointer"
            >
              <ShieldAlert class="w-3.5 h-3.5 text-amber-500" />
              <span>Bloqueios / Folga</span>
            </button>
          </div>

          <div class="flex items-center gap-2">
            <button
              @click="openEditModal(pro)"
              class="flex-1 py-2 px-3 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-300 dark:hover:text-white text-xs font-bold transition cursor-pointer"
            >
              Editar Perfil
            </button>
            <button
              @click="deleteProfessional(pro.id)"
              class="py-2 px-3 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 text-xs font-bold transition cursor-pointer"
            >
              Excluir
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar / Editar Profissional -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 w-full max-w-md rounded-2xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-surface-800 pb-3">
          <h3 class="font-bold text-lg text-[#202224] dark:text-white font-['Outfit']">
            {{ editingId ? 'Editar Profissional' : 'Novo Profissional' }}
          </h3>
          <button @click="showModal = false" class="text-gray-400 hover:text-gray-700 dark:text-surface-400 dark:hover:text-white p-1 rounded-lg">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="saveProfessional">
          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Nome Completo *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="Ex: Carlos Navalha"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Cargo / Título</label>
              <input
                v-model="form.title"
                type="text"
                placeholder="Barbeiro Master"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">WhatsApp</label>
              <input
                v-model="form.phone"
                type="text"
                placeholder="(11) 99999-8888"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Especialidades</label>
            <input
              v-model="form.specialty"
              type="text"
              placeholder="Ex: Degradê navalhado, barboterapia"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Foto / Avatar URL</label>
            <input
              v-model="form.avatar_url"
              type="url"
              placeholder="https://..."
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Biografia Curta</label>
            <textarea
              v-model="form.bio"
              rows="2"
              placeholder="Breve resumo da trajetória profissional..."
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition resize-none"
            ></textarea>
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-gray-100 dark:border-surface-800">
            <button
              type="button"
              @click="showModal = false"
              class="px-4 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-xs font-medium transition"
            >
              Cancelar
            </button>
            <button
              type="submit"
              class="px-5 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
            >
              Salvar Profissional
            </button>
          </div>
        </form>
      </div>
    </div>

    <!-- MODAL: Configuração de Horários de Trabalho -->
    <div
      v-if="showWorkingHoursModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 w-full max-w-2xl rounded-2xl p-6 shadow-2xl space-y-4 max-h-[90vh] overflow-y-auto my-8">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-surface-800 pb-3">
          <div>
            <h3 class="font-bold text-lg text-[#202224] dark:text-white font-['Outfit']">Grade Semanal de Trabalho</h3>
            <p class="text-xs text-[#4880FF] font-bold">{{ selectedProForHours?.name }}</p>
          </div>
          <button @click="showWorkingHoursModal = false" class="text-gray-400 hover:text-gray-700 dark:text-surface-400 dark:hover:text-white p-1 rounded-lg">✕</button>
        </div>

        <div class="space-y-3">
          <div
            v-for="wh in workingHoursForm"
            :key="wh.day_of_week"
            class="p-3.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800 flex flex-col sm:flex-row sm:items-center justify-between gap-3"
          >
            <div class="flex items-center gap-2.5 w-36 shrink-0">
              <input
                v-model="wh.is_active"
                type="checkbox"
                :id="`day_${wh.day_of_week}`"
                class="rounded bg-white dark:bg-surface-900 border-gray-300 dark:border-zinc-700 text-[#4880FF] focus:ring-blue-500"
              />
              <label :for="`day_${wh.day_of_week}`" class="text-xs font-bold text-[#202224] dark:text-white cursor-pointer">
                {{ getDayName(wh.day_of_week) }}
              </label>
            </div>

            <div v-if="wh.is_active" class="flex flex-wrap items-center gap-2.5 text-xs">
              <div class="flex items-center gap-1.5">
                <span class="text-gray-500 dark:text-surface-400">Expediente:</span>
                <input
                  v-model="wh.start_time"
                  type="time"
                  class="px-2.5 py-1 rounded-lg bg-white dark:bg-surface-900 border border-gray-200 dark:border-zinc-700 text-[#202224] dark:text-white text-xs"
                />
                <span class="text-gray-400 dark:text-surface-400">às</span>
                <input
                  v-model="wh.end_time"
                  type="time"
                  class="px-2.5 py-1 rounded-lg bg-white dark:bg-surface-900 border border-gray-200 dark:border-zinc-700 text-[#202224] dark:text-white text-xs"
                />
              </div>

              <div class="flex items-center gap-1.5 sm:ml-2">
                <span class="text-gray-500 dark:text-surface-400">Almoço:</span>
                <input
                  v-model="wh.break_start"
                  type="time"
                  class="px-2.5 py-1 rounded-lg bg-white dark:bg-surface-900 border border-gray-200 dark:border-zinc-700 text-[#202224] dark:text-white text-xs"
                />
                <span class="text-gray-400 dark:text-surface-400">às</span>
                <input
                  v-model="wh.break_end"
                  type="time"
                  class="px-2.5 py-1 rounded-lg bg-white dark:bg-surface-900 border border-gray-200 dark:border-zinc-700 text-[#202224] dark:text-white text-xs"
                />
              </div>
            </div>
            <div v-else class="text-xs text-gray-400 dark:text-surface-500 italic">
              Folga (Sem atendimento)
            </div>
          </div>
        </div>

        <div class="flex items-center justify-end gap-3 pt-4 border-t border-gray-100 dark:border-surface-800">
          <button
            type="button"
            @click="showWorkingHoursModal = false"
            class="px-4 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-xs font-medium transition"
          >
            Cancelar
          </button>
          <button
            type="button"
            @click="saveWorkingHours"
            class="px-5 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
          >
            Salvar Grade Semanal
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Bloqueios e Exceções de Agenda -->
    <div
      v-if="showExceptionsModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 w-full max-w-lg rounded-2xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-surface-800 pb-3">
          <div>
            <h3 class="font-bold text-lg text-[#202224] dark:text-white font-['Outfit']">Bloqueios & Exceções de Agenda</h3>
            <p class="text-xs text-[#4880FF] font-bold">{{ selectedProForHours?.name }}</p>
          </div>
          <button @click="showExceptionsModal = false" class="text-gray-400 hover:text-gray-700 dark:text-surface-400 dark:hover:text-white p-1 rounded-lg">✕</button>
        </div>

        <!-- Formulário para adicionar bloqueio -->
        <form class="p-4 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800 space-y-3.5" @submit.prevent="addException">
          <h4 class="text-xs font-bold uppercase tracking-wider text-gray-700 dark:text-surface-300">Novo Bloqueio</h4>
          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-bold text-gray-600 dark:text-surface-400 mb-1">Data *</label>
              <input
                v-model="exceptionForm.date"
                type="date"
                required
                class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-gray-200 dark:border-surface-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-blue-500/20"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-gray-600 dark:text-surface-400 mb-1">Motivo</label>
              <input
                v-model="exceptionForm.reason"
                type="text"
                placeholder="Ex: Consulta médica"
                class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-gray-200 dark:border-surface-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-blue-500/20"
              />
            </div>
          </div>

          <div class="flex items-center gap-2">
            <input
              v-model="exceptionForm.is_full_day"
              type="checkbox"
              id="full_day_exp"
              class="rounded bg-white dark:bg-surface-900 border-gray-300 dark:border-surface-700 text-[#4880FF]"
            />
            <label for="full_day_exp" class="text-xs text-gray-700 dark:text-surface-300 font-semibold cursor-pointer">Bloquear o dia inteiro (Folga/Férias)</label>
          </div>

          <div v-if="!exceptionForm.is_full_day" class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-[11px] font-bold text-gray-600 dark:text-surface-400 mb-1">Hora Início</label>
              <input
                v-model="exceptionForm.start_time"
                type="time"
                class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-gray-200 dark:border-surface-700 text-[#202224] dark:text-white text-xs"
              />
            </div>
            <div>
              <label class="block text-[11px] font-bold text-gray-600 dark:text-surface-400 mb-1">Hora Fim</label>
              <input
                v-model="exceptionForm.end_time"
                type="time"
                class="w-full px-3 py-2 rounded-xl bg-white dark:bg-surface-900 border border-gray-200 dark:border-surface-700 text-[#202224] dark:text-white text-xs"
              />
            </div>
          </div>

          <button
            type="submit"
            class="w-full py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
          >
            Adicionar Bloqueio
          </button>
        </form>

        <!-- Lista de Bloqueios Existentes -->
        <div class="space-y-2 max-h-48 overflow-y-auto">
          <div v-if="exceptionsList.length === 0" class="text-xs text-gray-400 dark:text-surface-500 text-center py-4">
            Nenhum bloqueio cadastrado para este profissional.
          </div>
          <div
            v-for="exp in exceptionsList"
            :key="exp.id"
            class="p-3 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800 flex items-center justify-between text-xs"
          >
            <div>
              <span class="font-bold text-[#202224] dark:text-white block">{{ exp.date }}</span>
              <span class="text-gray-500 dark:text-surface-400">
                {{ exp.is_full_day ? 'Dia Inteiro' : `${exp.start_time} às ${exp.end_time}` }}
                {{ exp.reason ? `— ${exp.reason}` : '' }}
              </span>
            </div>
            <button
              @click="deleteException(exp.id)"
              class="text-rose-500 hover:text-rose-700 dark:text-rose-400 dark:hover:text-rose-300 font-bold p-1.5"
            >
              ✕
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Plus, Clock, ShieldAlert } from 'lucide-vue-next'
import api from '../../services/api'

const professionals = ref<any[]>([])
const showModal = ref(false)
const editingId = ref<string | null>(null)

// Horários de trabalho
const showWorkingHoursModal = ref(false)
const selectedProForHours = ref<any>(null)
const workingHoursForm = ref<any[]>([])

// Bloqueios / Exceções
const showExceptionsModal = ref(false)
const exceptionsList = ref<any[]>([])
const exceptionForm = reactive({
  date: new Date().toISOString().split('T')[0],
  start_time: '14:00',
  end_time: '16:00',
  reason: '',
  is_full_day: false,
})

const form = reactive({
  name: '',
  title: '',
  phone: '',
  specialty: '',
  bio: '',
  avatar_url: '',
})

onMounted(() => {
  loadProfessionals()
})

async function loadProfessionals() {
  try {
    const res = await api.get('/admin/professionals')
    if (res.data.success) {
      professionals.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  }
}

function openCreateModal() {
  editingId.value = null
  form.name = ''
  form.title = ''
  form.phone = ''
  form.specialty = ''
  form.bio = ''
  form.avatar_url = ''
  showModal.value = true
}

function openEditModal(pro: any) {
  editingId.value = pro.id
  form.name = pro.name
  form.title = pro.title
  form.phone = pro.phone
  form.specialty = pro.specialty
  form.bio = pro.bio
  form.avatar_url = pro.avatar_url
  showModal.value = true
}

async function saveProfessional() {
  try {
    if (editingId.value) {
      await api.put(`/admin/professionals/${editingId.value}`, form)
    } else {
      await api.post('/admin/professionals', form)
    }
    showModal.value = false
    await loadProfessionals()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao salvar profissional')
  }
}

async function deleteProfessional(id: string) {
  if (!confirm('Deseja realmente remover este profissional?')) return
  try {
    await api.delete(`/admin/professionals/${id}`)
    await loadProfessionals()
  } catch (err) {
    alert('Erro ao excluir profissional')
  }
}

async function openWorkingHoursModal(pro: any) {
  selectedProForHours.value = pro
  try {
    const res = await api.get(`/admin/professionals/${pro.id}/working-hours`)
    if (res.data.success && res.data.data?.length > 0) {
      workingHoursForm.value = res.data.data
    } else {
      const defaultList = []
      for (let day = 0; day <= 6; day++) {
        defaultList.push({
          day_of_week: day,
          start_time: '08:00',
          end_time: '18:00',
          break_start: '12:00',
          break_end: '13:00',
          is_active: day >= 1 && day <= 5,
        })
      }
      workingHoursForm.value = defaultList
    }
    showWorkingHoursModal.value = true
  } catch (err) {
    console.error(err)
  }
}

async function saveWorkingHours() {
  if (!selectedProForHours.value) return
  try {
    await api.put(`/admin/professionals/${selectedProForHours.value.id}/working-hours`, workingHoursForm.value)
    alert('Grade horária salva com sucesso!')
    showWorkingHoursModal.value = false
  } catch (err) {
    alert('Erro ao salvar horários de trabalho')
  }
}

async function openExceptionsModal(pro: any) {
  selectedProForHours.value = pro
  await loadExceptions(pro.id)
  showExceptionsModal.value = true
}

async function loadExceptions(proId: string) {
  try {
    const res = await api.get(`/admin/professionals/${proId}/exceptions`)
    if (res.data.success) {
      exceptionsList.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  }
}

async function addException() {
  if (!selectedProForHours.value) return
  try {
    await api.post(`/admin/professionals/${selectedProForHours.value.id}/exceptions`, exceptionForm)
    exceptionForm.reason = ''
    await loadExceptions(selectedProForHours.value.id)
  } catch (err) {
    alert('Erro ao adicionar bloqueio')
  }
}

async function deleteException(expId: string) {
  if (!selectedProForHours.value) return
  try {
    await api.delete(`/admin/professionals/${selectedProForHours.value.id}/exceptions/${expId}`)
    await loadExceptions(selectedProForHours.value.id)
  } catch (err) {
    alert('Erro ao remover bloqueio')
  }
}

function getDayName(day: number) {
  const days = ['Domingo', 'Segunda-feira', 'Terça-feira', 'Quarta-feira', 'Quinta-feira', 'Sexta-feira', 'Sábado']
  return days[day] || `Dia ${day}`
}
</script>
