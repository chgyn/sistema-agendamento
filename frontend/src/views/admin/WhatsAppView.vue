<template>
  <div class="space-y-6 max-w-7xl mx-auto">
    <!-- Cabeçalho -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl sm:text-3xl font-black text-[#202224] dark:text-white tracking-tight flex items-center gap-3 font-['Outfit']">
          <div class="w-10 h-10 rounded-xl bg-blue-50 dark:bg-blue-950/60 border border-blue-200 dark:border-blue-900/40 flex items-center justify-center text-[#4880FF] shadow-sm">
            <Bot class="w-6 h-6" />
          </div>
          <span>WhatsApp & Atendimento com IA</span>
        </h1>
        <p class="text-xs sm:text-sm text-[#718096] dark:text-surface-400 mt-1">
          Conecte o número de WhatsApp do estabelecimento e configure o agente inteligente para conduzir agendamentos.
        </p>
      </div>

      <!-- Badge de Status Geral -->
      <div class="flex items-center gap-2 self-start sm:self-auto">
        <span
          v-if="statusData.status === 'CONNECTED'"
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 border border-emerald-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-500 dark:bg-emerald-400 animate-pulse"></span>
          WhatsApp Conectado
        </span>
        <span
          v-else-if="statusData.status === 'QRCODE'"
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-amber-500/10 text-amber-600 dark:text-amber-400 border border-amber-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-amber-500 dark:bg-amber-400 animate-ping"></span>
          Aguardando Leitura do QR Code
        </span>
        <span
          v-else
          class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full text-xs font-bold bg-gray-100 dark:bg-surface-800 text-gray-500 dark:text-zinc-400 border border-gray-200 dark:border-zinc-750"
        >
          <span class="w-2 h-2 rounded-full bg-gray-400 dark:bg-zinc-500"></span>
          WhatsApp Desconectado
        </span>
      </div>
    </div>

    <!-- Navegação por Abas -->
    <div class="flex border-b border-gray-200 dark:border-zinc-800 gap-6">
      <button
        @click="activeTab = 'connection'"
        class="pb-3.5 px-2 text-sm font-bold transition border-b-2 flex items-center gap-2 cursor-pointer"
        :class="activeTab === 'connection' ? 'border-[#4880FF] text-[#4880FF] font-black' : 'border-transparent text-gray-500 dark:text-zinc-400 hover:text-gray-800 dark:hover:text-zinc-200'"
      >
        <Smartphone class="w-4 h-4" />
        <span>1. Conexão do WhatsApp</span>
      </button>

      <button
        @click="activeTab = 'ai'"
        class="pb-3.5 px-2 text-sm font-bold transition border-b-2 flex items-center gap-2 cursor-pointer"
        :class="activeTab === 'ai' ? 'border-[#4880FF] text-[#4880FF] font-black' : 'border-transparent text-gray-500 dark:text-zinc-400 hover:text-gray-800 dark:hover:text-zinc-200'"
      >
        <Sparkles class="w-4 h-4" />
        <span>2. Configuração do Atendente com IA</span>
      </button>
    </div>

    <!-- Mensagens de Alerta Global -->
    <div v-if="successMessage" class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-600 dark:text-emerald-300 text-sm flex items-center justify-between">
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-5 h-5 shrink-0 text-emerald-500 dark:text-emerald-400" />
        <span>{{ successMessage }}</span>
      </div>
      <button @click="successMessage = ''" class="text-xs text-emerald-500 hover:text-emerald-700 dark:text-emerald-400 dark:hover:text-white">✕</button>
    </div>

    <div v-if="errorMessage" class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-600 dark:text-rose-300 text-sm flex items-center justify-between">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-5 h-5 shrink-0 text-rose-500 dark:text-rose-400" />
        <span>{{ errorMessage }}</span>
      </div>
      <button @click="errorMessage = ''" class="text-xs text-rose-500 hover:text-rose-700 dark:text-rose-400 dark:hover:text-white">✕</button>
    </div>

    <!-- ============================================================== -->
    <!-- ABA 1: CONEXÃO DO WHATSAPP (WUZAPI)                           -->
    <!-- ============================================================== -->
    <div v-if="activeTab === 'connection'" class="space-y-6">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Card Principal de Pareamento -->
        <div class="lg:col-span-7 saas-card p-6 sm:p-8 space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
              <QrCode class="w-5 h-5 text-[#4880FF]" />
              Pareamento do Dispositivo
            </h2>
            <button
              @click="fetchStatus"
              :disabled="loadingStatus"
              class="px-3 py-1.5 rounded-xl bg-white dark:bg-[#121318] hover:bg-gray-50 dark:hover:bg-zinc-800 text-gray-700 dark:text-zinc-300 transition text-xs font-bold flex items-center gap-1.5 border border-gray-200 dark:border-zinc-700/60 cursor-pointer shadow-sm"
              title="Atualizar Status"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingStatus }" />
              <span>Sincronizar</span>
            </button>
          </div>

          <!-- ESTADO 1: CONECTADO -->
          <div v-if="statusData.status === 'CONNECTED'" class="p-6 rounded-2xl bg-emerald-50 dark:bg-emerald-950/20 border border-emerald-200 dark:border-emerald-500/30 space-y-4">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-2xl bg-emerald-500/20 border border-emerald-500/40 flex items-center justify-center text-emerald-600 dark:text-emerald-400 shrink-0">
                <CheckCircle2 class="w-7 h-7" />
              </div>
              <div>
                <h3 class="font-bold text-[#202224] dark:text-white text-base font-['Outfit']">WhatsApp Conectado com Sucesso!</h3>
                <p class="text-xs text-gray-600 dark:text-zinc-400 mt-0.5">
                  Número Vinculado: <span class="text-emerald-600 dark:text-emerald-400 font-mono font-bold">{{ statusData.phone_number || 'Conectado via Web' }}</span>
                </p>
                <p class="text-[11px] text-gray-400 dark:text-zinc-500 mt-0.5">Instância: {{ statusData.instance_name }}</p>
              </div>
            </div>

            <div class="pt-4 border-t border-emerald-200/80 dark:border-zinc-800 flex flex-wrap gap-3">
              <button
                @click="disconnect"
                :disabled="loadingAction"
                class="py-2.5 px-4 rounded-xl bg-white hover:bg-gray-50 text-gray-700 dark:bg-[#181922] dark:hover:bg-zinc-800 dark:text-zinc-300 text-xs font-bold transition flex items-center gap-2 border border-gray-200 dark:border-zinc-700/60 cursor-pointer shadow-sm"
              >
                <Smartphone class="w-3.5 h-3.5" />
                <span>Pausar Conexão</span>
              </button>

              <button
                @click="logout"
                :disabled="loadingAction"
                class="py-2.5 px-4 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 border border-rose-500/30 text-rose-600 dark:text-rose-400 text-xs font-bold transition flex items-center gap-2 cursor-pointer"
              >
                <LogOut class="w-3.5 h-3.5" />
                <span>Desconectar e Limpar Sessão</span>
              </button>
            </div>
          </div>

          <!-- ESTADO 2: EXIBINDO QR CODE -->
          <div v-else-if="statusData.status === 'QRCODE' && qrCodeImage" class="flex flex-col items-center p-6 rounded-2xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 space-y-4 text-center">
            <div class="p-3 bg-white rounded-2xl shadow-xl shadow-black/10 dark:shadow-black/60 border border-gray-100 dark:border-transparent">
              <img :src="qrCodeImage" alt="QR Code WhatsApp" class="w-64 h-64 object-contain rounded-lg" />
            </div>

            <div class="space-y-1">
              <p class="text-sm font-bold text-[#202224] dark:text-white font-['Outfit']">Escaneie o código com o seu WhatsApp</p>
              <p class="text-xs text-gray-500 dark:text-zinc-400">O código é atualizado automaticamente a cada 60 segundos.</p>
            </div>

            <button
              @click="getQRCode"
              :disabled="loadingAction"
              class="py-2.5 px-4 rounded-xl bg-blue-50 hover:bg-blue-100 text-[#4880FF] dark:bg-blue-950/40 dark:hover:bg-blue-900/50 dark:text-blue-300 border border-blue-200 dark:border-blue-900/40 text-xs font-bold transition flex items-center gap-2 cursor-pointer"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingAction }" />
              <span>Gerar Novo QR Code</span>
            </button>
          </div>

          <!-- ESTADO 3: DESCONECTADO OU AGUARDANDO INÍCIO -->
          <div v-else class="p-8 rounded-2xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800 text-center space-y-4">
            <div class="w-14 h-14 rounded-2xl bg-blue-50 dark:bg-blue-950/60 border border-blue-200 dark:border-blue-900/40 mx-auto flex items-center justify-center text-[#4880FF]">
              <Smartphone class="w-7 h-7" />
            </div>
            <div class="max-w-md mx-auto space-y-1">
              <h3 class="font-bold text-[#202224] dark:text-white text-base font-['Outfit']">Nenhum WhatsApp Conectado</h3>
              <p class="text-xs text-gray-500 dark:text-zinc-400">
                Inicie uma sessão para gerar o QR Code e vincular o WhatsApp oficial da sua barbearia ou salão.
              </p>
            </div>
            <button
              @click="connect"
              :disabled="loadingAction"
              class="py-3 px-6 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white font-bold text-sm shadow-md shadow-blue-500/25 transition flex items-center gap-2 mx-auto cursor-pointer"
            >
              <QrCode class="w-4 h-4" />
              <span>Conectar WhatsApp (Gerar QR Code)</span>
            </button>
          </div>
        </div>

        <!-- Instruções Passo a Passo -->
        <div class="lg:col-span-5 saas-card p-6 sm:p-8 space-y-4">
          <h3 class="text-xs font-bold text-gray-500 dark:text-zinc-400 uppercase tracking-wider flex items-center gap-2">
            <ShieldCheck class="w-4 h-4 text-[#4880FF]" />
            Como Conectar seu Aparelho
          </h3>

          <ol class="space-y-3 text-xs text-gray-700 dark:text-zinc-300">
            <li class="flex items-start gap-3 p-3.5 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">1</span>
              <span>Abra o aplicativo do <strong>WhatsApp</strong> no seu celular principal.</span>
            </li>
            <li class="flex items-start gap-3 p-3.5 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">2</span>
              <span>Acesse as configurações tocando nos <strong>três pontos (Android)</strong> ou <strong>Configurações (iPhone)</strong>.</span>
            </li>
            <li class="flex items-start gap-3 p-3.5 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">3</span>
              <span>Toque em <strong>Aparelhos Conectados</strong> e depois em <strong>Conectar Aparelho</strong>.</span>
            </li>
            <li class="flex items-start gap-3 p-3.5 rounded-xl bg-gray-50 dark:bg-[#14151c] border border-gray-200 dark:border-zinc-800">
              <span class="w-5 h-5 rounded-full bg-blue-50 dark:bg-blue-950/60 text-[#4880FF] font-bold flex items-center justify-center shrink-0">4</span>
              <span>Aponte a câmera para o <strong>QR Code</strong> exibido nesta tela até a confirmação.</span>
            </li>
          </ol>

          <div class="p-3.5 rounded-xl bg-blue-50 dark:bg-blue-950/30 border border-blue-200 dark:border-blue-900/40 text-blue-900 dark:text-blue-300 text-[11px] leading-relaxed">
            💡 <strong>Dica de Segurança:</strong> O WUZAPI funciona como um aparelho web conectado oficial, permitindo responder seus clientes sem intermediários e com total privacidade.
          </div>
        </div>
      </div>
    </div>

    <!-- ============================================================== -->
    <!-- ABA 2: CONFIGURAÇÃO DO ATENDENTE COM IA                        -->
    <!-- ============================================================== -->
    <div v-else-if="activeTab === 'ai'" class="space-y-6">
      <form @submit.prevent="saveAIConfig" class="space-y-6">
        <!-- Switch Geral de Ativação -->
        <div class="saas-card p-6 sm:p-8 flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h3 class="text-base font-bold text-[#202224] dark:text-white flex items-center gap-2 font-['Outfit']">
              <Bot class="w-5 h-5 text-[#4880FF]" />
              Atendimento Automatizado com IA
            </h3>
            <p class="text-xs text-[#718096] dark:text-surface-400 mt-1">
              Quando ativado, a IA responde automaticamente as mensagens recebidas, tirando dúvidas e realizando agendamentos na sua agenda.
            </p>
          </div>

          <label class="relative inline-flex items-center cursor-pointer shrink-0">
            <input type="checkbox" v-model="aiForm.is_ai_enabled" class="sr-only peer" />
            <div class="w-12 h-6 bg-gray-200 dark:bg-surface-800 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 dark:after:border-surface-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-[#4880FF]"></div>
          </label>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <!-- Coluna 1: Provedor, Chaves e Modelo -->
          <div class="lg:col-span-7 space-y-6">
            <div class="saas-card p-6 sm:p-8 space-y-5">
              <h3 class="text-xs font-bold text-gray-500 dark:text-zinc-400 uppercase tracking-wider flex items-center gap-2">
                <Key class="w-4 h-4 text-[#4880FF]" />
                Provedor de Inteligência Artificial
              </h3>

              <!-- Seleção do Provedor -->
              <div class="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  @click="aiForm.ai_provider = 'GEMINI'"
                  class="p-4 rounded-xl border text-left transition relative cursor-pointer"
                  :class="aiForm.ai_provider === 'GEMINI' ? 'bg-blue-50/50 border-[#4880FF] text-[#202224] dark:bg-blue-950/30 dark:border-blue-500 dark:text-white font-bold ring-1 ring-[#4880FF]' : 'bg-gray-50 dark:bg-[#14151c] border-gray-200 dark:border-zinc-800 text-gray-600 dark:text-zinc-400 hover:border-gray-300 dark:hover:border-zinc-700'"
                >
                  <div class="text-sm font-bold font-['Outfit']">Google Gemini</div>
                  <div class="text-[11px] text-gray-500 dark:text-zinc-400 mt-1">Recomendado para custo-benefício e velocidade</div>
                  <span v-if="aiConfigData.has_gemini_api_key" class="inline-block mt-2 text-[10px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-300 px-2 py-0.5 rounded-md font-semibold border border-emerald-500/20">
                    ✓ Chave configurada
                  </span>
                </button>

                <button
                  type="button"
                  @click="aiForm.ai_provider = 'OPENAI'"
                  class="p-4 rounded-xl border text-left transition relative cursor-pointer"
                  :class="aiForm.ai_provider === 'OPENAI' ? 'bg-blue-50/50 border-[#4880FF] text-[#202224] dark:bg-blue-950/30 dark:border-blue-500 dark:text-white font-bold ring-1 ring-[#4880FF]' : 'bg-gray-50 dark:bg-[#14151c] border-gray-200 dark:border-zinc-800 text-gray-600 dark:text-zinc-400 hover:border-gray-300 dark:hover:border-zinc-700'"
                >
                  <div class="text-sm font-bold font-['Outfit']">OpenAI (ChatGPT)</div>
                  <div class="text-[11px] text-gray-500 dark:text-zinc-400 mt-1">Modelos GPT-4o e GPT-4o mini</div>
                  <span v-if="aiConfigData.has_openai_api_key" class="inline-block mt-2 text-[10px] bg-emerald-500/10 text-emerald-600 dark:text-emerald-300 px-2 py-0.5 rounded-md font-semibold border border-emerald-500/20">
                    ✓ Chave configurada
                  </span>
                </button>
              </div>

              <!-- Input da API Key do Gemini -->
              <div v-if="aiForm.ai_provider === 'GEMINI'" class="space-y-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider">
                  Google Gemini API Key
                  <span v-if="aiConfigData.has_gemini_api_key" class="text-emerald-600 dark:text-emerald-400 text-[11px] ml-2 font-normal">
                    (Já configurada com segurança. Preencha apenas se desejar substituir)
                  </span>
                </label>
                <div class="relative">
                  <input
                    type="password"
                    v-model="aiForm.gemini_api_key"
                    :placeholder="aiConfigData.has_gemini_api_key ? '••••••••••••••••••••••••••••••••' : 'Cole sua Gemini API Key aqui (AIzaSy...)'"
                    class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl px-4 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:ring-2 focus:ring-blue-500/20 focus:border-[#4880FF] transition font-mono"
                  />
                </div>
                <p class="text-[11px] text-gray-500 dark:text-zinc-500">
                  Obtenha gratuitamente em <a href="https://aistudio.google.com/app/apikey" target="_blank" class="text-[#4880FF] underline">Google AI Studio</a>.
                </p>
              </div>

              <!-- Input da API Key da OpenAI -->
              <div v-if="aiForm.ai_provider === 'OPENAI'" class="space-y-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider">
                  OpenAI API Key
                  <span v-if="aiConfigData.has_openai_api_key" class="text-emerald-600 dark:text-emerald-400 text-[11px] ml-2 font-normal">
                    (Já configurada com segurança. Preencha apenas se desejar substituir)
                  </span>
                </label>
                <div class="relative">
                  <input
                    type="password"
                    v-model="aiForm.openai_api_key"
                    :placeholder="aiConfigData.has_openai_api_key ? '••••••••••••••••••••••••••••••••' : 'Cole sua OpenAI API Key aqui (sk-...)'"
                    class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl px-4 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:ring-2 focus:ring-blue-500/20 focus:border-[#4880FF] transition font-mono"
                  />
                </div>
                <p class="text-[11px] text-gray-500 dark:text-zinc-500">
                  Obtenha no painel da <a href="https://platform.openai.com/api-keys" target="_blank" class="text-[#4880FF] underline">OpenAI Platform</a>.
                </p>
              </div>

              <!-- Seleção do Modelo -->
              <div class="space-y-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider">Modelo de Inteligência Artificial</label>
                <select
                  v-model="aiForm.ai_model"
                  class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl px-4 py-2.5 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:ring-2 focus:ring-blue-500/20 focus:border-[#4880FF] transition"
                >
                  <template v-if="aiForm.ai_provider === 'GEMINI'">
                    <option value="gemini-2.5-flash">Gemini 2.5 Flash (Ultrarrápido & Mais recente)</option>
                    <option value="gemini-1.5-flash">Gemini 1.5 Flash (Econômico & Estável)</option>
                    <option value="gemini-1.5-pro">Gemini 1.5 Pro (Raciocínio complexo)</option>
                  </template>
                  <template v-else>
                    <option value="gpt-4o-mini">GPT-4o Mini (Rápido, preciso e econômico)</option>
                    <option value="gpt-4o">GPT-4o (Máxima capacidade)</option>
                  </template>
                </select>
              </div>

              <!-- Instruções Personalizadas / System Prompt -->
              <div class="space-y-2">
                <label class="block text-xs font-bold text-gray-700 dark:text-surface-300 uppercase tracking-wider">
                  Instruções Adicionais do Estabelecimento (Prompt)
                </label>
                <textarea
                  v-model="aiForm.system_prompt_custom"
                  rows="4"
                  placeholder="Ex: Não atendemos sem agendamento prévio. Temos estacionamento conveniado gratuito na rua lateral. Em caso de atraso superior a 15 minutos, o horário poderá ser remarcado."
                  class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl p-3 text-sm text-[#202224] dark:text-white focus:outline-none focus:bg-white dark:focus:bg-surface-950 focus:ring-2 focus:ring-blue-500/20 focus:border-[#4880FF] transition leading-relaxed resize-y"
                ></textarea>
                <p class="text-[11px] text-gray-500 dark:text-zinc-500">
                  A IA já sabe automaticamente o nome dos seus serviços, preços, profissionais e horários livres. Use este campo para regras comerciais, tom de voz ou recados.
                </p>
              </div>
            </div>
          </div>

          <!-- Coluna 2: Ritmo Humanizado e Salvar -->
          <div class="lg:col-span-5 space-y-6">
            <div class="saas-card p-6 sm:p-8 space-y-4">
              <h3 class="text-xs font-bold text-gray-500 dark:text-zinc-400 uppercase tracking-wider flex items-center gap-2">
                <Sliders class="w-4 h-4 text-[#4880FF]" />
                Controle de Ritmo Humanizado
              </h3>

              <p class="text-xs text-gray-600 dark:text-zinc-400 leading-relaxed">
                Evite respostas robóticas instantâneas. O sistema simula o comportamento natural de uma pessoa digitando.
              </p>

              <!-- Tempo de espera inicial -->
              <div class="space-y-2 pt-2 border-t border-gray-100 dark:border-zinc-800">
                <div class="flex justify-between text-xs">
                  <span class="text-gray-700 dark:text-zinc-300 font-bold">Pausa antes de responder</span>
                  <span class="text-[#4880FF] font-mono font-bold">{{ aiForm.humanized_min_delay_sec }}s a {{ aiForm.humanized_max_delay_sec }}s</span>
                </div>
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <span class="text-[11px] text-gray-500 dark:text-zinc-500">Mínimo (seg):</span>
                    <input
                      type="number"
                      min="1"
                      max="10"
                      v-model.number="aiForm.humanized_min_delay_sec"
                      class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl p-2 text-xs text-[#202224] dark:text-white"
                    />
                  </div>
                  <div>
                    <span class="text-[11px] text-gray-500 dark:text-zinc-500">Máximo (seg):</span>
                    <input
                      type="number"
                      min="1"
                      max="15"
                      v-model.number="aiForm.humanized_max_delay_sec"
                      class="w-full bg-gray-50 dark:bg-surface-950/80 border border-gray-200 dark:border-surface-700/60 rounded-xl p-2 text-xs text-[#202224] dark:text-white"
                    />
                  </div>
                </div>
              </div>

              <!-- Velocidade de digitação -->
              <div class="space-y-2 pt-2 border-t border-gray-100 dark:border-zinc-800">
                <div class="flex justify-between text-xs">
                  <span class="text-gray-700 dark:text-zinc-300 font-bold">Velocidade do "Digitando..."</span>
                  <span class="text-[#4880FF] font-mono font-bold">{{ aiForm.typing_speed_chars_sec }} carac/seg</span>
                </div>
                <input
                  type="range"
                  min="15"
                  max="60"
                  v-model.number="aiForm.typing_speed_chars_sec"
                  class="w-full accent-[#4880FF]"
                />
                <p class="text-[10px] text-gray-400 dark:text-zinc-500">
                  Respostas mais longas exibirão "digitando..." por mais tempo (delimitado entre 1.5s e 7.5s).
                </p>
              </div>

              <!-- Janela de agrupamento (Debounce) -->
              <div class="space-y-2 pt-2 border-t border-gray-100 dark:border-zinc-800">
                <div class="flex justify-between text-xs">
                  <span class="text-gray-700 dark:text-zinc-300 font-bold">Janela de agrupamento</span>
                  <span class="text-[#4880FF] font-mono font-bold">{{ aiForm.debounce_window_sec }} segundos</span>
                </div>
                <input
                  type="range"
                  min="2"
                  max="8"
                  v-model.number="aiForm.debounce_window_sec"
                  class="w-full accent-[#4880FF]"
                />
                <p class="text-[10px] text-gray-400 dark:text-zinc-500">
                  Se o cliente mandar várias mensagens curtas seguidas, a IA agrupa todas antes de responder.
                </p>
              </div>
            </div>

            <!-- Botão Salvar -->
            <button
              type="submit"
              :disabled="savingAI"
              class="w-full py-3.5 px-6 rounded-xl bg-[#4880FF] hover:bg-[#386FF0] text-white font-bold text-sm shadow-md shadow-blue-500/25 transition flex items-center justify-center gap-2 cursor-pointer"
            >
              <RefreshCw v-if="savingAI" class="w-4 h-4 animate-spin" />
              <CheckCircle2 v-else class="w-4 h-4" />
              <span>Salvar Configurações de IA</span>
            </button>
          </div>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, watch } from 'vue'
import api from '../../services/api'
import {
  Bot,
  Smartphone,
  Sparkles,
  QrCode,
  CheckCircle2,
  AlertCircle,
  RefreshCw,
  LogOut,
  Key,
  Sliders,
  ShieldCheck,
} from 'lucide-vue-next'

type TabType = 'connection' | 'ai'
const activeTab = ref<TabType>('connection')

const successMessage = ref('')
const errorMessage = ref('')
const loadingStatus = ref(false)
const loadingAction = ref(false)
const savingAI = ref(false)

// Estado da Conexão
const statusData = reactive({
  status: 'DISCONNECTED',
  phone_number: '',
  is_connected: false,
  is_logged_in: false,
  has_qr_code: false,
  qr_code_base64: '',
  instance_name: '',
})

// Estado de Configurações de IA
const aiConfigData = reactive({
  is_ai_enabled: false,
  ai_provider: 'GEMINI',
  ai_model: 'gemini-2.5-flash',
  has_gemini_api_key: false,
  has_openai_api_key: false,
  system_prompt_custom: '',
  humanized_min_delay_sec: 2,
  humanized_max_delay_sec: 5,
  typing_speed_chars_sec: 35,
  debounce_window_sec: 4,
})

// Formulário de Edição de IA
const aiForm = reactive({
  is_ai_enabled: false,
  ai_provider: 'GEMINI',
  ai_model: 'gemini-2.5-flash',
  gemini_api_key: '',
  openai_api_key: '',
  system_prompt_custom: '',
  humanized_min_delay_sec: 2,
  humanized_max_delay_sec: 5,
  typing_speed_chars_sec: 35,
  debounce_window_sec: 4,
})

// QR code data URL
const qrCodeImage = computed(() => {
  if (!statusData.qr_code_base64) return ''
  if (statusData.qr_code_base64.startsWith('data:image')) {
    return statusData.qr_code_base64
  }
  return `data:image/png;base64,${statusData.qr_code_base64}`
})

// Polling interval para QR code
let pollingInterval: any = null

const startPolling = () => {
  stopPolling()
  pollingInterval = setInterval(async () => {
    if (statusData.status === 'QRCODE' || statusData.status === 'CONNECTING') {
      await fetchStatusSilently()
    } else {
      stopPolling()
    }
  }, 3500)
}

const stopPolling = () => {
  if (pollingInterval) {
    clearInterval(pollingInterval)
    pollingInterval = null
  }
}

// Watcher do Provedor de IA para ajustar modelo padrão caso não selecionado
watch(() => aiForm.ai_provider, (newProvider) => {
  if (newProvider === 'GEMINI' && !aiForm.ai_model.startsWith('gemini')) {
    aiForm.ai_model = 'gemini-2.5-flash'
  } else if (newProvider === 'OPENAI' && !aiForm.ai_model.startsWith('gpt')) {
    aiForm.ai_model = 'gpt-4o-mini'
  }
})

// 1. Consulta Status
const fetchStatus = async () => {
  loadingStatus.value = true
  try {
    const res = await api.get('/admin/whatsapp/status')
    if (res.data?.data) {
      Object.assign(statusData, res.data.data)
      if (statusData.status === 'QRCODE' || statusData.status === 'CONNECTING') {
        startPolling()
      } else {
        stopPolling()
      }
    }
  } catch (err: any) {
    console.error('Falha ao consultar status do WhatsApp:', err)
  } finally {
    loadingStatus.value = false
  }
}

const fetchStatusSilently = async () => {
  try {
    const res = await api.get('/admin/whatsapp/status')
    if (res.data?.data) {
      const oldStatus = statusData.status
      Object.assign(statusData, res.data.data)
      if (oldStatus !== 'CONNECTED' && statusData.status === 'CONNECTED') {
        successMessage.value = 'WhatsApp conectado com sucesso!'
        stopPolling()
      }
    }
  } catch (err) {
    // Silencioso no polling
  }
}

// 2. Conectar WhatsApp / Gerar QR
const connect = async () => {
  loadingAction.value = true
  errorMessage.value = ''
  try {
    const res = await api.post('/admin/whatsapp/connect')
    if (res.data?.data?.qr_code_base64) {
      statusData.status = 'QRCODE'
      statusData.qr_code_base64 = res.data.data.qr_code_base64
      startPolling()
    } else if (res.data?.data?.status === 'CONNECTED') {
      statusData.status = 'CONNECTED'
      successMessage.value = 'WhatsApp já estava conectado!'
    } else {
      await fetchStatus()
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Falha ao iniciar conexão com WhatsApp'
  } finally {
    loadingAction.value = false
  }
}

// 3. Obter novo QR Code
const getQRCode = async () => {
  loadingAction.value = true
  try {
    const res = await api.get('/admin/whatsapp/qr')
    if (res.data?.data?.qr_code_base64) {
      statusData.qr_code_base64 = res.data.data.qr_code_base64
      statusData.status = 'QRCODE'
      startPolling()
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Falha ao atualizar QR Code'
  } finally {
    loadingAction.value = false
  }
}

// 4. Desconectar
const disconnect = async () => {
  loadingAction.value = true
  try {
    await api.post('/admin/whatsapp/disconnect')
    statusData.status = 'DISCONNECTED'
    stopPolling()
    successMessage.value = 'WhatsApp desconectado.'
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Erro ao desconectar'
  } finally {
    loadingAction.value = false
  }
}

// 5. Logout
const logout = async () => {
  if (!confirm('Deseja realmente desconectar e remover a sessão do WhatsApp? Será necessário escanear o QR Code novamente.')) {
    return
  }
  loadingAction.value = true
  try {
    await api.post('/admin/whatsapp/logout')
    statusData.status = 'LOGGED_OUT'
    statusData.phone_number = ''
    statusData.qr_code_base64 = ''
    stopPolling()
    successMessage.value = 'Sessão encerrada com sucesso.'
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Erro ao encerrar sessão'
  } finally {
    loadingAction.value = false
  }
}

// 6. Consultar Configuração de IA
const fetchAIConfig = async () => {
  try {
    const res = await api.get('/admin/whatsapp/ai-config')
    if (res.data?.data) {
      Object.assign(aiConfigData, res.data.data)
      aiForm.is_ai_enabled = aiConfigData.is_ai_enabled
      aiForm.ai_provider = aiConfigData.ai_provider
      aiForm.ai_model = aiConfigData.ai_model
      aiForm.system_prompt_custom = aiConfigData.system_prompt_custom
      aiForm.humanized_min_delay_sec = aiConfigData.humanized_min_delay_sec
      aiForm.humanized_max_delay_sec = aiConfigData.humanized_max_delay_sec
      aiForm.typing_speed_chars_sec = aiConfigData.typing_speed_chars_sec
      aiForm.debounce_window_sec = aiConfigData.debounce_window_sec
    }
  } catch (err) {
    console.error('Falha ao carregar configurações de IA:', err)
  }
}

// 7. Salvar Configuração de IA
const saveAIConfig = async () => {
  savingAI.value = true
  errorMessage.value = ''
  successMessage.value = ''
  try {
    const payload: any = {
      is_ai_enabled: aiForm.is_ai_enabled,
      ai_provider: aiForm.ai_provider,
      ai_model: aiForm.ai_model,
      system_prompt_custom: aiForm.system_prompt_custom,
      humanized_min_delay_sec: aiForm.humanized_min_delay_sec,
      humanized_max_delay_sec: aiForm.humanized_max_delay_sec,
      typing_speed_chars_sec: aiForm.typing_speed_chars_sec,
      debounce_window_sec: aiForm.debounce_window_sec,
    }

    if (aiForm.gemini_api_key.trim()) {
      payload.gemini_api_key = aiForm.gemini_api_key.trim()
    }
    if (aiForm.openai_api_key.trim()) {
      payload.openai_api_key = aiForm.openai_api_key.trim()
    }

    const res = await api.put('/admin/whatsapp/ai-config', payload)
    if (res.data?.data) {
      Object.assign(aiConfigData, res.data.data)
      aiForm.gemini_api_key = ''
      aiForm.openai_api_key = ''
      successMessage.value = 'Configurações de IA salvas e ativadas com sucesso!'
    }
  } catch (err: any) {
    errorMessage.value = err.response?.data?.message || 'Falha ao salvar configurações de IA'
  } finally {
    savingAI.value = false
  }
}

onMounted(async () => {
  await Promise.all([fetchStatus(), fetchAIConfig()])
})

onUnmounted(() => {
  stopPolling()
})
</script>
