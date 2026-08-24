<template>
  <div class="space-y-6 max-w-4xl mx-auto">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">Configurações do Estabelecimento</h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Personalize as informações públicas, canais de contato e identidade da sua marca
        </p>
      </div>
    </div>

    <!-- Link Público de Agendamento em Destaque -->
    <div class="p-5 rounded-2xl bg-gradient-to-r from-emerald-950/40 via-teal-950/20 to-slate-900 border border-emerald-500/30 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <span class="text-xs uppercase font-bold text-emerald-400 tracking-wider">Seu Link Público de Agendamento:</span>
        <p class="text-sm font-mono text-white mt-1">
          {{ publicUrl }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <button
          @click="copyPublicUrl"
          class="px-3.5 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition flex items-center gap-1.5 shadow-md"
        >
          <Copy class="w-3.5 h-3.5" />
          <span>{{ copied ? 'Copiado!' : 'Copiar Link' }}</span>
        </button>
        <a
          :href="`/agendamento/${form.slug}`"
          target="_blank"
          class="px-3.5 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 text-xs font-semibold transition flex items-center gap-1.5"
        >
          <ExternalLink class="w-3.5 h-3.5" />
          <span>Abrir Página</span>
        </a>
      </div>
    </div>

    <!-- Formulário de Configurações -->
    <div class="glass-panel p-6 rounded-2xl border border-slate-800 space-y-6">
      <form class="space-y-4" @submit.prevent="saveSettings">
        <div v-if="successMessage" class="p-3 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs flex items-center gap-2">
          <CheckCircle2 class="w-4 h-4" />
          <span>{{ successMessage }}</span>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
            <input
              v-model="form.name"
              type="text"
              required
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
            <input
              v-model="form.document"
              type="text"
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
            <input
              v-model="form.phone"
              type="text"
              required
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">E-mail de Contato</label>
            <input
              v-model="form.email"
              type="email"
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Endereço Completo</label>
          <input
            v-model="form.address"
            type="text"
            placeholder="Rua, Número, Bairro"
            class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
          />
        </div>

        <div class="grid grid-cols-3 gap-4">
          <div class="col-span-2">
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Cidade</label>
            <input
              v-model="form.city"
              type="text"
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
            />
          </div>
          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">UF</label>
            <input
              v-model="form.state"
              type="text"
              maxlength="2"
              class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs uppercase focus:ring-2 focus:ring-emerald-500"
            />
          </div>
        </div>

        <div>
          <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Logo URL</label>
          <input
            v-model="form.logo_url"
            type="url"
            placeholder="https://..."
            class="w-full px-3.5 py-2 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs focus:ring-2 focus:ring-emerald-500"
          />
        </div>

        <div class="pt-4 border-t border-slate-800 flex justify-end">
          <button
            type="submit"
            :disabled="isSaving"
            class="px-6 py-2.5 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 text-xs font-bold transition shadow-md disabled:opacity-50"
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
  primary_color: '#10b981',
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
      form.primary_color = data.primary_color || '#10b981'
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
