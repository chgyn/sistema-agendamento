<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Cabeçalho -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight flex items-center gap-3 font-['Outfit']">
          <div class="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-950/60 border border-blue-200 dark:border-blue-900/40 flex items-center justify-center text-[#4880FF] shadow-sm">
            <CalendarCheck2 class="w-6 h-6" />
          </div>
          <span>Google Calendar</span>
        </h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-zinc-400 mt-1">
          Configure suas credenciais Google OAuth (BYOK), conecte sua agenda e sincronize agendamentos automaticamente.
        </p>
      </div>

      <!-- Badge de Status Geral -->
      <div class="flex items-center gap-2 self-start sm:self-auto">
        <span
          v-if="statusData.status === 'CONNECTED'"
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-500 dark:bg-emerald-400 animate-pulse"></span>
          Conectado
        </span>
        <span
          v-else-if="statusData.status === 'EXPIRED'"
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-amber-500 dark:bg-amber-400"></span>
          Autorização Expirada
        </span>
        <span
          v-else-if="statusData.status === 'ERROR'"
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-rose-500 dark:bg-rose-400"></span>
          Erro
        </span>
        <span
          v-else
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-gray-100 dark:bg-zinc-800 text-gray-500 dark:text-zinc-400 border border-gray-200 dark:border-zinc-700"
        >
          <span class="w-2 h-2 rounded-full bg-gray-400 dark:bg-zinc-500"></span>
          Desconectado
        </span>
      </div>
    </div>

    <!-- Mensagens de Alerta Global -->
    <div v-if="successMessage" class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-300 text-sm flex items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-5 h-5 shrink-0 text-emerald-500 dark:text-emerald-400" />
        <span>{{ successMessage }}</span>
      </div>
      <button @click="successMessage = ''" class="text-xs text-emerald-500 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-white cursor-pointer">✕</button>
    </div>

    <div v-if="errorMessage" class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-300 text-sm flex items-center justify-between gap-3">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-5 h-5 shrink-0 text-rose-500 dark:text-rose-400" />
        <span>{{ errorMessage }}</span>
      </div>
      <button @click="errorMessage = ''" class="text-xs text-rose-500 hover:text-rose-700 dark:text-rose-400 dark:hover:text-white cursor-pointer">✕</button>
    </div>

    <div v-if="loadingStatus" class="saas-card p-8 flex items-center justify-center gap-3 text-sm text-gray-500 dark:text-zinc-400">
      <RefreshCw class="w-5 h-5 animate-spin text-[#4880FF]" />
      <span>Carregando configurações do Google Calendar...</span>
    </div>

    <template v-else>
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- ============================================================== -->
        <!-- BLOCO 1: CREDENCIAIS GOOGLE OAUTH (BYOK)                        -->
        <!-- ============================================================== -->
        <div class="lg:col-span-7 saas-card p-6 sm:p-8 space-y-5">
          <div class="flex items-center justify-between gap-3">
            <div>
              <h2 class="text-base font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
                <Key class="w-5 h-5 text-[#4880FF]" />
                1. Credenciais Google OAuth (BYOK)
              </h2>
              <p class="text-xs text-[#718096] dark:text-zinc-400 mt-1">
                Cada estabelecimento utiliza seu próprio aplicativo OAuth no Google Cloud Console.
              </p>
            </div>
            <span
              v-if="statusData.has_credentials"
              class="text-[11px] px-2.5 py-1 rounded-lg bg-emerald-500/10 text-emerald-600 dark:text-emerald-300 border border-emerald-500/20 font-semibold shrink-0"
            >
              ✓ Credenciais salvas
            </span>
          </div>

          <!-- Redirect URI -->
          <div class="p-3.5 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 space-y-1.5">
            <div class="flex items-center justify-between gap-2">
              <span class="text-[11px] font-bold uppercase tracking-wider text-gray-500 dark:text-zinc-400">
                Redirect URI Autorizada (Google Cloud Console)
              </span>
              <button
                type="button"
                @click="copyRedirectURI"
                class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#4880FF] hover:text-[#386FF0] cursor-pointer"
              >
                <Check v-if="copiedRedirect" class="w-3.5 h-3.5 text-emerald-500" />
                <Copy v-else class="w-3.5 h-3.5" />
                <span>{{ copiedRedirect ? 'Copiada!' : 'Copiar' }}</span>
              </button>
            </div>
            <code class="block text-xs font-mono text-[#202224] dark:text-zinc-200 break-all select-all">
              {{ effectiveRedirectURI }}
            </code>
          </div>

          <form @submit.prevent="saveCredentials" class="space-y-4">
            <div class="space-y-1.5">
              <label class="block text-xs font-bold text-gray-700 dark:text-zinc-300">
                Google Client ID
              </label>
              <input
                v-model="credentialsForm.google_client_id"
                type="text"
                placeholder="1234567890-abcdef.apps.googleusercontent.com"
                class="w-full bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 rounded-xl px-3.5 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:ring-2 focus:ring-[#4880FF]/40 focus:border-[#4880FF]"
              />
            </div>

            <div class="space-y-1.5">
              <div class="flex items-center justify-between gap-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-zinc-300">
                  Google Client Secret
                </label>
                <span v-if="statusData.has_client_secret" class="text-[11px] text-emerald-600 dark:text-emerald-400 font-medium">
                  Secret salvo (deixe em branco para manter o atual)
                </span>
              </div>
              <input
                v-model="credentialsForm.google_client_secret"
                type="password"
                autocomplete="new-password"
                :placeholder="statusData.has_client_secret ? '••••••••••••••••••••••••' : 'GOCSPX-...'"
                class="w-full bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 rounded-xl px-3.5 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:ring-2 focus:ring-[#4880FF]/40 focus:border-[#4880FF]"
              />
              <p class="text-[11px] text-gray-500 dark:text-zinc-500">
                Armazenado criptografado com AES-256-GCM. Alterar Client ID ou Secret desconecta a sessão atual até nova autorização.
              </p>
            </div>

            <button
              type="submit"
              :disabled="savingCredentials"
              class="py-2.5 px-5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] disabled:opacity-60 text-white font-bold text-sm shadow-md shadow-blue-500/20 transition inline-flex items-center gap-2 cursor-pointer"
            >
              <RefreshCw v-if="savingCredentials" class="w-4 h-4 animate-spin" />
              <Save v-else class="w-4 h-4" />
              <span>Salvar Credenciais</span>
            </button>
          </form>
        </div>

        <!-- Passo a Passo Google Cloud Console -->
        <div class="lg:col-span-5 saas-card p-6 sm:p-8 space-y-4">
          <h3 class="text-xs font-bold text-gray-500 dark:text-zinc-400 uppercase tracking-wider flex items-center gap-2">
            <ShieldCheck class="w-4 h-4 text-[#4880FF]" />
            Como configurar no Google Cloud
          </h3>

          <ol class="space-y-3 text-xs text-gray-700 dark:text-zinc-300">
            <li class="flex items-start gap-3 p-3 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">1</span>
              <span>Acesse o <strong>Google Cloud Console</strong>, crie ou selecione um projeto e ative a <strong>Google Calendar API</strong>.</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">2</span>
              <span>Configure a <strong>Tela de permissão OAuth</strong> (adicione os usuários de teste caso esteja em modo de teste).</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">3</span>
              <span>Em <strong>Credenciais &gt; ID do cliente OAuth &gt; Aplicativo da Web</strong>, cadastre exatamente a <strong>Redirect URI</strong> ao lado.</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">4</span>
              <span>Cole o <strong>Client ID</strong> e o <strong>Client Secret</strong>, salve e clique em <strong>Conectar com Google</strong>.</span>
            </li>
          </ol>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- BLOCO 2 & 3: CONEXÃO DA CONTA E CONFIGURAÇÕES DE SINCRONIZAÇÃO -->
      <!-- ============================================================== -->
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Bloco 2: Conexão da Conta Google -->
        <div class="lg:col-span-5 saas-card p-6 sm:p-8 space-y-5">
          <h2 class="text-base font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
            <Link2 class="w-5 h-5 text-[#4880FF]" />
            2. Conexão da Conta Google
          </h2>

          <div v-if="statusData.status === 'CONNECTED'" class="space-y-4">
            <div class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/20 space-y-2">
              <div class="text-xs font-semibold text-emerald-700 dark:text-emerald-300">Conta conectada</div>
              <div class="text-sm font-bold text-[#202224] dark:text-white break-all">
                {{ statusData.google_account_email || 'Conta Google vinculada' }}
              </div>
              <div class="text-[11px] text-gray-600 dark:text-zinc-400">
                Última sincronização:
                <strong>{{ statusData.last_sync_at ? formatDateTime(statusData.last_sync_at) : 'Nenhuma ainda' }}</strong>
              </div>
            </div>

            <div v-if="statusData.last_error" class="p-3 rounded-xl bg-amber-500/10 border border-amber-500/20 text-xs text-amber-700 dark:text-amber-300">
              Último aviso: {{ statusData.last_error }}
            </div>

            <div v-if="!confirmDisconnect" class="pt-1">
              <button
                type="button"
                @click="confirmDisconnect = true"
                :disabled="disconnecting"
                class="py-2.5 px-4 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 text-rose-600 dark:text-rose-400 border border-rose-500/30 font-bold text-xs transition inline-flex items-center gap-2 cursor-pointer"
              >
                <Unlink class="w-4 h-4" />
                <span>Desconectar conta</span>
              </button>
            </div>
            <div v-else class="p-3.5 rounded-xl bg-rose-500/10 border border-rose-500/30 space-y-2.5">
              <p class="text-xs text-rose-700 dark:text-rose-300 font-medium">
                Confirmar desconexão? A sincronização automática será pausada, mantendo suas credenciais OAuth salvas.
              </p>
              <div class="flex items-center gap-2">
                <button
                  type="button"
                  @click="disconnectAccount"
                  :disabled="disconnecting"
                  class="py-1.5 px-3.5 rounded-lg bg-rose-600 hover:bg-rose-700 text-white text-xs font-bold transition inline-flex items-center gap-1.5 cursor-pointer"
                >
                  <RefreshCw v-if="disconnecting" class="w-3.5 h-3.5 animate-spin" />
                  <span>Confirmar desconexão</span>
                </button>
                <button
                  type="button"
                  @click="confirmDisconnect = false"
                  class="py-1.5 px-3 rounded-lg bg-white dark:bg-zinc-800 border border-gray-200 dark:border-zinc-700 text-xs font-semibold text-gray-700 dark:text-zinc-300 cursor-pointer"
                >
                  Cancelar
                </button>
              </div>
            </div>
          </div>

          <div v-else class="space-y-4">
            <div class="p-4 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 space-y-2">
              <p class="text-xs text-gray-600 dark:text-zinc-400 leading-relaxed">
                <template v-if="statusData.status === 'EXPIRED'">
                  A autorização da conta Google expirou ou foi revogada. Reconecte para retomar a sincronização automática.
                </template>
                <template v-else-if="statusData.status === 'ERROR'">
                  Ocorreu um erro na conexão anterior: {{ statusData.last_error || 'tente conectar novamente.' }}
                </template>
                <template v-else>
                  Salve as credenciais OAuth acima e autorize o acesso ao calendário do seu estabelecimento.
                </template>
              </p>
            </div>

            <button
              type="button"
              @click="connectGoogle"
              :disabled="!statusData.has_credentials || connecting"
              class="w-full py-3 px-5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] disabled:opacity-50 disabled:cursor-not-allowed text-white font-bold text-sm shadow-md shadow-blue-500/25 transition inline-flex items-center justify-center gap-2 cursor-pointer"
            >
              <RefreshCw v-if="connecting" class="w-4 h-4 animate-spin" />
              <ExternalLink v-else class="w-4 h-4" />
              <span>{{ statusData.status === 'EXPIRED' || statusData.status === 'ERROR' ? 'Reconectar conta' : 'Conectar com Google' }}</span>
            </button>
            <p v-if="!statusData.has_credentials" class="text-[11px] text-amber-600 dark:text-amber-400">
              Configure e salve o Client ID e Client Secret antes de conectar.
            </p>
          </div>
        </div>

        <!-- Bloco 3: Configurações de Sincronização -->
        <div class="lg:col-span-7 saas-card p-6 sm:p-8 space-y-5">
          <h2 class="text-base font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
            <Sliders class="w-5 h-5 text-[#4880FF]" />
            3. Configurações de Sincronização
          </h2>

          <div class="p-4 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 flex items-center justify-between gap-4">
            <div>
              <div class="text-sm font-bold text-[#202224] dark:text-white">
                Habilitar sincronização automática de agendamentos
              </div>
              <p class="text-xs text-[#718096] dark:text-zinc-400 mt-0.5">
                Exige que uma conta Google esteja conectada.
              </p>
            </div>

            <label class="relative inline-flex items-center cursor-pointer shrink-0">
              <input
                type="checkbox"
                v-model="syncForm.is_enabled"
                :disabled="statusData.status !== 'CONNECTED' || savingSyncSettings"
                @change="saveSyncSettings"
                class="sr-only peer"
              />
              <div class="w-12 h-6 bg-gray-200 dark:bg-zinc-800 peer-disabled:opacity-50 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#4880FF]"></div>
            </label>
          </div>

          <div class="space-y-2">
            <div class="flex items-center justify-between gap-2">
              <label class="block text-xs font-bold text-gray-700 dark:text-zinc-300">
                Calendário de Destino
              </label>
              <button
                type="button"
                @click="loadCalendars"
                :disabled="statusData.status !== 'CONNECTED' || loadingCalendars"
                class="inline-flex items-center gap-1.5 text-xs font-semibold text-[#4880FF] hover:text-[#386FF0] disabled:opacity-50 cursor-pointer"
              >
                <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingCalendars }" />
                <span>Atualizar lista</span>
              </button>
            </div>

            <div class="flex flex-col sm:flex-row gap-3">
              <select
                v-model="syncForm.calendar_id"
                :disabled="statusData.status !== 'CONNECTED'"
                class="flex-1 bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 rounded-xl px-3.5 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:ring-2 focus:ring-[#4880FF]/40 disabled:opacity-50"
              >
                <option value="primary">Agenda Principal (primary)</option>
                <option
                  v-for="cal in calendars"
                  :key="cal.id"
                  :value="cal.id"
                >
                  {{ cal.summary }}{{ cal.primary ? ' (Principal)' : '' }}
                </option>
              </select>

              <button
                type="button"
                @click="saveSyncSettings"
                :disabled="statusData.status !== 'CONNECTED' || savingSyncSettings"
                class="py-2.5 px-5 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] disabled:opacity-50 text-white font-bold text-sm transition inline-flex items-center justify-center gap-2 cursor-pointer shrink-0"
              >
                <RefreshCw v-if="savingSyncSettings" class="w-4 h-4 animate-spin" />
                <Save v-else class="w-4 h-4" />
                <span>Salvar</span>
              </button>
            </div>
          </div>

          <div class="p-3.5 rounded-xl bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900/40 text-blue-900 dark:text-blue-300 text-xs leading-relaxed">
            <strong>Como funciona:</strong> Agendamentos criados ou reagendados (pelo site público, painel administrativo e atendimento IA no WhatsApp) são sincronizados automaticamente no calendário escolhido. Cancelamentos removem o evento correspondente do Google Calendar.
          </div>
        </div>
      </div>

      <!-- ============================================================== -->
      <!-- BLOCO 4: FALHAS DE SINCRONIZAÇÃO                                -->
      <!-- ============================================================== -->
      <div class="saas-card p-6 sm:p-8 space-y-4">
        <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-3">
          <div>
            <h2 class="text-base font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
              <AlertTriangle class="w-5 h-5 text-amber-500" />
              4. Falhas de Sincronização
              <span
                v-if="failures.length > 0"
                class="text-xs px-2 py-0.5 rounded-full bg-rose-500/10 text-rose-600 dark:text-rose-400 border border-rose-500/20 font-bold"
              >
                {{ failures.length }}
              </span>
            </h2>
            <p class="text-xs text-[#718096] dark:text-zinc-400 mt-0.5">
              Agendamentos que falharam na sincronização com o Google Calendar e podem ser reprocessados manualmente.
            </p>
          </div>

          <button
            type="button"
            @click="loadFailures"
            :disabled="loadingFailures"
            class="py-2 px-3.5 rounded-xl bg-gray-100 dark:bg-zinc-800 hover:bg-gray-200 dark:hover:bg-zinc-700 text-xs font-bold text-gray-700 dark:text-zinc-200 inline-flex items-center gap-1.5 self-start sm:self-auto cursor-pointer"
          >
            <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingFailures }" />
            <span>Atualizar falhas</span>
          </button>
        </div>

        <div v-if="failures.length === 0" class="p-6 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 text-center text-xs text-gray-500 dark:text-zinc-400">
          Nenhuma falha de sincronização pendente no momento.
        </div>

        <div v-else class="overflow-x-auto">
          <table class="w-full text-left border-collapse text-xs">
            <thead>
              <tr class="border-b border-gray-200 dark:border-zinc-800 text-gray-500 dark:text-zinc-400 uppercase tracking-wider text-[11px]">
                <th class="py-3 pr-4 font-bold">Cliente</th>
                <th class="py-3 pr-4 font-bold">Serviço / Profissional</th>
                <th class="py-3 pr-4 font-bold">Data/Hora</th>
                <th class="py-3 pr-4 font-bold">Tentativas</th>
                <th class="py-3 pr-4 font-bold">Erro</th>
                <th class="py-3 text-right font-bold">Ação</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-100 dark:divide-zinc-800/80">
              <tr v-for="item in failures" :key="item.appointment_id">
                <td class="py-3.5 pr-4 font-semibold text-[#202224] dark:text-white">
                  {{ item.customer_name || '-' }}
                </td>
                <td class="py-3.5 pr-4 text-gray-600 dark:text-zinc-300">
                  <div>{{ item.service_name || '-' }}</div>
                  <div class="text-[11px] text-gray-400 dark:text-zinc-500">{{ item.professional_name || '-' }}</div>
                </td>
                <td class="py-3.5 pr-4 text-gray-600 dark:text-zinc-300 whitespace-nowrap">
                  {{ item.start_at ? formatDateTime(item.start_at) : '-' }}
                </td>
                <td class="py-3.5 pr-4 font-mono text-gray-600 dark:text-zinc-300">
                  {{ item.attempt_count }}
                </td>
                <td class="py-3.5 pr-4 text-rose-600 dark:text-rose-400 max-w-xs truncate" :title="item.last_error">
                  {{ item.last_error }}
                </td>
                <td class="py-3.5 text-right whitespace-nowrap">
                  <button
                    type="button"
                    @click="retrySync(item.appointment_id)"
                    :disabled="retryingId === item.appointment_id"
                    class="py-1.5 px-3 rounded-lg bg-[#4880FF] hover:bg-[#386FF0] disabled:opacity-50 text-white font-bold text-[11px] inline-flex items-center gap-1.5 cursor-pointer"
                  >
                    <RefreshCw class="w-3 h-3" :class="{ 'animate-spin': retryingId === item.appointment_id }" />
                    <span>Tentar sincronizar novamente</span>
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '../../services/api'
import {
  CalendarCheck2,
  CheckCircle2,
  AlertCircle,
  RefreshCw,
  Key,
  Copy,
  Check,
  Save,
  ShieldCheck,
  Link2,
  Unlink,
  ExternalLink,
  Sliders,
  AlertTriangle,
} from 'lucide-vue-next'

