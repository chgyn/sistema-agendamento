<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <Users class="w-7 h-7" :class="authStore.isGlobalAdmin ? 'text-purple-400' : 'text-emerald-400'" />
          <span>{{ authStore.isGlobalAdmin ? 'Usuários da Plataforma' : 'Gestão de Usuários' }}</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          {{ authStore.isGlobalAdmin
            ? 'Gerencie operadores e administradores cadastrados em todos os estabelecimentos da plataforma.'
            : 'Cadastre e administre os operadores e administradores vinculados à sua barbearia/salão.'
          }}
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl text-white text-sm font-bold shadow-lg transition active:scale-95 shrink-0"
        :class="authStore.isGlobalAdmin
          ? 'bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 shadow-purple-600/30'
          : 'bg-gradient-to-r from-emerald-500 to-teal-500 hover:from-emerald-400 hover:to-teal-400 shadow-emerald-500/30 text-slate-950 font-black'"
      >
        <UserPlus class="w-4 h-4" />
        <span>Novo Usuário</span>
      </button>
    </div>

    <!-- Filtros & Busca -->
    <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-4 flex flex-col sm:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-slate-400 absolute left-3 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome ou e-mail..."
          class="w-full pl-9 pr-4 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-slate-100 placeholder-slate-500 focus:outline-none focus:border-emerald-500 transition"
        />
      </div>

      <!-- Filtro de Tenant (Apenas para Global Admin) -->
      <div v-if="authStore.isGlobalAdmin" class="w-full sm:w-60">
        <select
          v-model="selectedTenantFilter"
          @change="fetchUsers"
          class="w-full py-2 px-3 rounded-xl bg-slate-950 border border-slate-800 text-sm text-slate-200 focus:outline-none focus:border-purple-500 transition"
        >
          <option value="">Todos os Estabelecimentos</option>
          <option v-for="t in tenantsList" :key="t.id" :value="t.id">
            {{ t.name }}
          </option>
        </select>
      </div>

      <!-- Filtro de Perfil -->
      <div class="w-full sm:w-44">
        <select
          v-model="selectedRoleFilter"
          class="w-full py-2 px-3 rounded-xl bg-slate-950 border border-slate-800 text-sm text-slate-200 focus:outline-none focus:border-emerald-500 transition"
        >
          <option value="ALL">Todos os Perfis</option>
          <option value="ADMIN_TENANT">Administradores</option>
          <option value="OPERATOR">Operadores</option>
        </select>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-16">
      <Loader2 class="w-8 h-8 text-emerald-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-slate-400">Carregando usuários...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredUsers.length === 0" class="bg-slate-900/50 border border-slate-800/80 rounded-2xl p-12 text-center">
      <Users class="w-12 h-12 text-slate-600 mx-auto mb-3" />
      <h3 class="text-base font-bold text-white">Nenhum usuário encontrado</h3>
      <p class="text-xs text-slate-400 mt-1 max-w-sm mx-auto">
        {{ searchTerm ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Cadastre os primeiros operadores ou administradores do estabelecimento.' }}
      </p>
    </div>

    <!-- Tabela de Usuários -->
    <div v-else class="bg-slate-900/80 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-950/80 border-b border-slate-800 text-xs text-slate-400 uppercase tracking-wider font-semibold">
            <tr>
              <th class="py-3.5 px-4">Usuário</th>
              <th v-if="authStore.isGlobalAdmin" class="py-3.5 px-4">Estabelecimento</th>
              <th class="py-3.5 px-4">Perfil / Permissão</th>
              <th class="py-3.5 px-4">Status</th>
              <th class="py-3.5 px-4">Cadastro</th>
              <th class="py-3.5 px-4 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60">
            <tr v-for="u in filteredUsers" :key="u.id" class="hover:bg-slate-800/40 transition">
              <td class="py-4 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-9 h-9 rounded-xl bg-slate-800 border border-slate-700/60 flex items-center justify-center font-bold text-sm text-slate-200 shrink-0">
                    {{ u.name.charAt(0) }}
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-white text-sm truncate">{{ u.name }}</p>
                    <p class="text-xs text-slate-400 truncate">{{ u.email }}</p>
                  </div>
                </div>
              </td>

              <td v-if="authStore.isGlobalAdmin" class="py-4 px-4 text-xs">
                <div v-if="u.tenant" class="flex items-center gap-1.5 text-slate-200">
                  <Store class="w-3.5 h-3.5 text-purple-400" />
                  <span>{{ u.tenant.name }}</span>
                </div>
                <div v-else class="text-purple-300 font-semibold flex items-center gap-1">
                  <ShieldCheck class="w-3.5 h-3.5" />
                  <span>Administração Geral</span>
                </div>
              </td>

              <td class="py-4 px-4">
                <span
                  class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold border"
                  :class="getRoleBadgeClass(u.role)"
                >
                  <ShieldCheck v-if="u.role === 'ADMIN_GLOBAL'" class="w-3 h-3 text-purple-400" />
                  <Crown v-else-if="u.role === 'ADMIN_TENANT' || u.role === 'ADMIN'" class="w-3 h-3 text-emerald-400" />
                  <UserCheck v-else class="w-3 h-3 text-sky-400" />
                  <span>{{ getRoleLabel(u.role) }}</span>
                </span>
              </td>

              <td class="py-4 px-4">
                <span
                  class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold border"
                  :class="u.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-rose-950/80 text-rose-300 border-rose-700/60'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="u.is_active ? 'bg-emerald-400' : 'bg-rose-400'"></span>
                  {{ u.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </td>

              <td class="py-4 px-4 text-xs text-slate-400">
                {{ formatDate(u.created_at) }}
              </td>

              <td class="py-4 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="openEditModal(u)"
                    title="Editar Usuário"
                    class="p-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <button
                    @click="toggleStatus(u)"
                    :title="u.is_active ? 'Inativar Usuário' : 'Ativar Usuário'"
                    class="p-2 rounded-lg transition"
                    :class="u.is_active ? 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-400' : 'bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400'"
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

    <!-- MODAL: Novo Usuário -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-lg shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400">
              <UserPlus class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base">Novo Usuário</h3>
              <p class="text-xs text-slate-400">Cadastre um operador ou administrador.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div v-if="createError" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ createError }}
          </div>

          <!-- Estabelecimento (se for Global Admin) -->
          <div v-if="authStore.isGlobalAdmin">
            <label class="block text-xs font-medium text-slate-300 mb-1">Estabelecimento (Tenant) *</label>
            <select
              v-model="createForm.tenant_id"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              required
            >
              <option value="" disabled>Selecione o estabelecimento</option>
              <option v-for="t in tenantsList" :key="t.id" :value="t.id">
                {{ t.name }} ({{ t.slug }})
              </option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nome Completo *</label>
            <input
              v-model="createForm.name"
              type="text"
              placeholder="Ex: João Silva"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-emerald-500"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">E-mail de Acesso *</label>
            <input
              v-model="createForm.email"
              type="email"
              placeholder="usuario@estabelecimento.com"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-emerald-500"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Perfil de Acesso *</label>
            <div class="grid grid-cols-2 gap-3">
              <label
                class="flex flex-col p-3 rounded-xl border cursor-pointer transition"
                :class="createForm.role === 'OPERATOR' ? 'bg-sky-500/10 border-sky-500/40 text-white' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
              >
                <div class="flex items-center justify-between">
                  <span class="text-xs font-bold">Operador</span>
                  <input type="radio" value="OPERATOR" v-model="createForm.role" class="text-sky-500" />
                </div>
                <span class="text-[11px] text-slate-400 mt-1 leading-snug">Acesso à agenda, clientes e atendimentos.</span>
              </label>

              <label
                class="flex flex-col p-3 rounded-xl border cursor-pointer transition"
                :class="createForm.role === 'ADMIN_TENANT' ? 'bg-emerald-500/10 border-emerald-500/40 text-white' : 'bg-slate-950 border-slate-800 text-slate-400 hover:border-slate-700'"
              >
                <div class="flex items-center justify-between">
                  <span class="text-xs font-bold">Administrador</span>
                  <input type="radio" value="ADMIN_TENANT" v-model="createForm.role" class="text-emerald-500" />
                </div>
                <span class="text-[11px] text-slate-400 mt-1 leading-snug">Gestão de usuários, serviços e configurações.</span>
              </label>
            </div>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Senha Inicial *</label>
            <input
              v-model="createForm.password"
              type="password"
              placeholder="Mínimo de 6 caracteres"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-emerald-500"
              required
            />
          </div>
        </div>

        <div class="p-5 border-t border-slate-800 flex items-center justify-end gap-3 bg-slate-950/60 rounded-b-2xl">
          <button @click="showCreateModal = false" class="px-4 py-2 rounded-xl text-slate-400 hover:text-white text-sm font-medium">
            Cancelar
          </button>
          <button
            @click="submitCreate"
            :disabled="saving"
            class="px-5 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm transition disabled:opacity-50 inline-flex items-center gap-2"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>Cadastrar Usuário</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Editar Usuário -->
    <div
      v-if="showEditModal && editingUser"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-lg shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <h3 class="font-bold text-white text-base">Editar Usuário</h3>
          <button @click="showEditModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nome Completo</label>
            <input v-model="editingUser.name" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">E-mail</label>
            <input v-model="editingUser.email" type="email" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Perfil de Acesso</label>
            <select v-model="editingUser.role" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white">
              <option value="OPERATOR">Operador</option>
              <option value="ADMIN_TENANT">Administrador do Estabelecimento</option>
            </select>
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Redefinir Senha (opcional)</label>
            <input
              v-model="editPassword"
              type="password"
              placeholder="Deixe em branco para manter a senha atual"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Status</label>
            <select v-model="editingUser.is_active" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white">
              <option :value="true">Ativo</option>
              <option :value="false">Inativo</option>
            </select>
          </div>
        </div>

        <div class="p-5 border-t border-slate-800 flex items-center justify-end gap-3 bg-slate-950/60 rounded-b-2xl">
          <button @click="showEditModal = false" class="px-4 py-2 rounded-xl text-slate-400 hover:text-white text-sm font-medium">
            Cancelar
          </button>
          <button
            @click="submitEdit"
            :disabled="saving"
            class="px-5 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-black text-sm transition disabled:opacity-50"
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
  Users, UserPlus, Search, Loader2, Pencil, Power, X,
  ShieldCheck, Crown, UserCheck, Store
} from 'lucide-vue-next'
import api from '../../services/api'
import { useAuthStore, type Tenant } from '../../stores/auth'

interface UserItem {
  id: string
  tenant_id?: string | null
  name: string
  email: string
  role: string
  is_active: boolean
  created_at?: string
  tenant?: Tenant
}

const authStore = useAuthStore()

const users = ref<UserItem[]>([])
const tenantsList = ref<Tenant[]>([])
const loading = ref(true)
const saving = ref(false)

const searchTerm = ref('')
const selectedTenantFilter = ref('')
const selectedRoleFilter = ref<'ALL' | 'ADMIN_TENANT' | 'OPERATOR'>('ALL')

const showCreateModal = ref(false)
const createError = ref('')
const createForm = ref({
  tenant_id: '',
  name: '',
  email: '',
  password: '',
  role: 'OPERATOR',
})

const showEditModal = ref(false)
const editingUser = ref<UserItem | null>(null)
const editPassword = ref('')

const filteredUsers = computed(() => {
  return users.value.filter(u => {
    const matchesSearch =
      !searchTerm.value ||
      u.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      u.email.toLowerCase().includes(searchTerm.value.toLowerCase())

    let matchesRole = true
    if (selectedRoleFilter.value === 'ADMIN_TENANT') {
      matchesRole = u.role === 'ADMIN_TENANT' || u.role === 'ADMIN'
    } else if (selectedRoleFilter.value === 'OPERATOR') {
      matchesRole = u.role === 'OPERATOR'
    }

    return matchesSearch && matchesRole
  })
})

function getRoleLabel(role: string) {
  if (role === 'ADMIN_GLOBAL') return 'Admin Geral'
  if (role === 'ADMIN_TENANT' || role === 'ADMIN') return 'Administrador'
  return 'Operador'
}

function getRoleBadgeClass(role: string) {
  if (role === 'ADMIN_GLOBAL') {
    return 'bg-purple-950/80 text-purple-300 border-purple-700/60'
  }
  if (role === 'ADMIN_TENANT' || role === 'ADMIN') {
    return 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60'
  }
  return 'bg-sky-950/80 text-sky-300 border-sky-700/60'
}

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('pt-BR')
  } catch {
    return dateStr
  }
}

