<template>
  <div class="space-y-6 max-w-5xl mx-auto">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight flex items-center gap-2.5 font-display">
          <CreditCard class="w-7 h-7 text-orange-400" />
          <span>Minha Assinatura & Cobrança</span>
        </h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1">
          Acompanhe o status do seu plano contratado, faturas geradas e opções de pagamento.
        </p>
      </div>

      <button
        @click="fetchSubscription"
        class="inline-flex items-center gap-2 px-3.5 py-2.5 rounded-xl bg-[#181922] border border-zinc-700/60 hover:bg-zinc-800 text-zinc-300 hover:text-white text-xs font-bold transition cursor-pointer"
      >
        <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loading }" />
        <span>Atualizar Status</span>
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-16">
      <Loader2 class="w-8 h-8 text-orange-500 animate-spin mx-auto mb-3" />
      <p class="text-sm text-zinc-400">Carregando detalhes da assinatura...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="!subscription" class="bg-[#14151c] border border-zinc-800 rounded-2xl p-12 text-center">
      <AlertTriangle class="w-12 h-12 text-amber-400 mx-auto mb-3" />
      <h3 class="text-base font-bold text-white font-display">Nenhuma assinatura ativa encontrada</h3>
      <p class="text-xs text-zinc-400 mt-1 max-w-sm mx-auto">
        Entre em contato com o suporte ou selecione um plano comercial para ativar seu estabelecimento.
      </p>
    </div>

    <div v-else class="space-y-6">
      <!-- Status Alert Banner -->
      <div
        class="p-5 rounded-2xl border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 shadow-lg transition-all"
        :class="statusBannerClass"
      >
        <div class="flex items-center gap-3.5">
          <div class="w-11 h-11 rounded-2xl flex items-center justify-center shrink-0" :class="statusIconContainerClass">
            <CheckCircle2 v-if="subscription.status === 'ACTIVE'" class="w-6 h-6 text-emerald-400" />
            <Clock v-else-if="subscription.status === 'PENDING'" class="w-6 h-6 text-amber-400 animate-pulse" />
            <AlertCircle v-else class="w-6 h-6 text-rose-400" />
          </div>
          <div>
            <div class="flex items-center gap-2 flex-wrap">
              <h3 class="font-bold text-white text-base font-display">
                Status da Assinatura: {{ statusLabel }}
              </h3>
              <span class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider border" :class="statusBadgeClass">
                {{ subscription.status }}
              </span>
              <span
                v-if="subscription.origin === 'MANUAL'"
                class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-purple-500/10 text-purple-300 border border-purple-500/30"
              >
                Concessão Manual (Admin)
              </span>
              <span
                v-else-if="subscription.origin === 'FREE_PLAN'"
                class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-emerald-500/10 text-emerald-300 border border-emerald-500/30"
              >
                Plano Gratuito
              </span>
              <span
                v-else
                class="px-2.5 py-0.5 rounded text-[10px] font-black uppercase tracking-wider bg-blue-500/10 text-blue-300 border border-blue-500/30"
              >
                Asaas Recorrente
              </span>
            </div>
            <p class="text-xs text-zinc-300 mt-1">
              {{ statusDescription }}
            </p>
          </div>
        </div>

        <!-- Botão de Ação Imediata (Pagar Fatura) - Apenas se plano Asaas -->
        <a
          v-if="subscription.payment_url && subscription.status !== 'ACTIVE' && subscription.origin === 'ASAAS'"
          :href="subscription.payment_url"
          target="_blank"
          class="w-full sm:w-auto inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-orange-500 hover:bg-orange-400 text-zinc-950 font-bold text-xs shadow-glow-sm transition shrink-0"
        >
          <ExternalLink class="w-4 h-4" />
          <span>Efetuar Pagamento da Fatura</span>
        </a>
      </div>

      <!-- Detalhes do Plano Contratado -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-5">
        <!-- Card 1: Informações do Plano -->
        <div class="glass-panel border border-zinc-800 rounded-2xl p-6 space-y-4 md:col-span-2">
          <div class="flex items-center justify-between pb-3.5 border-b border-zinc-800">
            <div class="flex items-center gap-2.5">
              <div class="w-8 h-8 rounded-xl bg-orange-500/10 border border-orange-500/20 flex items-center justify-center text-orange-400">
                <Sparkles class="w-4 h-4" />
              </div>
              <div>
                <h3 class="font-bold text-white text-sm font-display">Plano Contratado</h3>
                <p class="text-xs text-zinc-400">Configuração comercial vigente</p>
              </div>
            </div>
            <span class="px-3 py-1 rounded-full text-xs font-bold bg-orange-950/80 text-orange-300 border border-orange-700/60 font-display">
              {{ subscription.plan?.name || 'Plano Personalizado' }}
            </span>
          </div>

          <div class="grid grid-cols-2 sm:grid-cols-3 gap-3.5 text-xs">
            <div class="p-3.5 rounded-xl bg-[#14151c] border border-zinc-800">
              <span class="text-zinc-400 block font-semibold">Valor Recorrente</span>
              <span class="text-lg font-black text-white mt-1 block font-display">
                {{ subscription.origin === 'FREE_PLAN' || subscription.price === 0 ? 'Gratuito (R$ 0,00)' : `R$ ${subscription.price.toFixed(2).replace('.', ',')}` }}
              </span>
            </div>

            <div class="p-3.5 rounded-xl bg-[#14151c] border border-zinc-800">
              <span class="text-zinc-400 block font-semibold">Periodicidade</span>
              <span class="text-sm font-bold text-white mt-1 block">
                {{ subscription.origin === 'FREE_PLAN' ? 'Sem cobrança' : formatCycle(subscription.billing_cycle) }}
              </span>
            </div>

            <div class="p-3.5 rounded-xl bg-[#14151c] border border-zinc-800 col-span-2 sm:col-span-1">
              <span class="text-zinc-400 block font-semibold">
                {{ subscription.origin === 'MANUAL' ? 'Validade da Liberação' : 'Próximo Vencimento' }}
              </span>
              <span class="text-sm font-bold text-orange-400 mt-1 block font-display">
                {{ subscription.current_period_end ? formatDate(subscription.current_period_end) : (subscription.origin === 'FREE_PLAN' ? 'Acesso Contínuo' : formatDate(subscription.next_due_date)) }}
              </span>
            </div>
          </div>

          <div v-if="subscription.plan?.description" class="text-xs text-zinc-400 leading-relaxed pt-1">
            {{ subscription.plan.description }}
          </div>
        </div>

        <!-- Card 2: Modalidade / Integração Asaas -->
        <div class="glass-panel border border-zinc-800 rounded-2xl p-6 space-y-3 flex flex-col justify-between">
          <div>
            <div class="flex items-center gap-2 text-zinc-300 pb-3 border-b border-zinc-800">
              <ShieldCheck class="w-4 h-4 text-orange-400" />
              <h4 class="text-xs font-bold uppercase tracking-wider">
                {{ subscription.origin === 'ASAAS' ? 'Identificador Asaas' : 'Modalidade de Acesso' }}
              </h4>
            </div>

            <div class="mt-3 space-y-2.5 text-xs">
              <template v-if="subscription.origin === 'MANUAL'">
                <div>
                  <span class="text-zinc-400 block font-semibold">Origem do Contrato:</span>
                  <span class="text-purple-300 font-bold block mt-0.5">Liberado Manualmente pelo Suporte</span>
                </div>
                <div v-if="subscription.manual_grant_reason">
                  <span class="text-zinc-400 block font-semibold">Justificativa:</span>
                  <span class="text-zinc-200 text-xs block mt-0.5">{{ subscription.manual_grant_reason }}</span>
                </div>
                <div>
                  <span class="text-zinc-400 block font-semibold">Status de Cobrança:</span>
                  <span class="text-emerald-400 font-bold block mt-0.5">Isento de faturas automáticas</span>
                </div>
              </template>

              <template v-else-if="subscription.origin === 'FREE_PLAN'">
                <div>
                  <span class="text-zinc-400 block font-semibold">Origem do Contrato:</span>
                  <span class="text-emerald-400 font-bold block mt-0.5">Plano Gratuito da Plataforma</span>
                </div>
                <div>
                  <span class="text-zinc-400 block font-semibold">Status de Cobrança:</span>
                  <span class="text-zinc-300 block mt-0.5">100% gratuito, sem gateway de pagamento</span>
                </div>
              </template>

              <template v-else>
                <div>
                  <span class="text-zinc-400 block font-semibold">ID Assinatura:</span>
                  <span class="font-mono text-[11px] text-zinc-200 truncate block mt-0.5">
                    {{ subscription.asaas_subscription_id || 'sub_mock_local' }}
                  </span>
                </div>
                <div>
                  <span class="text-zinc-400 block font-semibold">ID Cliente:</span>
                  <span class="font-mono text-[11px] text-zinc-200 truncate block mt-0.5">
                    {{ subscription.asaas_customer_id || 'cus_mock_local' }}
                  </span>
                </div>
              </template>
            </div>
          </div>

          <div v-if="subscription.origin === 'ASAAS' && subscription.payment_url" class="pt-2">
            <a
              :href="subscription.payment_url"
              target="_blank"
              class="w-full inline-flex items-center justify-center gap-2 py-2.5 px-3 rounded-xl bg-[#181922] hover:bg-zinc-800 text-orange-400 text-xs font-bold transition border border-orange-500/30"
            >
              <ExternalLink class="w-3.5 h-3.5" />
              <span>Abrir Portal de Faturas</span>
            </a>
          </div>
        </div>
      </div>

      <!-- Histórico de Faturas e Pagamentos -->
      <div class="glass-panel border border-zinc-800 rounded-2xl overflow-hidden shadow-xl">
        <div class="p-5 border-b border-zinc-800 flex items-center justify-between">
          <div class="flex items-center gap-2.5">
            <FileText class="w-4 h-4 text-orange-400" />
            <h3 class="font-bold text-white text-sm font-display">Histórico de Cobranças Recorrentes</h3>
          </div>
          <span class="text-xs text-zinc-400">Sincronizado via Webhook</span>
        </div>

        <div v-if="!subscription.invoices || subscription.invoices.length === 0" class="p-8 text-center text-xs text-zinc-400">
          Nenhuma fatura registrada no histórico até o momento. As cobranças aparecerão aqui assim que forem geradas no Asaas.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="bg-[#121318]/90 border-b border-zinc-800 text-xs text-zinc-400 uppercase tracking-wider font-bold">
              <tr>
                <th class="py-3.5 px-4">Fatura</th>
                <th class="py-3.5 px-4">Vencimento</th>
                <th class="py-3.5 px-4">Valor</th>
                <th class="py-3.5 px-4">Método</th>
                <th class="py-3.5 px-4">Status</th>
                <th class="py-3.5 px-4 text-right">Comprovante / Ação</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-zinc-800/60 text-xs">
              <tr v-for="inv in subscription.invoices" :key="inv.id" class="hover:bg-[#16171e]/80 transition">
                <td class="py-3.5 px-4 font-mono text-zinc-300">
                  {{ inv.asaas_payment_id }}
                </td>
                <td class="py-3.5 px-4 text-zinc-300 font-display">
                  {{ formatDate(inv.due_date) }}
                </td>
                <td class="py-3.5 px-4 font-bold text-white font-display">
                  R$ {{ inv.value.toFixed(2).replace('.', ',') }}
                </td>
                <td class="py-3.5 px-4 text-zinc-400 uppercase">
                  {{ inv.billing_type || 'PIX / Boleto' }}
                </td>
                <td class="py-3.5 px-4">
                  <span
                    class="px-2.5 py-0.5 rounded-full text-[10px] font-bold border"
                    :class="getInvoiceBadgeClass(inv.status)"
                  >
                    {{ inv.status }}
                  </span>
                </td>
                <td class="py-3.5 px-4 text-right">
                  <a
                    v-if="inv.invoice_url"
                    :href="inv.invoice_url"
                    target="_blank"
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-200 hover:text-white transition font-bold border border-zinc-750"
                  >
                    <ExternalLink class="w-3 h-3" />
                    <span>Ver Fatura</span>
                  </a>
                  <span v-else class="text-zinc-500">-</span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  CreditCard, RefreshCw, Loader2, AlertTriangle, CheckCircle2,
  Clock, AlertCircle, ExternalLink, Sparkles, ShieldCheck, FileText
} from 'lucide-vue-next'
import api from '../../services/api'
import type { Subscription, PlanBillingCycle } from '../../stores/auth'

