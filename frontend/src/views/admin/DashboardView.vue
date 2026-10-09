<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- ======================================================== -->
    <!-- 1. DASHBOARD GLOBAL DA PLATAFORMA (SUPER ADMIN) -->
    <!-- ======================================================== -->
    <template v-if="authStore.isGlobalAdmin">
      <!-- Header Global -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 class="text-xl sm:text-2xl font-black text-[#202224] dark:text-white font-display tracking-tight flex items-center gap-2.5">
            <ShieldCheck class="w-7 h-7 text-[#4880FF]" />
            <span>Visão Global da Rede Multi-Tenant</span>
          </h2>
          <p class="text-xs sm:text-sm text-[#718096] dark:text-zinc-400 mt-1">
            Métricas operacionais consolidadas e faturamento em tempo real de toda a plataforma.
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadGlobalKPIs"
            class="flex items-center gap-1.5 px-3.5 py-2.5 rounded-xl bg-white dark:bg-[#121318] border border-gray-200/80 dark:border-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#4880FF] dark:hover:text-white text-xs font-bold transition shadow-sm cursor-pointer"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/estabelecimentos"
            class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Estabelecimento</span>
          </RouterLink>
        </div>
      </div>

      <!-- 4 Cards Superiores (Figma Standard) -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Card 1: Estabelecimentos (Azul) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#E9F0FE] text-[#4880FF] flex items-center justify-center shrink-0 shadow-sm">
            <Building2 class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ globalStats.total_tenants || 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Estabelecimentos Cadastrados</p>
            <span class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold mt-0.5 block">
              {{ globalStats.active_tenants || 0 }} ativos na rede
            </span>
          </div>
        </div>

        <!-- Card 2: Faturamento Total (Coral) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#FFEFE7] text-[#FF6647] flex items-center justify-center shrink-0 shadow-sm">
            <DollarSign class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              R$ {{ (globalStats.total_revenue || 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Faturamento Estimado</p>
            <span class="text-[11px] text-[#718096] dark:text-zinc-500 font-semibold mt-0.5 block">
              Volume geral transacionado
            </span>
          </div>
        </div>

        <!-- Card 3: Agendamentos Totais (Amarelo) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#FFF7E6] text-[#FEC53D] flex items-center justify-center shrink-0 shadow-sm">
            <Calendar class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ globalStats.total_appointments || 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Atendimentos Totais</p>
            <span class="text-[11px] text-[#4880FF] font-semibold mt-0.5 block">
              {{ globalStats.completed_appointments || 0 }} concluídos
            </span>
          </div>
        </div>

        <!-- Card 4: Usuários & Clientes (Roxo) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#ECEBFE] text-[#8280FF] flex items-center justify-center shrink-0 shadow-sm">
            <Users class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ globalStats.total_customers || 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Base de Clientes</p>
            <span class="text-[11px] text-purple-600 dark:text-purple-400 font-semibold mt-0.5 block">
              {{ globalStats.total_users || 0 }} operadores do sistema
            </span>
          </div>
        </div>
      </div>

      <!-- Tabelas Recentes da Plataforma -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Estabelecimentos Recentes -->
        <div class="saas-card overflow-hidden">
          <div class="p-5 border-b border-gray-100 dark:border-zinc-800 flex items-center justify-between">
            <div class="flex items-center space-x-2.5">
              <Store class="w-5 h-5 text-[#4880FF]" />
              <h3 class="font-bold text-base text-[#202224] dark:text-white font-display">Estabelecimentos Recentes</h3>
            </div>
            <RouterLink to="/admin/estabelecimentos" class="text-xs text-[#4880FF] hover:underline font-bold">
              Ver Todos →
            </RouterLink>
          </div>

          <div v-if="!globalStats.recent_tenants || globalStats.recent_tenants.length === 0" class="p-8 text-center text-[#718096] text-xs">
            Nenhum estabelecimento cadastrado.
          </div>
          <div v-else class="divide-y divide-gray-100 dark:divide-zinc-800/80">
            <div v-for="t in globalStats.recent_tenants" :key="t.id" class="p-4 flex items-center justify-between hover:bg-gray-50/60 dark:hover:bg-zinc-800/40 transition">
              <div class="flex items-center gap-3">
                <div class="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-950/40 border border-blue-100 dark:border-blue-900/40 flex items-center justify-center font-bold text-sm text-[#4880FF] shrink-0 font-display">
                  {{ t.name.charAt(0) }}
                </div>
                <div>
                  <h4 class="font-bold text-sm text-[#202224] dark:text-white font-display">{{ t.name }}</h4>
                  <p class="text-xs text-[#718096] dark:text-zinc-400">/{{ t.slug }} • {{ t.city || 'Brasil' }}</p>
                </div>
              </div>
              <span
                class="px-2.5 py-1 rounded-full text-[10px] font-bold border"
                :class="t.is_active ? 'bg-emerald-50 text-emerald-600 border-emerald-200 dark:bg-emerald-950/60 dark:text-emerald-300 dark:border-emerald-800' : 'bg-rose-50 text-rose-600 border-rose-200 dark:bg-rose-950/60 dark:text-rose-300 dark:border-rose-800'"
              >
                {{ t.is_active ? '● Ativo' : '○ Inativo' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Atividade Recente -->
        <div class="saas-card overflow-hidden">
          <div class="p-5 border-b border-gray-100 dark:border-zinc-800 flex items-center justify-between">
            <div class="flex items-center space-x-2.5">
              <Clock class="w-5 h-5 text-[#FEC53D]" />
              <h3 class="font-bold text-base text-[#202224] dark:text-white font-display">Atividade Recente na Plataforma</h3>
            </div>
          </div>

          <div v-if="!globalStats.recent_appointments || globalStats.recent_appointments.length === 0" class="p-8 text-center text-[#718096] text-xs">
            Nenhuma atividade recente registrada.
          </div>
          <div v-else class="divide-y divide-gray-100 dark:divide-zinc-800/80">
            <div v-for="apt in globalStats.recent_appointments" :key="apt.id" class="p-4 flex items-center justify-between hover:bg-gray-50/60 dark:hover:bg-zinc-800/40 transition">
              <div>
                <h4 class="font-bold text-sm text-[#202224] dark:text-white font-display">{{ apt.customer?.name || 'Cliente' }}</h4>
                <p class="text-xs text-[#718096] dark:text-zinc-400 mt-0.5">
                  <span class="text-[#4880FF] font-semibold">{{ apt.service?.name }}</span> com {{ apt.professional?.name }}
                </p>
              </div>
              <div class="text-right">
                <p class="text-xs font-black text-[#202224] dark:text-white font-display">R$ {{ apt.total_price?.toFixed(2).replace('.', ',') }}</p>
                <span :class="getStatusBadgeClass(apt.status)" class="inline-block mt-1">
                  {{ formatStatus(apt.status) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- ======================================================== -->
    <!-- 2. DASHBOARD DO TENANT (PADRÃO FIGMA SAAS DASHBOARD) -->
    <!-- ======================================================== -->
    <template v-else>
      <!-- Sub-header com Ações Rápidas -->
      <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
        <div>
          <h2 class="text-xl sm:text-2xl font-black text-[#202224] dark:text-white font-display tracking-tight">
            Dashboard
          </h2>
          <p class="text-xs sm:text-sm text-[#718096] dark:text-zinc-400 mt-0.5">
            Acompanhe métricas, tendências de atendimento e faturamento em tempo real.
          </p>
        </div>

        <div class="flex items-center gap-3">
          <button
            @click="loadAnalytics"
            class="flex items-center gap-1.5 px-3.5 py-2.5 rounded-xl bg-white dark:bg-[#121318] border border-gray-200/80 dark:border-zinc-800 text-[#718096] dark:text-zinc-300 hover:text-[#4880FF] dark:hover:text-white text-xs font-bold transition shadow-sm cursor-pointer"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': isLoading }" />
            <span>Atualizar</span>
          </button>
          <RouterLink
            to="/admin/agenda"
            class="flex items-center gap-1.5 px-4 py-2.5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white text-xs font-bold transition shadow-md shadow-blue-500/25 cursor-pointer"
          >
            <Plus class="w-4 h-4" />
            <span>Novo Agendamento</span>
          </RouterLink>
        </div>
      </div>

      <!-- ======================================================== -->
      <!-- LINHA 1: 4 KPI CARDS (FIGMA: Save Products, Stock, Sales, Job Application) -->
      <!-- ======================================================== -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- KPI 1: Agendamentos Hoje (Azul Pastel / #4880FF) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#E9F0FE] text-[#4880FF] flex items-center justify-center shrink-0 shadow-sm">
            <CalendarCheck class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ analyticsData?.kpis?.today_appointments_count ?? 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Agendamentos Hoje</p>
            <span class="text-[11px] text-[#4880FF] font-semibold mt-0.5 block">
              {{ analyticsData?.kpis?.confirmed_count ?? 0 }} confirmados
            </span>
          </div>
        </div>

        <!-- KPI 2: Atendimentos Pendentes (Amarelo Pastel / #FFB800) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#FFF7E6] text-[#FFB800] flex items-center justify-center shrink-0 shadow-sm">
            <Clock class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ analyticsData?.kpis?.pending_count ?? 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Atendimentos Pendentes</p>
            <span class="text-[11px] text-amber-600 dark:text-amber-400 font-semibold mt-0.5 block">
              Aguardando confirmação
            </span>
          </div>
        </div>

        <!-- KPI 3: Faturamento Estimado (Coral Pastel / #FF6647) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#FFEFE7] text-[#FF6647] flex items-center justify-center shrink-0 shadow-sm">
            <DollarSign class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              R$ {{ (analyticsData?.kpis?.today_revenue ?? 0).toLocaleString('pt-BR', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Faturamento Estimado</p>
            <span class="text-[11px] text-emerald-600 dark:text-emerald-400 font-semibold mt-0.5 block">
              Previsão do dia
            </span>
          </div>
        </div>

        <!-- KPI 4: Clientes na Base (Roxo Pastel / #8280FF) -->
        <div class="saas-card p-6 flex items-center gap-4">
          <div class="w-14 h-14 rounded-full bg-[#ECEBFE] text-[#8280FF] flex items-center justify-center shrink-0 shadow-sm">
            <Users class="w-7 h-7" />
          </div>
          <div class="min-w-0">
            <h3 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
              {{ analyticsData?.kpis?.total_customers_count ?? 0 }}
            </h3>
            <p class="text-xs text-[#718096] dark:text-zinc-400 font-medium">Total de Clientes</p>
            <span class="text-[11px] text-purple-600 dark:text-purple-400 font-semibold mt-0.5 block">
              {{ analyticsData?.kpis?.active_professionals_count ?? 0 }} profissionais ativos
            </span>
          </div>
        </div>
      </div>

      <!-- ======================================================== -->
      <!-- LINHA 2: GRÁFICO REPORTS (SPLINE) + GRÁFICO ANALYTICS (DONUT) -->
      <!-- ======================================================== -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- CARD REPORTS: Curva Spline com Gradiente (Figma 4:117) -->
        <div class="lg:col-span-8 saas-card p-6 flex flex-col justify-between">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="font-bold text-lg text-[#202224] dark:text-white font-display">Reports</h3>
              <p class="text-xs text-[#718096] dark:text-zinc-400">Fluxo e volume de atendimentos ao longo do dia</p>
            </div>
            <button class="p-1.5 text-gray-400 hover:text-[#4880FF] rounded-lg transition" title="Mais Opções">
              <MoreHorizontal class="w-5 h-5" />
            </button>
          </div>

          <!-- Área do Gráfico Spline em SVG -->
          <div class="relative w-full h-64 sm:h-72 mt-2">
            <!-- Grid Lines de fundo estilo Figma -->
            <div class="absolute inset-0 flex flex-col justify-between pointer-events-none opacity-40">
              <div class="border-b border-gray-100 dark:border-zinc-800 w-full"></div>
              <div class="border-b border-gray-100 dark:border-zinc-800 w-full"></div>
              <div class="border-b border-gray-100 dark:border-zinc-800 w-full"></div>
              <div class="border-b border-gray-100 dark:border-zinc-800 w-full"></div>
              <div class="border-b border-gray-100 dark:border-zinc-800 w-full"></div>
            </div>

            <!-- SVG Spline Curve -->
            <svg class="w-full h-full overflow-visible" viewBox="0 0 700 240" preserveAspectRatio="none">
              <defs>
                <!-- Gradiente de Preenchimento Translúcido (Azul -> Violeta) -->
                <linearGradient id="splineGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stop-color="#4880FF" stop-opacity="0.35" />
                  <stop offset="100%" stop-color="#8280FF" stop-opacity="0.10" />
                </linearGradient>
                <!-- Gradiente do Traço da Linha -->
                <linearGradient id="strokeGradient" x1="0%" y1="0%" x2="100%" y2="0%">
                  <stop offset="0%" stop-color="#4880FF" />
                  <stop offset="50%" stop-color="#5B8EFF" />
                  <stop offset="100%" stop-color="#A855F7" />
                </linearGradient>
              </defs>

              <!-- Área sob a curva -->
              <path
                :d="splineAreaPath"
                fill="url(#splineGradient)"
                class="transition-all duration-500"
              />

              <!-- Linha Bezier Principal -->
              <path
                :d="splineLinePath"
                fill="none"
                stroke="url(#strokeGradient)"
                stroke-width="3.5"
                stroke-linecap="round"
                class="transition-all duration-500"
              />

              <!-- Linha Pontilhada Vertical Guia no Ponto de Pico (Figma Line 52) -->
              <line
                v-if="highlightPoint"
                :x1="highlightPoint.x"
                y1="30"
                :x2="highlightPoint.x"
                y2="210"
                stroke="#4880FF"
                stroke-width="1.5"
                stroke-dasharray="4 4"
                opacity="0.8"
              />

              <!-- Pontos de dados interativos -->
              <g v-for="(pt, idx) in splinePoints" :key="idx">
                <circle
                  :cx="pt.x"
                  :cy="pt.y"
                  r="5"
                  class="fill-white dark:fill-[#121318] stroke-[#4880FF] stroke-[3] cursor-pointer hover:r-7 transition-all"
                  @mouseenter="activePoint = pt"
                />
              </g>
            </svg>

            <!-- Tooltip Flutuante estilo Figma: "Sales / Atendimentos: 2,678" -->
            <div
              v-if="activePoint"
              class="absolute z-20 px-3 py-1.5 rounded-xl bg-[#202224] text-white text-[11px] font-bold shadow-xl pointer-events-none transform -translate-x-1/2 -translate-y-full transition-all"
              :style="{ left: `${(activePoint.x / 700) * 100}%`, top: `${(activePoint.y / 240) * 100}%` }"
            >
              <div class="text-[9px] text-gray-300 uppercase tracking-wider font-semibold">Atendimentos: {{ activePoint.data.appointments }}</div>
              <div class="text-emerald-400 font-mono">R$ {{ activePoint.data.revenue.toFixed(2).replace('.', ',') }}</div>
            </div>
          </div>

          <!-- Rótulos do Eixo X (Horários) -->
          <div class="flex items-center justify-between pt-3 text-[11px] text-[#718096] dark:text-zinc-500 font-medium">
            <span v-for="(lbl, idx) in chartLabels" :key="idx">{{ lbl }}</span>
          </div>
        </div>

        <!-- CARD ANALYTICS: Donut Chart de Status (Figma 4:92) -->
        <div class="lg:col-span-4 saas-card p-6 flex flex-col justify-between">
          <div class="flex items-center justify-between mb-2">
            <div>
              <h3 class="font-bold text-lg text-[#202224] dark:text-white font-display">Analytics</h3>
              <p class="text-xs text-[#718096] dark:text-zinc-400">Distribuição de Status dos Atendimentos</p>
            </div>
            <button class="p-1.5 text-gray-400 hover:text-[#4880FF] rounded-lg transition" title="Mais Opções">
              <MoreHorizontal class="w-5 h-5" />
            </button>
          </div>

          <!-- Gráfico de Rosca (Donut SVG) -->
          <div class="relative w-48 h-48 sm:w-52 sm:h-52 mx-auto my-4 flex items-center justify-center">
            <svg class="w-full h-full -rotate-90" viewBox="0 0 100 100">
              <!-- Anel Cinza de Fundo -->
              <circle
                cx="50"
                cy="50"
                r="38"
                fill="none"
                class="stroke-gray-100 dark:stroke-zinc-800"
                stroke-width="12"
              />

              <!-- Segmento Confirmados (#4880FF) -->
              <circle
                cx="50"
                cy="50"
                r="38"
                fill="none"
                stroke="#4880FF"
                stroke-width="12"
                stroke-linecap="round"
                :stroke-dasharray="donutSegments.confirmed.dash"
                :stroke-dashoffset="donutSegments.confirmed.offset"
                class="transition-all duration-700"
              />

              <!-- Segmento Concluídos (#FEC53D) -->
              <circle
                cx="50"
                cy="50"
                r="38"
                fill="none"
                stroke="#FEC53D"
                stroke-width="12"
                stroke-linecap="round"
                :stroke-dasharray="donutSegments.completed.dash"
                :stroke-dashoffset="donutSegments.completed.offset"
                class="transition-all duration-700"
              />

              <!-- Segmento Cancelados (#FF6647) -->
              <circle
                cx="50"
                cy="50"
                r="38"
                fill="none"
                stroke="#FF6647"
                stroke-width="12"
                stroke-linecap="round"
                :stroke-dasharray="donutSegments.cancelled.dash"
                :stroke-dashoffset="donutSegments.cancelled.offset"
                class="transition-all duration-700"
              />
            </svg>

            <!-- Rótulo Central do Donut: "80% Transactions" -->
            <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none text-center">
              <span class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white font-display tracking-tight">
                {{ conversionRate }}%
              </span>
              <span class="text-[10px] uppercase font-bold text-[#718096] dark:text-zinc-400">
                Concluídos
              </span>
            </div>
          </div>

          <!-- Legenda do Gráfico (Padrão Figma: Bullets redondos) -->
          <div class="flex items-center justify-around pt-3 border-t border-gray-100 dark:border-zinc-800 text-xs">
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full bg-[#4880FF]"></span>
              <span class="text-[#718096] dark:text-zinc-400 font-medium">Confirmados</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full bg-[#FEC53D]"></span>
              <span class="text-[#718096] dark:text-zinc-400 font-medium">Concluídos</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-3 h-3 rounded-full bg-[#FF6647]"></span>
              <span class="text-[#718096] dark:text-zinc-400 font-medium">Cancelados</span>
            </div>
          </div>
        </div>
      </div>

      <!-- ======================================================== -->
      <!-- LINHA 3: TABELA RECENT ORDERS + CARD TOP SELLING PRODUCTS -->
      <!-- ======================================================== -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- TABELA RECENT ORDERS / ÚLTIMOS AGENDAMENTOS (Figma 4:36) -->
        <div class="lg:col-span-8 saas-card overflow-hidden">
          <div class="p-6 border-b border-gray-100 dark:border-zinc-800 flex items-center justify-between">
            <div>
              <h3 class="font-bold text-lg text-[#202224] dark:text-white font-display">Recent Orders (Atendimentos Recentes)</h3>
              <p class="text-xs text-[#718096] dark:text-zinc-400">Últimos agendamentos efetuados pelos clientes</p>
            </div>
            <RouterLink to="/admin/agenda" class="text-xs text-[#4880FF] hover:underline font-bold">
              Ver Todos →
            </RouterLink>
          </div>

          <div v-if="!analyticsData?.recent_appointments || analyticsData.recent_appointments.length === 0" class="p-10 text-center text-[#718096] text-xs">
            Nenhum atendimento recente registrado.
          </div>
          <div v-else class="overflow-x-auto">
            <table class="w-full text-left text-xs text-[#202224] dark:text-zinc-300">
              <thead class="bg-gray-50/70 dark:bg-zinc-800/50 text-[#718096] dark:text-zinc-400 uppercase font-bold text-[11px] tracking-wider border-b border-gray-100 dark:border-zinc-800">
                <tr>
                  <th class="p-4">Cliente</th>
                  <th class="p-4">Serviço / Profissional</th>
                  <th class="p-4">Horário</th>
                  <th class="p-4 text-center">Status</th>
                  <th class="p-4 text-right">Valor</th>
                  <th class="p-4 text-center">Ações</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-100 dark:divide-zinc-800/80">
                <tr
                  v-for="apt in analyticsData.recent_appointments"
                  :key="apt.id"
                  class="hover:bg-gray-50/60 dark:hover:bg-zinc-800/40 transition"
                >
                  <!-- Cliente com Avatar -->
                  <td class="p-4">
                    <div class="flex items-center gap-3">
                      <div class="w-8 h-8 rounded-full bg-[#E9F0FE] text-[#4880FF] font-bold text-xs flex items-center justify-center shrink-0 font-display">
                        {{ (apt.customer?.name || 'C').charAt(0).toUpperCase() }}
                      </div>
                      <span class="font-bold text-[#202224] dark:text-white truncate max-w-[140px]">
                        {{ apt.customer?.name || 'Cliente' }}
                      </span>
                    </div>
                  </td>

                  <!-- Serviço & Profissional -->
                  <td class="p-4">
                    <div class="font-medium text-[#202224] dark:text-white truncate max-w-[160px]">
                      {{ apt.service?.name }}
                    </div>
                    <div class="text-[11px] text-[#718096] dark:text-zinc-400">
                      com {{ apt.professional?.name }}
                    </div>
                  </td>

                  <!-- Horário -->
                  <td class="p-4 font-mono font-semibold text-[#718096] dark:text-zinc-300">
                    {{ formatTime(apt.start_at) }}
                  </td>

                  <!-- Status em Pílula Suave (Figma Total Order Badge) -->
                  <td class="p-4 text-center">
                    <span :class="getStatusBadgeClass(apt.status)">
                      {{ formatStatus(apt.status) }}
                    </span>
                  </td>

                  <!-- Valor -->
                  <td class="p-4 text-right font-black text-[#202224] dark:text-white font-display">
                    R$ {{ (apt.total_price || 0).toFixed(2).replace('.', ',') }}
                  </td>

                  <!-- Ações Rápidas -->
                  <td class="p-4 text-center">
                    <div v-if="apt.status === 'CONFIRMED'" class="flex items-center justify-center gap-1.5">
                      <button
                        @click="completeAppointment(apt.id)"
                        class="px-2 py-1 rounded-lg bg-emerald-50 dark:bg-emerald-950/40 text-emerald-600 dark:text-emerald-400 hover:bg-emerald-100 text-[11px] font-bold transition"
                        title="Concluir"
                      >
                        ✓
                      </button>
                      <button
                        @click="cancelAppointment(apt.id)"
                        class="px-2 py-1 rounded-lg bg-rose-50 dark:bg-rose-950/40 text-rose-600 dark:text-rose-400 hover:bg-rose-100 text-[11px] font-bold transition"
                        title="Cancelar"
                      >
                        ✕
                      </button>
                    </div>
                    <span v-else class="text-[11px] text-[#718096] dark:text-zinc-500">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- CARD TOP SELLING PRODUCTS / SERVIÇOS MAIS PROCURADOS (Figma 4:3) -->
        <div class="lg:col-span-4 saas-card p-6 flex flex-col justify-between">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="font-bold text-lg text-[#202224] dark:text-white font-display">Top Selling Products</h3>
              <p class="text-xs text-[#718096] dark:text-zinc-400">Serviços com maior volume de agendamento</p>
            </div>
            <button class="p-1.5 text-gray-400 hover:text-[#4880FF] rounded-lg transition" title="Mais Opções">
              <MoreHorizontal class="w-5 h-5" />
            </button>
          </div>

          <div v-if="!analyticsData?.top_services || analyticsData.top_services.length === 0" class="py-8 text-center text-[#718096] text-xs">
            Nenhum serviço registrado com agendamentos.
          </div>
          <div v-else class="space-y-4 flex-1">
            <div
              v-for="svc in analyticsData.top_services"
              :key="svc.service_id"
              class="flex items-center gap-4 p-3 rounded-2xl hover:bg-gray-50/80 dark:hover:bg-zinc-800/40 transition"
            >
              <!-- Thumbnail arredondado (Figma Blue Box) -->
              <div class="w-14 h-14 rounded-2xl bg-[#E9F0FE] text-[#4880FF] flex items-center justify-center shrink-0 shadow-sm font-display font-black text-lg">
                <Sparkles class="w-6 h-6" />
              </div>

              <!-- Informações do Serviço & Estrelas -->
              <div class="min-w-0 flex-1">
                <h4 class="font-bold text-sm text-[#202224] dark:text-white truncate font-display">
                  {{ svc.service_name }}
                </h4>

                <!-- Avaliação em Estrelas Douradas -->
                <div class="flex items-center gap-1 mt-1 text-[#FEC53D]">
                  <Star v-for="i in 5" :key="i" class="w-3.5 h-3.5 fill-[#FEC53D]" />
                </div>

                <div class="flex items-center justify-between mt-1">
                  <span class="text-xs font-black text-[#202224] dark:text-white font-display">
                    R$ {{ svc.price.toFixed(2).replace('.', ',') }}
                  </span>
                  <span class="text-[11px] text-[#718096] dark:text-zinc-400">
                    {{ svc.bookings_count }} pedidos
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import {
  Calendar, DollarSign, Users, Clock, Plus, RefreshCw,
  Store, Building2, ShieldCheck, CalendarCheck, MoreHorizontal,
  Sparkles, Star
} from 'lucide-vue-next'
import { useAuthStore } from '../../stores/auth'
import api from '../../services/api'

const authStore = useAuthStore()
const isLoading = ref(false)

// Dados do Tenant (SaaS Figma Dashboard)
interface AnalyticsPayload {
  kpis: {
    today_appointments_count: number
    confirmed_count: number
    completed_count: number
    pending_count: number
    today_revenue: number
    total_customers_count: number
    active_professionals_count: number
  }
  reports_chart: Array<{ label: string; appointments: number; revenue: number }>
  status_distribution: Array<{ status: string; label: string; count: number; percentage: number; color: string }>
  top_services: Array<{ service_id: string; service_name: string; price: number; bookings_count: number; rating: number }>
  recent_appointments: any[]
}

const analyticsData = ref<AnalyticsPayload | null>(null)
const globalStats = ref<any>({})

// Ponto ativo do gráfico com tooltip
const activePoint = ref<{ x: number; y: number; data: any } | null>(null)

onMounted(async () => {
  if (authStore.isGlobalAdmin) {
    await loadGlobalKPIs()
  } else {
    await loadAnalytics()
  }
})

async function loadAnalytics() {
  isLoading.value = true
  try {
    const res = await api.get('/admin/dashboard/analytics')
    if (res.data.success && res.data.data) {
      analyticsData.value = res.data.data
    }
  } catch (err) {
    console.error('Falha ao carregar dashboard analytics:', err)
  } finally {
    isLoading.value = false
  }
}

async function loadGlobalKPIs() {
  isLoading.value = true
  try {
    const res = await api.get('/admin/global-dashboard')
    if (res.data.success && res.data.data) {
      globalStats.value = res.data.data
    }
  } catch (err) {
    console.error('Falha ao carregar estatísticas globais:', err)
  } finally {
    isLoading.value = false
  }
}

// -------------------------------------------------------------
// CÁLCULOS DO GRÁFICO SPLINE REPORTS (FIGMA BEZIER CURVE)
// -------------------------------------------------------------
const chartLabels = computed(() => {
  if (analyticsData.value?.reports_chart?.length) {
    return analyticsData.value.reports_chart.map(p => p.label)
  }
  return ['08:00', '10:00', '12:00', '14:00', '16:00', '18:00', '20:00']
})

const splinePoints = computed(() => {
  const chart = analyticsData.value?.reports_chart || []
  if (chart.length === 0) return []

  const maxVal = Math.max(...chart.map(p => p.appointments), 5)
  const width = 700
  const height = 240
  const paddingX = 40
  const paddingY = 40

  const availableW = width - (paddingX * 2)
  const availableH = height - (paddingY * 2)
  const stepX = availableW / (chart.length - 1 || 1)

  return chart.map((pt, i) => {
    const x = paddingX + (i * stepX)
    const ratio = pt.appointments / maxVal
    const y = (height - paddingY) - (ratio * availableH)
    return { x, y, data: pt }
  })
})

const highlightPoint = computed(() => {
  const pts = splinePoints.value
  if (!pts.length) return null
  return pts.reduce((max, p) => p.data.appointments > (max?.data.appointments || 0) ? p : max, pts[0])
})

const splineLinePath = computed(() => {
  const pts = splinePoints.value
  if (pts.length === 0) return ''
  if (pts.length === 1) return `M ${pts[0].x} ${pts[0].y}`

  let d = `M ${pts[0].x} ${pts[0].y}`
  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[i]
    const p1 = pts[i + 1]
    const cpX = (p0.x + p1.x) / 2
    d += ` C ${cpX} ${p0.y}, ${cpX} ${p1.y}, ${p1.x} ${p1.y}`
  }
  return d
})

const splineAreaPath = computed(() => {
  const line = splineLinePath.value
  const pts = splinePoints.value
  if (!line || !pts.length) return ''
  const first = pts[0]
  const last = pts[pts.length - 1]
  return `${line} L ${last.x} 230 L ${first.x} 230 Z`
})

// -------------------------------------------------------------
// CÁLCULOS DO GRÁFICO DONUT ANALYTICS (STATUS DISTRIBUTION)
// -------------------------------------------------------------
const conversionRate = computed(() => {
  const kpis = analyticsData.value?.kpis
  if (!kpis || !kpis.today_appointments_count) return 80
  const rate = Math.round((kpis.completed_count / kpis.today_appointments_count) * 100)
  return rate || 0
})

const donutSegments = computed(() => {
  const circumference = 2 * Math.PI * 38 // ~238.76
  const dist = analyticsData.value?.status_distribution || []

  let confPct = 45
  let compPct = 35
  let cancPct = 20

  dist.forEach(d => {
    if (d.status === 'CONFIRMED') confPct = d.percentage || 10
    if (d.status === 'COMPLETED') compPct = d.percentage || 10
    if (d.status === 'CANCELLED') cancPct = d.percentage || 10
  })

  const confDash = (confPct / 100) * circumference
  const compDash = (compPct / 100) * circumference
  const cancDash = (cancPct / 100) * circumference

  return {
    confirmed: {
      dash: `${confDash} ${circumference}`,
      offset: 0
    },
    completed: {
      dash: `${compDash} ${circumference}`,
      offset: -confDash
    },
    cancelled: {
      dash: `${cancDash} ${circumference}`,
      offset: -(confDash + compDash)
    }
  }
})

// -------------------------------------------------------------
// FORMATADORES E BADGES
// -------------------------------------------------------------
function formatTime(isoString: string) {
  if (!isoString) return '--:--'
  const d = new Date(isoString)
  return d.toLocaleTimeString('pt-BR', { hour: '2-digit', minute: '2-digit' })
}

function formatStatus(status: string) {
  switch (status) {
    case 'CONFIRMED':
      return 'Confirmado'
    case 'COMPLETED':
      return 'Concluído'
    case 'CANCELLED':
      return 'Cancelado'
    case 'PENDING':
      return 'Pendente'
    default:
      return status
  }
}

function getStatusBadgeClass(status: string) {
  switch (status) {
    case 'CONFIRMED':
      return 'px-2.5 py-1 rounded-full text-[10px] font-bold bg-[#E9F0FE] text-[#4880FF] border border-blue-200/80 dark:bg-blue-950/60 dark:text-blue-300 dark:border-blue-800'
    case 'COMPLETED':
      return 'px-2.5 py-1 rounded-full text-[10px] font-bold bg-[#FFF7E6] text-[#D97706] border border-amber-200/80 dark:bg-amber-950/60 dark:text-amber-300 dark:border-amber-800'
    case 'CANCELLED':
      return 'px-2.5 py-1 rounded-full text-[10px] font-bold bg-[#FFEFE7] text-[#DC2626] border border-rose-200/80 dark:bg-rose-950/60 dark:text-rose-300 dark:border-rose-800'
    default:
      return 'px-2.5 py-1 rounded-full text-[10px] font-bold bg-gray-100 text-gray-700 border border-gray-200 dark:bg-zinc-800 dark:text-zinc-300'
  }
}

async function completeAppointment(id: string) {
  try {
    await api.patch(`/admin/appointments/${id}/complete`)
    await loadAnalytics()
  } catch (err) {
    console.error('Falha ao concluir agendamento:', err)
  }
}

async function cancelAppointment(id: string) {
  if (!confirm('Deseja realmente cancelar este agendamento?')) return
  try {
    await api.patch(`/admin/appointments/${id}/cancel`)
    await loadAnalytics()
  } catch (err) {
    console.error('Falha ao cancelar agendamento:', err)
  }
}
</script>
