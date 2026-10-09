<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-purple-500/10 border border-purple-500/20 text-purple-400 text-xs font-semibold mb-2">
          <Layers class="w-3.5 h-3.5" />
          <span>Catálogo Comercial da Plataforma</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-['Outfit'] flex items-center gap-2.5">
          <CreditCard class="w-8 h-8 text-purple-400" />
          <span>Planos de Assinatura</span>
        </h1>
        <p class="text-xs sm:text-sm text-surface-400 mt-1">
          Gerencie a grade comercial de planos, periodicidades e limites para contratação dos estabelecimentos.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 via-indigo-600 to-purple-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-900/30 hover:shadow-purple-700/40 transition active:scale-95 shrink-0"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Plano Comercial</span>
      </button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="glass-card p-5 flex items-center gap-4 group hover:border-purple-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 group-hover:scale-105 transition-transform shrink-0">
          <Package class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Total de Planos</p>
          <p class="text-2xl font-black text-white font-['Outfit'] mt-0.5">{{ plans.length }}</p>
        </div>
      </div>

      <div class="glass-card p-5 flex items-center gap-4 group hover:border-emerald-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 group-hover:scale-105 transition-transform shrink-0">
          <CheckCircle2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Planos Ativos para Venda</p>
          <p class="text-2xl font-black text-emerald-400 font-['Outfit'] mt-0.5">{{ activeCount }}</p>
        </div>
      </div>

      <div class="glass-card p-5 flex items-center gap-4 group hover:border-indigo-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-indigo-500/10 border border-indigo-500/20 flex items-center justify-center text-indigo-400 group-hover:scale-105 transition-transform shrink-0">
          <Layers class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Integração Asaas Recorrente</p>
          <p class="text-sm font-bold text-indigo-300 mt-1 flex items-center gap-2">
            <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
            <span>API v3 Conectada</span>
          </p>
        </div>
      </div>
    </div>

    <!-- Filtros e Busca -->
    <div class="glass-card p-3 sm:p-4 flex flex-col sm:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-surface-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome ou descrição do plano..."
          class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-surface-100 placeholder-surface-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        />
      </div>

      <div class="flex items-center gap-2 w-full sm:w-auto">
        <select
          v-model="statusFilter"
          class="w-full sm:w-48 py-2.5 px-3.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-surface-200 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        >
          <option value="ALL">Todos os status</option>
          <option value="ACTIVE">Apenas Ativos</option>
          <option value="INACTIVE">Apenas Inativos</option>
          <option value="FREE">Apenas Gratuitos</option>
          <option value="PAID">Apenas Pagos (Asaas)</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-20">
      <Loader2 class="w-9 h-9 text-purple-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-surface-400">Carregando planos de assinatura...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredPlans.length === 0" class="glass-card p-12 text-center">
      <div class="w-14 h-14 rounded-2xl bg-surface-800/80 border border-surface-700/60 flex items-center justify-center mx-auto mb-4 text-surface-500">
        <Package class="w-7 h-7" />
      </div>
      <h3 class="text-base font-bold text-white font-['Outfit']">Nenhum plano comercial encontrado</h3>
      <p class="text-xs text-surface-400 mt-1.5 max-w-sm mx-auto">
        {{ searchTerm ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Cadastre o primeiro plano para disponibilizar no cadastro dos estabelecimentos.' }}
      </p>
    </div>

    <!-- Grade de Cards de Planos -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-5">
      <div
        v-for="plan in filteredPlans"
        :key="plan.id"
        class="glass-card p-6 flex flex-col justify-between transition-all duration-300 relative group overflow-hidden"
        :class="plan.is_active ? 'hover:border-purple-500/50 hover:shadow-xl hover:shadow-purple-950/30 hover:-translate-y-1' : 'opacity-70 border-surface-800'"
      >
        <div class="space-y-4">
          <!-- Top Row: Nome, Status, Tipo (Gratuito/Pago) e Ciclo -->
          <div class="flex items-start justify-between gap-2">
            <div>
              <div class="flex items-center gap-1.5 flex-wrap">
                <span
                  v-if="plan.is_free"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-emerald-500/10 text-emerald-400 border border-emerald-500/30 shadow-sm"
                >
                  Gratuito
                </span>
                <span
                  v-else
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider bg-purple-500/10 text-purple-300 border border-purple-500/30"
                >
                  Asaas Recorrente
                </span>
                <span
                  v-if="!plan.is_free"
                  class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border"
                  :class="getCycleBadgeClass(plan.billing_cycle)"
                >
                  {{ formatCycle(plan.billing_cycle) }}
                </span>
              </div>
              <h3 class="text-lg font-black text-white font-['Outfit'] mt-2 group-hover:text-purple-300 transition">
                {{ plan.name }}
              </h3>
            </div>

            <button
              @click="toggleStatus(plan)"
              :title="plan.is_active ? 'Desativar contratação' : 'Ativar contratação'"
              class="px-2.5 py-1 rounded-full text-xs font-bold border flex items-center gap-1.5 transition"
              :class="plan.is_active ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20 hover:bg-emerald-500/20' : 'bg-surface-800 text-surface-400 border-surface-700 hover:bg-surface-750'"
            >
              <span class="w-1.5 h-1.5 rounded-full" :class="plan.is_active ? 'bg-emerald-400' : 'bg-surface-500'"></span>
              <span>{{ plan.is_active ? 'Ativo' : 'Inativo' }}</span>
            </button>
          </div>

          <!-- Preço -->
          <div class="pt-3 border-t border-surface-800/80">
            <div v-if="plan.is_free" class="flex items-baseline gap-1.5">
              <span class="text-2xl font-black text-emerald-400 font-['Outfit'] tracking-tight">
                Gratuito
              </span>
              <span class="text-xs text-surface-400 font-medium">(R$ 0,00)</span>
            </div>
            <div v-else class="flex items-baseline gap-1">
              <span class="text-xs text-surface-400 font-semibold">R$</span>
              <span class="text-3xl font-black text-white font-['Outfit'] tracking-tight">
                {{ plan.price.toFixed(2).replace('.', ',') }}
              </span>
              <span class="text-xs text-surface-400">/ {{ formatCycleShort(plan.billing_cycle) }}</span>
            </div>
            <p v-if="plan.description" class="text-xs text-surface-400 mt-2 leading-relaxed line-clamp-2">
              {{ plan.description }}
            </p>
          </div>

          <!-- Limites e Recursos -->
          <div class="pt-3 border-t border-surface-800/80 space-y-2 text-xs">
            <div class="flex items-center justify-between text-surface-300">
              <span class="text-surface-400">Profissionais:</span>
              <span class="font-bold text-white">
                {{ plan.max_professionals === 0 ? 'Ilimitados' : `Até ${plan.max_professionals}` }}
              </span>
            </div>
            <div class="flex items-center justify-between text-surface-300">
              <span class="text-surface-400">Serviços:</span>
              <span class="font-bold text-white">
                {{ plan.max_services === 0 ? 'Ilimitados' : `Até ${plan.max_services}` }}
              </span>
            </div>

            <!-- Features Chips -->
            <div v-if="parseFeatures(plan.features).length > 0" class="pt-2 space-y-1.5">
              <div
                v-for="(feat, idx) in parseFeatures(plan.features)"
                :key="idx"
                class="flex items-center gap-2 text-surface-300 text-xs"
              >
                <div class="w-4 h-4 rounded-full bg-emerald-500/10 flex items-center justify-center shrink-0">
                  <Check class="w-3 h-3 text-emerald-400" />
                </div>
                <span class="truncate">{{ feat }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer Card Actions -->
        <div class="pt-4 mt-5 border-t border-surface-800/80 flex items-center justify-between gap-2">
          <span class="text-[11px] text-surface-500 font-mono">Ordem: #{{ plan.sort_order }}</span>

          <div class="flex items-center gap-1.5">
            <button
              @click="openEditModal(plan)"
              class="p-2 rounded-xl bg-surface-800 hover:bg-surface-700 text-surface-300 hover:text-white transition"
              title="Editar plano"
            >
              <Pencil class="w-4 h-4" />
            </button>

            <button
              @click="deletePlan(plan)"
              class="p-2 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 transition"
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
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <CreditCard class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base font-['Outfit']">
                {{ isEditing ? 'Editar Plano Comercial' : 'Novo Plano de Assinatura' }}
              </h3>
              <p class="text-xs text-surface-400">Configure os limites e termos para este plano.</p>
            </div>
          </div>
          <button @click="showModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-4">
          <div v-if="errorMessage" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ errorMessage }}
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Toggle Plano Gratuito -->
            <div class="sm:col-span-2 p-3.5 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-between">
              <div>
                <span class="text-xs font-bold text-white flex items-center gap-1.5">
                  <Sparkles class="w-3.5 h-3.5 text-emerald-400" />
                  <span>Plano de Assinatura Gratuito</span>
                </span>
                <span class="text-[11px] text-surface-400 block mt-0.5">
                  Estabelecimentos cadastrados neste plano têm ativação imediata sem cobranças ou integração no Asaas.
                </span>
              </div>
              <label class="relative inline-flex items-center cursor-pointer shrink-0">
                <input type="checkbox" v-model="form.is_free" @change="onToggleFreePlan" class="sr-only peer" />
                <div class="w-11 h-6 bg-surface-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-surface-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-emerald-500"></div>
              </label>
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome do Plano *</label>
              <input
                v-model="form.name"
                type="text"
                required
                placeholder="Ex: Plano Profissional"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Descrição Comercial</label>
              <textarea
                v-model="form.description"
                rows="2"
                placeholder="Breve resumo da proposta de valor do plano..."
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              ></textarea>
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">
                Valor (R$) <span v-if="!form.is_free">*</span><span v-else class="text-emerald-400 lowercase font-medium">(grátis)</span>
              </label>
              <input
                v-if="!form.is_free"
                v-model.number="form.price"
                type="number"
                step="0.01"
                min="0"
                required
                placeholder="99.90"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition font-mono font-bold"
              />
              <div v-else class="w-full px-3.5 py-2.5 rounded-xl bg-surface-900 border border-emerald-500/30 text-emerald-400 text-sm font-bold font-mono">
                R$ 0,00 (Gratuito)
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">
                {{ form.is_free ? 'Periodicidade (Simbólica)' : 'Periodicidade (Asaas) *' }}
              </label>
              <select
                v-model="form.billing_cycle"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              >
                <option value="MONTHLY">Mensal</option>
                <option value="QUARTERLY">Trimestral</option>
                <option value="SEMIANNUALLY">Semestral</option>
                <option value="YEARLY">Anual</option>
              </select>
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Máx. Profissionais (0 = Ilimitado)</label>
              <input
                v-model.number="form.max_professionals"
                type="number"
                min="0"
                placeholder="0"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Máx. Serviços (0 = Ilimitado)</label>
              <input
                v-model.number="form.max_services"
                type="number"
                min="0"
                placeholder="0"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Ordem de Exibição</label>
              <input
                v-model.number="form.sort_order"
                type="number"
                min="0"
                placeholder="1"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Status de Venda</label>
              <select
                v-model="form.is_active"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              >
                <option :value="true">Ativo (Visível no Onboarding)</option>
                <option :value="false">Inativo (Oculto)</option>
              </select>
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">
                Benefícios & Recursos (um por linha)
              </label>
              <textarea
                v-model="rawFeatures"
                rows="3"
                placeholder="Ex:&#10;Agendamento Online 24/7&#10;Lembretes Automáticos&#10;Relatórios Financeiros"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end gap-3 bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showModal = false"
            class="px-4 py-2.5 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 text-sm font-medium transition"
          >
            Cancelar
          </button>
          <button
            @click="savePlan"
            :disabled="saving"
            class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold transition disabled:opacity-50 shadow-lg shadow-purple-900/30"
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
  Package, Layers, Check, Pencil, Trash2, X, Sparkles
} from 'lucide-vue-next'
import api from '../../services/api'
import type { Plan, PlanBillingCycle } from '../../stores/auth'

