<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-display">Catálogo de Serviços</h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1">
          Configure os serviços prestados, durações, preços e profissionais habilitados
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-zinc-950 text-xs font-bold transition shadow-glow-sm self-start sm:self-auto cursor-pointer"
      >
        <Plus class="w-4 h-4" />
        <span>Cadastrar Novo Serviço</span>
      </button>
    </div>

    <!-- Lista de Serviços -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="svc in services"
        :key="svc.id"
        class="glass-panel p-5 rounded-2xl border border-zinc-800 flex flex-col justify-between hover:border-zinc-700 hover:bg-[#16171e]/90 transition"
      >
        <div>
          <div class="flex items-start justify-between gap-2 mb-2.5">
            <h3 class="font-bold text-base text-white font-display">{{ svc.name }}</h3>
            <span class="text-xs font-black text-orange-400 bg-orange-950/60 px-3 py-1 rounded-xl border border-orange-800/40 shrink-0 font-display">
              R$ {{ svc.price.toFixed(2).replace('.', ',') }}
            </span>
          </div>

          <p class="text-xs text-zinc-400 mb-4 line-clamp-2 leading-relaxed">
            {{ svc.description || 'Sem descrição cadastrada.' }}
          </p>
        </div>

        <div>
          <div class="flex items-center justify-between py-2.5 border-t border-zinc-800 text-xs text-zinc-400">
            <span class="flex items-center gap-1.5 font-semibold">
              <Clock class="w-3.5 h-3.5 text-orange-400" />
              {{ svc.duration_minutes }} minutos
            </span>
            <span :class="svc.is_active ? 'text-emerald-400 font-bold' : 'text-zinc-500'">
              {{ svc.is_active ? '● Ativo' : '○ Inativo' }}
            </span>
          </div>

          <div class="flex items-center gap-2 pt-3 border-t border-zinc-800">
            <button
              @click="openEditModal(svc)"
              class="flex-1 py-2 px-3 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-200 text-xs font-bold transition border border-zinc-700/60 cursor-pointer"
            >
              Editar
            </button>
            <button
              @click="deleteService(svc.id)"
              class="py-2 px-3 rounded-xl bg-rose-500/15 hover:bg-rose-500/25 text-rose-400 text-xs font-bold transition cursor-pointer"
            >
              Excluir
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar / Editar Serviço -->
    <div v-if="showModal" class="fixed inset-0 z-50 bg-black/80 backdrop-blur-md flex items-center justify-center p-4">
      <div class="glass-panel w-full max-w-md rounded-2xl border border-zinc-800 p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <h3 class="font-bold text-lg text-white font-display">
            {{ editingId ? 'Editar Serviço' : 'Novo Serviço' }}
          </h3>
          <button @click="showModal = false" class="text-zinc-400 hover:text-white text-lg font-bold">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="saveService">
          <div>
            <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1">Nome do Serviço *</label>
            <input
              v-model="form.name"
              type="text"
              required
              placeholder="Ex: Corte Degradê Navalhado"
              class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white text-xs focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500"
            />
          </div>

          <div class="grid grid-cols-2 gap-3">
            <div>
              <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1">Duração (Minutos) *</label>
              <input
                v-model.number="form.duration_minutes"
                type="number"
                min="5"
                step="5"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white text-xs focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500"
              />
            </div>
            <div>
              <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1">Preço (R$) *</label>
              <input
                v-model.number="form.price"
                type="number"
                min="0"
                step="0.50"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white text-xs focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1">Descrição</label>
            <textarea
              v-model="form.description"
              rows="2"
              placeholder="Descreva detalhes do serviço..."
              class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white text-xs focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 resize-none"
            ></textarea>
          </div>

          <div class="flex items-center gap-2 pt-2">
            <input
              v-model="form.is_active"
              type="checkbox"
              id="is_active_svc"
              class="rounded bg-[#14151c] border-zinc-750 text-orange-500 focus:ring-orange-500"
            />
            <label for="is_active_svc" class="text-xs text-zinc-300 font-semibold">Serviço ativo para agendamento público</label>
          </div>

          <div class="flex items-center justify-end gap-3 pt-4 border-t border-zinc-800">
            <button
              type="button"
              @click="showModal = false"
              class="px-4 py-2.5 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-300 text-xs font-bold transition border border-zinc-700/60"
            >
              Cancelar
            </button>
            <button
              type="submit"
              class="px-4 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-zinc-950 text-xs font-bold transition shadow-glow-sm cursor-pointer"
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
