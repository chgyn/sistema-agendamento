<template>
  <div class="min-h-screen bg-slate-950 flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <div class="sm:mx-auto sm:w-full sm:max-w-3xl relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-12 h-12 rounded-2xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-glow">
          <Scissors class="w-6 h-6 text-slate-950 font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl sm:text-3xl font-black tracking-tight text-white">
        Cadastre seu Estabelecimento
      </h2>
      <p class="mt-1 text-center text-sm text-slate-400">
        Escolha o plano ideal e comece a receber agendamentos online em minutos
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-3xl relative z-10 px-4">
      <div class="glass-panel py-8 px-6 shadow-2xl rounded-2xl sm:px-10 border border-slate-800 space-y-6">
        <form class="space-y-6" @submit.prevent="handleRegister">
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- SEÇÃO 1: Escolha do Plano Comercial -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-emerald-400 uppercase tracking-wider flex items-center gap-1.5">
                <Sparkles class="w-3.5 h-3.5" />
                <span>1. Selecione o Plano de Assinatura *</span>
              </label>
              <span v-if="loadingPlans" class="text-xs text-slate-500 flex items-center gap-1">
                <Loader2 class="w-3 h-3 animate-spin" /> Carregando planos...
              </span>
            </div>

            <div v-if="plans.length === 0 && !loadingPlans" class="p-4 rounded-xl bg-slate-900 border border-slate-800 text-center text-xs text-slate-400">
              Nenhum plano disponível no momento. Entre em contato com o suporte.
            </div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-3.5">
              <div
                v-for="plan in plans"
                :key="plan.id"
                @click="form.plan_id = plan.id"
                class="p-4 rounded-xl border cursor-pointer transition-all duration-200 relative flex flex-col justify-between"
                :class="form.plan_id === plan.id ? 'bg-emerald-950/40 border-emerald-500 shadow-lg shadow-emerald-500/10 ring-1 ring-emerald-500' : 'bg-slate-900/80 border-slate-800 hover:border-slate-700'"
              >
                <!-- Radio check indicator -->
                <div class="flex items-start justify-between">
                  <span class="text-xs font-black text-white">{{ plan.name }}</span>
                  <div
                    class="w-4 h-4 rounded-full border flex items-center justify-center transition"
                    :class="form.plan_id === plan.id ? 'border-emerald-400 bg-emerald-500' : 'border-slate-600'"
                  >
                    <div v-if="form.plan_id === plan.id" class="w-1.5 h-1.5 rounded-full bg-slate-950 font-bold"></div>
                  </div>
                </div>

                <div class="my-2.5">
                  <span class="text-xs text-slate-400">R$</span>
                  <span class="text-2xl font-black text-white ml-0.5">{{ plan.price.toFixed(2).replace('.', ',') }}</span>
                  <span class="text-[10px] text-slate-400 block">/ {{ formatCycle(plan.billing_cycle) }}</span>
                </div>

                <p v-if="plan.description" class="text-[11px] text-slate-400 leading-tight line-clamp-2">
                  {{ plan.description }}
                </p>
              </div>
            </div>
          </div>

          <!-- SEÇÃO 2: Dados do Estabelecimento -->
          <div class="space-y-4 pt-4 border-t border-slate-800">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-300 flex items-center gap-1.5">
              <Scissors class="w-3.5 h-3.5 text-emerald-400" />
              <span>2. Dados do Estabelecimento</span>
            </h4>

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

              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Cidade</label>
                <input
                  v-model="form.city"
                  type="text"
                  placeholder="São Paulo"
                  class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Estado (UF)</label>
                <input
                  v-model="form.state"
                  type="text"
                  placeholder="SP"
                  maxlength="2"
                  class="w-full px-3.5 py-2 rounded-xl bg-slate-900/90 border border-slate-700/80 text-white placeholder-slate-500 text-sm uppercase focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                />
              </div>
            </div>
          </div>

          <!-- SEÇÃO 3: Dados do Administrador -->
          <div class="space-y-4 pt-4 border-t border-slate-800">
            <h4 class="text-xs font-bold uppercase tracking-wider text-slate-300 flex items-center gap-1.5">
              <ShieldCheck class="w-3.5 h-3.5 text-emerald-400" />
              <span>3. Dados de Acesso do Administrador</span>
            </h4>
            
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
            :disabled="isLoading || !form.plan_id"
            class="w-full mt-6 flex justify-center items-center py-3.5 px-4 rounded-xl shadow-lg text-sm font-bold text-slate-950 bg-emerald-500 hover:bg-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition disabled:opacity-50"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Criando Conta e Assinatura...' : 'Concluir Cadastro e Iniciar' }}</span>
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
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Scissors, AlertCircle, Loader2, Sparkles, ShieldCheck } from 'lucide-vue-next'
import { useAuthStore, type Plan, type PlanBillingCycle } from '../../stores/auth'
import api from '../../services/api'

const router = useRouter()
const authStore = useAuthStore()

const plans = ref<Plan[]>([])
const loadingPlans = ref(true)

const form = reactive({
  plan_id: '',
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

function formatCycle(cycle: PlanBillingCycle) {
  const map: Record<PlanBillingCycle, string> = {
    MONTHLY: 'mês',
    QUARTERLY: 'trimestre',
    SEMIANNUALLY: 'semestre',
    YEARLY: 'ano',
  }
  return map[cycle] || 'mês'
}

function generateSlug() {
  form.slug = form.tenant_name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/(^-|-$)+/g, '')
}

async function fetchPlans() {
  loadingPlans.value = true
  try {
    const res = await api.get('/public/plans')
    if (res.data.success && res.data.data.length > 0) {
      plans.value = res.data.data
      // Seleciona o primeiro ou o plano intermediário por padrão
      const defaultPlan = plans.value.find(p => p.sort_order === 2) || plans.value[0]
      form.plan_id = defaultPlan.id
    }
  } catch (err: any) {
    console.error('Erro ao buscar planos:', err)
  } finally {
    loadingPlans.value = false
  }
}

async function handleRegister() {
  if (!form.plan_id) {
    errorMessage.value = 'Por favor, selecione um plano comercial para continuar.'
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    await authStore.registerTenant(form)
    router.push('/admin/minha-assinatura')
  } catch (err: any) {
    errorMessage.value = err.message || 'Erro ao registrar estabelecimento'
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchPlans()
})
</script>
