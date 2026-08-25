<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <CreditCard class="w-7 h-7 text-purple-400" />
          <span>Planos de Assinatura</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Gerencie a grade comercial de planos, periodicidades e limites para contratação dos estabelecimentos.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-600/30 transition active:scale-95 shrink-0"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Plano Comercial</span>
      </button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400">
          <Package class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Total de Planos</p>
          <p class="text-2xl font-black text-white mt-0.5">{{ plans.length }}</p>
        </div>
      </div>

      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400">
          <CheckCircle2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Planos Ativos para Venda</p>
          <p class="text-2xl font-black text-emerald-400 mt-0.5">{{ activeCount }}</p>
        </div>
      </div>

      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400">
          <Layers class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Integração Asaas Recorrente</p>
          <p class="text-sm font-bold text-indigo-300 mt-1 flex items-center gap-1.5">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>API v3 Conectada</span>
          </p>
        </div>
      </div>
    </div>

    <!-- Filtros e Busca -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex flex-col sm:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome ou descrição do plano..."
          class="w-full pl-9 pr-4 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-purple-500 transition"
        />
      </div>

      <div class="flex items-center gap-2 w-full sm:w-auto">
        <select
          v-model="statusFilter"
          class="w-full sm:w-44 py-2 px-3 rounded-xl bg-slate-950 border border-slate-800 text-sm text-slate-200 focus:outline-none focus:border-purple-500 transition"
        >
          <option value="ALL">Todos os status</option>
          <option value="ACTIVE">Apenas Ativos</option>
          <option value="INACTIVE">Apenas Inativos</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-16">
      <Loader2 class="w-8 h-8 text-purple-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-slate-400">Carregando planos de assinatura...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredPlans.length === 0" class="bg-slate-900/50 border border-slate-800/80 rounded-2xl p-12 text-center">
      <Package class="w-12 h-12 text-slate-600 mx-auto mb-3" />
      <h3 class="text-base font-bold text-white">Nenhum plano comercial encontrado</h3>
      <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
        {{ searchTerm ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Cadastre o primeiro plano para disponibilizar no cadastro dos estabelecimentos.' }}
      </p>
    </div>

    <!-- Grade de Cards de Planos -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="plan in filteredPlans"
        :key="plan.id"
        class="bg-slate-900/90 border rounded-2xl p-5 flex flex-col justify-between transition-all duration-300 relative group overflow-hidden"
        :class="plan.is_active ? 'border-slate-800 hover:border-purple-500/50 shadow-xl' : 'border-slate-800/40 opacity-75'"
      >
        <div class="space-y-4">
          <!-- Top Row: Nome, Status e Ciclo -->
          <div class="flex items-start justify-between gap-2">
            <div>
              <span
                class="px-2 py-0.5 rounded text-[10px] font-black uppercase tracking-wider border"
                :class="getCycleBadgeClass(plan.billing_cycle)"
              >
                {{ formatCycle(plan.billing_cycle) }}
              </span>
              <h3 class="text-lg font-black text-white mt-1.5 group-hover:text-purple-300 transition">
                {{ plan.name }}
              </h3>
            </div>

            <button
              @click="toggleStatus(plan)"
              :title="plan.is_active ? 'Desativar contratação' : 'Ativar contratação'"
              class="px-2.5 py-1 rounded-full text-xs font-bold border flex items-center gap-1.5 transition"
              :class="plan.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-slate-800 text-slate-400 border-slate-700'"
            >
              <span class="w-1.5 h-1.5 rounded-full" :class="plan.is_active ? 'bg-emerald-400' : 'bg-slate-500'"></span>
              <span>{{ plan.is_active ? 'Ativo' : 'Inativo' }}</span>
            </button>
          </div>

          <!-- Preço -->
          <div class="pt-2 border-t border-slate-800/60">
            <div class="flex items-baseline gap-1">
              <span class="text-xs text-slate-400 font-semibold">R$</span>
              <span class="text-3xl font-black text-white tracking-tight">
                {{ plan.price.toFixed(2).replace('.', ',') }}
              </span>
              <span class="text-xs text-slate-400">/ {{ formatCycleShort(plan.billing_cycle) }}</span>
            </div>
            <p v-if="plan.description" class="text-xs text-slate-400 mt-2 leading-relaxed line-clamp-2">
              {{ plan.description }}
            </p>
          </div>

          <!-- Limites e Recursos -->
          <div class="pt-3 border-t border-slate-800/60 space-y-2 text-xs">
            <div class="flex items-center justify-between text-slate-300">
              <span class="text-slate-400">Profissionais:</span>
              <span class="font-bold text-white">
                {{ plan.max_professionals === 0 ? 'Ilimitados' : `Até ${plan.max_professionals}` }}
              </span>
            </div>
            <div class="flex items-center justify-between text-slate-300">
              <span class="text-slate-400">Serviços:</span>
              <span class="font-bold text-white">
                {{ plan.max_services === 0 ? 'Ilimitados' : `Até ${plan.max_services}` }}
              </span>
            </div>

            <!-- Features Chips -->
            <div v-if="parseFeatures(plan.features).length > 0" class="pt-2 space-y-1.5">
              <div
                v-for="(feat, idx) in parseFeatures(plan.features)"
                :key="idx"
                class="flex items-center gap-1.5 text-slate-300"
              >
                <Check class="w-3.5 h-3.5 text-emerald-400 shrink-0" />
                <span class="truncate">{{ feat }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Card Actions -->
        <div class="pt-4 mt-5 border-t border-slate-800/80 flex items-center justify-between gap-2">
          <span class="text-[11px] text-slate-500">Ordem: #{{ plan.sort_order }}</span>

          <div class="flex items-center gap-1.5">
            <button
              @click="openEditModal(plan)"
              class="p-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition"
              title="Editar plano"
            >
              <Pencil class="w-4 h-4" />
            </button>

            <button
              @click="deletePlan(plan)"
              class="p-2 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 transition"
              title="Excluir plano"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar / Editar Plano -->
    <div
      v-if="showModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400">
              <CreditCard class="w-4 h-4" />
            </div>
            <h3 class="font-bold text-white text-base">
              {{ isEditing ? 'Editar Plano Comercial' : 'Novo Plano de Assinatura' }}
            </h3>
          </div>
          <button @click="showModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-4">
          <div v-if="errorMessage" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ errorMessage }}
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Nome do Plano *</label>
              <input
                v-model="form.name"
                type="text"
                required
                placeholder="Ex: Plano Profissional"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              />
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Descrição Comercial</label>
              <textarea
                v-model="form.description"
                rows="2"
                placeholder="Breve resumo da proposta de valor do plano..."
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              ></textarea>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Valor (R$) *</label>
              <input
                v-model.number="form.price"
                type="number"
                step="0.01"
                min="1"
                required
                placeholder="99.90"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500 font-mono font-bold"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Periodicidade (Asaas) *</label>
              <select
                v-model="form.billing_cycle"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              >
                <option value="MONTHLY">Mensal</option>
                <option value="QUARTERLY">Trimestral</option>
                <option value="SEMIANNUALLY">Semestral</option>
                <option value="YEARLY">Anual</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Máx. Profissionais (0 = Ilimitado)</label>
              <input
                v-model.number="form.max_professionals"
                type="number"
                min="0"
                placeholder="0"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Máx. Serviços (0 = Ilimitado)</label>
              <input
                v-model.number="form.max_services"
                type="number"
                min="0"
                placeholder="0"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Ordem de Exibição</label>
              <input
                v-model.number="form.sort_order"
                type="number"
                min="0"
                placeholder="1"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              />
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Status de Venda</label>
              <select
                v-model="form.is_active"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              >
                <option :value="true">Ativo (Visível no Onboarding)</option>
                <option :value="false">Inativo (Oculto)</option>
              </select>
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">
                Benefícios & Recursos (um por linha)
              </label>
              <textarea
                v-model="rawFeatures"
                rows="3"
                placeholder="Ex:&#10;Agendamento Online 24/7&#10;Lembretes Automáticos&#10;Relatórios Financeiros"
                class="w-full px-3.5 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 border-t border-slate-800 flex items-center justify-end gap-3 bg-slate-950/60 rounded-b-2xl">
          <button
            @click="showModal = false"
            class="px-4 py-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 text-sm font-medium transition"
          >
            Cancelar
          </button>
          <button
            @click="savePlan"
            :disabled="saving"
            class="inline-flex items-center gap-2 px-5 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold transition disabled:opacity-50"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>{{ isEditing ? 'Salvar Alterações' : 'Criar Plano' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  CreditCard, Plus, Search, CheckCircle2, Loader2,
  Package, Layers, Check, Pencil, Trash2, X
} from 'lucide-vue-next'
import api from '../../services/api'
import type { Plan, PlanBillingCycle } from '../../stores/auth'

const plans = ref<Plan[]>([])
const loading = ref(true)
const saving = ref(false)
const searchTerm = ref('')
const statusFilter = ref<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL')

const showModal = ref(false)
const isEditing = ref(false)
const editingID = ref<string | null>(null)
const errorMessage = ref('')
const rawFeatures = ref('')

const form = ref<{
  name: string
  description: string
  price: number
  billing_cycle: PlanBillingCycle
  max_professionals: number
  max_services: number
  sort_order: number
  is_active: boolean
}>({
  name: '',
  description: '',
  price: 49.90,
  billing_cycle: 'MONTHLY',
  max_professionals: 0,
  max_services: 0,
  sort_order: 0,
  is_active: true,
})

const activeCount = computed(() => plans.value.filter(p => p.is_active).length)

const filteredPlans = computed(() => {
  return plans.value.filter(p => {
    const matchesSearch =
      !searchTerm.value ||
      p.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      (p.description && p.description.toLowerCase().includes(searchTerm.value.toLowerCase()))

    const matchesStatus =
      statusFilter.value === 'ALL' ||
      (statusFilter.value === 'ACTIVE' && p.is_active) ||
      (statusFilter.value === 'INACTIVE' && !p.is_active)

    return matchesSearch && matchesStatus
  })
})

function formatCycle(cycle: PlanBillingCycle) {
  const map: Record<PlanBillingCycle, string> = {
    MONTHLY: 'Mensal',
    QUARTERLY: 'Trimestral',
    SEMIANNUALLY: 'Semestral',
    YEARLY: 'Anual',
  }
  return map[cycle] || cycle
}

function formatCycleShort(cycle: PlanBillingCycle) {
  const map: Record<PlanBillingCycle, string> = {
    MONTHLY: 'mês',
    QUARTERLY: 'trimestre',
    SEMIANNUALLY: 'semestre',
    YEARLY: 'ano',
  }
  return map[cycle] || 'mês'
}

function getCycleBadgeClass(cycle: PlanBillingCycle) {
  const map: Record<PlanBillingCycle, string> = {
    MONTHLY: 'bg-purple-950/80 text-purple-300 border-purple-700/60',
    QUARTERLY: 'bg-indigo-950/80 text-indigo-300 border-indigo-700/60',
    SEMIANNUALLY: 'bg-sky-950/80 text-sky-300 border-sky-700/60',
    YEARLY: 'bg-amber-950/80 text-amber-300 border-amber-700/60',
  }
  return map[cycle] || 'bg-slate-800 text-slate-300 border-slate-700'
}

function parseFeatures(featStr?: string): string[] {
  if (!featStr) return []
  try {
    const parsed = JSON.parse(featStr)
    if (Array.isArray(parsed)) return parsed
  } catch {
    // String dividida por quebras de linha
    return featStr.split('\n').filter(s => s.trim().length > 0)
  }
  return []
}

async function fetchPlans() {
  loading.value = true
  try {
    const res = await api.get('/admin/plans')
    if (res.data.success) {
      plans.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao carregar planos:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  isEditing.value = false
  editingID.value = null
  errorMessage.value = ''
  rawFeatures.value = ''
  form.value = {
    name: '',
    description: '',
    price: 49.90,
    billing_cycle: 'MONTHLY',
    max_professionals: 0,
    max_services: 0,
    sort_order: plans.value.length + 1,
    is_active: true,
  }
  showModal.value = true
}

function openEditModal(plan: Plan) {
  isEditing.value = true
  editingID.value = plan.id
  errorMessage.value = ''
  const featuresList = parseFeatures(plan.features)
  rawFeatures.value = featuresList.join('\n')
  form.value = {
    name: plan.name,
    description: plan.description || '',
    price: plan.price,
    billing_cycle: plan.billing_cycle,
    max_professionals: plan.max_professionals,
    max_services: plan.max_services,
    sort_order: plan.sort_order,
    is_active: plan.is_active,
  }
  showModal.value = true
}

async function savePlan() {
  if (!form.value.name || form.value.price <= 0) {
    errorMessage.value = 'Por favor, informe o nome e um valor válido para o plano.'
    return
  }

  saving.value = true
  errorMessage.value = ''

  // Formata features em JSON Array
  const featuresArray = rawFeatures.value
    .split('\n')
    .map(s => s.trim())
    .filter(s => s.length > 0)

  const payload = {
    ...form.value,
    features: JSON.stringify(featuresArray),
  }

  try {
    if (isEditing.value && editingID.value) {
      const res = await api.put(`/admin/plans/${editingID.value}`, payload)
      if (res.data.success) {
        showModal.value = false
        await fetchPlans()
      }
    } else {
      const res = await api.post('/admin/plans', payload)
      if (res.data.success) {
        showModal.value = false
        await fetchPlans()
      }
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || err.message || 'Erro ao salvar plano comercial.'
  } finally {
    saving.value = false
  }
}

async function toggleStatus(plan: Plan) {
  const nextStatus = !plan.is_active
  try {
    const res = await api.patch(`/admin/plans/${plan.id}/status`, { is_active: nextStatus })
    if (res.data.success) {
      plan.is_active = nextStatus
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao alterar status do plano.')
  }
}

async function deletePlan(plan: Plan) {
  if (!confirm(`Tem certeza que deseja excluir o plano "${plan.name}"?`)) {
    return
  }

  try {
    const res = await api.delete(`/admin/plans/${plan.id}`)
    if (res.data.success) {
      await fetchPlans()
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Não foi possível excluir o plano.')
  }
}

onMounted(() => {
  fetchPlans()
})
</script>