interface CalendarItem {
  id: string
  summary: string
  primary: boolean
  access_role: string
}

interface FailureItem {
  appointment_id: string
  customer_name: string
  service_name: string
  professional_name: string
  start_at: string
  sync_status: string
  last_action: string
  attempt_count: number
  last_error: string
  updated_at: string
}

const route = useRoute()
const router = useRouter()

const loadingStatus = ref(true)
const savingCredentials = ref(false)
const savingSyncSettings = ref(false)
const connecting = ref(false)
const disconnecting = ref(false)
const confirmDisconnect = ref(false)
const loadingCalendars = ref(false)
const loadingFailures = ref(false)
const retryingId = ref<string | null>(null)
const copiedRedirect = ref(false)

const successMessage = ref('')
const errorMessage = ref('')

const statusData = reactive({
  google_client_id: '',
  has_client_secret: false,
  has_credentials: false,
  redirect_uri: '',
  is_enabled: false,
  status: 'DISCONNECTED',
  google_account_email: '',
  calendar_id: 'primary',
  calendar_summary: '',
  last_sync_at: '' as string | null,
  last_error: '',
  failed_syncs_count: 0,
})

const credentialsForm = reactive({
  google_client_id: '',
  google_client_secret: '',
})

const syncForm = reactive({
  is_enabled: false,
  calendar_id: 'primary',
})

