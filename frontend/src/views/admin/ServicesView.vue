<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight font-['Outfit']">Catálogo de Serviços</h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-surface-400 mt-1">
          Configure os serviços prestados, durações, preços e profissionais habilitados
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 self-start sm:self-auto cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>Cadastrar Novo Serviço</span>
      </button>
    </div>

    <!-- Lista de Serviços -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="svc in services"
        :key="svc.id"
        class="saas-card p-6 flex flex-col justify-between hover:border-[#4880FF]/40 transition group"
      >
        <div>
          <div class="flex items-start justify-between gap-2 mb-2.5">
            <h3 class="font-bold text-base text-[#202224] dark:text-white font-['Outfit'] group-hover:text-[#4880FF] transition">{{ svc.name }}</h3>
            <span class="text-xs font-black text-[#4880FF] dark:text-blue-400 bg-blue-50 dark:bg-blue-950/60 px-3 py-1 rounded-xl border border-blue-200 dark:border-blue-900/40 shrink-0 font-mono">
              R$ {{ svc.price.toFixed(2).replace('.', ',') }}
            </span>
          </div>

          <p class="text-xs text-[#718096] dark:text-surface-400 mb-4 line-clamp-2 leading-relaxed">
            {{ svc.description || 'Sem descrição cadastrada.' }}
          </p>
        </div>

        <div>
          <div class="flex items-center justify-between py-2.5 border-t border-gray-100 dark:border-surface-800 text-xs text-gray-500 dark:text-surface-400">
            <span class="flex items-center gap-1.5 font-medium">
              <Clock class="w-3.5 h-3.5 text-[#4880FF]" />
              {{ svc.duration_minutes }} minutos
            </span>
            <span :class="svc.is_active ? 'text-emerald-600 dark:text-emerald-400 font-bold' : 'text-gray-400 dark:text-surface-500'">
              {{ svc.is_active ? '● Ativo' : '○ Inativo' }}
            </span>
          </div>

          <div class="flex items-center gap-2 pt-3 border-t border-gray-100 dark:border-surface-800">
            <button
              @click="openEditModal(svc)"
              class="flex-1 py-2 px-3 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-300 dark:hover:text-white text-xs font-bold transition cursor-pointer"
            >
              Editar
            </button>
            <button
              @click="deleteService(svc.id)"
              class="py-2 px-3 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 text-xs font-bold transition cursor-pointer"
            >
              Excluir
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar / Editar Serviço -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 w-full max-w-md rounded-2xl p-6 shadow-2xl space-y-4 my-8">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-surface-800 pb-3">
          <h3 class="font-bold text-lg text-[#202224] dark:text-white font-['Outfit']">
            {{ editingId ? 'Editar Serviço' : 'Novo Serviço' }}
          </h3>
          <button @click="showModal = false" class="text-gray-400 hover:text-gray-700 dark:text-surface-400 dark:hover:text-white p-1 rounded-lg">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="saveService">
          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Nome do Serviço *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="Ex: Corte Degradê Navalhado"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Duração (Minutos) *</label>
              <input
                v-model.number="form.duration_minutes"
                type="number"
                min="5"
                step="5"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Preço (R$) *</label>
              <input
                v-model.number="form.price"
                type="number"
                min="0"
                step="0.50"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs font-mono font-bold focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1">Descrição</label>
            <textarea
              v-model="form.description"
              rows="2"
              placeholder="Descreva detalhes do serviço..."
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition resize-none"
            ></textarea>
          </div>

          <div class="flex items-center gap-2 pt-2">
            <input
              v-model="form.is_active"
              type="checkbox"
              id="is_active_svc"
              class="rounded bg-gray-50 dark:bg-surface-950/80 border-gray-300 dark:border-surface-700 text-[#4880FF] focus:ring-blue-500"
            />
            <label for="is_active_svc" class="text-xs text-gray-700 dark:text-surface-300 font-semibold cursor-pointer">Serviço ativo para agendamento público</label>
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
              Salvar Serviço
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Plus, Clock } from 'lucide-vue-next'
import api from '../../services/api'

const services = ref<any[]>([])
const showModal = ref(false)
const editingId = ref<string | null>(null)

const form = reactive({
  name: '',
  description: '',
  duration_minutes: 30,
  price: 45.0,
  is_active: true,
})

onMounted(() => {
  loadServices()
})

async function loadServices() {
  try {
    const res = await api.get('/admin/services')
    if (res.data.success) {
      services.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  }
}

function openCreateModal() {
  editingId.value = null
  form.name = ''
  form.description = ''
  form.duration_minutes = 30
  form.price = 45.0
  form.is_active = true
  showModal.value = true
}

function openEditModal(svc: any) {
  editingId.value = svc.id
  form.name = svc.name
  form.description = svc.description
  form.duration_minutes = svc.duration_minutes
  form.price = svc.price
  form.is_active = svc.is_active
  showModal.value = true
}

async function saveService() {
  try {
    if (editingId.value) {
      await api.put(`/admin/services/${editingId.value}`, form)
    } else {
      await api.post('/admin/services', form)
    }
    showModal.value = false
    await loadServices()
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao salvar serviço')
  }
}

async function deleteService(id: string) {
  if (!confirm('Deseja realmente remover este serviço?')) return
  try {
    await api.delete(`/admin/services/${id}`)
    await loadServices()
  } catch (err) {
    alert('Erro ao excluir serviço')
  }
}
</script>
