<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header & Busca -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-black text-white tracking-tight">Base de Clientes</h1>
        <p class="text-xs sm:text-sm text-slate-400 mt-1">
          Acompanhe o histórico de visitas, faturamento e dados de contato dos clientes
        </p>
      </div>

      <div class="w-full sm:w-72">
        <input
          v-model="search"
          @input="loadCustomers"
          type="text"
          placeholder="Buscar por nome, telefone ou email..."
          class="w-full px-4 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-white text-xs placeholder-slate-500 focus:ring-2 focus:ring-emerald-500"
        />
      </div>
    </div>

    <!-- Tabela de Clientes -->
    <div class="glass-panel rounded-2xl border border-slate-800 overflow-hidden">
      <div v-if="isLoading" class="text-center py-16">
        <Loader2 class="w-8 h-8 animate-spin text-emerald-500 mx-auto mb-2" />
        <p class="text-xs text-slate-400">Carregando clientes...</p>
      </div>

      <div v-else-if="customers.length === 0" class="text-center py-16 px-4">
        <Users class="w-10 h-10 text-slate-600 mx-auto mb-2" />
        <p class="text-sm font-semibold text-slate-300">Nenhum cliente encontrado</p>
        <p class="text-xs text-slate-500 mt-1">Os clientes são cadastrados automaticamente ao realizar agendamentos.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-slate-300">
          <thead class="bg-slate-900/80 text-slate-400 uppercase font-semibold border-b border-slate-800">
            <tr>
              <th class="p-4">Cliente</th>
              <th class="p-4">WhatsApp / Telefone</th>
              <th class="p-4">E-mail</th>
              <th class="p-4 text-center">Total de Visitas</th>
              <th class="p-4 text-right">Total Investido</th>
              <th class="p-4 text-center">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-slate-800/80">
            <tr v-for="c in customers" :key="c.id" class="hover:bg-slate-900/50 transition">
              <td class="p-4 font-bold text-white flex items-center space-x-2">
                <div class="w-8 h-8 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center justify-center">
                  {{ c.name.charAt(0) }}
                </div>
                <span>{{ c.name }}</span>
              </td>
              <td class="p-4 font-mono text-emerald-400">{{ c.phone }}</td>
              <td class="p-4 text-slate-400">{{ c.email || '—' }}</td>
              <td class="p-4 text-center font-bold text-white">
                <span class="px-2 py-1 rounded-md bg-slate-800 border border-slate-700">
                  {{ c.total_visits }} atendimentos
                </span>
              </td>
              <td class="p-4 text-right font-extrabold text-emerald-400">
                R$ {{ (c.total_spent || 0).toFixed(2) }}
              </td>
              <td class="p-4 text-center">
                <button
                  @click="openCustomerHistory(c)"
                  class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-200 font-semibold transition"
                >
                  Ver Histórico
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- MODAL: Histórico de Atendimentos do Cliente -->
    <div v-if="showHistoryModal" class="fixed inset-0 z-50 bg-black/70 backdrop-blur-sm flex items-center justify-center p-4">
      <div class="glass-panel w-full max-w-lg rounded-2xl border border-slate-800 p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-slate-800 pb-3">
          <div>
            <h3 class="font-bold text-lg text-white">Histórico do Cliente</h3>
            <p class="text-xs text-emerald-400">{{ selectedCustomer?.name }} ({{ selectedCustomer?.phone }})</p>
          </div>
          <button @click="showHistoryModal = false" class="text-slate-400 hover:text-white text-lg">✕</button>
        </div>

        <div class="space-y-3 max-h-72 overflow-y-auto">
          <div v-if="customerAppointments.length === 0" class="text-center py-6 text-xs text-slate-500">
            Nenhum histórico registrado para este cliente.
          </div>
          <div
            v-for="apt in customerAppointments"
            :key="apt.id"
            class="p-3 rounded-xl bg-slate-900 border border-slate-800 flex items-center justify-between text-xs"
          >
            <div>
              <span class="font-bold text-white block">{{ apt.service?.name }}</span>
              <span class="text-slate-400">{{ formatDate(apt.start_at) }} com {{ apt.professional?.name }}</span>
            </div>
            <div class="text-right">
              <span class="font-bold text-emerald-400 block">R$ {{ apt.total_price?.toFixed(2) }}</span>
              <span class="text-[10px] text-slate-400">{{ apt.status }}</span>
            </div>
          </div>
        </div>

        <div class="pt-3 border-t border-slate-800 text-right">
          <button
            @click="showHistoryModal = false"
            class="px-4 py-2 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold"
          >
            Fechar
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { Users, Loader2 } from 'lucide-vue-next'
import { format } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import api from '../../services/api'

const customers = ref<any[]>([])
const search = ref('')
const isLoading = ref(true)

const showHistoryModal = ref(false)
const selectedCustomer = ref<any>(null)
const customerAppointments = ref<any[]>([])

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

async function openCustomerHistory(customer: any) {
  selectedCustomer.value = customer
  try {
    const res = await api.get(`/admin/customers/${customer.id}`)
    if (res.data.success) {
      customerAppointments.value = res.data.data.appointments || []
    }
    showHistoryModal.value = true
  } catch (err) {
    alert('Erro ao carregar histórico')
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return ''
  try {
    return format(new Date(dateStr), "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })
  } catch {
    return dateStr
  }
}
</script>