const calendars = ref<CalendarItem[]>([])
const failures = ref<FailureItem[]>([])

const effectiveRedirectURI = computed(() => {
  if (statusData.redirect_uri) return statusData.redirect_uri
  return `${window.location.origin}/api/v1/integrations/google-calendar/callback`
})

function applyStatusResponse(data: any) {
  if (!data) return
  statusData.google_client_id = data.google_client_id || ''
  statusData.has_client_secret = !!data.has_client_secret
  statusData.has_credentials = !!data.has_credentials
  statusData.redirect_uri = data.redirect_uri || ''
  statusData.is_enabled = !!data.is_enabled
  statusData.status = data.status || 'DISCONNECTED'
  statusData.google_account_email = data.google_account_email || ''
  statusData.calendar_id = data.calendar_id || 'primary'
  statusData.calendar_summary = data.calendar_summary || ''
  statusData.last_sync_at = data.last_sync_at || null
  statusData.last_error = data.last_error || ''
  statusData.failed_syncs_count = data.failed_syncs_count || 0

  credentialsForm.google_client_id = statusData.google_client_id
  syncForm.is_enabled = statusData.is_enabled
  syncForm.calendar_id = statusData.calendar_id
}

function handleOAuthQueryParams() {
  const oauthStatus = route.query.status as string | undefined
  const reason = route.query.reason as string | undefined
  if (!oauthStatus) return

  if (oauthStatus === 'connected') {
    successMessage.value = 'Conta Google conectada com sucesso! Selecione o calendário e ative a sincronização automática.'
  } else if (oauthStatus === 'error') {
    const reasonMessages: Record<string, string> = {
      access_denied: 'Autorização cancelada na tela de consentimento do Google.',
      invalid_state: 'Sessão de autorização inválida ou expirada. Tente conectar novamente.',
      not_configured: 'Credenciais OAuth não encontradas. Salve o Client ID e Secret primeiro.',
      oauth_exchange_failed: 'Não foi possível concluir a autenticação com o Google. Verifique Client ID, Secret e Redirect URI.',
    }
    errorMessage.value = reasonMessages[reason || ''] || 'Falha ao conectar conta do Google Calendar.'
  }

  router.replace({ query: {} })
}

