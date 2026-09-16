<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-purple-500/10 border border-purple-500/20 text-purple-400 text-xs font-semibold mb-2">
          <ShieldCheck class="w-3.5 h-3.5" />
          <span>Administração Multi-tenant Global</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-['Outfit'] flex items-center gap-2.5">
          <Building2 class="w-8 h-8 text-purple-400" />
          <span>Gestão de Estabelecimentos</span>
        </h1>
        <p class="text-xs sm:text-sm text-surface-400 mt-1">
          Cadastre e gerencie todos os estabelecimentos, planos contratados e o status financeiro de suas assinaturas.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 via-indigo-600 to-purple-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-900/30 hover:shadow-purple-700/40 transition active:scale-95 shrink-0"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Estabelecimento</span>
      </button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="glass-card p-5 flex items-center gap-4 group hover:border-purple-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 group-hover:scale-105 transition-transform shrink-0">
          <Building2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Total Estabelecimentos</p>
          <p class="text-2xl font-black text-white font-['Outfit'] mt-0.5">{{ tenants.length }}</p>
        </div>
      </div>

      <div class="glass-card p-5 flex items-center gap-4 group hover:border-emerald-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 group-hover:scale-105 transition-transform shrink-0">
          <CheckCircle2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Assinaturas Ativas</p>
          <p class="text-2xl font-black text-emerald-400 font-['Outfit'] mt-0.5">{{ activeSubscriptionsCount }}</p>
        </div>
      </div>

      <div class="glass-card p-5 flex items-center gap-4 group hover:border-amber-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400 group-hover:scale-105 transition-transform shrink-0">
          <Clock class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Pagamento Pendente</p>
          <p class="text-2xl font-black text-amber-400 font-['Outfit'] mt-0.5">{{ pendingSubscriptionsCount }}</p>
        </div>
      </div>

      <div class="glass-card p-5 flex items-center gap-4 group hover:border-rose-500/30 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-rose-500/10 border border-rose-500/20 flex items-center justify-center text-rose-400 group-hover:scale-105 transition-transform shrink-0">
          <AlertTriangle class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-surface-400 font-medium">Inadimplentes / Atrasadas</p>
          <p class="text-2xl font-black text-rose-400 font-['Outfit'] mt-0.5">{{ overdueSubscriptionsCount }}</p>
        </div>
      </div>
    </div>

    <!-- Filtros & Busca -->
    <div class="glass-card p-3 sm:p-4 flex flex-col md:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-surface-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome, slug, cidade ou e-mail..."
          class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-surface-100 placeholder-surface-500 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        />
      </div>

      <div class="flex flex-col sm:flex-row items-center gap-2 w-full md:w-auto">
        <!-- Filtro por Status da Assinatura -->
        <select
          v-model="subscriptionFilter"
          class="w-full sm:w-56 py-2.5 px-3.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-surface-200 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        >
          <option value="ALL">Todas as Assinaturas</option>
          <option value="ACTIVE">Assinatura Ativa 🟢</option>
          <option value="PENDING">Aguardando Pagamento 🟡</option>
          <option value="OVERDUE">Vencida / Atrasada 🔴</option>
          <option value="CANCELLED">Cancelada ⚫</option>
          <option value="NO_SUB">Sem Assinatura</option>
        </select>

        <!-- Filtro por Status do Estabelecimento -->
        <select
          v-model="statusFilter"
          class="w-full sm:w-44 py-2.5 px-3.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-surface-200 focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        >
          <option value="ALL">Status: Todos</option>
          <option value="ACTIVE">Apenas Ativos</option>
          <option value="INACTIVE">Apenas Inativos</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-20">
      <Loader2 class="w-9 h-9 text-purple-400 animate-spin mx-auto mb-3" />
      <p class="text-sm text-surface-400">Carregando estabelecimentos...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredTenants.length === 0" class="glass-card p-12 text-center">
      <div class="w-14 h-14 rounded-2xl bg-surface-800/80 border border-surface-700/60 flex items-center justify-center mx-auto mb-4 text-surface-500">
        <Building2 class="w-7 h-7" />
      </div>
      <h3 class="text-base font-bold text-white font-['Outfit']">Nenhum estabelecimento encontrado</h3>
      <p class="text-xs text-surface-400 mt-1.5 max-w-sm mx-auto">
        {{ searchTerm || subscriptionFilter !== 'ALL' || statusFilter !== 'ALL' ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Comece cadastrando o primeiro estabelecimento da plataforma.' }}
      </p>
    </div>

    <!-- Tabela de Estabelecimentos com Assinatura -->
    <div v-else class="glass-card overflow-hidden shadow-2xl">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-surface-950/90 border-b border-surface-800 text-[11px] text-surface-400 uppercase tracking-wider font-bold">
            <tr>
              <th class="py-4 px-5">Estabelecimento</th>
              <th class="py-4 px-5">Plano & Assinatura</th>
              <th class="py-4 px-5">Contato</th>
              <th class="py-4 px-5">Localização</th>
              <th class="py-4 px-5">Status Conta</th>
              <th class="py-4 px-5 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-surface-800/60">
            <tr v-for="t in filteredTenants" :key="t.id" class="hover:bg-surface-800/40 transition">
              <!-- Estabelecimento -->
              <td class="py-4 px-5">
                <div class="flex items-center gap-3.5">
                  <div
                    class="w-11 h-11 rounded-2xl flex items-center justify-center font-black text-sm text-white shrink-0 overflow-hidden border border-surface-700/60 bg-surface-800 shadow-sm"
                    :style="t.logo_url ? '' : { backgroundColor: t.primary_color || '#8b5cf6' }"
                  >
                    <img v-if="t.logo_url" :src="t.logo_url" :alt="t.name" class="w-full h-full object-cover" />
                    <span v-else class="font-['Outfit'] text-base">{{ t.name.charAt(0) }}</span>
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-white text-sm truncate font-['Outfit']">{{ t.name }}</p>
                    <p class="text-xs text-purple-400 font-mono truncate">/{{ t.slug }}</p>
                  </div>
                </div>
              </td>

              <!-- Plano & Assinatura -->
              <td class="py-4 px-5">
                <div v-if="t.subscription" class="space-y-1">
                  <div class="flex items-center gap-2">
                    <span class="font-bold text-white text-xs">
                      {{ t.subscription.plan?.name || 'Plano Personalizado' }}
                    </span>
                    <span
                      class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border inline-flex items-center gap-1.5"
                      :class="getSubscriptionBadgeClass(t.subscription.status)"
                    >
                      <span class="w-1.5 h-1.5 rounded-full" :class="getSubscriptionDotClass(t.subscription.status)"></span>
                      <span>{{ formatSubStatus(t.subscription.status) }}</span>
                    </span>
                  </div>
                  <div class="text-[11px] text-surface-400 flex items-center gap-1.5 font-medium">
                    <span class="font-mono text-surface-300">R$ {{ t.subscription.price.toFixed(2).replace('.', ',') }}</span>
                    <span>•</span>
                    <span>{{ formatCycle(t.subscription.billing_cycle) }}</span>
                  </div>
                </div>
                <div v-else class="text-xs text-surface-500 italic">
                  Sem assinatura ativa
                </div>
              </td>

              <!-- Contato -->
              <td class="py-4 px-5 text-xs text-surface-300">
                <div class="font-medium text-white">{{ t.phone || '-' }}</div>
                <div class="text-surface-400 truncate mt-0.5">{{ t.email || '-' }}</div>
              </td>

              <!-- Localização -->
              <td class="py-4 px-5 text-xs text-surface-300">
                <div v-if="t.city || t.state" class="font-medium">{{ t.city }} - {{ t.state }}</div>
                <div v-else class="text-surface-500">Não informado</div>
              </td>

              <!-- Status Conta -->
              <td class="py-4 px-5">
                <span
                  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold border"
                  :class="t.is_active ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-400 border-rose-500/20'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="t.is_active ? 'bg-emerald-400' : 'bg-rose-400'"></span>
                  {{ t.is_active ? 'Ativo' : 'Inativo' }}
                </span>
              </td>

              <!-- Ações -->
              <td class="py-4 px-5 text-right">
                <div class="flex items-center justify-end gap-1.5">
                  <!-- Ver Detalhamento da Assinatura -->
                  <button
                    @click="openSubscriptionModal(t)"
                    title="Ver detalhamento da assinatura e faturas"
                    class="p-2 rounded-xl bg-purple-500/10 hover:bg-purple-500/20 text-purple-300 hover:text-white transition"
                  >
                    <CreditCard class="w-4 h-4" />
                  </button>

                  <!-- Link Público -->
                  <a
                    :href="`/agendamento/${t.slug}`"
                    target="_blank"
                    title="Ver página de agendamento público"
                    class="p-2 rounded-xl bg-surface-800 hover:bg-surface-700 text-surface-300 hover:text-white transition"
                  >
                    <ExternalLink class="w-4 h-4" />
                  </a>

                  <!-- Editar Estabelecimento -->
                  <button
                    @click="openEditModal(t)"
                    title="Editar informações"
                    class="p-2 rounded-xl bg-surface-800 hover:bg-surface-700 text-surface-300 hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <!-- Alternar Status -->
                  <button
                    @click="toggleStatus(t)"
                    :title="t.is_active ? 'Inativar estabelecimento' : 'Ativar estabelecimento'"
                    class="p-2 rounded-xl transition"
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

    <!-- MODAL: Detalhamento Completo da Assinatura do Estabelecimento -->
    <div
      v-if="showSubscriptionModal && selectedTenant"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-3xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <CreditCard class="w-5 h-5" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-white text-base font-['Outfit']">{{ selectedTenant.name }}</h3>
                <span class="text-xs text-purple-400 font-mono">/{{ selectedTenant.slug }}</span>
              </div>
              <p class="text-xs text-surface-400">Detalhamento contratual, financeiro e faturas Asaas</p>
            </div>
          </div>
          <button @click="showSubscriptionModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-6">
          <div v-if="loadingSubscriptionDetails" class="py-12 text-center">
            <Loader2 class="w-8 h-8 text-purple-400 animate-spin mx-auto mb-3" />
            <p class="text-xs text-surface-400">Carregando dados da assinatura...</p>
          </div>

          <div v-else-if="!selectedSubscription" class="bg-surface-950/60 border border-surface-800 rounded-2xl p-8 text-center space-y-2">
            <AlertTriangle class="w-8 h-8 text-amber-400 mx-auto" />
            <p class="text-sm font-bold text-white">Nenhuma assinatura vinculada a este estabelecimento</p>
            <p class="text-xs text-surface-400">Este estabelecimento não possui registro de assinatura recorrente no momento.</p>
          </div>

          <div v-else class="space-y-6">
            <!-- Banner de Status da Assinatura -->
            <div
              class="p-4 rounded-2xl border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 shadow-md"
              :class="getStatusBannerClass(selectedSubscription.status)"
            >
              <div class="flex items-center gap-3.5">
                <div class="w-10 h-10 rounded-2xl flex items-center justify-center shrink-0" :class="getStatusIconBoxClass(selectedSubscription.status)">
                  <CheckCircle2 v-if="selectedSubscription.status === 'ACTIVE'" class="w-5 h-5 text-emerald-400" />
                  <Clock v-else-if="selectedSubscription.status === 'PENDING'" class="w-5 h-5 text-amber-400 animate-pulse" />
                  <AlertCircle v-else class="w-5 h-5 text-rose-400" />
                </div>
                <div>
                  <div class="flex items-center gap-2">
                    <span class="text-xs font-bold text-white">Status da Assinatura:</span>
                    <span class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border" :class="getSubscriptionBadgeClass(selectedSubscription.status)">
                      {{ formatSubStatus(selectedSubscription.status) }}
                    </span>
                  </div>
                  <p class="text-xs text-surface-300 mt-0.5">
                    {{ getStatusExplanation(selectedSubscription.status) }}
                  </p>
                </div>
              </div>

              <!-- Ação de Abrir Checkout/Fatura no Asaas -->
              <a
                v-if="selectedSubscription.payment_url"
                :href="selectedSubscription.payment_url"
                target="_blank"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-surface-950 font-bold text-xs shadow-lg shadow-emerald-500/20 transition shrink-0"
              >
                <ExternalLink class="w-3.5 h-3.5" />
                <span>Abrir Fatura Asaas</span>
              </a>
            </div>

            <!-- Grid de Informações Contratuais -->
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Card 1: Detalhes do Plano -->
              <div class="bg-surface-950/80 border border-surface-800/80 rounded-2xl p-4 space-y-3">
                <div class="flex items-center justify-between pb-2.5 border-b border-surface-800">
                  <div class="flex items-center gap-2">
                    <Sparkles class="w-4 h-4 text-purple-400" />
                    <h4 class="text-xs font-bold uppercase tracking-wider text-purple-300">Plano Contratado</h4>
                  </div>
                  <span class="px-2.5 py-0.5 rounded-full bg-purple-500/10 text-purple-300 border border-purple-500/30 text-[10px] font-bold">
                    {{ selectedSubscription.plan?.name || 'Personalizado' }}
                  </span>
                </div>

                <div class="grid grid-cols-2 gap-3 text-xs">
                  <div>
                    <span class="text-surface-400 block text-[11px]">Valor Recorrente:</span>
                    <span class="text-sm font-black text-white font-mono">
                      R$ {{ selectedSubscription.price.toFixed(2).replace('.', ',') }}
                    </span>
                  </div>
                  <div>
                    <span class="text-surface-400 block text-[11px]">Periodicidade:</span>
                    <span class="text-sm font-bold text-white">
                      {{ formatCycle(selectedSubscription.billing_cycle) }}
                    </span>
                  </div>
                  <div>
                    <span class="text-surface-400 block text-[11px]">Data de Início:</span>
                    <span class="text-xs font-medium text-surface-200">
                      {{ formatDate(selectedSubscription.created_at) }}
                    </span>
                  </div>
                  <div>
                    <span class="text-surface-400 block text-[11px]">Próxima Cobrança:</span>
                    <span class="text-xs font-bold text-emerald-400">
                      {{ formatDate(selectedSubscription.next_due_date) }}
                    </span>
                  </div>
                </div>
              </div>

              <!-- Card 2: Integração com Asaas -->
              <div class="bg-surface-950/80 border border-surface-800/80 rounded-2xl p-4 space-y-3 flex flex-col justify-between">
                <div>
                  <div class="flex items-center gap-2 pb-2.5 border-b border-surface-800">
                    <ShieldCheck class="w-4 h-4 text-indigo-400" />
                    <h4 class="text-xs font-bold uppercase tracking-wider text-indigo-300">Identificadores Asaas</h4>
                  </div>

                  <div class="mt-3 space-y-2 text-xs">
                    <div>
                      <span class="text-surface-400 block text-[11px]">Asaas Subscription ID:</span>
                      <span class="font-mono text-xs text-indigo-200 select-all">
                        {{ selectedSubscription.asaas_subscription_id || 'Não integrado' }}
                      </span>
                    </div>
                    <div>
                      <span class="text-surface-400 block text-[11px]">Asaas Customer ID:</span>
                      <span class="font-mono text-xs text-indigo-200 select-all">
                        {{ selectedSubscription.asaas_customer_id || 'Não integrado' }}
                      </span>
                    </div>
                  </div>
                </div>

                <div class="text-[11px] text-surface-400 flex items-center gap-1.5 pt-2 border-t border-surface-800/60">
                  <span>Forma de Pagamento:</span>
                  <span class="font-bold text-white uppercase">{{ selectedSubscription.payment_method || 'PIX / Cartão' }}</span>
                </div>
              </div>
            </div>

            <!-- Histórico de Faturas e Cobranças Sincronizadas -->
            <div class="bg-surface-950/90 border border-surface-800/80 rounded-2xl overflow-hidden shadow">
              <div class="p-4 border-b border-surface-800 flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <FileText class="w-4 h-4 text-emerald-400" />
                  <h4 class="text-xs font-bold uppercase tracking-wider text-surface-200">Histórico de Cobranças & Faturas</h4>
                </div>
                <span class="text-[11px] text-surface-500">Sincronizado via Webhook</span>
              </div>

              <div v-if="!selectedSubscription.invoices || selectedSubscription.invoices.length === 0" class="p-6 text-center text-xs text-surface-400">
                Nenhuma fatura registrada no histórico até o momento.
              </div>

              <div v-else class="overflow-x-auto max-h-48">
                <table class="w-full text-left text-xs">
                  <thead class="bg-surface-900 border-b border-surface-800 text-[11px] text-surface-400 uppercase font-semibold sticky top-0">
                    <tr>
                      <th class="py-2.5 px-3.5">Fatura</th>
                      <th class="py-2.5 px-3.5">Vencimento</th>
                      <th class="py-2.5 px-3.5">Valor</th>
                      <th class="py-2.5 px-3.5">Método</th>
                      <th class="py-2.5 px-3.5">Status</th>
                      <th class="py-2.5 px-3.5 text-right">Comprovante</th>
                    </tr>
                  </thead>
                  <tbody class="divide-y divide-surface-800/60">
                    <tr v-for="inv in selectedSubscription.invoices" :key="inv.id" class="hover:bg-surface-800/40 transition">
                      <td class="py-2.5 px-3.5 font-mono text-surface-300">
                        {{ inv.asaas_payment_id }}
                      </td>
                      <td class="py-2.5 px-3.5 text-surface-300">
                        {{ formatDate(inv.due_date) }}
                      </td>
                      <td class="py-2.5 px-3.5 font-bold text-white font-mono">
                        R$ {{ inv.value.toFixed(2).replace('.', ',') }}
                      </td>
                      <td class="py-2.5 px-3.5 text-surface-400 uppercase">
                        {{ inv.billing_type || 'PIX' }}
                      </td>
                      <td class="py-2.5 px-3.5">
                        <span class="px-2 py-0.5 rounded-full text-[10px] font-bold border" :class="getInvoiceBadgeClass(inv.status)">
                          {{ inv.status }}
                        </span>
                      </td>
                      <td class="py-2.5 px-3.5 text-right">
                        <a
                          v-if="inv.invoice_url"
                          :href="inv.invoice_url"
                          target="_blank"
                          class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-surface-800 hover:bg-surface-700 text-surface-300 hover:text-white transition font-medium"
                        >
                          <ExternalLink class="w-3 h-3" />
                          <span>Abrir</span>
                        </a>
                        <span v-else class="text-surface-500">-</span>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Ajuste Manual de Status (Exclusivo Admin Geral) -->
            <div class="p-4 rounded-2xl bg-purple-500/5 border border-purple-500/20 space-y-3">
              <div class="flex items-center justify-between">
                <div class="flex items-center gap-2">
                  <SlidersHorizontal class="w-4 h-4 text-purple-400" />
                  <h4 class="text-xs font-bold uppercase tracking-wider text-purple-300">Ajuste Manual de Status (Admin Geral)</h4>
                </div>
                <span class="text-[10px] text-surface-400">Controle direto de liberação de acesso</span>
              </div>

              <div class="flex flex-col sm:flex-row items-center gap-3">
                <select
                  v-model="overrideStatusForm"
                  class="w-full sm:w-64 py-2.5 px-3.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-xs text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20"
                >
                  <option value="ACTIVE">ACTIVE (Liberado / Ativo 🟢)</option>
                  <option value="PENDING">PENDING (Aguardando Pagamento 🟡)</option>
                  <option value="OVERDUE">OVERDUE (Inadimplente / Bloqueado 🔴)</option>
                  <option value="CANCELLED">CANCELLED (Cancelado ⚫)</option>
                  <option value="TRIAL">TRIAL (Período de Testes 🚀)</option>
                </select>

                <button
                  @click="applyStatusOverride"
                  :disabled="updatingStatus || overrideStatusForm === selectedSubscription.status"
                  class="w-full sm:w-auto inline-flex items-center justify-center gap-1.5 px-4 py-2.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white text-xs font-bold transition disabled:opacity-50 shadow-md"
                >
                  <Loader2 v-if="updatingStatus" class="w-3.5 h-3.5 animate-spin" />
                  <span>Salvar Alteração de Status</span>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showSubscriptionModal = false"
            class="px-5 py-2.5 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 text-sm font-medium transition"
          >
            Fechar
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar Novo Estabelecimento -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-2xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <Plus class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base font-['Outfit']">Novo Estabelecimento & Administrador</h3>
              <p class="text-xs text-surface-400">Cadastre a empresa e credenciais de acesso inicial.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-6">
          <div v-if="createError" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-300 text-xs">
            {{ createError }}
          </div>

          <!-- 1. Dados do Estabelecimento -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold uppercase tracking-wider text-purple-400 flex items-center gap-1.5">
              <Building2 class="w-3.5 h-3.5" />
              <span>1. Dados do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
                <input
                  v-model="createForm.name"
                  @input="autoGenerateSlug"
                  type="text"
                  required
                  placeholder="Ex: Barbearia Prime"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Slug de Acesso *</label>
                <input
                  v-model="createForm.slug"
                  type="text"
                  required
                  placeholder="barbearia-prime"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white font-mono focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
                <input
                  v-model="createForm.phone"
                  type="text"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
                <input
                  v-model="createForm.document"
                  type="text"
                  placeholder="00.000.000/0001-00"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Cidade</label>
                <input
                  v-model="createForm.city"
                  type="text"
                  placeholder="São Paulo"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">UF</label>
                <input
                  v-model="createForm.state"
                  type="text"
                  placeholder="SP"
                  maxlength="2"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white uppercase focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>
            </div>
          </div>

          <!-- 2. Administrador do Estabelecimento -->
          <div class="space-y-3 pt-4 border-t border-surface-800">
            <h4 class="text-xs font-bold uppercase tracking-wider text-purple-400 flex items-center gap-1.5">
              <ShieldCheck class="w-3.5 h-3.5" />
              <span>2. Administrador do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome do Administrador *</label>
                <input
                  v-model="createForm.admin_name"
                  type="text"
                  required
                  placeholder="Carlos Silva"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Acesso *</label>
                <input
                  v-model="createForm.admin_email"
                  type="email"
                  required
                  placeholder="admin@barbeariaprime.com"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Senha de Acesso *</label>
                <input
                  v-model="createForm.admin_password"
                  type="password"
                  required
                  minlength="6"
                  placeholder="••••••••"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end gap-3 bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showCreateModal = false"
            class="px-4 py-2.5 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 text-sm font-medium transition"
          >
            Cancelar
          </button>
          <button
            @click="submitCreate"
            :disabled="saving"
            class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold transition disabled:opacity-50 shadow-lg shadow-purple-900/30"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>Cadastrar Estabelecimento</span>
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Editar Estabelecimento -->
    <div
      v-if="showEditModal && editingTenant"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-surface-950/80 backdrop-blur-md overflow-y-auto"
    >
      <div class="glass-card-elevated w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-400 shrink-0">
              <Pencil class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-white text-base font-['Outfit']">Editar Estabelecimento</h3>
              <p class="text-xs text-surface-400">Atualize dados cadastrais e de contato.</p>
            </div>
          </div>
          <button @click="showEditModal = false" class="p-1.5 rounded-lg text-surface-400 hover:text-white hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
              <input
                v-model="editingTenant.name"
                type="text"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone</label>
              <input
                v-model="editingTenant.phone"
                type="text"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Contato</label>
              <input
                v-model="editingTenant.email"
                type="email"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Cidade</label>
              <input
                v-model="editingTenant.city"
                type="text"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-surface-300 uppercase tracking-wider mb-1.5">Estado (UF)</label>
              <input
                v-model="editingTenant.state"
                type="text"
                maxlength="2"
                class="w-full px-3.5 py-2.5 rounded-xl bg-surface-950/80 border border-surface-700/60 text-sm text-white uppercase focus:outline-none focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-surface-800/80 flex items-center justify-end gap-3 bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showEditModal = false"
            class="px-4 py-2.5 rounded-xl text-surface-400 hover:text-white hover:bg-surface-800 text-sm font-medium transition"
          >
            Cancelar
          </button>
          <button
            @click="submitEdit"
            :disabled="saving"
            class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 to-indigo-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold transition disabled:opacity-50 shadow-lg shadow-purple-900/30"
          >
            <Loader2 v-if="saving" class="w-4 h-4 animate-spin" />
            <span>Salvar Alterações</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Building2, Plus, Search, CheckCircle2, Loader2,
  ExternalLink, Pencil, Power, X, ShieldCheck, CreditCard,
  Clock, AlertTriangle, AlertCircle, Sparkles, FileText, SlidersHorizontal
} from 'lucide-vue-next'
import api from '../../services/api'
import type { Tenant, Subscription, SubscriptionStatus, PlanBillingCycle } from '../../stores/auth'

const tenants = ref<Tenant[]>([])
const loading = ref(true)
const saving = ref(false)
const searchTerm = ref('')
const statusFilter = ref<'ALL' | 'ACTIVE' | 'INACTIVE'>('ALL')
const subscriptionFilter = ref<'ALL' | 'ACTIVE' | 'PENDING' | 'OVERDUE' | 'CANCELLED' | 'NO_SUB'>('ALL')

// Modais
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showSubscriptionModal = ref(false)

const createError = ref('')
const editingTenant = ref<Tenant | null>(null)
const selectedTenant = ref<Tenant | null>(null)
const selectedSubscription = ref<Subscription | null>(null)
const loadingSubscriptionDetails = ref(false)
const updatingStatus = ref(false)
const overrideStatusForm = ref<SubscriptionStatus>('ACTIVE')

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

// KPIs
const activeSubscriptionsCount = computed(() =>
  tenants.value.filter(t => t.subscription?.status === 'ACTIVE' || t.subscription?.status === 'TRIAL').length
)
const pendingSubscriptionsCount = computed(() =>
  tenants.value.filter(t => t.subscription?.status === 'PENDING').length
)
const overdueSubscriptionsCount = computed(() =>
  tenants.value.filter(t => t.subscription?.status === 'OVERDUE').length
)

// Filtro Composto
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

    let matchesSub = true
    if (subscriptionFilter.value === 'ACTIVE') {
      matchesSub = t.subscription?.status === 'ACTIVE' || t.subscription?.status === 'TRIAL'
    } else if (subscriptionFilter.value === 'PENDING') {
      matchesSub = t.subscription?.status === 'PENDING'
    } else if (subscriptionFilter.value === 'OVERDUE') {
      matchesSub = t.subscription?.status === 'OVERDUE'
    } else if (subscriptionFilter.value === 'CANCELLED') {
      matchesSub = t.subscription?.status === 'CANCELLED'
    } else if (subscriptionFilter.value === 'NO_SUB') {
      matchesSub = !t.subscription
    }

    return matchesSearch && matchesStatus && matchesSub
  })
})

