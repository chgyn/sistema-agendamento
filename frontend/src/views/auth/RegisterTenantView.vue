<template>
  <div class="min-h-screen bg-[#0d0e12] flex flex-col justify-center py-12 sm:px-6 lg:px-8 relative overflow-hidden">
    <!-- Efeito Luz Ambiente -->
    <div class="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[350px] bg-orange-500/10 rounded-full blur-[130px] pointer-events-none"></div>

    <div class="sm:mx-auto sm:w-full sm:max-w-3xl relative z-10">
      <div class="flex justify-center">
        <RouterLink to="/" class="w-13 h-13 rounded-2xl bg-gradient-to-tr from-orange-500 to-amber-500 flex items-center justify-center shadow-glow-sm">
          <Scissors class="w-7 h-7 text-zinc-950 font-bold" />
        </RouterLink>
      </div>
      <h2 class="mt-4 text-center text-2xl sm:text-3xl font-black tracking-tight text-white font-display">
        Cadastre seu Estabelecimento
      </h2>
      <p class="mt-1 text-center text-sm text-zinc-400">
        Escolha o plano ideal e comece a receber agendamentos online em minutos
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-3xl relative z-10 px-4">
      <div class="glass-panel py-8 px-6 shadow-2xl rounded-2xl sm:px-10 border border-zinc-800/80 space-y-6">
        <form class="space-y-6" @submit.prevent="handleRegister">
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-xs flex items-center gap-2">
            <AlertCircle class="w-4 h-4 shrink-0" />
            <span>{{ errorMessage }}</span>
          </div>

          <!-- SEÇÃO 1: Escolha do Plano Comercial -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <label class="text-xs font-bold text-orange-400 uppercase tracking-wider flex items-center gap-1.5">
                <Sparkles class="w-3.5 h-3.5" />
                <span>1. Selecione o Plano de Assinatura *</span>
              </label>
              <span v-if="loadingPlans" class="text-xs text-zinc-500 flex items-center gap-1">
                <Loader2 class="w-3 h-3 animate-spin text-orange-500" /> Carregando planos...
              </span>
            </div>

            <div v-if="plans.length === 0 && !loadingPlans" class="p-4 rounded-xl bg-[#14151c] border border-zinc-800 text-center text-xs text-zinc-400">
              Nenhum plano disponível no momento. Entre em contato com o suporte.
            </div>

            <div v-else class="grid grid-cols-1 sm:grid-cols-3 gap-3.5">
              <div
                v-for="plan in plans"
                :key="plan.id"
                @click="form.plan_id = plan.id"
                class="p-4 rounded-xl border cursor-pointer transition-all duration-200 relative flex flex-col justify-between"
                :class="form.plan_id === plan.id ? 'bg-orange-950/40 border-orange-500 shadow-glow-sm ring-1 ring-orange-500' : 'bg-[#14151c]/80 border-zinc-800 hover:border-zinc-700'"
              >
                <!-- Radio check indicator & Badge Gratuito -->
                <div>
                  <div class="flex items-start justify-between gap-1.5">
                    <div>
                      <span v-if="plan.is_free" class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-emerald-500/20 text-emerald-400 border border-emerald-500/30 inline-block mb-1">
                        Gratuito
                      </span>
                      <h4 class="text-xs font-black text-white font-display block">{{ plan.name }}</h4>
                    </div>
                    <div
                      class="w-4 h-4 rounded-full border flex items-center justify-center transition shrink-0 mt-0.5"
                      :class="form.plan_id === plan.id ? 'border-orange-400 bg-orange-500' : 'border-zinc-600'"
                    >
                      <div v-if="form.plan_id === plan.id" class="w-1.5 h-1.5 rounded-full bg-zinc-950 font-bold"></div>
                    </div>
                  </div>

                  <div class="my-2.5">
                    <template v-if="plan.is_free">
                      <span class="text-xl font-black text-emerald-400 font-display">100% Grátis</span>
                      <span class="text-[10px] text-zinc-400 block">Sem mensalidade ou Asaas</span>
                    </template>
                    <template v-else>
                      <span class="text-xs text-zinc-400">R$</span>
                      <span class="text-2xl font-black text-white ml-0.5 font-display">{{ plan.price.toFixed(2).replace('.', ',') }}</span>
                      <span class="text-[10px] text-zinc-400 block">/ {{ formatCycle(plan.billing_cycle) }}</span>
                    </template>
                  </div>
                </div>

                <p v-if="plan.description" class="text-[11px] text-zinc-400 leading-tight line-clamp-2">
                  {{ plan.description }}
                </p>
              </div>
            </div>
          </div>

          <!-- SEÇÃO 2: Dados do Estabelecimento -->
          <div class="space-y-4 pt-4 border-t border-zinc-800/80">
            <h4 class="text-xs font-bold uppercase tracking-wider text-zinc-300 flex items-center gap-1.5">
              <Scissors class="w-3.5 h-3.5 text-orange-400" />
              <span>2. Dados do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
                <input
                  v-model="form.tenant_name"
                  @input="generateSlug"
                  type="text"
                  required
                  placeholder="Ex: Barbearia Imperial"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Link de Agendamento (Slug) *</label>
                <div class="flex items-center rounded-xl bg-[#14151c] border border-zinc-750 px-3 py-2.5 text-sm text-zinc-400">
                  <span class="text-xs text-zinc-500">/agendamento/</span>
                  <input
                    v-model="form.slug"
                    type="text"
                    required
                    placeholder="barbearia-imperial"
                    class="w-full bg-transparent text-orange-400 font-mono text-sm focus:outline-none"
                  />
                </div>
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
                <input
                  v-model="form.phone"
                  type="text"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
                <input
                  v-model="form.document"
                  type="text"
                  placeholder="Opcional"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Cidade</label>
                <input
                  v-model="form.city"
                  type="text"
                  placeholder="São Paulo"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Estado (UF)</label>
                <input
                  v-model="form.state"
                  type="text"
                  placeholder="SP"
                  maxlength="2"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm uppercase focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>
            </div>
          </div>

          <!-- SEÇÃO 3: Dados do Administrador -->
          <div class="space-y-4 pt-4 border-t border-zinc-800/80">
            <h4 class="text-xs font-bold uppercase tracking-wider text-zinc-300 flex items-center gap-1.5">
              <ShieldCheck class="w-3.5 h-3.5 text-orange-400" />
              <span>3. Dados de Acesso do Administrador</span>
            </h4>

            <div class="space-y-3">
              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Seu Nome Completo *</label>
                <input
                  v-model="form.admin_name"
                  type="text"
                  required
                  placeholder="Nome do proprietário ou gerente"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                />
              </div>

              <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                <div>
                  <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">E-mail de Login *</label>
                  <input
                    v-model="form.admin_email"
                    type="email"
                    required
                    placeholder="seuemail@exemplo.com"
                    class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                  />
                </div>

                <div>
                  <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Senha (Mín. 6 caracteres) *</label>
                  <input
                    v-model="form.password"
                    type="password"
                    required
                    minlength="6"
                    placeholder="••••••••"
                    class="w-full px-3.5 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
                  />
                </div>
              </div>
            </div>
          </div>

          <button
            type="submit"
            :disabled="isLoading || !form.plan_id"
            class="w-full mt-6 flex justify-center items-center py-4 px-4 rounded-xl shadow-glow text-sm font-bold text-zinc-950 bg-orange-500 hover:bg-orange-400 focus:outline-none focus:ring-2 focus:ring-orange-500 transition active:scale-[0.98] disabled:opacity-50 cursor-pointer"
          >
            <Loader2 v-if="isLoading" class="w-4 h-4 animate-spin mr-2" />
            <span>{{ isLoading ? 'Criando Conta e Assinatura...' : 'Concluir Cadastro e Iniciar' }}</span>
          </button>
        </form>
      </div>

      <div class="mt-6 text-center text-xs text-zinc-400">
        Já possui cadastro?
        <RouterLink to="/login" class="font-bold text-orange-400 hover:text-orange-300 underline ml-1">
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
    const selectedPlan = plans.value.find(p => p.id === form.plan_id)
    await authStore.registerTenant(form)

    // Se plano gratuito, o estabelecimento já está 100% ativo e vai direto para o dashboard
    if (selectedPlan?.is_free || authStore.subscription?.status === 'ACTIVE') {
      router.push('/admin/dashboard')
    } else {
      router.push('/admin/minha-assinatura')
    }
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
