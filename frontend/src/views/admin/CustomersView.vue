<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header & Ações -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h2 class="text-xl sm:text-2xl font-black text-[#202224] dark:text-white font-display tracking-tight">
          Customer List (Base de Clientes)
        </h2>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-zinc-400 mt-0.5">
          Acompanhe o histórico de visitas, faturamento e dados cadastrais dos clientes.
        </p>
      </div>

      <div class="flex items-center gap-3">
        <div class="w-full sm:w-72 relative">
          <Search class="w-3.5 h-3.5 text-gray-400 absolute left-3 top-3 pointer-events-none" />
          <input
            v-model="search"
            @input="loadCustomers"
            type="text"
            placeholder="Buscar nome, telefone ou e-mail..."
            class="w-full pl-8 pr-3 py-2 rounded-xl bg-white dark:bg-[#121318] border border-gray-200/80 dark:border-zinc-800 text-xs text-[#202224] dark:text-white placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF] transition shadow-sm"
          />
        </div>

        <button
          @click="openNewCustomerModal"
          class="flex items-center gap-1.5 px-4 py-2 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 shrink-0 cursor-pointer"
        >
          <Plus class="w-4 h-4" />
          <span>+ Add Customer</span>
        </button>
      </div>
    </div>

    <!-- Layout Grid: Tabela de Clientes + Painel Lateral de Detalhes (Figma 4:2803) -->
    <div class="grid grid-cols-1" :class="selectedCustomer ? 'lg:grid-cols-12 gap-6' : ''">
      <!-- Tabela Principal (7-8 colunas quando selecionado, ou 12 quando fechado) -->
      <div :class="selectedCustomer ? 'lg:col-span-8' : 'w-full'" class="saas-card overflow-hidden">
        <div v-if="isLoading" class="text-center py-16">
          <Loader2 class="w-8 h-8 animate-spin text-[#4880FF] mx-auto mb-2" />
          <p class="text-xs text-[#718096]">Carregando base de clientes...</p>
        </div>

        <div v-else-if="customers.length === 0" class="text-center py-16 px-4">
          <Users class="w-10 h-10 text-gray-300 dark:text-zinc-600 mx-auto mb-2" />
          <h4 class="font-bold text-sm text-[#202224] dark:text-white font-display">Nenhum cliente encontrado</h4>
          <p class="text-xs text-[#718096] dark:text-zinc-400 mt-1">
            Os clientes são cadastrados automaticamente ao realizar agendamentos ou manualmente pelo botão acima.
          </p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left text-xs text-[#202224] dark:text-zinc-300">
            <thead class="bg-gray-50/70 dark:bg-zinc-800/50 text-[#718096] dark:text-zinc-400 uppercase font-bold text-[11px] tracking-wider border-b border-gray-100 dark:border-zinc-800">
              <tr>
                <th class="p-4">Name</th>
                <th class="p-4">Email</th>
                <th class="p-4">Phone number</th>
                <th class="p-4 text-center">Status / Visitas</th>
                <th class="p-4 text-right">Total Spent</th>
                <th class="p-4 text-center">Action</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-zinc-800/80">
              <tr
                v-for="c in customers"
                :key="c.id"
                @click="inspectCustomer(c)"
                class="hover:bg-gray-50/60 dark:hover:bg-zinc-800/40 transition cursor-pointer"
                :class="selectedCustomer?.id === c.id ? 'bg-[#E9F0FE]/40 dark:bg-blue-950/20' : ''"
              >
                <!-- Avatar & Nome -->
                <td class="p-4 font-bold text-[#202224] dark:text-white flex items-center space-x-3">
                  <div class="w-9 h-9 rounded-full bg-[#E9F0FE] text-[#4880FF] font-black flex items-center justify-center font-display shrink-0 text-xs shadow-sm">
                    {{ c.name.charAt(0).toUpperCase() }}
                  </div>
                  <span class="font-display truncate max-w-[150px]">{{ c.name }}</span>
                </td>

                <!-- Email -->
                <td class="p-4 text-[#718096] dark:text-zinc-400 truncate max-w-[180px]">
                  {{ c.email || '—' }}
                </td>

                <!-- Telefone -->
                <td class="p-4 font-mono font-semibold text-[#202224] dark:text-zinc-200">
                  {{ c.phone }}
                </td>

                <!-- Pílula de Status / Frequência (Figma Pill) -->
                <td class="p-4 text-center">
                  <span class="px-2.5 py-1 rounded-full text-[10px] font-bold bg-[#E9F0FE] text-[#4880FF] border border-blue-200/80 dark:bg-blue-950/60 dark:text-blue-300">
                    {{ c.total_visits || 0 }} visitas
                  </span>
                </td>

                <!-- Total Investido -->
                <td class="p-4 text-right font-black text-[#202224] dark:text-white font-display">
                  R$ {{ (c.total_spent || 0).toFixed(2).replace('.', ',') }}
                </td>

                <!-- Botão 3 Pontos de Ação -->
                <td class="p-4 text-center" @click.stop>
                  <button
                    @click="inspectCustomer(c)"
                    class="p-1.5 rounded-lg text-gray-400 hover:text-[#4880FF] hover:bg-gray-100 dark:hover:bg-zinc-800 transition"
                    title="Ver Perfil"
                  >
                    <MoreHorizontal class="w-4 h-4" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- PAINEL LATERAL DE DETALHES DO CLIENTE (Figma 4:2803 Right Card) -->
      <div v-if="selectedCustomer" class="lg:col-span-4 saas-card p-6 space-y-6">
        <!-- Topo com botão fechar -->
        <div class="flex items-center justify-between">
          <span class="text-[11px] font-bold text-[#718096] uppercase tracking-wider">Perfil do Cliente</span>
          <button @click="selectedCustomer = null" class="text-gray-400 hover:text-gray-600 text-xs font-bold">✕ Fechar</button>
        </div>

        <!-- Avatar Central & Nome (Figma John Deo) -->
        <div class="text-center space-y-2">
          <div class="w-20 h-20 mx-auto rounded-full bg-gradient-to-tr from-[#4880FF] to-[#8280FF] text-white font-black text-2xl flex items-center justify-center shadow-lg shadow-blue-500/20">
            {{ selectedCustomer.name.charAt(0).toUpperCase() }}
          </div>
          <h3 class="font-black text-lg text-[#202224] dark:text-white font-display">{{ selectedCustomer.name }}</h3>
          <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Cliente Registrado</p>
        </div>

        <!-- Contact Info Section -->
        <div class="space-y-3 pt-3 border-t border-gray-100 dark:border-zinc-800">
          <h4 class="font-bold text-xs text-[#202224] dark:text-white uppercase tracking-wider font-display">Contact Info</h4>
          <div class="space-y-2 text-xs">
            <div class="flex items-center gap-3 text-[#718096] dark:text-zinc-300">
              <Mail class="w-4 h-4 text-[#4880FF] shrink-0" />
              <span class="truncate">{{ selectedCustomer.email || 'Não informado' }}</span>
            </div>
            <div class="flex items-center gap-3 text-[#718096] dark:text-zinc-300">
              <Phone class="w-4 h-4 text-[#4880FF] shrink-0" />
              <span>{{ selectedCustomer.phone }}</span>
            </div>
          </div>
        </div>

        <!-- Performance / Histórico Resumido (Figma Performance Box) -->
        <div class="pt-3 border-t border-gray-100 dark:border-zinc-800 space-y-3">
          <div class="flex items-center justify-between">
            <h4 class="font-bold text-xs text-[#202224] dark:text-white uppercase tracking-wider font-display">Performance</h4>
            <span class="text-[11px] text-[#4880FF] font-bold">{{ customerAppointments.length }} atendimentos</span>
          </div>

          <!-- Histórico em Lista -->
          <div class="space-y-2 max-h-52 overflow-y-auto pr-1">
            <div v-if="customerAppointments.length === 0" class="text-center py-4 text-xs text-[#718096]">
              Nenhum agendamento anterior registrado.
            </div>
            <div
              v-for="apt in customerAppointments"
              :key="apt.id"
              class="p-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800/60 border border-gray-100 dark:border-zinc-800 text-xs flex items-center justify-between"
            >
              <div>
                <span class="font-bold text-[#202224] dark:text-white block">{{ apt.service?.name }}</span>
                <span class="text-[10px] text-[#718096]">{{ formatDateShort(apt.start_at) }} com {{ apt.professional?.name }}</span>
              </div>
              <span class="font-black text-[#202224] dark:text-white font-display">R$ {{ (apt.total_price || 0).toFixed(2).replace('.', ',') }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- MODAL: Novo Cliente Manual -->
    <div v-if="showNewCustomerModal" class="fixed inset-0 z-50 bg-black/50 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="saas-card w-full max-w-md p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-gray-100 dark:border-zinc-800 pb-3">
          <h3 class="font-bold text-base text-[#202224] dark:text-white font-display">Novo Cliente</h3>
          <button @click="showNewCustomerModal = false" class="text-gray-400 hover:text-gray-600 text-lg font-bold">✕</button>
        </div>

        <form class="space-y-3.5" @submit.prevent="submitNewCustomer">
          <div>
            <label class="block text-xs font-bold text-[#718096] uppercase tracking-wider mb-1">Nome Completo *</label>
            <input
              v-model="newCustomerForm.name"
              type="text"
              required
              placeholder="Ex: João Silva"
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#718096] uppercase tracking-wider mb-1">WhatsApp / Telefone *</label>
            <input
              v-model="newCustomerForm.phone"
              type="text"
              required
              placeholder="(11) 99999-8888"
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div>
            <label class="block text-xs font-bold text-[#718096] uppercase tracking-wider mb-1">E-mail</label>
            <input
              v-model="newCustomerForm.email"
              type="email"
              placeholder="cliente@email.com"
              class="w-full px-3 py-2.5 rounded-xl bg-gray-50 dark:bg-zinc-800 border border-gray-200/80 dark:border-zinc-700 text-[#202224] dark:text-white text-xs focus:ring-2 focus:ring-[#4880FF]/20 focus:border-[#4880FF]"
            />
          </div>

          <div class="flex items-center justify-end gap-3 pt-3 border-t border-gray-100 dark:border-zinc-800">
            <button
              type="button"
              @click="showNewCustomerModal = false"
              class="px-4 py-2.5 rounded-xl bg-gray-100 hover:bg-gray-200 dark:bg-zinc-800 text-[#718096] text-xs font-bold transition"
            >
              Cancelar
            </button>
            <button
              type="submit"
              class="px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
            >
              Cadastrar Cliente
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref } from 'vue'
import { Users, Loader2, Search, Plus, MoreHorizontal, Mail, Phone } from 'lucide-vue-next'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '../../services/api'

const customers = ref<any[]>([])
const search = ref('')
const isLoading = ref(true)

const selectedCustomer = ref<any>(null)
const customerAppointments = ref<any[]>([])
const showNewCustomerModal = ref(false)

const newCustomerForm = reactive({
  name: '',
  phone: '',
  email: '',
})

onMounted(() => {
  loadCustomers()
})

async function loadCustomers() {
  isLoading.value = true
  try {
    const res = await api.get('/admin/customers', {
      params: { search: search.value },
    })
    if (res.data.success) {
      customers.value = res.data.data
    }
  } catch (err) {
    console.error(err)
  } finally {
    isLoading.value = false
  }
}

async function inspectCustomer(customer: any) {
  selectedCustomer.value = customer
  try {
    const res = await api.get(`/admin/customers/${customer.id}`)
    if (res.data.success) {
      customerAppointments.value = res.data.data.appointments || []
    }
  } catch (err) {
    console.error('Falha ao carregar atendimentos do cliente:', err)
  }
}

function openNewCustomerModal() {
  newCustomerForm.name = ''
  newCustomerForm.phone = ''
  newCustomerForm.email = ''
  showNewCustomerModal.value = true
}

async function submitNewCustomer() {
  showNewCustomerModal.value = false
  await loadCustomers()
}

function formatDateShort(dateStr: string) {
  if (!dateStr) return ''
  try {
    return format(new Date(dateStr), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })
  } catch {
    return dateStr
  }
}
</script>
