<template>
  <div class="min-h-screen bg-slate-950 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <div class="sm:mx-auto sm:w-full sm:max-w-xl relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-glow">
          <Scissors class="w-6 h-6 text-slate-950 font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl font-bold tracking-tight text-white">
        Cadastre sua Barbearia ou Salão
      </h2>
      <p class="mt-1 text-center text-sm text-slate-400">
        Comece a receber agendamentos online em minutos
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-xl relative z-10 px-4">
      <div class="glass-panel py-8 px-6 shadow-2xl rounded-2xl sm:px-10 border border-slate-800">
        <form class="space-y-4" @submit.prevent="handleRegister">
          <div v-if="errorMessage" class="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-400 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
              <input
                v-model="form.tenant_name"
                @input="generateSlug"
                type="text"
                required
                placeholder="Ex: Barbearia Imperial"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Link de Agendamento (Slug) *</label>
              <div class="flex items-center rounded-xl bg-slate-900/90 border border-slate-700/80 px-3 py-2 text-sm text-slate-400">
                <span class="text-xs text-slate-500">/agendamento/</span>
                <input
                  v-model="form.slug"
                  type="text"
                  required
                  placeholder="barbearia-imperial"
                  class="w-full bg-transparent text-emerald-400 font-mono text-sm focus:outline-none"
                />
              </div>
            </div>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
              <input
                v-model="form.phone"
                type="text"
                required
                placeholder="(11) 98765-4321"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
              <input
                v-model="form.document"
                type="text"
                placeholder="Opcional"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>
          </div>

          <div class="grid grid-cols-3 gap-4">
            <div class="col-span-2">
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Cidade</label>
              <input
                v-model="form.city"
                type="text"
                placeholder="São Paulo"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>
            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">UF</label>
              <input
                v-model="form.state"
                type="text"
                placeholder="SP"
                maxlength="2"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm uppercase focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>
          </div>

          <div class="pt-3 border-t border-slate-800">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-400 mb-3">Dados de Acesso do Administrador</h4>
            
            <div class="space-y-3">
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Seu Nome Completo *</label>
                <input
                  v-model="form.admin_name"
                  type="text"
                  required
                  placeholder="Nome do proprietário ou gerente"
                  class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-medium text-slate-300 mb-1">E-mail de Login *</label>
                  <input
                    v-model="form.admin_email"
                    type="email"
                    required
                    placeholder="seuemail@exemplo.com"
                    class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                  />
                </div>

                <div>
                  <label class="block text-xs font-medium text-slate-300 mb-1">Senha (Mín. 6 caracteres) *</label>
                  <input
                    v-model="form.password"
                    type="password"
                    required
                    minlength="6"
                    placeholder="••••••••"
                    class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                  />
                </div>
              </div>
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full mt-4 flex justify-center items-center py-3 px-4 rounded-xl shadow-md text-sm font-semibold text-slate-950 bg-emerald-500 hover:bg-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition disabled:opacity-50"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Cadastrando...' : 'Criar Minha Conta e Iniciar' }}</span>
          </button>
        </form>
      </div>

      <div class="mt-6 text-center text-xs text-slate-400">
        Já possui cadastro?
        <RouterLink to="/login" class="font-medium text-emerald-400 hover:text-emerald-300 underline ml-1">
          Fazer Login
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Scissors, AlertCircle, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const form = reactive({
  tenant_name: '',
  slug: '',
  document: '',
  phone: '',
  city: '',
  state: '',
  admin_name: '',
  admin_email: '',
  password: '',
})

const isLoading = ref(false)
const errorMessage = ref('')

function generateSlug() {
  form.slug = form.tenant_name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)+/g, '')
}

async function handleRegister() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.registerTenant(form)
    router.push('/admin/dashboard')
  } catch (err: any) {
    errorMessage.value = err.message || 'Erro ao registrar estabelecimento'
  } finally {
    isLoading.value = false
  }
}
</script>
