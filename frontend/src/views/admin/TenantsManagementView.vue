<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <Building2 class="w-7 h-7 text-purple-400" />
          <span>Gestão de Estabelecimentos</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Cadastre e gerencie todos os tenants (barbearias e salões) e seus administradores iniciais na plataforma.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-600/30 transition active:scale-95 shrink-0"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Estabelecimento</span>
      </button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400">
          <Building2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Total de Estabelecimentos</p>
          <p class="text-2xl font-black text-white mt-0.5">{{ tenants.length }}</p>
        </div>
      </div>

      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400">
          <CheckCircle2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Estabelecimentos Ativos</p>
          <p class="text-2xl font-black text-emerald-400 mt-0.5">{{ activeCount }}</p>
        </div>
      </div>

      <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex items-center gap-4">
        <div class="w-12 h-12 rounded-xl bg-rose-500/10 border border-rose-500/20 flex items-center justify-center text-rose-400">
          <XCircle class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-slate-400 font-medium">Inativos / Pausados</p>
          <p class="text-2xl font-black text-rose-400 mt-0.5">{{ inactiveCount }}</p>
        </div>
      </div>
    </div>

    <!-- Filtros & Busca -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex flex-col sm:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome, slug, cidade ou e-mail..."
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
      <p class="text-sm text-slate-400">Carregando estabelecimentos...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredTenants.length === 0" class="bg-slate-900/50 border border-slate-800/80 rounded-2xl p-12 text-center">
      <Building2 class="w-12 h-12 text-slate-600 mx-auto mb-3" />
      <h3 class="text-base font-bold text-white">Nenhum estabelecimento encontrado</h3>
      <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
        {{ searchTerm ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Comece cadastrando o primeiro estabelecimento da plataforma.' }}
      </p>
    </div>

    <!-- Tabela de Estabelecimentos -->
    <div v-else class="bg-slate-900/80 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-950/80 border-b border-slate-800 text-xs text-slate-400 uppercase tracking-wider font-semibold">
            <tr>
              <th class="py-3.5 px-4">Estabelecimento</th>
              <th class="py-3.5 px-4">Contato</th>
              <th class="py-3.5 px-4">Localização</th>
              <th class="py-3.5 px-4">Status</th>
              <th class="py-3.5 px-4">Cadastro</th>
              <th class="py-3.5 px-4 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60">
            <tr v-for="t in filteredTenants" :key="t.id" class="hover:bg-slate-800/40 transition">
              <td class="py-4 px-4">
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-xl flex items-center justify-center font-black text-sm text-white shrink-0 overflow-hidden border border-slate-700/60 bg-slate-800"
                    :style="t.logo_url ? '' : { backgroundColor: t.primary_color || '#8b5cf6' }"
                  >
                    <img v-if="t.logo_url" :src="t.logo_url" :alt="t.name" class="w-full h-full object-cover" />
                    <span v-else>{{ t.name.charAt(0) }}</span>
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-white text-sm truncate">{{ t.name }}</p>
                    <p class="text-xs text-purple-400 font-mono truncate">/{{ t.slug }}</p>
                  </div>
                </div>
              </td>
              <td class="py-4 px-4 text-xs text-slate-300">
                <div>{{ t.phone || '-' }}</div>
                <div class="text-slate-400 truncate">{{ t.email || '-' }}</div>
              </td>
              <td class="py-4 px-4 text-xs text-slate-300">
                <div v-if="t.city || t.state">{{ t.city }} - {{ t.state }}</div>
                <div v-else class="text-slate-500">Não informado</div>
              </td>
              <td class="py-4 px-4">
                <span
                  class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold border"
                  :class="t.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-rose-950/80 text-rose-300 border-rose-700/60'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="t.is_active ? 'bg-emerald-400' : 'bg-rose-400'"></span>
                  {{ t.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </td>
              <td class="py-4 px-4 text-xs text-slate-400">
                {{ formatDate(t.created_at) }}
              </td>
              <td class="py-4 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <a
                    :href="`/agendamento/${t.slug}`"
                    target="_blank"
                    title="Ver página de agendamento público"
                    class="p-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                  >
                    <ExternalLink class="w-4 h-4" />
                  </a>

                  <button
                    @click="openEditModal(t)"
                    title="Editar informações"
                    class="p-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <button
                    @click="toggleStatus(t)"
                    :title="t.is_active ? 'Inativar estabelecimento' : 'Ativar estabelecimento'"
                    class="p-2 rounded-lg transition"
                    :class="t.is_active ? 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-400' : 'bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400'"
                  >
                    <Power class="w-4 h-4" />
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: Cadastrar Novo Estabelecimento -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-2xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400">
              <Building2 class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base">Novo Estabelecimento (Tenant)</h3>
              <p class="text-xs text-slate-400">Crie o estabelecimento e configure o administrador inicial.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-6">
          <div v-if="createError" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ createError }}
          </div>

          <!-- Seção 1: Dados do Estabelecimento -->
          <div class="space-y-4">
            <h4 class="text-xs font-bold text-purple-400 uppercase tracking-wider flex items-center gap-1.5">
              <Store class="w-3.5 h-3.5" />
              <span>1. Dados do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Nome do Estabelecimento *</label>
                <input
                  v-model="createForm.name"
                  @input="autoGenerateSlug"
                  type="text"
                  placeholder="Ex: Barbearia Dom Navalha"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Identificador URL (Slug) *</label>
                <input
                  v-model="createForm.slug"
                  type="text"
                  placeholder="Ex: dom-navalha"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-purple-400 font-mono focus:outline-none focus:border-purple-500"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Telefone / WhatsApp *</label>
                <input
                  v-model="createForm.phone"
                  type="text"
                  placeholder="(11) 98765-4321"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">E-mail de Contato</label>
                <input
                  v-model="createForm.email"
                  type="email"
                  placeholder="contato@estabelecimento.com"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Cidade</label>
                <input
                  v-model="createForm.city"
                  type="text"
                  placeholder="São Paulo"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Estado (UF)</label>
                <input
                  v-model="createForm.state"
                  type="text"
                  maxlength="2"
                  placeholder="SP"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                />
              </div>

              <div class="sm:col-span-2">
                <label class="block text-xs font-medium text-slate-300 mb-1">Endereço Completo</label>
                <input
                  v-model="createForm.address"
                  type="text"
                  placeholder="Av. Paulista, 1578 - Bela Vista"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Intervalo de Horários (minutos)</label>
                <select
                  v-model="createForm.slot_interval_minutes"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                >
                  <option :value="15">15 minutos</option>
                  <option :value="30">30 minutos (Recomendado)</option>
                  <option :value="45">45 minutos</option>
                  <option :value="60">60 minutos</option>
                </select>
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Cor Primária (Tema)</label>
                <div class="flex items-center gap-2">
                  <input
                    v-model="createForm.primary_color"
                    type="color"
                    class="w-9 h-9 rounded-lg bg-transparent border-0 cursor-pointer"
                  />
                  <input
                    v-model="createForm.primary_color"
                    type="text"
                    class="flex-1 px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white font-mono"
                  />
                </div>
              </div>
            </div>
          </div>

          <!-- Seção 2: Administrador Inicial -->
          <div class="space-y-4 pt-4 border-t border-slate-800">
            <h4 class="text-xs font-bold text-emerald-400 uppercase tracking-wider flex items-center gap-1.5">
              <ShieldCheck class="w-3.5 h-3.5" />
              <span>2. Administrador Inicial do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="sm:col-span-2">
                <label class="block text-xs font-medium text-slate-300 mb-1">Nome do Administrador *</label>
                <input
                  v-model="createForm.admin_name"
                  type="text"
                  placeholder="Ex: Carlos Administrador"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">E-mail de Login *</label>
                <input
                  v-model="createForm.admin_email"
                  type="email"
                  placeholder="admin@estabelecimento.com"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                  required
                />
              </div>

              <div>
                <label class="block text-xs font-medium text-slate-300 mb-1">Senha de Acesso Inicial *</label>
                <input
                  v-model="createForm.admin_password"
                  type="password"
                  placeholder="Mínimo 6 caracteres"
                  class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
                  required
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 border-t border-slate-800 flex items-center justify-end gap-3 bg-slate-950/60 rounded-b-2xl">
          <button
            @click="showCreateModal = false"
            class="px-4 py-2 rounded-xl text-slate-400 hover:text-white hover:bg-slate-800 text-sm font-medium transition"
          >
            Cancelar
          </button>
          <button
            @click="submitCreate"
            :disabled="saving"
            class="inline-flex items-center gap-2 px-5 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold transition disabled:opacity-50"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>Criar Estabelecimento</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Editar Estabelecimento -->
    <div
      v-if="showEditModal && editingTenant"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <h3 class="font-bold text-white text-base">Editar Estabelecimento</h3>
          <button @click="showEditModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 overflow-y-auto space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="sm:col-span-2">
              <label class="block text-xs font-medium text-slate-300 mb-1">Nome</label>
              <input v-model="editingTenant.name" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Slug</label>
              <input v-model="editingTenant.slug" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-purple-400 font-mono" />
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Telefone</label>
              <input v-model="editingTenant.phone" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">E-mail</label>
              <input v-model="editingTenant.email" type="email" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Cidade</label>
              <input v-model="editingTenant.city" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
            </div>

            <div class="sm:col-span-2">
              <label class="block text-xs font-medium text-slate-300 mb-1">Endereço</label>
              <input v-model="editingTenant.address" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Cor Primária</label>
              <div class="flex items-center gap-2">
                <input v-model="editingTenant.primary_color" type="color" class="w-9 h-9 rounded-lg bg-transparent border-0 cursor-pointer" />
                <input v-model="editingTenant.primary_color" type="text" class="flex-1 px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white font-mono" />
              </div>
            </div>

            <div>
              <label class="block text-xs font-medium text-slate-300 mb-1">Status</label>
              <select v-model="editingTenant.is_active" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white">
                <option :value="true">Ativo</option>
                <option :value="false">Inativo</option>
              </select>
            </div>
          </div>
        </div>

        <div class="p-5 border-t border-slate-800 flex items-center justify-end gap-3 bg-slate-950/60 rounded-b-2xl">
          <button @click="showEditModal = false" class="px-4 py-2 rounded-xl text-slate-400 hover:text-white text-sm font-medium">
            Cancelar
          </button>
          <button
            @click="submitEdit"
            :disabled="saving"
            class="px-5 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white text-sm font-bold transition disabled:opacity-50"
          >
            Salvar Alterações
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Building2, Plus, Search, CheckCircle2, XCircle, Loader2,
  ExternalLink, Pencil, Power, X, Store, ShieldCheck
} from 'lucide-vue-next'
import api from '../../services/api'
import type { Tenant } from '../../stores/auth'

const tenants = ref<Tenant[]>([])
const loading = ref(true)
const saving = ref(false)
const searchTerm = ref('')
const statusFilter = ref<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL')

const showCreateModal = ref(false)
const createError = ref('')
const createForm = ref({
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
  slot_interval_minutes: 30,
  admin_name: '',
  admin_email: '',
  admin_password: '',
})

const showEditModal = ref(false)
const editingTenant = ref<Tenant | null>(null)

const activeCount = computed(() => tenants.value.filter(t => t.is_active).length)
const inactiveCount = computed(() => tenants.value.filter(t => !t.is_active).length)

const filteredTenants = computed(() => {
  return tenants.value.filter(t => {
    const matchesSearch =
      !searchTerm.value ||
      t.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      t.slug.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      (t.city && t.city.toLowerCase().includes(searchTerm.value.toLowerCase())) ||
      (t.email && t.email.toLowerCase().includes(searchTerm.value.toLowerCase()))

    const matchesStatus =
      statusFilter.value === 'ALL' ||
      (statusFilter.value === 'ACTIVE' && t.is_active) ||
      (statusFilter.value === 'INACTIVE' && !t.is_active)

    return matchesSearch && matchesStatus
  })
})

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('pt-BR')
  } catch {
    return dateStr
  }
}