async function fetchUsers() {
  loading.value = true
  try {
    let url = '/admin/users'
    if (authStore.isGlobalAdmin && selectedTenantFilter.value) {
      url += `?tenant_id=${selectedTenantFilter.value}`
    }
    const res = await api.get(url)
    if (res.data.success) {
      users.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao buscar usuários:', err)
  } finally {
    loading.value = false
  }
}

async function fetchTenantsList() {
  if (!authStore.isGlobalAdmin) return
  try {
    const res = await api.get('/admin/tenants')
    if (res.data.success) {
      tenantsList.value = res.data.data
    }
  } catch (err) {
    console.error('Erro ao carregar lista de estabelecimentos:', err)
  }
}

function openCreateModal() {
  createError.value = ''
  createForm.value = {
    tenant_id: selectedTenantFilter.value || '',
    name: '',
    email: '',
    password: '',
    role: 'OPERATOR',
  }
  showCreateModal.value = true
}

async function submitCreate() {
  if (!createForm.value.name || !createForm.value.email || !createForm.value.password) {
    createError.value = 'Por favor, preencha todos os campos obrigatórios (*).'
    return
  }

  if (authStore.isGlobalAdmin && !createForm.value.tenant_id) {
    createError.value = 'Selecione um estabelecimento para vincular o usuário.'
    return
  }

  saving.value = true
  createError.value = ''
  try {
    const payload: any = {
      name: createForm.value.name,
      email: createForm.value.email,
      password: createForm.value.password,
      role: createForm.value.role,
    }
    if (authStore.isGlobalAdmin) {
      payload.tenant_id = createForm.value.tenant_id
    }

    const res = await api.post('/admin/users', payload)
    if (res.data.success) {
      showCreateModal.value = false
      await fetchUsers()
    }
  } catch (err: any) {
    createError.value = err.response?.data?.error || 'Erro ao cadastrar usuário.'
  } finally {
    saving.value = false
  }
}

function openEditModal(user: UserItem) {
  editingUser.value = { ...user }
  editPassword.value = ''
  showEditModal.value = true
}

async function submitEdit() {
  if (!editingUser.value) return
  saving.value = true
  try {
    const payload: any = {
      name: editingUser.value.name,
      email: editingUser.value.email,
      role: editingUser.value.role,
      is_active: editingUser.value.is_active,
    }
    if (editPassword.value) {
      payload.password = editPassword.value
    }

    const res = await api.put(`/admin/users/${editingUser.value.id}`, payload)
    if (res.data.success) {
      showEditModal.value = false
      await fetchUsers()
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao atualizar usuário.')
  } finally {
    saving.value = false
  }
}

async function toggleStatus(user: UserItem) {
  const nextStatus = !user.is_active
  const actionName = nextStatus ? 'ativar' : 'inativar'
  if (!confirm(`Deseja realmente ${actionName} o usuário "${user.name}"?`)) {
    return
  }

  try {
    const res = await api.patch(`/admin/users/${user.id}/status`, { is_active: nextStatus })
    if (res.data.success) {
      user.is_active = nextStatus
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao alterar status.')
  }
}

onMounted(() => {
  fetchUsers()
  fetchTenantsList()
})
</script>