function formatSubStatus(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'Ativa'
    case 'PENDING':
      return 'Pendente'
    case 'OVERDUE':
      return 'Vencida'
    case 'CANCELLED':
      return 'Cancelada'
    case 'TRIAL':
      return 'Trial'
    default:
      return status || 'Pendente'
  }
}

function getSubscriptionBadgeClass(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
    case 'PENDING':
      return 'bg-amber-500/10 text-amber-400 border-amber-500/20'
    case 'OVERDUE':
      return 'bg-rose-500/10 text-rose-400 border-rose-500/20'
    case 'CANCELLED':
      return 'bg-surface-800 text-surface-400 border-surface-700'
    case 'TRIAL':
      return 'bg-sky-500/10 text-sky-400 border-sky-500/20'
    default:
      return 'bg-surface-800 text-surface-400 border-surface-700'
  }
}

function getSubscriptionDotClass(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'bg-emerald-400'
    case 'PENDING':
      return 'bg-amber-400'
    case 'OVERDUE':
      return 'bg-rose-400'
    case 'CANCELLED':
      return 'bg-surface-500'
    case 'TRIAL':
      return 'bg-sky-400'
    default:
      return 'bg-surface-500'
  }
}

function getStatusBannerClass(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'bg-emerald-500/10 border-emerald-500/30'
    case 'PENDING':
      return 'bg-amber-500/10 border-amber-500/30'
    case 'OVERDUE':
      return 'bg-rose-500/10 border-rose-500/30'
    default:
      return 'bg-surface-900 border-surface-800'
  }
}