function autoGenerateSlug() {
  if (!createForm.value.name) return
  createForm.value.slug = createForm.value.name
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
    .replace(/[^a-z0-9]+/g, '-')
    .replace(/^-+|-+$/g, '')
}

async function fetchTenants() {
  loading.value = true
  try {
    const res = await api.get('/admin/tenants')
    if (res.data.success) {
      tenants.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao buscar estabelecimentos:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  createError.value = ''
  createForm.value = {
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
    slot_interval_minutes: 30,
    admin_name: '',
    admin_email: '',
    admin_password: '',
  }
  showCreateModal.value = true
}

async function submitCreate() {
  if (!createForm.value.name || !createForm.value.slug || !createForm.value.phone || !createForm.value.admin_email || !createForm.value.admin_password) {
    createError.value = 'Por favor, preencha todos os campos obrigatórios (*).'
    return
  }

  saving.value = true
  createError.value = ''
  try {
    const res = await api.post('/admin/tenants', createForm.value)
    if (res.data.success) {
      showCreateModal.value = false
      await fetchTenants()
    }
  } catch (err: any) {
    createError.value = err.response?.data?.error || err.message || 'Erro ao cadastrar estabelecimento.'
  } finally {
    saving.value = false
  }
}

function openEditModal(tenant: Tenant) {
  editingTenant.value = { ...tenant }
  showEditModal.value = true
}

async function submitEdit() {
  if (!editingTenant.value) return
  saving.value = true
  try {
    const res = await api.put(`/admin/tenants/${editingTenant.value.id}`, editingTenant.value)
    if (res.data.success) {
      showEditModal.value = false
      await fetchTenants()
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao atualizar estabelecimento.')
  } finally {
    saving.value = false
  }
}

async function toggleStatus(tenant: Tenant) {
  const nextStatus = !tenant.is_active
  const actionName = nextStatus ? 'ativar' : 'inativar'
  if (!confirm(`Deseja realmente ${actionName} o estabelecimento "${tenant.name}"?`)) {
    return
  }

  try {
    const res = await api.patch(`/admin/tenants/${tenant.id}/status`, { is_active: nextStatus })
    if (res.data.success) {
      tenant.is_active = nextStatus
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao alterar status.')
  }
}

onMounted(() => {
  fetchTenants()
})
</script>