async function loadStatus() {
  loadingStatus.value = true
  try {
    const res = await api.get('/admin/integrations/google-calendar/status')
    applyStatusResponse(res.data?.data)
    if (statusData.status === 'CONNECTED') {
      await loadCalendars()
    }
    await loadFailures()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao carregar status do Google Calendar.'
  } finally {
    loadingStatus.value = false
  }
}

async function loadCalendars() {
  if (statusData.status !== 'CONNECTED') return
  loadingCalendars.value = true
  try {
    const res = await api.get('/admin/integrations/google-calendar/calendars')
    calendars.value = res.data?.data || []
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao carregar calendários da conta Google.'
  } finally {
    loadingCalendars.value = false
  }
}

async function loadFailures() {
  loadingFailures.value = true
  try {
    const res = await api.get('/admin/integrations/google-calendar/failures')
    failures.value = res.data?.data || []
  } catch {
    // Silencioso no carregamento inicial caso não haja falhas
  } finally {
    loadingFailures.value = false
  }
}

async function copyRedirectURI() {
  try {
    await navigator.clipboard.writeText(effectiveRedirectURI.value)
    copiedRedirect.value = true
    setTimeout(() => {
      copiedRedirect.value = false
    }, 2000)
  } catch {
    errorMessage.value = 'Não foi possível copiar automaticamente. Copie o endereço manualmente.'
  }
}

