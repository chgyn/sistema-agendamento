<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Header & Busca -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-display">Base de Clientes</h1>
        <p class="text-xs sm:text-sm text-zinc-400 mt-1">
          Acompanhe o histórico de visitas, faturamento e dados de contato dos clientes
        </p>
      </div>

      <div class="w-full sm:w-80">
        <input
          v-model="search"
          @input="loadCustomers"
          type="text"
          placeholder="Buscar por nome, telefone ou email..."
          class="w-full px-4 py-2.5 rounded-xl bg-[#14151c] border border-zinc-750 text-white text-xs placeholder-zinc-500 focus:outline-none focus:ring-2 focus:ring-orange-500/25 focus:border-orange-500 transition"
        />
      </div>
    </div>

    <!-- Tabela de Clientes -->
    <div class="glass-panel rounded-2xl border border-zinc-800 overflow-hidden">
      <div v-if="isLoading" class="text-center py-16">
        <Loader2 class="w-8 h-8 animate-spin text-orange-500 mx-auto mb-2" />
        <p class="text-xs text-zinc-400">Carregando clientes...</p>
      </div>

      <div v-else-if="customers.length === 0" class="text-center py-16 px-4">
        <Users class="w-10 h-10 text-zinc-600 mx-auto mb-2" />
        <p class="text-sm font-bold text-zinc-300">Nenhum cliente encontrado</p>
        <p class="text-xs text-zinc-500 mt-1">Os clientes são cadastrados automaticamente ao realizar agendamentos.</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="w-full text-left text-xs text-zinc-300">
          <thead class="bg-[#121318]/90 text-zinc-400 uppercase font-bold text-[11px] tracking-wider border-b border-zinc-800">
            <tr>
              <th class="p-4">Cliente</th>
              <th class="p-4">WhatsApp / Telefone</th>
              <th class="p-4">E-mail</th>
              <th class="p-4 text-center">Total de Visitas</th>
              <th class="p-4 text-right">Total Investido</th>
              <th class="p-4 text-center">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-zinc-800/70">
            <tr v-for="c in customers" :key="c.id" class="hover:bg-[#16171e]/80 transition">
              <td class="p-4 font-bold text-white flex items-center space-x-3">
                <div class="w-9 h-9 rounded-xl bg-orange-500/15 border border-orange-500/30 text-orange-400 font-bold flex items-center justify-center font-display shrink-0">
                  {{ c.name.charAt(0) }}
                </div>
                <span class="font-display">{{ c.name }}</span>
              </td>
              <td class="p-4 font-mono font-semibold text-orange-400">{{ c.phone }}</td>
              <td class="p-4 text-zinc-400">{{ c.email || '—' }}</td>
              <td class="p-4 text-center font-bold text-white">
                <span class="px-2.5 py-1 rounded-lg bg-[#181922] border border-zinc-750 font-display">
                  {{ c.total_visits }} atendimentos
                </span>
              </td>
              <td class="p-4 text-right font-black text-emerald-400 font-display">
                R$ {{ (c.total_spent || 0).toFixed(2).replace('.', ',') }}
              </td>
              <td class="p-4 text-center">
                <button
                  @click="openCustomerHistory(c)"
                  class="px-3.5 py-1.5 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-200 hover:text-white border border-zinc-700/60 font-bold transition cursor-pointer"
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
    <div v-if="showHistoryModal" class="fixed inset-0 z-50 bg-black/80 backdrop-blur-md flex items-center justify-center p-4">
      <div class="glass-panel w-full max-w-lg rounded-2xl border border-zinc-800 p-6 shadow-2xl space-y-4">
        <div class="flex items-center justify-between border-b border-zinc-800 pb-3">
          <div>
            <h3 class="font-bold text-lg text-white font-display">Histórico do Cliente</h3>
            <p class="text-xs text-orange-400 font-bold">{{ selectedCustomer?.name }} ({{ selectedCustomer?.phone }})</p>
          </div>
          <button @click="showHistoryModal = false" class="text-zinc-400 hover:text-white text-lg font-bold">✕</button>
        </div>

        <div class="space-y-3 max-h-72 overflow-y-auto">
          <div v-if="customerAppointments.length === 0" class="text-center py-6 text-xs text-zinc-500">
            Nenhum histórico registrado para este cliente.
          </div>
          <div
            v-for="apt in customerAppointments"
            :key="apt.id"
            class="p-3.5 rounded-xl bg-[#14151c] border border-zinc-800 flex items-center justify-between text-xs"
          >
            <div>
              <span class="font-bold text-white block font-display">{{ apt.service?.name }}</span>
              <span class="text-zinc-400">{{ formatDate(apt.start_at) }} com {{ apt.professional?.name }}</span>
            </div>
            <div class="text-right">
              <span class="font-black text-orange-400 block font-display">R$ {{ apt.total_price?.toFixed(2).replace('.', ',') }}</span>
              <span class="text-[10px] text-zinc-400 font-semibold">{{ apt.status }}</span>
            </div>
          </div>
        </div>

        <div class="pt-3 border-t border-zinc-800 text-right">
          <button
            @click="showHistoryModal = false"
            class="px-4 py-2 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-300 text-xs font-bold border border-zinc-700/60"
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
