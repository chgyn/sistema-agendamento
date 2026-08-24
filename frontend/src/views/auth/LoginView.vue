<template>
  <div class="min-h-screen bg-slate-950 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Efeitos de Fundo -->
    <div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-emerald-500/10 rounded-full blur-3xl pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-glow">
          <Scissors class="w-6 h-6 text-slate-950 font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl font-bold tracking-tight text-white">
        Painel de Gestão
      </h2>
      <p class="mt-1 text-center text-sm text-slate-400">
        Acesse sua barbearia ou salão de beleza
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md relative z-10 px-4">
      <div class="glass-panel py-8 px-6 shadow-2xl rounded-2xl sm:px-10 border border-slate-800">
        <form class="space-y-5" @submit.prevent="handleLogin">
          <div v-if="errorMessage" class="p-3 rounded-xl bg-red-500/10 border border-red-500/30 text-red-400 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">E-mail</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-500">
                <Mail class="w-4 h-4" />
              </div>
              <input
                v-model="email"
                type="email"
                required
                placeholder="exemplo@barbearia.com"
                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-2">Senha</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-slate-500">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:border-transparent transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center items-center py-3 px-4 rounded-xl shadow-md text-sm font-semibold text-slate-950 bg-emerald-500 hover:bg-emerald-400 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-emerald-500 transition disabled:opacity-50"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Entrando...' : 'Entrar no Sistema' }}</span>
          </button>
        </form>

        <!-- Atalhos Demo de 1 Clique -->
        <div class="mt-8 pt-6 border-t border-slate-800">
          <p class="text-xs text-center text-slate-400 font-medium mb-3">
            Acesso Rápido de Demonstração (1 Clique):
          </p>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="fillDemo('admin@domnavalha.com', 'admin123')"
              class="py-2 px-3 rounded-lg bg-emerald-950/40 hover:bg-emerald-900/50 border border-emerald-800/40 text-emerald-300 text-xs font-medium text-center transition"
            >
              💈 Barbearia Dom Navalha
            </button>
            <button
              type="button"
              @click="fillDemo('admin@bellavista.com', 'admin123')"
              class="py-2 px-3 rounded-lg bg-pink-950/40 hover:bg-pink-900/50 border border-pink-800/40 text-pink-300 text-xs font-medium text-center transition"
            >
              💇 Salão Bella Vista
            </button>
          </div>
        </div>
      </div>

      <div class="mt-6 text-center text-xs text-slate-400">
        Não tem uma conta ainda?
        <RouterLink to="/register" class="font-medium text-emerald-400 hover:text-emerald-300 underline ml-1">
          Cadastre seu estabelecimento
        </RouterLink>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Scissors, Mail, Lock, AlertCircle, Loader2 } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('admin@domnavalha.com')
const password = ref('admin123')
const isLoading = ref(false)
const errorMessage = ref('')

function fillDemo(demoEmail: string, demoPass: string) {
  email.value = demoEmail
  password.value = demoPass
}

async function handleLogin() {
  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.login(email.value, password.value)
    router.push('/admin/dashboard')
  } catch (err: any) {
    errorMessage.value = err.message || 'Credenciais inválidas'
  } finally {
    isLoading.value = false
  }
}
</script>
