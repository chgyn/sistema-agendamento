<template>
  <div class="min-h-screen bg-[#F8F9FD] flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden font-sans">
    <!-- Efeitos de Luz Ambiente -->
    <div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-96 h-96 bg-orange-500/5 rounded-full blur-[140px] pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-md relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-14 h-14 rounded-2xl bg-gradient-to-tr from-orange-500 to-amber-500 flex items-center justify-center shadow-md shadow-orange-500/25 transition transform hover:scale-105">
          <Scissors class="w-7 h-7 text-white font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl sm:text-3xl font-black tracking-tight text-gray-900 font-display">
        Painel de Gestão
      </h2>
      <p class="mt-1 text-center text-sm text-gray-500">
        Acesse sua barbearia ou salão de beleza
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md relative z-10 px-4">
      <div class="bg-white py-8 px-6 shadow-xl rounded-2xl sm:px-10 border border-gray-200/90">
        <form class="space-y-5" @submit.prevent="handleLogin">
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-2">E-mail</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
                <Mail class="w-4 h-4" />
              </div>
              <input
                v-model="email"
                type="email"
                required
                placeholder="exemplo@barbearia.com"
                class="w-full pl-10 pr-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
              />
            </div>
          </div>

          <div>
            <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-2">Senha</label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3.5 flex items-center pointer-events-none text-gray-400">
                <Lock class="w-4 h-4" />
              </div>
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••"
                class="w-full pl-10 pr-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
              />
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading"
            class="w-full flex justify-center items-center py-3.5 px-4 rounded-xl shadow-md shadow-orange-500/20 text-sm font-bold text-white bg-orange-500 hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500/30 transition active:scale-[0.98] disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Entrando...' : 'Entrar no Sistema' }}</span>
          </button>
        </form>

        <!-- Atalhos Demo de 1 Clique -->
        <div class="mt-8 pt-6 border-t border-gray-200/80 space-y-2.5">
          <p class="text-xs text-center text-gray-500 font-bold uppercase tracking-wider mb-2">
            Acesso Rápido de Demonstração (1 Clique):
          </p>
          <button
            type="button"
            @click="fillDemo('admin@plataforma.com', 'admin123')"
            class="w-full py-2.5 px-3 rounded-xl bg-purple-50 hover:bg-purple-100 border border-purple-200 text-purple-700 text-xs font-bold text-center transition flex items-center justify-center gap-1.5 shadow-sm cursor-pointer"
          >
            <span>👑 Administrador Geral da Plataforma</span>
          </button>
          <div class="grid grid-cols-2 gap-2">
            <button
              type="button"
              @click="fillDemo('admin@domnavalha.com', 'admin123')"
              class="py-2.5 px-3 rounded-xl bg-orange-50 hover:bg-orange-100 border border-orange-200 text-orange-700 text-xs font-bold text-center transition cursor-pointer"
            >
              💈 Dom Navalha
            </button>
            <button
              type="button"
              @click="fillDemo('admin@bellavista.com', 'admin123')"
              class="py-2.5 px-3 rounded-xl bg-pink-50 hover:bg-pink-100 border border-pink-200 text-pink-700 text-xs font-bold text-center transition cursor-pointer"
            >
              💇 Salão Bella Vista
            </button>
          </div>
        </div>
      </div>

      <p class="mt-6 text-center text-xs text-gray-500">
        Não tem uma conta ainda?
        <RouterLink to="/register" class="font-bold text-orange-600 hover:text-orange-500 hover:underline">
          Cadastre seu estabelecimento
        </RouterLink>
      </p>

      <p class="mt-4 text-center text-[11px] text-gray-400">
        Projeto Open-Source sob <a href="https://github.com/chgyn/sistema-agendamento/blob/main/LICENSE" target="_blank" rel="noopener noreferrer" class="hover:text-gray-600 underline">Licença MIT</a>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Scissors, Lock, Mail, Loader2, AlertCircle } from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'

const router = useRouter()
const authStore = useAuthStore()

const email = ref('')
const password = ref('')
const errorMessage = ref('')
const isLoading = ref(false)

function fillDemo(e: string, p: string) {
  email.value = e
  password.value = p
  errorMessage.value = ''
}

async function handleLogin() {
  errorMessage.value = ''
  isLoading.value = true

  try {
    await authStore.login(email.value, password.value)
    router.push('/admin/dashboard')
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || err.message || 'Erro ao realizar login. Verifique suas credenciais.'
  } finally {
    isLoading.value = false
  }
}
</script>
