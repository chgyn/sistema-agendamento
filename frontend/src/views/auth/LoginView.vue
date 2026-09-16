<template>
  <div class="min-h-screen bg-[#0d0e12] flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Efeitos de Luz Ambiente -->
    <div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-orange-500/10 rounded-full blur-[120px] pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-13 h-13 rounded-2xl bg-gradient-to-tr from-orange-500 to-amber-500 flex items-center justify-center shadow-glow-sm">
          <Scissors class="w-7 h-7 text-zinc-950 font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl sm:text-3xl font-black tracking-tight text-white font-display">
        Painel de Gestão
      </h2>
      <p class="mt-1 text-center text-sm text-zinc-400">
        Acesse sua barbearia ou salão de beleza
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md relative z-10 px-4">
      <div class="glass-panel py-8 px-6 shadow-2xl rounded-2xl sm:px-10 border border-zinc-800/80">
        <form class="space-y-5" @submit.prevent="handleLogin">
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <div>
            <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-2">E-mail</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-500">
                <Mail class="w-4 h-4" />
              </div>
              <input
                v-model="email"
                type="email"
                required
                placeholder="exemplo@barbearia.com"
                class="w-full pl-10 pr-4 py-3 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-2">Senha</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-zinc-500">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-10 pr-4 py-3 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center items-center py-3.5 px-4 rounded-xl shadow-glow text-sm font-bold text-zinc-950 bg-orange-500 hover:bg-orange-400 focus:outline-none focus:ring-2 focus:ring-orange-500 transition active:scale-[0.98] disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Entrando...' : 'Entrar no Sistema' }}</span>
          </button>
        </form>

        <!-- Atalhos Demo de 1 Clique -->
        <div class="mt-8 pt-6 border-t border-zinc-800/80 space-y-2.5">
          <p class="text-xs text-center text-zinc-400 font-bold uppercase tracking-wider mb-2">
            Acesso Rápido de Demonstração (1 Clique):
          </p>
          <button
            type="button"
            @click="fillDemo('admin@plataforma.com', 'admin123')"
            class="w-full py-2.5 px-3 rounded-xl bg-purple-950/60 hover:bg-purple-900/70 border border-purple-700/50 text-purple-200 text-xs font-bold text-center transition flex items-center justify-center gap-1.5 shadow-sm cursor-pointer"
          >
            <span>👑 Administrador Geral da Plataforma</span>
          </button>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="fillDemo('admin@domnavalha.com', 'admin123')"
              class="py-2.5 px-3 rounded-xl bg-orange-950/40 hover:bg-orange-900/50 border border-orange-800/40 text-orange-200 text-xs font-bold text-center transition cursor-pointer"
            >
              💈 Dom Navalha
            </button>
            <button
              type="button"
              @click="fillDemo('admin@bellavista.com', 'admin123')"
              class="py-2.5 px-3 rounded-xl bg-pink-950/40 hover:bg-pink-900/50 border border-pink-800/40 text-pink-200 text-xs font-bold text-center transition cursor-pointer"
            >
              💇 Salão Bella Vista
            </button>
          </div>
        </div>
      </div>

      <div class="mt-6 text-center text-xs text-zinc-400">
        Não tem uma conta ainda?
        <RouterLink to="/register" class="font-bold text-orange-400 hover:text-orange-300 underline ml-1">
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