const subscription = ref<Subscription | null>(null)
const loading = ref(true)

const statusLabel = computed(() => {
  switch (subscription.value?.status) {
    case 'ACTIVE':
      return 'Ativa & Liberada'
    case 'PENDING':
      return 'Aguardando Pagamento'
    case 'OVERDUE':
      return 'Vencida / Pendente'
    case 'CANCELLED':
      return 'Cancelada'
    case 'TRIAL':
      return 'Período de Testes (Trial)'
    default:
      return subscription.value?.status || 'Desconhecido'
  }
})

const statusDescription = computed(() => {
  if (subscription.value?.origin === 'MANUAL') {
    return 'Assinatura concedida manualmente pela administração da plataforma. Recursos liberados sem cobranças recorrentes.'
  }
  if (subscription.value?.origin === 'FREE_PLAN') {
    return 'Plano gratuito ativo. Você possui acesso aos recursos inclusos sem necessidade de mensalidade ou gateway de pagamento.'
  }
  switch (subscription.value?.status) {
    case 'ACTIVE':
      return 'Todos os recursos de agendamento, agenda e profissionais estão disponíveis normalmente.'
    case 'PENDING':
      return 'Sua conta foi criada! Realize o pagamento da fatura para desbloquear o acesso total aos agendamentos.'
    case 'OVERDUE':
      return 'Identificamos uma fatura pendente. Regularize o pagamento para evitar o bloqueio da agenda.'
    case 'CANCELLED':
      return 'Sua assinatura foi descontinuada. Entre em contato com a equipe para reativação.'
    default:
      return ''
  }
})

