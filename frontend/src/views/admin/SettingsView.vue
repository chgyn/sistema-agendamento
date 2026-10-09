<template>
  <div class="space-y-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight font-['Outfit']">Configurações do Estabelecimento</h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-surface-400 mt-1">
          Personalize as informações públicas, canais de contato e identidade da sua marca
        </p>
      </div>
    </div>

    <!-- Link Público de Agendamento em Destaque -->
    <div class="saas-card p-6 flex flex-col sm:flex-row sm:items-center justify-between gap-4 border-l-4 border-l-[#4880FF]">
      <div>
        <span class="text-xs uppercase font-bold text-[#4880FF] tracking-wider block">Seu Link Público de Agendamento:</span>
        <p class="text-sm font-mono text-[#202224] dark:text-white mt-1 font-bold">
          {{ publicUrl }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="copyPublicUrl"
          class="px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition flex items-center gap-1.5 shadow-md shadow-blue-500/25 cursor-pointer"
        >
          <Copy class="w-3.5 h-3.5" />
          <span>{{ copied ? 'Copiado!' : 'Copiar Link' }}</span>
        </button>
        <a
          :href="`/agendamento/${form.slug}`"
          target="_blank"
          class="px-4 py-2.5 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-200 dark:hover:text-white text-xs font-bold transition flex items-center gap-1.5"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Abrir Página</span>
        </a>
      </div>
    </div>

    <!-- Formulário de Configurações -->
    <div class="saas-card p-6 sm:p-8 space-y-6">
      <form class="space-y-4.5" @submit.prevent="saveSettings">
        <div v-if="successMessage" class="p-3.5 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-400 text-xs flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4 shrink-0" />
          <span>{{ successMessage }}</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
            <input
              v-model="form.document"
              type="text"
              placeholder="00.000.000/0001-00"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
            <input
              v-model="form.phone"
              type="text"
              required
              placeholder="(11) 98765-4321"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Contato</label>
            <input
              v-model="form.email"
              type="email"
              placeholder="contato@empresa.com"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Endereço Completo</label>
          <input
            v-model="form.address"
            type="text"
            placeholder="Rua, Número, Bairro"
            class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
          />
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="col-span-2">
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Cidade</label>
            <input
              v-model="form.city"
              type="text"
              placeholder="São Paulo"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>
          <div>
            <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">UF</label>
            <input
              v-model="form.state"
              type="text"
              maxlength="2"
              placeholder="SP"
              class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs uppercase focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Logo URL</label>
          <input
            v-model="form.logo_url"
            type="url"
            placeholder="https://..."
            class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-[#202224] dark:text-white text-xs focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-[#4880FF] focus:ring-2 focus:ring-blue-500/20 transition"
          />
        </div>

        <div class="pt-4 border-t border-gray-100 dark:border-surface-800 flex justify-end">
          <button
            type="submit"
            :disabled="isSaving"
            class="px-6 py-3 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 disabled:opacity-50 cursor-pointer"
          >
            {{ isSaving ? 'Salvando...' : 'Salvar Alterações' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { Copy, ExternalLink, CheckCircle2 } from 'lucide-vue-next'
import api from '../../services/api'
import { useAuthStore } from '../../stores/auth'

const authStore = useAuthStore()
const isSaving = ref(false)
const copied = ref(false)
const successMessage = ref('')

const form = reactive({
  name: '',
  slug: '',
  document: '',
  phone: '',
  email: '',
  address: '',
  city: '',
  state: '',
  logo_url: '',
  primary_color: '#4880FF',
})

const publicUrl = computed(() => {
  return `${window.location.origin}/agendamento/${form.slug}`
})

onMounted(() => {
  loadSettings()
})

async function loadSettings() {
  try {
    const res = await api.get('/admin/settings')
    if (res.data.success) {
      const data = res.data.data
      form.name = data.name || ''
      form.slug = data.slug || ''
      form.document = data.document || ''
      form.phone = data.phone || ''
      form.email = data.email || ''
      form.address = data.address || ''
      form.city = data.city || ''
      form.state = data.state || ''
      form.logo_url = data.logo_url || ''
      form.primary_color = data.primary_color || '#4880FF'
    }
  } catch (err) {
    console.error(err)
  }
}

async function saveSettings() {
  isSaving.value = true
  successMessage.value = ''
  try {
    const res = await api.put('/admin/settings', form)
    if (res.data.success) {
      authStore.updateTenant(res.data.data)
      successMessage.value = 'Configurações salvas com sucesso!'
      setTimeout(() => { successMessage.value = '' }, 4000)
    }
  } catch (err) {
    alert('Erro ao salvar configurações')
  } finally {
    isSaving.value = false
  }
}

function copyPublicUrl() {
  navigator.clipboard.writeText(publicUrl.value)
  copied.value = true
  setTimeout(() => { copied.value = false }, 3000)
}
</script>
