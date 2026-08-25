<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight flex items-center gap-2.5">
          <ShieldCheck class="w-7 h-7 text-purple-400" />
          <span>Administradores Gerais</span>
        </h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Usuários com permissão administrativa irrestrita sobre toda a plataforma multi-tenant.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-4 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-600/30 transition active:scale-95 shrink-0"
      >
        <UserPlus class="w-4 h-4" />
        <span>Novo Administrador Geral</span>
      </button>
    </div>

    <!-- Alert de Segurança -->
    <div class="p-4 rounded-2xl bg-purple-950/40 border border-purple-800/60 flex items-start gap-3.5">
      <ShieldAlert class="w-5 h-5 text-purple-400 shrink-0 mt-0.5" />
      <div class="text-xs text-purple-200 leading-relaxed">
        <span class="font-bold text-white">Atenção ao controle de acesso:</span> Os Administradores Gerais têm acesso total a todos os estabelecimentos, relatórios financeiros e configurações da plataforma. Conceda esse acesso apenas a membros confiáveis da equipe de administração central.
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-16">
      <Loader2 class="w-8 h-8 text-purple-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-slate-400">Carregando administradores gerais...</p>
    </div>

    <!-- Tabela de Administradores Gerais -->
    <div v-else class="bg-slate-900/80 border border-slate-800 rounded-2xl overflow-hidden shadow-xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-slate-950/80 border-b border-slate-800 text-xs text-slate-400 uppercase tracking-wider font-semibold">
            <tr>
              <th class="py-3.5 px-4">Administrador</th>
              <th class="py-3.5 px-4">Escopo</th>
              <th class="py-3.5 px-4">Status</th>
              <th class="py-3.5 px-4">Data de Cadastro</th>
              <th class="py-3.5 px-4 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/60">
            <tr v-for="admin in admins" :key="admin.id" class="hover:bg-slate-800/40 transition">
              <td class="py-4 px-4">
                <div class="flex items-center gap-3">
                  <div class="w-10 h-10 rounded-xl bg-purple-600/20 border border-purple-500/30 flex items-center justify-center font-bold text-sm text-purple-300 shrink-0">
                    {{ admin.name.charAt(0) }}
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-white text-sm truncate flex items-center gap-2">
                      <span>{{ admin.name }}</span>
                      <span v-if="admin.id === authStore.user?.id" class="text-[10px] px-2 py-0.2 rounded bg-purple-500/20 text-purple-300 font-bold border border-purple-500/40">Você</span>
                    </p>
                    <p class="text-xs text-slate-400 truncate">{{ admin.email }}</p>
                  </div>
                </div>
              </td>

              <td class="py-4 px-4">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-0.5 rounded-full text-[11px] font-bold bg-purple-950/80 text-purple-300 border border-purple-700/60">
                  <Crown class="w-3 h-3 text-purple-400" />
                  <span>Acesso Global</span>
                </span>
              </td>

              <td class="py-4 px-4">
                <span
                  class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded text-[10px] font-bold border"
                  :class="admin.is_active ? 'bg-emerald-950/80 text-emerald-300 border-emerald-700/60' : 'bg-rose-950/80 text-rose-300 border-rose-700/60'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="admin.is_active ? 'bg-emerald-400' : 'bg-rose-400'"></span>
                  {{ admin.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </td>

              <td class="py-4 px-4 text-xs text-slate-400">
                {{ formatDate(admin.created_at) }}
              </td>

              <td class="py-4 px-4 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="openEditModal(admin)"
                    title="Editar Administrador"
                    class="p-2 rounded-lg bg-slate-800/80 hover:bg-slate-700 text-slate-300 hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <button
                    v-if="admin.id !== authStore.user?.id"
                    @click="toggleStatus(admin)"
                    :title="admin.is_active ? 'Inativar Administrador' : 'Ativar Administrador'"
                    class="p-2 rounded-lg transition"
                    :class="admin.is_active ? 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-400' : 'bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-400'"
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

    <!-- MODAL: Novo Administrador Geral -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <div class="flex items-center gap-2.5">
            <div class="w-8 h-8 rounded-lg bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400">
              <ShieldCheck class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base">Novo Administrador Geral</h3>
              <p class="text-xs text-slate-400">Cadastre outro administrador global da plataforma.</p>
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

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nome Completo *</label>
            <input
              v-model="createForm.name"
              type="text"
              placeholder="Ex: Roberto Administrador"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">E-mail de Acesso *</label>
            <input
              v-model="createForm.email"
              type="email"
              placeholder="roberto@plataforma.com"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Senha de Acesso *</label>
            <input
              v-model="createForm.password"
              type="password"
              placeholder="Mínimo 6 caracteres"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white focus:outline-none focus:border-purple-500"
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
            class="px-5 py-2 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white font-bold text-sm transition disabled:opacity-50 inline-flex items-center gap-2"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>Cadastrar Administrador</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Editar Administrador Geral -->
    <div
      v-if="showEditModal && editingAdmin"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-slate-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-slate-900 border border-slate-800 rounded-2xl w-full max-w-md shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 border-b border-slate-800">
          <h3 class="font-bold text-white text-base">Editar Administrador Geral</h3>
          <button @click="showEditModal = false" class="p-1 rounded-lg text-slate-400 hover:text-white hover:bg-slate-800">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Nome Completo</label>
            <input v-model="editingAdmin.name" type="text" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">E-mail</label>
            <input v-model="editingAdmin.email" type="email" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white" />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Redefinir Senha (opcional)</label>
            <input
              v-model="editPassword"
              type="password"
              placeholder="Deixe em branco para manter a atual"
              class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white"
            />
          </div>

          <div>
            <label class="block text-xs font-medium text-slate-300 mb-1">Status</label>
            <select v-model="editingAdmin.is_active" class="w-full px-3 py-2 rounded-xl bg-slate-950 border border-slate-800 text-sm text-white">
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
            class="px-5 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-sm transition disabled:opacity-50"
          >
            Salvar Alterações
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import {
  ShieldCheck, ShieldAlert, UserPlus, Loader2, Pencil, Power, X, Crown
} from 'lucide-vue-next'
import api from '../../services/api'
import { useAuthStore, type User } from '../../stores/auth'

const authStore = useAuthStore()

const admins = ref<User[]>([])
const loading = ref(true)
const saving = ref(false)

const showCreateModal = ref(false)
const createError = ref('')
const createForm = ref({
  name: '',
  email: '',
  password: '',
})

const showEditModal = ref(false)
const editingAdmin = ref<User | null>(null)
const editPassword = ref('')

function formatDate(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleDateString('pt-BR')
  } catch {
    return dateStr
  }
}

async function fetchAdmins() {
  loading.value = true
  try {
    const res = await api.get('/admin/global-admins')
    if (res.data.success) {
      admins.value = res.data.data
    }
  } catch (err) {
    console.error('Erro ao carregar administradores gerais:', err)
  } finally {
    loading.value = false
  }
}

function openCreateModal() {
  createError.value = ''
  createForm.value = {
    name: '',
    email: '',
    password: '',
  }
  showCreateModal.value = true
}

async function submitCreate() {
  if (!createForm.value.name || !createForm.value.email || !createForm.value.password) {
    createError.value = 'Por favor, preencha todos os campos obrigatórios (*).'
    return
  }

  saving.value = true
  createError.value = ''
  try {
    const res = await api.post('/admin/global-admins', createForm.value)
    if (res.data.success) {
      showCreateModal.value = false
      await fetchAdmins()
    }
  } catch (err: any) {
    createError.value = err.response?.data?.error || 'Erro ao cadastrar administrador geral.'
  } finally {
    saving.value = false
  }
}

function openEditModal(admin: User) {
  editingAdmin.value = { ...admin }
  editPassword.value = ''
  showEditModal.value = true
}

async function submitEdit() {
  if (!editingAdmin.value) return
  saving.value = true
  try {
    const payload: any = {
      name: editingAdmin.value.name,
      email: editingAdmin.value.email,
      is_active: editingAdmin.value.is_active,
    }
    if (editPassword.value) {
      payload.password = editPassword.value
    }

    const res = await api.put(`/admin/users/${editingAdmin.value.id}`, payload)
    if (res.data.success) {
      showEditModal.value = false
      await fetchAdmins()
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao atualizar administrador.')
  } finally {
    saving.value = false
  }
}

async function toggleStatus(admin: User) {
  const nextStatus = !admin.is_active
  const actionName = nextStatus ? 'ativar' : 'inativar'
  if (!confirm(`Deseja realmente ${actionName} o administrador "${admin.name}"?`)) {
    return
  }

  try {
    const res = await api.patch(`/admin/users/${admin.id}/status`, { is_active: nextStatus })
    if (res.data.success) {
      admin.is_active = nextStatus
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao alterar status.')
  }
}

onMounted(() => {
  fetchAdmins()
})
</script>
