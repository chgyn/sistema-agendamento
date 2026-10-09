<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-full bg-purple-500/10 border border-purple-500/20 text-purple-600 dark:text-purple-400 text-xs font-semibold mb-2">
          <ShieldCheck class="w-3.5 h-3.5" />
          <span>Administração Multi-tenant Global</span>
        </div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight font-['Outfit'] flex items-center gap-2.5">
          <Building2 class="w-8 h-8 text-purple-600 dark:text-purple-400" />
          <span>Gestão de Estabelecimentos</span>
        </h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-surface-400 mt-1">
          Cadastre e gerencie todos os estabelecimentos, planos contratados e o status financeiro de suas assinaturas.
        </p>
      </div>

      <button
        @click="openCreateModal"
        class="inline-flex items-center justify-center gap-2 px-5 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 via-indigo-600 to-purple-600 hover:from-purple-500 hover:to-indigo-500 text-white text-sm font-bold shadow-lg shadow-purple-900/20 hover:shadow-purple-700/30 transition active:scale-95 shrink-0"
      >
        <Plus class="w-4 h-4" />
        <span>Novo Estabelecimento</span>
      </button>
    </div>

    <!-- KPI Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="saas-card p-5 flex items-center gap-4 group hover:border-purple-500/40 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-600 dark:text-purple-400 group-hover:scale-105 transition-transform shrink-0">
          <Building2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-[#718096] dark:text-surface-400 font-medium">Total Estabelecimentos</p>
          <p class="text-2xl font-black text-[#202224] dark:text-white font-['Outfit'] mt-0.5">{{ tenants.length }}</p>
        </div>
      </div>

      <div class="saas-card p-5 flex items-center gap-4 group hover:border-emerald-500/40 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-600 dark:text-emerald-400 group-hover:scale-105 transition-transform shrink-0">
          <CheckCircle2 class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-[#718096] dark:text-surface-400 font-medium">Assinaturas Ativas</p>
          <p class="text-2xl font-black text-emerald-600 dark:text-emerald-400 font-['Outfit'] mt-0.5">{{ activeSubscriptionsCount }}</p>
        </div>
      </div>

      <div class="saas-card p-5 flex items-center gap-4 group hover:border-amber-500/40 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-600 dark:text-amber-400 group-hover:scale-105 transition-transform shrink-0">
          <Clock class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-[#718096] dark:text-surface-400 font-medium">Pagamento Pendente</p>
          <p class="text-2xl font-black text-amber-600 dark:text-amber-400 font-['Outfit'] mt-0.5">{{ pendingSubscriptionsCount }}</p>
        </div>
      </div>

      <div class="saas-card p-5 flex items-center gap-4 group hover:border-rose-500/40 transition-colors">
        <div class="w-12 h-12 rounded-2xl bg-rose-500/10 border border-rose-500/20 flex items-center justify-center text-rose-600 dark:text-rose-400 group-hover:scale-105 transition-transform shrink-0">
          <AlertTriangle class="w-6 h-6" />
        </div>
        <div>
          <p class="text-xs text-[#718096] dark:text-surface-400 font-medium">Inadimplentes / Atrasadas</p>
          <p class="text-2xl font-black text-rose-600 dark:text-rose-400 font-['Outfit'] mt-0.5">{{ overdueSubscriptionsCount }}</p>
        </div>
      </div>
    </div>

    <!-- Filtros & Busca -->
    <div class="saas-card p-3 sm:p-4 flex flex-col md:flex-row items-center gap-3">
      <div class="relative flex-1 w-full">
        <Search class="w-4 h-4 text-gray-400 dark:text-surface-400 absolute left-3.5 top-1/2 -translate-y-1/2" />
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Buscar por nome, slug, cidade ou e-mail..."
          class="w-full pl-10 pr-4 py-2.5 rounded-xl bg-gray-50/80 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-surface-100 placeholder-gray-400 dark:placeholder-surface-500 focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        />
      </div>

      <div class="flex flex-col sm:flex-row items-center gap-2 w-full md:w-auto">
        <!-- Filtro por Status da Assinatura -->
        <select
          v-model="subscriptionFilter"
          class="w-full sm:w-56 py-2.5 px-3.5 rounded-xl bg-gray-50/80 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-surface-200 focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
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
          class="w-full sm:w-44 py-2.5 px-3.5 rounded-xl bg-gray-50/80 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-surface-200 focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
        >
          <option value="ALL">Status: Todos</option>
          <option value="ACTIVE">Apenas Ativos</option>
          <option value="INACTIVE">Apenas Inativos</option>
        </select>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="text-center py-20">
      <Loader2 class="w-9 h-9 text-purple-500 animate-spin mx-auto mb-3" />
      <p class="text-sm text-[#718096] dark:text-surface-400">Carregando estabelecimentos...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredTenants.length === 0" class="saas-card p-12 text-center">
      <div class="w-14 h-14 rounded-2xl bg-gray-100 dark:bg-surface-800/80 border border-gray-200 dark:border-surface-700/60 flex items-center justify-center mx-auto mb-4 text-gray-400 dark:text-surface-500">
        <Building2 class="w-7 h-7" />
      </div>
      <h3 class="text-base font-bold text-[#202224] dark:text-white font-['Outfit']">Nenhum estabelecimento encontrado</h3>
      <p class="text-xs text-[#718096] dark:text-surface-400 mt-1.5 max-w-sm mx-auto">
        {{ searchTerm || subscriptionFilter !== 'ALL' || statusFilter !== 'ALL' ? 'Nenhum resultado corresponde aos filtros aplicados.' : 'Comece cadastrando o primeiro estabelecimento da plataforma.' }}
      </p>
    </div>

    <!-- Tabela de Estabelecimentos com Assinatura -->
    <div v-else class="saas-card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left text-sm">
          <thead class="bg-gray-50/90 dark:bg-surface-950/90 border-b border-gray-200/80 dark:border-surface-800 text-[11px] text-gray-500 dark:text-surface-400 uppercase tracking-wider font-bold">
            <tr>
              <th class="py-4 px-5">Estabelecimento</th>
              <th class="py-4 px-5">Plano & Assinatura</th>
              <th class="py-4 px-5">Contato</th>
              <th class="py-4 px-5">Localização</th>
              <th class="py-4 px-5">Status Conta</th>
              <th class="py-4 px-5 text-right">Ações</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-gray-100 dark:divide-surface-800/60">
            <tr v-for="t in filteredTenants" :key="t.id" class="hover:bg-gray-50/70 dark:hover:bg-surface-800/40 transition">
              <!-- Estabelecimento -->
              <td class="py-4 px-5">
                <div class="flex items-center gap-3.5">
                  <div
                    class="w-11 h-11 rounded-2xl flex items-center justify-center font-black text-sm text-white shrink-0 overflow-hidden border border-gray-200 dark:border-surface-700/60 bg-gray-100 dark:bg-surface-800 shadow-sm"
                    :style="t.logo_url ? '' : { backgroundColor: t.primary_color || '#8b5cf6' }"
                  >
                    <img v-if="t.logo_url" :src="t.logo_url" :alt="t.name" class="w-full h-full object-cover" />
                    <span v-else class="font-['Outfit'] text-base">{{ t.name.charAt(0) }}</span>
                  </div>
                  <div class="min-w-0">
                    <p class="font-bold text-[#202224] dark:text-white text-sm truncate font-['Outfit']">{{ t.name }}</p>
                    <p class="text-xs text-purple-600 dark:text-purple-400 font-mono truncate">/{{ t.slug }}</p>
                  </div>
                </div>
              </td>

              <!-- Plano & Assinatura -->
              <td class="py-4 px-5">
                <div v-if="t.subscription" class="space-y-1">
                  <div class="flex items-center gap-2 flex-wrap">
                    <span class="font-bold text-[#202224] dark:text-white text-xs">
                      {{ t.subscription.plan?.name || 'Plano Personalizado' }}
                    </span>
                    <span
                      class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border inline-flex items-center gap-1.5"
                      :class="getSubscriptionBadgeClass(t.subscription.status)"
                    >
                      <span class="w-1.5 h-1.5 rounded-full" :class="getSubscriptionDotClass(t.subscription.status)"></span>
                      <span>{{ formatSubStatus(t.subscription.status) }}</span>
                    </span>

                    <!-- Origin Badge -->
                    <span
                      v-if="t.subscription.origin === 'MANUAL'"
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-purple-500/10 text-purple-600 dark:text-purple-400 border border-purple-500/30 inline-flex items-center gap-1"
                      title="Liberado manualmente pelo administrador"
                    >
                      <Key class="w-2.5 h-2.5" />
                      <span>Manual</span>
                    </span>
                    <span
                      v-else-if="t.subscription.origin === 'FREE_PLAN' || t.subscription.plan?.is_free"
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 border border-cyan-500/30 inline-flex items-center gap-1"
                      title="Plano Gratuito"
                    >
                      <Gift class="w-2.5 h-2.5" />
                      <span>Grátis</span>
                    </span>
                    <span
                      v-else
                      class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/30 inline-flex items-center gap-1"
                      title="Gateway Asaas"
                    >
                      <span>Asaas</span>
                    </span>
                  </div>

                  <div class="text-[11px] text-[#718096] dark:text-surface-400 flex items-center gap-1.5 font-medium flex-wrap">
                    <span v-if="t.subscription.origin === 'FREE_PLAN' || t.subscription.plan?.is_free" class="text-cyan-600 dark:text-cyan-400 font-bold">100% Gratuito</span>
                    <span v-else class="font-mono text-gray-700 dark:text-surface-300">R$ {{ t.subscription.price.toFixed(2).replace('.', ',') }}</span>
                    <span>•</span>
                    <span>{{ formatCycle(t.subscription.billing_cycle) }}</span>
                    <template v-if="t.subscription.current_period_end">
                      <span>•</span>
                      <span class="text-amber-600 dark:text-amber-400 font-medium">Expira: {{ formatDate(t.subscription.current_period_end) }}</span>
                    </template>
                  </div>
                </div>
                <div v-else class="text-xs text-gray-400 dark:text-surface-500 italic">
                  Sem assinatura ativa
                </div>
              </td>

              <!-- Contato -->
              <td class="py-4 px-5 text-xs text-gray-600 dark:text-surface-300">
                <div class="font-medium text-[#202224] dark:text-white">{{ t.phone || '-' }}</div>
                <div class="text-gray-400 dark:text-surface-400 truncate mt-0.5">{{ t.email || '-' }}</div>
              </td>

              <!-- Localização -->
              <td class="py-4 px-5 text-xs text-gray-600 dark:text-surface-300">
                <div v-if="t.city || t.state" class="font-medium">{{ t.city }} - {{ t.state }}</div>
                <div v-else class="text-gray-400 dark:text-surface-500">Não informado</div>
              </td>

              <!-- Status Conta -->
              <td class="py-4 px-5">
                <span
                  class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-[11px] font-bold border"
                  :class="t.is_active ? 'bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border-emerald-500/20' : 'bg-rose-500/10 text-rose-600 dark:text-rose-400 border-rose-500/20'"
                >
                  <span class="w-1.5 h-1.5 rounded-full" :class="t.is_active ? 'bg-emerald-500 dark:bg-emerald-400' : 'bg-rose-500 dark:bg-rose-400'"></span>
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
                    class="p-2 rounded-xl bg-purple-500/10 hover:bg-purple-500/20 text-purple-600 dark:text-purple-300 transition"
                  >
                    <CreditCard class="w-4 h-4" />
                  </button>

                  <!-- Link Público -->
                  <a
                    :href="`/agendamento/${t.slug}`"
                    target="_blank"
                    title="Ver página de agendamento público"
                    class="p-2 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-600 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-300 dark:hover:text-white transition"
                  >
                    <ExternalLink class="w-4 h-4" />
                  </a>

                  <!-- Editar Estabelecimento -->
                  <button
                    @click="openEditModal(t)"
                    title="Editar informações"
                    class="p-2 rounded-xl bg-gray-100 hover:bg-gray-200 text-gray-600 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-300 dark:hover:text-white transition"
                  >
                    <Pencil class="w-4 h-4" />
                  </button>

                  <!-- Alternar Status -->
                  <button
                    @click="toggleStatus(t)"
                    :title="t.is_active ? 'Inativar estabelecimento' : 'Ativar estabelecimento'"
                    class="p-2 rounded-xl transition"
                    :class="t.is_active ? 'bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400' : 'bg-emerald-500/10 hover:bg-emerald-500/20 text-emerald-600 dark:text-emerald-400'"
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
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 rounded-2xl w-full max-w-3xl max-h-[92vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-gray-100 dark:border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-10 h-10 rounded-2xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-600 dark:text-purple-400 shrink-0">
              <CreditCard class="w-5 h-5" />
            </div>
            <div>
              <div class="flex items-center gap-2">
                <h3 class="font-bold text-[#202224] dark:text-white text-base font-['Outfit']">{{ selectedTenant.name }}</h3>
                <span class="text-xs text-purple-600 dark:text-purple-400 font-mono">/{{ selectedTenant.slug }}</span>
              </div>
              <p class="text-xs text-[#718096] dark:text-surface-400">Detalhamento contratual, liberação manual e auditoria</p>
            </div>
          </div>
          <button @click="showSubscriptionModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Subnavigation Tabs -->
        <div class="flex items-center gap-1 sm:gap-2 px-5 sm:px-6 pt-3 border-b border-gray-100 dark:border-surface-800/80 bg-gray-50/60 dark:bg-surface-950/40 overflow-x-auto">
          <button
            @click="activeSubTab = 'details'"
            class="pb-3 px-3 text-xs font-bold transition border-b-2 flex items-center gap-2 shrink-0"
            :class="activeSubTab === 'details' ? 'border-purple-600 text-purple-600 dark:border-purple-500 dark:text-purple-400' : 'border-transparent text-gray-500 dark:text-surface-400 hover:text-gray-800 dark:hover:text-surface-200'"
          >
            <CreditCard class="w-4 h-4" />
            <span>Visão Geral & Faturas</span>
          </button>

          <button
            @click="activeSubTab = 'manual_grant'"
            class="pb-3 px-3 text-xs font-bold transition border-b-2 flex items-center gap-2 shrink-0"
            :class="activeSubTab === 'manual_grant' ? 'border-purple-600 text-purple-600 dark:border-purple-500 dark:text-purple-400' : 'border-transparent text-gray-500 dark:text-surface-400 hover:text-gray-800 dark:hover:text-surface-200'"
          >
            <Key class="w-4 h-4" />
            <span>Liberação Manual</span>
          </button>

          <button
            @click="activeSubTab = 'audit_logs'"
            class="pb-3 px-3 text-xs font-bold transition border-b-2 flex items-center gap-2 shrink-0"
            :class="activeSubTab === 'audit_logs' ? 'border-purple-500 text-purple-400' : 'border-transparent text-surface-400 hover:text-surface-200'"
          >
            <History class="w-4 h-4" />
            <span>Trilha de Auditoria</span>
            <span v-if="auditLogs.length > 0" class="px-1.5 py-0.2 rounded-full text-[10px] bg-purple-500/20 text-purple-300 font-mono">
              {{ auditLogs.length }}
            </span>
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-6 flex-1">
          <div v-if="loadingSubscriptionDetails" class="py-12 text-center">
            <Loader2 class="w-8 h-8 text-purple-400 animate-spin mx-auto mb-3" />
            <p class="text-xs text-surface-400">Carregando dados da assinatura...</p>
          </div>

          <!-- TAB 1: VISÃO GERAL & DETALHES -->
          <div v-else-if="activeSubTab === 'details'" class="space-y-6">
            <div v-if="!selectedSubscription" class="bg-gray-50 dark:bg-surface-950/60 border border-gray-200 dark:border-surface-800 rounded-2xl p-8 text-center space-y-4">
              <div class="w-12 h-12 rounded-2xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-500 dark:text-amber-400 mx-auto">
                <AlertTriangle class="w-6 h-6" />
              </div>
              <div>
                <p class="text-sm font-bold text-[#202224] dark:text-white">Nenhuma assinatura vinculada a este estabelecimento</p>
                <p class="text-xs text-gray-500 dark:text-surface-400 mt-1 max-w-md mx-auto">
                  Este estabelecimento não possui assinatura ativa no momento. Você pode liberar uma assinatura manualmente agora mesmo sem cobrança no Asaas.
                </p>
              </div>
              <button
                @click="activeSubTab = 'manual_grant'"
                class="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-purple-600 hover:bg-purple-500 text-white font-bold text-xs shadow-lg shadow-purple-900/30 transition"
              >
                <Key class="w-4 h-4" />
                <span>Liberar Assinatura Manualmente</span>
              </button>
            </div>

            <div v-else class="space-y-6">
              <!-- Banner de Status da Assinatura -->
              <div
                class="p-4 rounded-2xl border flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 shadow-md"
                :class="getStatusBannerClass(selectedSubscription.status)"
              >
                <div class="flex items-center gap-3.5">
                  <div class="w-10 h-10 rounded-2xl flex items-center justify-center shrink-0" :class="getStatusIconBoxClass(selectedSubscription.status)">
                    <CheckCircle2 v-if="selectedSubscription.status === 'ACTIVE'" class="w-5 h-5 text-emerald-500 dark:text-emerald-400" />
                    <Clock v-else-if="selectedSubscription.status === 'PENDING'" class="w-5 h-5 text-amber-500 dark:text-amber-400 animate-pulse" />
                    <AlertCircle v-else class="w-5 h-5 text-rose-500 dark:text-rose-400" />
                  </div>
                  <div>
                    <div class="flex items-center gap-2 flex-wrap">
                      <span class="text-xs font-bold text-[#202224] dark:text-white">Status da Assinatura:</span>
                      <span class="px-2.5 py-0.5 rounded-full text-[10px] font-black uppercase tracking-wider border" :class="getSubscriptionBadgeClass(selectedSubscription.status)">
                        {{ formatSubStatus(selectedSubscription.status) }}
                      </span>
                      <!-- Modalidade / Origem -->
                      <span
                        v-if="selectedSubscription.origin === 'MANUAL'"
                        class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-purple-500/10 text-purple-600 dark:text-purple-400 border border-purple-500/30 inline-flex items-center gap-1"
                      >
                        <Key class="w-2.5 h-2.5" />
                        <span>Manual</span>
                      </span>
                      <span
                        v-else-if="selectedSubscription.origin === 'FREE_PLAN' || selectedSubscription.plan?.is_free"
                        class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-cyan-500/10 text-cyan-600 dark:text-cyan-400 border border-cyan-500/30 inline-flex items-center gap-1"
                      >
                        <Gift class="w-2.5 h-2.5" />
                        <span>Plano Gratuito</span>
                      </span>
                      <span
                        v-else
                        class="px-2 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider bg-blue-500/10 text-blue-600 dark:text-blue-400 border border-blue-500/30 inline-flex items-center gap-1"
                      >
                        <span>Asaas</span>
                      </span>
                    </div>
                    <p class="text-xs text-gray-600 dark:text-surface-300 mt-0.5">
                      {{ getStatusExplanation(selectedSubscription.status) }}
                    </p>
                  </div>
                </div>

                <!-- Ação de Abrir Checkout/Fatura no Asaas (apenas se origem Asaas) -->
                <a
                  v-if="selectedSubscription.origin === 'ASAAS' && selectedSubscription.payment_url"
                  :href="selectedSubscription.payment_url"
                  target="_blank"
                  class="inline-flex items-center gap-1.5 px-4 py-2 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-white font-bold text-xs shadow-lg shadow-emerald-500/20 transition shrink-0"
                >
                  <ExternalLink class="w-3.5 h-3.5" />
                  <span>Abrir Fatura Asaas</span>
                </a>
              </div>

              <!-- Grid de Informações Contratuais -->
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- Card 1: Detalhes do Plano -->
                <div class="bg-gray-50/70 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800/80 rounded-2xl p-4 space-y-3">
                  <div class="flex items-center justify-between pb-2.5 border-b border-gray-200/80 dark:border-surface-800">
                    <div class="flex items-center gap-2">
                      <Sparkles class="w-4 h-4 text-purple-600 dark:text-purple-400" />
                      <h4 class="text-xs font-bold uppercase tracking-wider text-purple-700 dark:text-purple-300">Plano Contratado</h4>
                    </div>
                    <span class="px-2.5 py-0.5 rounded-full bg-purple-500/10 text-purple-700 dark:text-purple-300 border border-purple-500/30 text-[10px] font-bold">
                      {{ selectedSubscription.plan?.name || 'Personalizado' }}
                    </span>
                  </div>

                  <div class="grid grid-cols-2 gap-3 text-xs">
                    <div>
                      <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Valor Recorrente:</span>
                      <span v-if="selectedSubscription.origin === 'FREE_PLAN' || selectedSubscription.plan?.is_free" class="text-sm font-black text-cyan-600 dark:text-cyan-400 font-mono">
                        R$ 0,00 (Grátis)
                      </span>
                      <span v-else class="text-sm font-black text-[#202224] dark:text-white font-mono">
                        R$ {{ selectedSubscription.price.toFixed(2).replace('.', ',') }}
                      </span>
                    </div>
                    <div>
                      <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Periodicidade:</span>
                      <span class="text-sm font-bold text-[#202224] dark:text-white">
                        {{ formatCycle(selectedSubscription.billing_cycle) }}
                      </span>
                    </div>
                    <div>
                      <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Data de Início:</span>
                      <span class="text-xs font-medium text-gray-700 dark:text-surface-200">
                        {{ formatDate(selectedSubscription.created_at) }}
                      </span>
                    </div>
                    <div>
                      <span class="text-gray-500 dark:text-surface-400 block text-[11px]">
                        {{ selectedSubscription.origin === 'MANUAL' ? 'Validade / Expiração:' : 'Próxima Cobrança:' }}
                      </span>
                      <span v-if="selectedSubscription.current_period_end" class="text-xs font-bold text-amber-600 dark:text-amber-400">
                        {{ formatDate(selectedSubscription.current_period_end) }}
                      </span>
                      <span v-else-if="selectedSubscription.next_due_date" class="text-xs font-bold text-emerald-600 dark:text-emerald-400">
                        {{ formatDate(selectedSubscription.next_due_date) }}
                      </span>
                      <span v-else class="text-xs text-gray-500 dark:text-surface-400">
                        Acesso contínuo / Vitalício
                      </span>
                    </div>
                  </div>
                </div>

                <!-- Card 2: Modalidade & Origem da Assinatura -->
                <div class="bg-gray-50/70 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800/80 rounded-2xl p-4 space-y-3 flex flex-col justify-between">
                  <div>
                    <div class="flex items-center gap-2 pb-2.5 border-b border-gray-200/80 dark:border-surface-800">
                      <ShieldCheck class="w-4 h-4 text-indigo-600 dark:text-indigo-400" />
                      <h4 class="text-xs font-bold uppercase tracking-wider text-indigo-700 dark:text-indigo-300">Modalidade & Origem</h4>
                    </div>

                    <!-- CASO 1: CONCESSÃO MANUAL -->
                    <div v-if="selectedSubscription.origin === 'MANUAL'" class="mt-3 space-y-2 text-xs">
                      <div>
                        <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Concedido por:</span>
                        <span class="font-medium text-[#202224] dark:text-white flex items-center gap-1.5 mt-0.5">
                          <User class="w-3.5 h-3.5 text-purple-600 dark:text-purple-400" />
                          <span>{{ selectedSubscription.granted_by_user?.name || 'Administrador Global' }}</span>
                        </span>
                      </div>
                      <div>
                        <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Motivo / Justificativa:</span>
                        <p class="text-xs text-purple-900 dark:text-purple-200 italic mt-0.5 bg-purple-50 dark:bg-surface-900/60 p-2 rounded-lg border border-purple-200 dark:border-purple-500/10">
                          "{{ selectedSubscription.manual_grant_reason || 'Concessão direta sem cobrança pelo administrador.' }}"
                        </p>
                      </div>
                    </div>

                    <!-- CASO 2: PLANO GRATUITO -->
                    <div v-else-if="selectedSubscription.origin === 'FREE_PLAN' || selectedSubscription.plan?.is_free" class="mt-3 space-y-2 text-xs">
                      <div class="p-2.5 rounded-xl bg-cyan-500/10 border border-cyan-500/20 text-cyan-700 dark:text-cyan-300 space-y-1">
                        <div class="font-bold flex items-center gap-1.5">
                          <Gift class="w-3.5 h-3.5" />
                          <span>Plano Gratuito Nativo</span>
                        </div>
                        <p class="text-[11px] text-gray-600 dark:text-surface-300">
                          Estabelecimento cadastrado diretamente no plano grátis. Não requer integração nem faturas do gateway Asaas.
                        </p>
                      </div>
                    </div>

                    <!-- CASO 3: GATEWAY ASAAS -->
                    <div v-else class="mt-3 space-y-2 text-xs">
                      <div>
                        <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Asaas Subscription ID:</span>
                        <span class="font-mono text-xs text-indigo-700 dark:text-indigo-200 select-all">
                          {{ selectedSubscription.asaas_subscription_id || 'Não integrado' }}
                        </span>
                      </div>
                      <div>
                        <span class="text-gray-500 dark:text-surface-400 block text-[11px]">Asaas Customer ID:</span>
                        <span class="font-mono text-xs text-indigo-700 dark:text-indigo-200 select-all">
                          {{ selectedSubscription.asaas_customer_id || 'Não integrado' }}
                        </span>
                      </div>
                    </div>
                  </div>

                  <div class="text-[11px] text-gray-500 dark:text-surface-400 flex items-center justify-between pt-2 border-t border-gray-200/80 dark:border-surface-800/60">
                    <span>Tipo de Acesso:</span>
                    <span class="font-bold text-[#202224] dark:text-white uppercase">{{ formatOrigin(selectedSubscription.origin) }}</span>
                  </div>
                </div>
              </div>

              <!-- Histórico de Faturas e Cobranças Sincronizadas (Apenas Asaas) -->
              <div v-if="selectedSubscription.origin === 'ASAAS'" class="bg-white dark:bg-surface-950/90 border border-gray-200 dark:border-surface-800/80 rounded-2xl overflow-hidden shadow-sm">
                <div class="p-4 border-b border-gray-200 dark:border-surface-800 flex items-center justify-between bg-gray-50/50 dark:bg-transparent">
                  <div class="flex items-center gap-2">
                    <FileText class="w-4 h-4 text-emerald-500 dark:text-emerald-400" />
                    <h4 class="text-xs font-bold uppercase tracking-wider text-gray-700 dark:text-surface-200">Histórico de Cobranças & Faturas Asaas</h4>
                  </div>
                  <span class="text-[11px] text-gray-400 dark:text-surface-500">Sincronizado via Webhook</span>
                </div>

                <div v-if="!selectedSubscription.invoices || selectedSubscription.invoices.length === 0" class="p-6 text-center text-xs text-gray-500 dark:text-surface-400">
                  Nenhuma fatura registrada no histórico até o momento.
                </div>

                <div v-else class="overflow-x-auto max-h-48">
                  <table class="w-full text-left text-xs">
                    <thead class="bg-gray-50 dark:bg-surface-900 border-b border-gray-200 dark:border-surface-800 text-[11px] text-gray-500 dark:text-surface-400 uppercase font-semibold sticky top-0">
                      <tr>
                        <th class="py-2.5 px-3.5">Fatura</th>
                        <th class="py-2.5 px-3.5">Vencimento</th>
                        <th class="py-2.5 px-3.5">Valor</th>
                        <th class="py-2.5 px-3.5">Método</th>
                        <th class="py-2.5 px-3.5">Status</th>
                        <th class="py-2.5 px-3.5 text-right">Comprovante</th>
                      </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-100 dark:divide-surface-800/60">
                      <tr v-for="inv in selectedSubscription.invoices" :key="inv.id" class="hover:bg-gray-50 dark:hover:bg-surface-800/40 transition">
                        <td class="py-2.5 px-3.5 font-mono text-gray-600 dark:text-surface-300">
                          {{ inv.asaas_payment_id }}
                        </td>
                        <td class="py-2.5 px-3.5 text-gray-600 dark:text-surface-300">
                          {{ formatDate(inv.due_date) }}
                        </td>
                        <td class="py-2.5 px-3.5 font-bold text-[#202224] dark:text-white font-mono">
                          R$ {{ inv.value.toFixed(2).replace('.', ',') }}
                        </td>
                        <td class="py-2.5 px-3.5 text-gray-500 dark:text-surface-400 uppercase">
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
                            class="inline-flex items-center gap-1 px-2.5 py-1 rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 dark:bg-surface-800 dark:hover:bg-surface-700 dark:text-surface-300 dark:hover:text-white transition font-medium"
                          >
                            <ExternalLink class="w-3 h-3" />
                            <span>Abrir</span>
                          </a>
                          <span v-else class="text-gray-400 dark:text-surface-500">-</span>
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
                    <SlidersHorizontal class="w-4 h-4 text-purple-600 dark:text-purple-400" />
                    <h4 class="text-xs font-bold uppercase tracking-wider text-purple-700 dark:text-purple-300">Ajuste Manual de Status</h4>
                  </div>
                  <span class="text-[10px] text-gray-400 dark:text-surface-400">Grava evento na auditoria</span>
                </div>

                <div class="space-y-3">
                  <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                    <div>
                      <label class="block text-[11px] text-gray-600 dark:text-surface-400 font-medium mb-1">Novo Status:</label>
                      <select
                        v-model="overrideStatusForm"
                        class="w-full py-2 px-3 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-xs text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                      >
                        <option value="ACTIVE">ACTIVE (Liberado / Ativo 🟢)</option>
                        <option value="PENDING">PENDING (Aguardando Pagamento 🟡)</option>
                        <option value="OVERDUE">OVERDUE (Inadimplente / Bloqueado 🔴)</option>
                        <option value="CANCELLED">CANCELLED (Cancelado ⚫)</option>
                        <option value="TRIAL">TRIAL (Período de Testes 🚀)</option>
                      </select>
                    </div>

                    <div>
                      <label class="block text-[11px] text-gray-600 dark:text-surface-400 font-medium mb-1">Motivo / Justificativa:</label>
                      <input
                        v-model="overrideReasonForm"
                        type="text"
                        placeholder="Ex: Regularização off-line de fatura"
                        class="w-full py-2 px-3 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-xs text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                      />
                    </div>
                  </div>

                  <div class="flex justify-end">
                    <button
                      @click="applyStatusOverride"
                      :disabled="updatingStatus || overrideStatusForm === selectedSubscription.status"
                      class="inline-flex items-center justify-center gap-1.5 px-4 py-2 rounded-xl bg-purple-600 hover:bg-purple-500 text-white text-xs font-bold transition disabled:opacity-50 shadow-md"
                    >
                      <Loader2 v-if="updatingStatus" class="w-3.5 h-3.5 animate-spin" />
                      <span>Salvar Alteração de Status</span>
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- TAB 2: LIBERAÇÃO MANUAL DE ASSINATURA -->
          <div v-else-if="activeSubTab === 'manual_grant'" class="space-y-6">
            <div class="p-4 rounded-2xl bg-gradient-to-r from-purple-900/30 via-indigo-900/20 to-purple-900/30 border border-purple-500/30 flex items-start gap-3.5">
              <div class="w-10 h-10 rounded-2xl bg-purple-500/20 border border-purple-500/30 flex items-center justify-center text-purple-300 shrink-0 mt-0.5">
                <Key class="w-5 h-5" />
              </div>
              <div class="space-y-1">
                <h4 class="text-sm font-bold text-[#202224] dark:text-white font-['Outfit']">Liberação Manual de Assinatura</h4>
                <p class="text-xs text-gray-600 dark:text-surface-300 leading-relaxed">
                  Conceda ou altere a assinatura deste estabelecimento diretamente, sem a necessidade de gateway de pagamento Asaas. O status será imediatamente ativado e o evento registrado na auditoria com o autor da concessão.
                </p>
              </div>
            </div>

            <!-- Feedback de Erro ou Sucesso -->
            <div v-if="manualGrantError" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-300 text-xs flex items-center gap-2">
              <AlertCircle class="w-4 h-4 shrink-0" />
              <span>{{ manualGrantError }}</span>
            </div>

            <div v-if="manualGrantSuccess" class="p-3.5 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-300 text-xs flex items-center gap-2">
              <Check class="w-4 h-4 shrink-0" />
              <span>{{ manualGrantSuccess }}</span>
            </div>

            <form @submit.prevent="submitManualGrant" class="space-y-5">
              <!-- Campo 1: Selecionar Plano -->
              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-2">
                  1. Selecione o Plano a Associar *
                </label>
                <select
                  v-model="manualGrantForm.plan_id"
                  required
                  class="w-full py-2.5 px-3.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                >
                  <option value="" disabled>Escolha um plano...</option>
                  <option v-for="plan in availablePlans" :key="plan.id" :value="plan.id">
                    {{ plan.name }} — {{ plan.is_free ? '100% Gratuito' : 'R$ ' + plan.price.toFixed(2).replace('.', ',') }}
                    (Até {{ plan.max_professionals }} prof., {{ plan.max_services }} serv.)
                  </option>
                </select>
              </div>

              <!-- Campo 2: Período de Validade -->
              <div class="space-y-3">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider">
                  2. Período de Validade da Assinatura
                </label>

                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
                  <label
                    class="p-3.5 rounded-xl border cursor-pointer transition flex items-center gap-3"
                    :class="!manualGrantForm.has_expiration ? 'bg-purple-500/10 border-purple-500/40 text-purple-700 dark:text-white' : 'bg-gray-50 dark:bg-surface-950/60 border-gray-200 dark:border-surface-800 text-gray-600 dark:text-surface-400 hover:border-purple-300 dark:hover:border-surface-700'"
                  >
                    <input
                      type="radio"
                      :value="false"
                      v-model="manualGrantForm.has_expiration"
                      class="text-purple-600 focus:ring-purple-500"
                    />
                    <div>
                      <p class="text-xs font-bold">Sem Expiração</p>
                      <p class="text-[11px] text-gray-500 dark:text-surface-400">Acesso contínuo / vitalício</p>
                    </div>
                  </label>

                  <label
                    class="p-3.5 rounded-xl border cursor-pointer transition flex items-center gap-3"
                    :class="manualGrantForm.has_expiration ? 'bg-purple-500/10 border-purple-500/40 text-purple-700 dark:text-white' : 'bg-gray-50 dark:bg-surface-950/60 border-gray-200 dark:border-surface-800 text-gray-600 dark:text-surface-400 hover:border-purple-300 dark:hover:border-surface-700'"
                  >
                    <input
                      type="radio"
                      :value="true"
                      v-model="manualGrantForm.has_expiration"
                      class="text-purple-600 focus:ring-purple-500"
                    />
                    <div>
                      <p class="text-xs font-bold">Definir Validade</p>
                      <p class="text-[11px] text-gray-500 dark:text-surface-400">Expira na data limite escolhida</p>
                    </div>
                  </label>
                </div>

                <!-- Input de Data quando Definir Validade está selecionado -->
                <div v-if="manualGrantForm.has_expiration" class="p-3 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-800 space-y-1.5">
                  <label class="block text-[11px] font-bold text-purple-700 dark:text-purple-300 uppercase tracking-wider">
                    Data Limite de Validade *
                  </label>
                  <div class="relative">
                    <Calendar class="w-4 h-4 text-gray-400 dark:text-surface-400 absolute left-3 top-1/2 -translate-y-1/2" />
                    <input
                      type="date"
                      v-model="manualGrantForm.expires_at"
                      :min="minExpirationDate"
                      required
                      class="w-full pl-9 pr-3 py-2 rounded-lg bg-white dark:bg-surface-900 border border-gray-200 dark:border-surface-700/60 text-xs text-[#202224] dark:text-white focus:outline-none focus:border-purple-500"
                    />
                  </div>
                  <p class="text-[10px] text-gray-500 dark:text-surface-400">
                    Após as 23:59 desta data, a assinatura será automaticamente marcada como expirada.
                  </p>
                </div>
              </div>

              <!-- Campo 3: Justificativa / Motivo -->
              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-2">
                  3. Justificativa / Motivo da Concessão * (Obrigatório para Auditoria)
                </label>
                <textarea
                  v-model="manualGrantForm.reason"
                  rows="3"
                  required
                  placeholder="Ex: Parceria comercial estratégica, concessão de cortesia para novos franqueados, teste VIP ou acordo comercial off-line..."
                  class="w-full p-3 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-xs text-[#202224] dark:text-white placeholder-gray-400 dark:placeholder-surface-500 focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition resize-none"
                ></textarea>
                <p class="text-[10px] text-gray-500 dark:text-surface-400 mt-1">Mínimo de 3 caracteres. Ficará registrado na trilha de auditoria do sistema.</p>
              </div>

              <!-- Botões de Ação do Formulário -->
              <div class="flex items-center justify-end gap-3 pt-3 border-t border-gray-100 dark:border-surface-800">
                <button
                  type="button"
                  @click="activeSubTab = 'details'"
                  class="px-4 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-xs font-medium transition"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  :disabled="grantingManual"
                  class="inline-flex items-center gap-2 px-6 py-2.5 rounded-xl bg-gradient-to-r from-purple-600 via-indigo-600 to-purple-600 hover:from-purple-500 hover:to-indigo-500 text-white text-xs font-bold shadow-lg shadow-purple-900/30 transition disabled:opacity-50"
                >
                  <Loader2 v-if="grantingManual" class="w-4 h-4 animate-spin" />
                  <Key v-else class="w-4 h-4" />
                  <span>Confirmar e Liberar Assinatura</span>
                </button>
              </div>
            </form>
          </div>

          <!-- TAB 3: TRILHA DE AUDITORIA -->
          <div v-else-if="activeSubTab === 'audit_logs'" class="space-y-4">
            <div class="flex items-center justify-between pb-2 border-b border-gray-100 dark:border-surface-800">
              <div class="flex items-center gap-2">
                <History class="w-4 h-4 text-purple-600 dark:text-purple-400" />
                <h4 class="text-xs font-bold uppercase tracking-wider text-gray-700 dark:text-surface-200">Trilha de Auditoria & Alterações</h4>
              </div>
              <button
                v-if="selectedSubscription?.id"
                @click="fetchAuditLogs(selectedSubscription.id)"
                :disabled="loadingAuditLogs"
                title="Atualizar histórico"
                class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-lg bg-gray-100 hover:bg-gray-200 text-gray-700 hover:text-gray-900 dark:bg-surface-900 dark:hover:bg-surface-800 dark:text-surface-300 dark:hover:text-white text-xs transition"
              >
                <RefreshCw class="w-3 h-3" :class="{ 'animate-spin': loadingAuditLogs }" />
                <span>Atualizar</span>
              </button>
            </div>

            <div v-if="loadingAuditLogs" class="py-12 text-center">
              <Loader2 class="w-7 h-7 text-purple-500 dark:text-purple-400 animate-spin mx-auto mb-2" />
              <p class="text-xs text-gray-500 dark:text-surface-400">Carregando eventos de auditoria...</p>
            </div>

            <div v-else-if="auditLogs.length === 0" class="p-8 text-center bg-gray-50 dark:bg-surface-950/60 rounded-2xl border border-gray-200 dark:border-surface-800 text-gray-500 dark:text-surface-400 space-y-2">
              <History class="w-7 h-7 text-gray-400 dark:text-surface-500 mx-auto" />
              <p class="text-xs font-medium">Nenhum evento registrado nesta assinatura até o momento.</p>
            </div>

            <div v-else class="space-y-3">
              <div
                v-for="log in auditLogs"
                :key="log.id"
                class="p-4 rounded-xl bg-gray-50/70 dark:bg-surface-950/80 border border-gray-200/80 dark:border-surface-800/80 space-y-2 hover:border-purple-300 dark:hover:border-surface-700 transition"
              >
                <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-1.5">
                  <div class="flex items-center gap-2">
                    <span
                      class="px-2.5 py-0.5 rounded-full text-[10px] font-bold uppercase tracking-wider border"
                      :class="getAuditActionInfo(log.action).class"
                    >
                      {{ getAuditActionInfo(log.action).label }}
                    </span>
                    <span v-if="log.plan" class="text-xs text-purple-700 dark:text-purple-300 font-semibold">
                      • Plano: {{ log.plan.name }}
                    </span>
                  </div>
                  <span class="text-[11px] text-gray-500 dark:text-surface-400 font-mono">
                    {{ formatDateTime(log.created_at) }}
                  </span>
                </div>

                <p v-if="log.reason" class="text-xs text-gray-700 dark:text-surface-200 bg-gray-100 dark:bg-surface-900/60 p-2.5 rounded-lg border border-gray-200 dark:border-surface-800">
                  <span class="text-gray-500 dark:text-surface-400 font-bold">Justificativa: </span>
                  {{ log.reason }}
                </p>

                <div class="flex items-center justify-between text-[11px] text-gray-500 dark:text-surface-400 pt-1 border-t border-gray-200/80 dark:border-surface-800/60">
                  <div class="flex items-center gap-1.5">
                    <User class="w-3 h-3 text-gray-400 dark:text-surface-400" />
                    <span>Realizado por: </span>
                    <span class="text-gray-800 dark:text-surface-200 font-medium">{{ log.performed_by?.name || 'Sistema Automático' }}</span>
                  </div>
                  <div v-if="log.previous_status || log.new_status" class="font-mono text-[10px]">
                    <span :class="log.previous_status ? 'text-gray-400 dark:text-surface-400' : 'text-gray-300 dark:text-surface-500'">{{ log.previous_status || '-' }}</span>
                    <span class="mx-1 text-purple-500 dark:text-purple-400">➔</span>
                    <span class="text-emerald-600 dark:text-emerald-400 font-bold">{{ log.new_status }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-gray-100 dark:border-surface-800/80 flex items-center justify-between bg-gray-50/70 dark:bg-surface-950/60 rounded-b-2xl">
          <div class="text-xs text-gray-500 dark:text-surface-400">
            <span v-if="selectedSubscription">Origem: <strong class="text-gray-900 dark:text-white">{{ formatOrigin(selectedSubscription.origin) }}</strong></span>
          </div>
          <button
            @click="showSubscriptionModal = false"
            class="px-5 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-sm font-medium transition"
          >
            Fechar
          </button>
        </div>
      </div>
    </div>

    <!-- MODAL: Criar Novo Estabelecimento -->
    <div
      v-if="showCreateModal"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 rounded-2xl w-full max-w-2xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-gray-100 dark:border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-600 dark:text-purple-400 shrink-0">
              <Plus class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-[#202224] dark:text-white text-base font-['Outfit']">Novo Estabelecimento & Administrador</h3>
              <p class="text-xs text-[#718096] dark:text-surface-400">Cadastre a empresa e credenciais de acesso inicial.</p>
            </div>
          </div>
          <button @click="showCreateModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-6">
          <div v-if="createError" class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/20 text-rose-600 dark:text-rose-300 text-xs">
            {{ createError }}
          </div>

          <!-- 1. Dados do Estabelecimento -->
          <div class="space-y-3">
            <h4 class="text-xs font-bold uppercase tracking-wider text-purple-600 dark:text-purple-400 flex items-center gap-1.5">
              <Building2 class="w-3.5 h-3.5" />
              <span>1. Dados do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
                <input
                  v-model="createForm.name"
                  @input="autoGenerateSlug"
                  type="text"
                  required
                  placeholder="Ex: Barbearia Prime"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Slug de Acesso *</label>
                <input
                  v-model="createForm.slug"
                  type="text"
                  required
                  placeholder="barbearia-prime"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white font-mono focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone *</label>
                <input
                  v-model="createForm.phone"
                  type="text"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">CNPJ / CPF</label>
                <input
                  v-model="createForm.document"
                  type="text"
                  placeholder="00.000.000/0001-00"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Cidade</label>
                <input
                  v-model="createForm.city"
                  type="text"
                  placeholder="São Paulo"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">UF</label>
                <input
                  v-model="createForm.state"
                  type="text"
                  placeholder="SP"
                  maxlength="2"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white uppercase focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>
            </div>
          </div>

          <!-- 2. Administrador do Estabelecimento -->
          <div class="space-y-3 pt-4 border-t border-gray-100 dark:border-surface-800">
            <h4 class="text-xs font-bold uppercase tracking-wider text-purple-600 dark:text-purple-400 flex items-center gap-1.5">
              <ShieldCheck class="w-3.5 h-3.5" />
              <span>2. Administrador do Estabelecimento</span>
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
              <div class="sm:col-span-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Nome do Administrador *</label>
                <input
                  v-model="createForm.admin_name"
                  type="text"
                  required
                  placeholder="Carlos Silva"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Acesso *</label>
                <input
                  v-model="createForm.admin_email"
                  type="email"
                  required
                  placeholder="admin@barbeariaprime.com"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Senha de Acesso *</label>
                <input
                  v-model="createForm.admin_password"
                  type="password"
                  required
                  minlength="6"
                  placeholder="••••••••"
                  class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-gray-100 dark:border-surface-800/80 flex items-center justify-end gap-3 bg-gray-50/70 dark:bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showCreateModal = false"
            class="px-4 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-sm font-medium transition"
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
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 dark:bg-surface-950/80 backdrop-blur-sm overflow-y-auto"
    >
      <div class="bg-white dark:bg-[#121318] border border-gray-200 dark:border-surface-800 rounded-2xl w-full max-w-xl max-h-[90vh] flex flex-col shadow-2xl my-8">
        <!-- Modal Header -->
        <div class="flex items-center justify-between p-5 sm:p-6 border-b border-gray-100 dark:border-surface-800/80">
          <div class="flex items-center gap-3">
            <div class="w-9 h-9 rounded-xl bg-purple-500/10 border border-purple-500/20 flex items-center justify-center text-purple-600 dark:text-purple-400 shrink-0">
              <Pencil class="w-4 h-4" />
            </div>
            <div>
              <h3 class="font-bold text-[#202224] dark:text-white text-base font-['Outfit']">Editar Estabelecimento</h3>
              <p class="text-xs text-[#718096] dark:text-surface-400">Atualize dados cadastrais e de contato.</p>
            </div>
          </div>
          <button @click="showEditModal = false" class="p-1.5 rounded-lg text-gray-400 hover:text-gray-700 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 transition">
            <X class="w-5 h-5" />
          </button>
        </div>

        <!-- Modal Body -->
        <div class="p-6 overflow-y-auto space-y-4">
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <div class="sm:col-span-2">
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Nome do Estabelecimento *</label>
              <input
                v-model="editingTenant.name"
                type="text"
                required
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">WhatsApp / Telefone</label>
              <input
                v-model="editingTenant.phone"
                type="text"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">E-mail de Contato</label>
              <input
                v-model="editingTenant.email"
                type="email"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Cidade</label>
              <input
                v-model="editingTenant.city"
                type="text"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider mb-1.5">Estado (UF)</label>
              <input
                v-model="editingTenant.state"
                type="text"
                maxlength="2"
                class="w-full px-3.5 py-2.5 rounded-xl bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 text-sm text-[#202224] dark:text-white uppercase focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:border-purple-500 focus:ring-2 focus:ring-purple-500/20 transition"
              />
            </div>
          </div>
        </div>

        <!-- Modal Footer -->
        <div class="p-5 sm:p-6 border-t border-gray-100 dark:border-surface-800/80 flex items-center justify-end gap-3 bg-gray-50/70 dark:bg-surface-950/60 rounded-b-2xl">
          <button
            @click="showEditModal = false"
            class="px-4 py-2.5 rounded-xl text-gray-600 hover:text-gray-900 hover:bg-gray-100 dark:text-surface-400 dark:hover:text-white dark:hover:bg-surface-800 text-sm font-medium transition"
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
  Clock, AlertTriangle, AlertCircle, Sparkles, FileText, SlidersHorizontal,
  Key, Gift, History, Calendar, Check, RefreshCw, User
} from 'lucide-vue-next'
import api from '../../services/api'
import type {
  Tenant, Subscription, SubscriptionStatus, PlanBillingCycle,
  Plan, SubscriptionAuditLog
} from '../../stores/auth'

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
const overrideReasonForm = ref('')

// Tabs e Liberação Manual na Modal de Assinatura
const activeSubTab = ref<'details' | 'manual_grant' | 'audit_logs'>('details')
const availablePlans = ref<Plan[]>([])
const loadingPlans = ref(false)
const auditLogs = ref<SubscriptionAuditLog[]>([])
const loadingAuditLogs = ref(false)
const grantingManual = ref(false)
const manualGrantError = ref('')
const manualGrantSuccess = ref('')
const manualGrantForm = ref({
  plan_id: '',
  has_expiration: false,
  expires_at: '',
  reason: '',
})

const minExpirationDate = computed(() => {
  const tomorrow = new Date()
  tomorrow.setDate(tomorrow.getDate() + 1)
  return tomorrow.toISOString().split('T')[0]
})

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

function formatDateTime(dateStr?: string) {
  if (!dateStr) return '-'
  try {
    return new Date(dateStr).toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    })
  } catch {
    return dateStr
  }
}

function formatOrigin(origin?: string) {
  switch (origin) {
    case 'MANUAL':
      return 'Concessão Manual (Sem Asaas)'
    case 'FREE_PLAN':
      return 'Plano 100% Gratuito'
    case 'ASAAS':
      return 'Gateway Asaas'
    default:
      return origin || 'Asaas'
  }
}

function getAuditActionInfo(action: string) {
  switch (action) {
    case 'MANUAL_GRANT':
      return { label: 'Concessão Manual', class: 'bg-purple-500/10 text-purple-400 border-purple-500/30' }
    case 'FREE_REGISTRATION':
      return { label: 'Cadastro Gratuito', class: 'bg-cyan-500/10 text-cyan-400 border-cyan-500/30' }
    case 'STATUS_OVERRIDE':
      return { label: 'Alteração de Status', class: 'bg-amber-500/10 text-amber-400 border-amber-500/30' }
    case 'PLAN_CHANGE':
      return { label: 'Mudança de Plano', class: 'bg-indigo-500/10 text-indigo-400 border-indigo-500/30' }
    case 'AUTO_EXPIRED':
      return { label: 'Expiração Automática', class: 'bg-rose-500/10 text-rose-400 border-rose-500/30' }
    case 'ASAAS_SYNC':
      return { label: 'Sincronização Asaas', class: 'bg-blue-500/10 text-blue-400 border-blue-500/30' }
    default:
      return { label: action, class: 'bg-surface-800 text-surface-400 border-surface-700' }
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

async function fetchAvailablePlans() {
  loadingPlans.value = true
  try {
    const res = await api.get('/admin/plans?is_active=true')
    if (res.data.success) {
      availablePlans.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao carregar planos disponíveis:', err)
  } finally {
    loadingPlans.value = false
  }
}

async function fetchAuditLogs(subscriptionId: string) {
  loadingAuditLogs.value = true
  try {
    const res = await api.get(`/admin/subscriptions/${subscriptionId}/audit-logs`)
    if (res.data.success) {
      auditLogs.value = res.data.data
    }
  } catch (err: any) {
    console.error('Erro ao carregar histórico de auditoria:', err)
    auditLogs.value = []
  } finally {
    loadingAuditLogs.value = false
  }
}

async function openSubscriptionModal(tenant: Tenant) {
  selectedTenant.value = tenant
  selectedSubscription.value = tenant.subscription || null
  overrideStatusForm.value = tenant.subscription?.status || 'ACTIVE'
  overrideReasonForm.value = ''
  manualGrantError.value = ''
  manualGrantSuccess.value = ''
  activeSubTab.value = tenant.subscription ? 'details' : 'manual_grant'
  showSubscriptionModal.value = true

  if (availablePlans.value.length === 0) {
    await fetchAvailablePlans()
  }

  manualGrantForm.value = {
    plan_id: tenant.subscription?.plan_id || (availablePlans.value[0]?.id || ''),
    has_expiration: !!tenant.subscription?.current_period_end,
    expires_at: tenant.subscription?.current_period_end
      ? tenant.subscription.current_period_end.split('T')[0]
      : '',
    reason: '',
  }

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
      if (!manualGrantForm.value.plan_id) {
        manualGrantForm.value.plan_id = res.data.data.subscription.plan_id
      }
    }
  } catch (err: any) {
    console.error('Erro ao carregar detalhes completos da assinatura:', err)
  } finally {
    loadingSubscriptionDetails.value = false
  }

  if (selectedSubscription.value?.id) {
    await fetchAuditLogs(selectedSubscription.value.id)
  } else {
    auditLogs.value = []
  }
}

async function submitManualGrant() {
  if (!selectedTenant.value) return
  if (!manualGrantForm.value.plan_id) {
    manualGrantError.value = 'Por favor, selecione um plano para associar.'
    return
  }
  if (!manualGrantForm.value.reason || manualGrantForm.value.reason.trim().length < 3) {
    manualGrantError.value = 'Por favor, informe a justificativa da concessão (mínimo de 3 caracteres).'
    return
  }
  if (manualGrantForm.value.has_expiration && !manualGrantForm.value.expires_at) {
    manualGrantError.value = 'Por favor, informe a data limite de validade da assinatura.'
    return
  }

  grantingManual.value = true
  manualGrantError.value = ''
  manualGrantSuccess.value = ''

  try {
    let expiresAtISO: string | null = null
    if (manualGrantForm.value.has_expiration && manualGrantForm.value.expires_at) {
      expiresAtISO = new Date(manualGrantForm.value.expires_at + 'T23:59:59').toISOString()
    }

    const payload: any = {
      plan_id: manualGrantForm.value.plan_id,
      reason: manualGrantForm.value.reason.trim(),
    }
    if (expiresAtISO) {
      payload.expires_at = expiresAtISO
    }

    const res = await api.post(`/admin/tenants/${selectedTenant.value.id}/subscriptions/grant-manual`, payload)
    if (res.data.success) {
      manualGrantSuccess.value = 'Assinatura concedida com sucesso!'
      selectedSubscription.value = res.data.data
      overrideStatusForm.value = res.data.data.status
      if (selectedTenant.value) {
        selectedTenant.value.subscription = res.data.data
      }
      await fetchTenants()
      if (res.data.data.id) {
        await fetchAuditLogs(res.data.data.id)
      }
      setTimeout(() => {
        activeSubTab.value = 'details'
        manualGrantSuccess.value = ''
      }, 1200)
    }
  } catch (err: any) {
    manualGrantError.value = err.response?.data?.error || err.message || 'Falha ao conceder assinatura manual.'
  } finally {
    grantingManual.value = false
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
      reason: overrideReasonForm.value.trim() || undefined,
    })
    if (res.data.success) {
      selectedSubscription.value.status = newStatus
      if (selectedTenant.value?.subscription) {
        selectedTenant.value.subscription.status = newStatus
      }
      overrideReasonForm.value = ''
      await fetchTenants()
      if (selectedSubscription.value.id) {
        await fetchAuditLogs(selectedSubscription.value.id)
      }
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
  fetchAvailablePlans()
})
</script>