const plans = ref<Plan[]>([])
const loading = ref(true)
const saving = ref(false)
const searchTerm = ref('')
const statusFilter = ref<'ALL' | 'ACTIVE' | 'INACTIVE' | 'FREE' | 'PAID'>('ALL')

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
  is_free: boolean
}>({
  name: '',
  description: '',
  price: 49.90,
  billing_cycle: 'MONTHLY',
  max_professionals: 0,
  max_services: 0,
  sort_order: 0,
  is_active: true,
  is_free: false,
})

function onToggleFreePlan() {
  if (form.value.is_free) {
    form.value.price = 0
  } else if (form.value.price === 0) {
    form.value.price = 49.90
  }
}

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
      (statusFilter.value === 'INACTIVE' && !p.is_active) ||
      (statusFilter.value === 'FREE' && p.is_free) ||
      (statusFilter.value === 'PAID' && !p.is_free)

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
    MONTHLY: 'bg-purple-500/10 text-purple-300 border-purple-500/30',
    QUARTERLY: 'bg-indigo-500/10 text-indigo-300 border-indigo-500/30',
    SEMIANNUALLY: 'bg-sky-500/10 text-sky-300 border-sky-500/30',
    YEARLY: 'bg-brand-500/10 text-brand-300 border-brand-500/30',
  }
  return map[cycle] || 'bg-surface-800 text-surface-300 border-surface-700'
}

function parseFeatures(featStr?: string): string[] {
  if (!featStr) return []
  try {
    const parsed = JSON.parse(featStr)
    if (Array.isArray(parsed)) return parsed
  } catch {
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
    is_free: false,
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
    is_free: Boolean(plan.is_free),
  }
  showModal.value = true
}

async function savePlan() {
  if (!form.value.name || (!form.value.is_free && form.value.price <= 0)) {
    errorMessage.value = 'Por favor, informe o nome e um valor válido para o plano comercial.'
    return
  }

  if (form.value.is_free) {
    form.value.price = 0
  }

  saving.value = true
  errorMessage.value = ''

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