async function saveCredentials() {
  errorMessage.value = ''
  successMessage.value = ''
  savingCredentials.value = true
  try {
    const payload: Record<string, any> = {
      google_client_id: credentialsForm.google_client_id.trim(),
    }
    if (credentialsForm.google_client_secret.trim() !== '') {
      payload.google_client_secret = credentialsForm.google_client_secret.trim()
    }
    const res = await api.put('/admin/integrations/google-calendar/settings', payload)
    applyStatusResponse(res.data?.data)
    credentialsForm.google_client_secret = ''
    successMessage.value = 'Credenciais OAuth salvas com sucesso!'
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao salvar credenciais OAuth.'
  } finally {
    savingCredentials.value = false
  }
}

async function connectGoogle() {
  errorMessage.value = ''
  connecting.value = true
  try {
    const res = await api.get('/admin/integrations/google-calendar/auth-url')
    const authUrl = res.data?.data?.auth_url
    if (!authUrl) {
      throw new Error('URL de autorização não retornada pela API.')
    }
    window.location.href = authUrl
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || err.message || 'Erro ao iniciar autenticação com Google.'
    connecting.value = false
  }
}

async function disconnectAccount() {
  errorMessage.value = ''
  successMessage.value = ''
  disconnecting.value = true
  try {
    await api.post('/admin/integrations/google-calendar/disconnect')
    confirmDisconnect.value = false
    calendars.value = []
    successMessage.value = 'Conta Google desconectada com sucesso.'
    await loadStatus()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao desconectar conta Google.'
  } finally {
    disconnecting.value = false
  }
}