const statusBannerClass = computed(() => {
  switch (subscription.value?.status) {
    case 'ACTIVE':
      return 'bg-emerald-950/30 border-emerald-800/60'
    case 'PENDING':
      return 'bg-amber-950/30 border-amber-800/60'
    case 'OVERDUE':
      return 'bg-rose-950/30 border-rose-800/60'
    default:
      return 'bg-[#14151c] border-zinc-800'
  }
})

const statusIconContainerClass = computed(() => {
  switch (subscription.value?.status) {
    case 'ACTIVE':
      return 'bg-emerald-500/15 border border-emerald-500/30'
    case 'PENDING':
      return 'bg-amber-500/15 border border-amber-500/30'
    default:
      return 'bg-rose-500/15 border border-rose-500/30'
  }
})

const statusBadgeClass = computed(() => {
  switch (subscription.value?.status) {
    case 'ACTIVE':
      return 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60'
    case 'PENDING':
      return 'bg-amber-950/80 text-amber-300 border-amber-700/60'
    default:
      return 'bg-rose-950/80 text-rose-300 border-rose-700/60'
  }
})

function formatCycle(cycle?: PlanBillingCycle) {
  if (!cycle) return 'Mensal'
  const map: Record<PlanBillingCycle, string> = {
    MONTHLY: 'Mensal',
    QUARTERLY: 'Trimestral',
    SEMIANNUALLY: 'Semestral',
    YEARLY: 'Anual',
  }
  return map[cycle] || cycle
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('pt-BR')
  } catch {
    return dateStr
  }
}

function getInvoiceBadgeClass(status: string) {
  switch (status) {
    case 'CONFIRMED':
    case 'RECEIVED':
      return 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60'
    case 'PENDING':
      return 'bg-amber-950/80 text-amber-300 border-amber-700/60'
    default:
      return 'bg-rose-950/80 text-rose-300 border-rose-700/60'
  }
}

async function fetchSubscription() {
  loading.value = true
  try {
    const res = await api.get('/admin/subscription')
    if (res.data.success) {
      subscription.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao carregar assinatura:', err)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchSubscription()
})
</script>
