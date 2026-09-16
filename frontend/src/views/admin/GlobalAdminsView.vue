<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-purple-500/10 border border-purple-500/20 text-purple-400 text-xs font-semibold mb-2">
          <ShieldCheck class="w-3.5 h-3.5" />
          <span>Super Administradores da Plataforma</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-['Outfit'] flex items-center gap-2.5">
          <ShieldCheck class="w-8 h-8 text-purple-400" />
          <span>Administradores Gerais</span>
        </h1>
        <p class="text-xs sm:text-sm text-surface-400 mt-1">
          Usuários com permissão administrativa irrestrita sobre toda a plataforma multi-tenant.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 via-indigo-600 to-purple-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-900/30 hover:shadow-purple-700/40 transition active:scale-95 shrink-0"
      >
        <UserPlus class="w-4 h-4" />
        <span>Novo Administrador Geral</span>
      </button>
    </div>

    <!-- Alert de Segurança -->
    <div class="p-4 sm:p-5 rounded-2xl bg-purple-500/10 border border-purple-500/30 flex items-start gap-3.5 shadow-lg shadow-purple-950/20">
      <div class="w-9 h-9 rounded-xl bg-purple-500/20 border border-purple-500/40 flex items-center justify-center text-purple-300 shrink-0 mt-0.5">
        <ShieldAlert class="w-5 h-5" />
      </div>
      <div class="text-xs text-purple-200 leading-relaxed">
        <span class="font-bold text-white block text-sm mb-0.5 font-['Outfit']">Atenção ao controle de acesso central:</span>
        Os Administradores Gerais têm acesso total a todos os estabelecimentos, relatórios financeiros e configurações da plataforma. Conceda esse acesso apenas a membros confiáveis da equipe de administração central.
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="text-center py-20">
      <Loader2 class="w-9 h-9 text-purple-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-surface-400">Carregando administradores gerais...</p>
    </div>

    <!-- Tabela de Administradores Gerais -->
    <div v-else class="glass-card overflow-hidden shadow-2xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-surface-950/90 border-b border-surface-800 text-[11px] text-surface-400 uppercase tracking-wider font-bold">
            <tr>
              <th class="py-4 px-5">Administrador</th>
              <th class="py-4 px-5">Escopo</th>
              <th class="py-4 px-5">Status</th>
              <th class="py-4 px-5">Data de Cadastro</th>
              <th class="py-4 px-5 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-surface-800/60">
            <tr v-for="admin in admins" :key="admin.id" class="hover:bg-surface-800/40 transition">
              <td class="py-4 px-5">
                <div class="flex items-center gap-3.5">
                  <div class="w-10 h-10 rounded-2xl bg-purple-500/15 border border-purple-500/30 flex items-center justify-center font-bold text-sm text-purple-300 shrink-0 font-['Outfit'] shadow-sm">
                    {{ admin.name.charAt(0) }}
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-white text-sm truncate flex items-center gap-2 font-['Outfit']">
                      <span>{{ admin.name }}</span>
                      <span v-if="admin.id === authStore.user?.id" class="text-[10px] px-2 py-0.5 rounded-full bg-purple-500/20 text-purple-300 font-bold border border-purple-500/40">Você</span>
                    </p>
                    <p class="text-xs text-surface-400 truncate">{{ admin.email }}</p>
                  </div>
                </div>
              </td>

              <td class="py-4 px-5">
                <span class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold bg-purple-500/10 text-purple-300 border border-purple-500/30">
                  <Crown class="w-3 h-3 text-purple-400" />
                  <span>Acesso Global</span>
                </span>
              </td>

              <td class="py-4 px-5">
                <span
                  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold border"
                  :class="admin.is_active ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border-rose-500/20'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="admin.is_active ? 'bg-emerald-400' : 'bg-rose-400'"></span>
                  {{ admin.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </td>

              <td class="py-4 px-5 text-xs text-surface-400">
                {{ formatDate(admin.created_at) }}
              </td>

              <td class="py-4 px-5 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <button
                    @click="openEditModal(admin)"
                    title="Editar Administrador"
                    class="p-2 rounded-xl bg-surface-800 hover:bg-surface-700 text-surface-300 hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <button
                    v-if="admin.id !== authStore.user?.id"
                    @click="toggleStatus(admin)"
                    :title="admin.is_active ? 'Inativar Administrador' : 'Ativar Administrador'"
                    class="p-2 rounded-xl transition"
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
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-md shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <ShieldCheck class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base font-['Outfit']">Novo Administrador Geral</h3>
              <p class="text-xs text-surface-400">Cadastre outro administrador global da plataforma.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div v-if="createError" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ createError }}
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome Completo *</label>
            <input
              v-model="createForm.name"
              type="text"
              placeholder="Ex: Roberto Administrador"
              class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Acesso *</label>
            <input
              v-model="createForm.email"
              type="email"
              placeholder="roberto@plataforma.com"
              class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              required
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Senha de Acesso *</label>
            <input
              v-model="createForm.password"
              type="password"
              placeholder="Mínimo 6 caracteres"
              class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              required
            />
          </div>
        </div>

        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end gap-3 bg-surface-950/60 rounded-b-2xl">
          <button @click="showCreateModal = false" class="px-4 py-2.5 rounded-xl text-surface-400 hover:text-white text-sm font-medium transition">
            Cancelar
          </button>
          <button
            @click="submitCreate"
            :disabled="saving"
            class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white font-bold text-sm transition disabled:opacity-50 inline-flex items-center gap-2 shadow-lg shadow-purple-900/30"
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
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-md shadow-2xl my-8 flex flex-col">
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <Pencil class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base font-['Outfit']">Editar Administrador Geral</h3>
              <p class="text-xs text-surface-400">Atualize dados cadastrais ou redefina a senha.</p>
            </div>
          </div>
          <button @click="showEditModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <div class="p-6 space-y-4">
          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome Completo</label>
            <input v-model="editingAdmin.name" type="text" class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition" />
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">E-mail</label>
            <input v-model="editingAdmin.email" type="email" class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition" />
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Redefinir Senha (opcional)</label>
            <input
              v-model="editPassword"
              type="password"
              placeholder="Deixe em branco para manter a atual"
              class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Status</label>
            <select v-model="editingAdmin.is_active" class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition">
              <option :value="true">Ativo</option>
              <option :value="false">Inativo</option>
            </select>
          </div>
        </div>

        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end gap-3 bg-surface-950/60 rounded-b-2xl">
          <button @click="showEditModal = false" class="px-4 py-2.5 rounded-xl text-surface-400 hover:text-white text-sm font-medium transition">
            Cancelar
          </button>
          <button
            @click="submitEdit"
            :disabled="saving"
            class="px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white font-bold text-sm transition disabled:opacity-50 shadow-lg shadow-purple-900/30"
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