async function saveSyncSettings() {
  errorMessage.value = ''
  successMessage.value = ''
  savingSyncSettings.value = true
  try {
    const selectedCal = calendars.value.find(c => c.id === syncForm.calendar_id)
    const summary = selectedCal ? selectedCal.summary : (syncForm.calendar_id === 'primary' ? 'Agenda Principal' : '')
    const res = await api.put('/admin/integrations/google-calendar/settings', {
      is_enabled: syncForm.is_enabled,
      calendar_id: syncForm.calendar_id,
      calendar_summary: summary,
    })
    applyStatusResponse(res.data?.data)
    successMessage.value = 'Configurações de sincronização salvas com sucesso!'
  } catch (err: any) {
    syncForm.is_enabled = statusData.is_enabled
    errorMessage.value = err.response?.data?.error || 'Erro ao salvar configurações de sincronização.'
  } finally {
    savingSyncSettings.value = false
  }
}

async function retrySync(appointmentId: string) {
  errorMessage.value = ''
  successMessage.value = ''
  retryingId.value = appointmentId
  try {
    await api.post(`/admin/integrations/google-calendar/syncs/${appointmentId}/retry`)
    successMessage.value = 'Sincronização reenfileirada! Aguarde alguns instantes e atualize a lista.'
    await loadFailures()
  } catch (err: any) {
    errorMessage.value = err.response?.data?.error || 'Erro ao reenfileirar sincronização.'
  } finally {
    retryingId.value = null
  }
}

function formatDateTime(iso: string): string {
  try {
    return new Intl.DateTimeFormat('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(new Date(iso))
  } catch {
    return iso
  }
}

onMounted(async () => {
  handleOAuthQueryParams()
  await loadStatus()
})
</script>