function getStatusIconBoxClass(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'bg-emerald-500/10 border border-emerald-500/20'
    case 'PENDING':
      return 'bg-amber-500/10 border border-amber-500/20'
    default:
      return 'bg-rose-500/10 border border-rose-500/20'
  }
}

function getStatusExplanation(status?: SubscriptionStatus) {
  switch (status) {
    case 'ACTIVE':
      return 'Assinatura regularizada. Estabelecimento possui acesso total a todos os recursos da plataforma.'
    case 'PENDING':
      return 'Assinatura aguardando confirmação do pagamento inicial no Asaas.'
    case 'OVERDUE':
      return 'Inadimplente. O estabelecimento possui cobrança vencida e o acesso operacional está bloqueado pelo Gatekeeper.'
    case 'CANCELLED':
      return 'Assinatura cancelada. O estabelecimento não tem mais permissão de uso.'
    case 'TRIAL':
      return 'Estabelecimento operando em período gratuito de avaliação (Trial).'
    default:
      return ''
  }
}

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
      return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
    case 'PENDING':
      return 'bg-amber-500/10 text-amber-400 border-amber-500/20'
    default:
      return 'bg-rose-500/10 text-rose-400 border-rose-500/20'
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

async function openSubscriptionModal(tenant: Tenant) {
  selectedTenant.value = tenant
  selectedSubscription.value = tenant.subscription || null
  overrideStatusForm.value = tenant.subscription?.status || 'ACTIVE'
  showSubscriptionModal.value = true

  loadingSubscriptionDetails.value = true
  try {
    const res = await api.get(`/admin/tenants/${tenant.id}`)
    if (res.data.success && res.data.data.subscription) {
      selectedSubscription.value = res.data.data.subscription
      overrideStatusForm.value = res.data.data.subscription.status
      const idx = tenants.value.findIndex(t => t.id === tenant.id)
      if (idx !== -1) {
        tenants.value[idx].subscription = res.data.data.subscription
      }
    }
  } catch (err: any) {
    console.error('Erro ao carregar detalhes completos da assinatura:', err)
  } finally {
    loadingSubscriptionDetails.value = false
  }
}

async function applyStatusOverride() {
  if (!selectedSubscription.value) return

  const newStatus = overrideStatusForm.value
  if (!confirm(`Confirmar alteração manual do status da assinatura para "${newStatus}"?`)) {
    return
  }

  updatingStatus.value = true
  try {
    const res = await api.patch(`/admin/subscriptions/${selectedSubscription.value.id}/status`, {
      status: newStatus,
    })
    if (res.data.success) {
      selectedSubscription.value.status = newStatus
      if (selectedTenant.value?.subscription) {
        selectedTenant.value.subscription.status = newStatus
      }
      await fetchTenants()
    }
  } catch (err: any) {
    alert(err.response?.data?.error || 'Erro ao alterar status da assinatura.')
  } finally {
    updatingStatus.value = false
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
