<template>
  <div class="space-y-6">
    <!-- Cabeçalho -->
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold text-white tracking-tight flex items-center gap-3">
          <div class="w-10 h-10 rounded-xl bg-emerald-500/20 border border-emerald-500/30 flex items-center justify-center text-emerald-400">
            <Bot class="w-6 h-6" />
          </div>
          WhatsApp & Atendimento com IA
        </h1>
        <p class="text-sm text-slate-400 mt-1">
          Conecte o número de WhatsApp do estabelecimento e configure o agente inteligente para conduzir agendamentos.
        </p>
      </div>

      <!-- Badge de Status Geral -->
      <div class="flex items-center gap-2 self-start sm:self-auto">
        <span
          v-if="statusData.status === 'CONNECTED'"
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
          WhatsApp Conectado
        </span>
        <span
          v-else-if="statusData.status === 'QRCODE'"
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/30"
        >
          <span class="w-2 h-2 rounded-full bg-amber-400 animate-ping"></span>
          Aguardando Leitura do QR Code
        </span>
        <span
          v-else
          class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-xs font-semibold bg-slate-800 text-slate-400 border border-slate-700"
        >
          <span class="w-2 h-2 rounded-full bg-slate-500"></span>
          WhatsApp Desconectado
        </span>
      </div>
    </div>

    <!-- Navegação por Abas -->
    <div class="flex border-b border-slate-800 gap-4">
      <button
        @click="activeTab = 'connection'"
        class="pb-3 px-2 text-sm font-semibold transition border-b-2 flex items-center gap-2"
        :class="activeTab === 'connection' ? 'border-emerald-500 text-emerald-400' : 'border-transparent text-slate-400 hover:text-slate-200'"
      >
        <Smartphone class="w-4 h-4" />
        <span>1. Conexão do WhatsApp</span>
      </button>

      <button
        @click="activeTab = 'ai'"
        class="pb-3 px-2 text-sm font-semibold transition border-b-2 flex items-center gap-2"
        :class="activeTab === 'ai' ? 'border-emerald-500 text-emerald-400' : 'border-transparent text-slate-400 hover:text-slate-200'"
      >
        <Sparkles class="w-4 h-4" />
        <span>2. Configuração do Atendente com IA</span>
      </button>
    </div>

    <!-- Mensagens de Alerta Global -->
    <div v-if="successMessage" class="p-4 rounded-xl bg-emerald-500/10 border border-emerald-500/30 text-emerald-300 text-sm flex items-center justify-between">
      <div class="flex items-center gap-2">
        <CheckCircle2 class="w-5 h-5 shrink-0 text-emerald-400" />
        <span>{{ successMessage }}</span>
      </div>
      <button @click="successMessage = ''" class="text-xs text-emerald-400 hover:text-white">✕</button>
    </div>

    <div v-if="errorMessage" class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-sm flex items-center justify-between">
      <div class="flex items-center gap-2">
        <AlertCircle class="w-5 h-5 shrink-0 text-rose-400" />
        <span>{{ errorMessage }}</span>
      </div>
      <button @click="errorMessage = ''" class="text-xs text-rose-400 hover:text-white">✕</button>
    </div>

    <!-- ============================================================== -->
    <!-- ABA 1: CONEXÃO DO WHATSAPP (WUZAPI)                           -->
    <!-- ============================================================== -->
    <div v-if="activeTab === 'connection'" class="space-y-6">
      <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
        <!-- Card Principal de Pareamento -->
        <div class="lg:col-span-7 bg-slate-900/80 border border-slate-800 rounded-2xl p-6 backdrop-blur-md space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <QrCode class="w-5 h-5 text-emerald-400" />
              Pareamento do Dispositivo
            </h2>
            <button
              @click="fetchStatus"
              :disabled="loadingStatus"
              class="p-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-slate-300 transition text-xs flex items-center gap-1.5"
              title="Atualizar Status"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingStatus }" />
              <span>Sincronizar</span>
            </button>
          </div>

          <!-- ESTADO 1: CONECTADO -->
          <div v-if="statusData.status === 'CONNECTED'" class="p-6 rounded-xl bg-emerald-950/20 border border-emerald-500/30 space-y-4">
            <div class="flex items-center gap-4">
              <div class="w-12 h-12 rounded-full bg-emerald-500/20 flex items-center justify-center text-emerald-400">
                <CheckCircle2 class="w-7 h-7" />
              </div>
              <div>
                <h3 class="font-bold text-white text-base">WhatsApp Conectado com Sucesso!</h3>
                <p class="text-xs text-slate-400 mt-0.5">
                  Número Vinculado: <span class="text-emerald-400 font-mono font-bold">{{ statusData.phone_number || 'Conectado via Web' }}</span>
                </p>
                <p class="text-[11px] text-slate-500 mt-0.5">Instância: {{ statusData.instance_name }}</p>
              </div>
            </div>

            <div class="pt-4 border-t border-slate-800 flex flex-wrap gap-3">
              <button
                @click="disconnect"
                :disabled="loadingAction"
                class="py-2 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-300 text-xs font-semibold transition flex items-center gap-2"
              >
                <Smartphone class="w-3.5 h-3.5" />
                <span>Pausar Conexão</span>
              </button>

              <button
                @click="logout"
                :disabled="loadingAction"
                class="py-2 px-4 rounded-xl bg-rose-500/10 hover:bg-rose-500/20 border border-rose-500/30 text-rose-400 text-xs font-semibold transition flex items-center gap-2"
              >
                <LogOut class="w-3.5 h-3.5" />
                <span>Desconectar e Limpar Sessão</span>
              </button>
            </div>
          </div>

          <!-- ESTADO 2: EXIBINDO QR CODE -->
          <div v-else-if="statusData.status === 'QRCODE' && qrCodeImage" class="flex flex-col items-center p-6 rounded-xl bg-slate-950/60 border border-slate-800 space-y-4 text-center">
            <div class="p-3 bg-white rounded-2xl shadow-xl shadow-black/40">
              <img :src="qrCodeImage" alt="QR Code WhatsApp" class="w-64 h-64 object-contain rounded-lg" />
            </div>

            <div class="space-y-1">
              <p class="text-sm font-semibold text-white">Escaneie o código com o seu WhatsApp</p>
              <p class="text-xs text-slate-400">O código é atualizado automaticamente a cada 60 segundos.</p>
            </div>

            <button
              @click="getQRCode"
              :disabled="loadingAction"
              class="py-2 px-4 rounded-xl bg-emerald-500/10 hover:bg-emerald-500/20 border border-emerald-500/30 text-emerald-400 text-xs font-semibold transition flex items-center gap-2"
            >
              <RefreshCw class="w-3.5 h-3.5" :class="{ 'animate-spin': loadingAction }" />
              <span>Gerar Novo QR Code</span>
            </button>
          </div>

          <!-- ESTADO 3: DESCONECTADO OU AGUARDANDO INÍCIO -->
          <div v-else class="p-8 rounded-xl bg-slate-950/40 border border-slate-800/80 text-center space-y-4">
            <div class="w-14 h-14 rounded-2xl bg-slate-800/80 mx-auto flex items-center justify-center text-slate-400">
              <Smartphone class="w-7 h-7" />
            </div>
            <div class="max-w-md mx-auto space-y-1">
              <h3 class="font-bold text-white text-base">Nenhum WhatsApp Conectado</h3>
              <p class="text-xs text-slate-400">
                Inicie uma sessão para gerar o QR Code e vincular o WhatsApp oficial da sua barbearia ou salão.
              </p>
            </div>
            <button
              @click="connect"
              :disabled="loadingAction"
              class="py-2.5 px-6 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-sm shadow-lg shadow-emerald-500/20 transition flex items-center gap-2 mx-auto"
            >
              <QrCode class="w-4 h-4" />
              <span>Conectar WhatsApp (Gerar QR Code)</span>
            </button>
          </div>
        </div>

        <!-- Instruções Passo a Passo -->
        <div class="lg:col-span-5 bg-slate-900/80 border border-slate-800 rounded-2xl p-6 backdrop-blur-md space-y-4">
          <h3 class="text-sm font-bold text-white uppercase tracking-wider text-slate-400 flex items-center gap-2">
            <ShieldCheck class="w-4 h-4 text-emerald-400" />
            Como Conectar seu Aparelho
          </h3>

          <ol class="space-y-3 text-xs text-slate-300">
            <li class="flex items-start gap-3 p-3 rounded-xl bg-slate-950/50 border border-slate-800">
              <span class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center justify-center shrink-0">1</span>
              <span>Abra o aplicativo do <strong>WhatsApp</strong> no seu celular principal.</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-slate-950/50 border border-slate-800">
              <span class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center justify-center shrink-0">2</span>
              <span>Acesse as configurações tocando nos <strong>três pontos (Android)</strong> ou <strong>Configurações (iPhone)</strong>.</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-slate-950/50 border border-slate-800">
              <span class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center justify-center shrink-0">3</span>
              <span>Toque em <strong>Aparelhos Conectados</strong> e depois em <strong>Conectar Aparelho</strong>.</span>
            </li>
            <li class="flex items-start gap-3 p-3 rounded-xl bg-slate-950/50 border border-slate-800">
              <span class="w-5 h-5 rounded-full bg-emerald-500/20 text-emerald-400 font-bold flex items-center justify-center shrink-0">4</span>
              <span>Aponte a câmera para o <strong>QR Code</strong> exibido nesta tela até a confirmação.</span>
            </li>
          </ol>

          <div class="p-3 rounded-xl bg-blue-950/20 border border-blue-500/30 text-blue-300 text-[11px] leading-relaxed">
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
        <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 backdrop-blur-md flex flex-col sm:flex-row sm:items-center justify-between gap-4">
          <div>
            <h3 class="text-base font-bold text-white flex items-center gap-2">
              <Bot class="w-5 h-5 text-emerald-400" />
              Atendimento Automatizado com IA
            </h3>
            <p class="text-xs text-slate-400 mt-0.5">
              Quando ativado, a IA responde automaticamente as mensagens recebidas, tirando dúvidas e realizando agendamentos na sua agenda.
            </p>
          </div>

          <label class="relative inline-flex items-center cursor-pointer shrink-0">
            <input type="checkbox" v-model="aiForm.is_ai_enabled" class="sr-only peer" />
            <div class="w-11 h-6 bg-slate-700 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-slate-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-emerald-500"></div>
          </label>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-12 gap-6">
          <!-- Coluna 1: Provedor, Chaves e Modelo -->
          <div class="lg:col-span-7 space-y-6">
            <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 backdrop-blur-md space-y-5">
              <h3 class="text-sm font-bold text-white uppercase tracking-wider text-slate-400 flex items-center gap-2">
                <Key class="w-4 h-4 text-emerald-400" />
                Provedor de Inteligência Artificial
              </h3>

              <!-- Seleção do Provedor -->
              <div class="grid grid-cols-2 gap-3">
                <button
                  type="button"
                  @click="aiForm.ai_provider = 'GEMINI'"
                  class="p-4 rounded-xl border text-left transition relative"
                  :class="aiForm.ai_provider === 'GEMINI' ? 'bg-emerald-500/10 border-emerald-500 text-white font-bold' : 'bg-slate-950/60 border-slate-800 text-slate-400 hover:border-slate-700'"
                >
                  <div class="text-sm font-semibold">Google Gemini</div>
                  <div class="text-[11px] text-slate-400 mt-1">Recomendado para custo-benefício e velocidade</div>
                  <span v-if="aiConfigData.has_gemini_api_key" class="inline-block mt-2 text-[10px] bg-emerald-500/20 text-emerald-300 px-2 py-0.5 rounded-md font-semibold">
                    ✓ Chave configurada
                  </span>
                </button>

                <button
                  type="button"
                  @click="aiForm.ai_provider = 'OPENAI'"
                  class="p-4 rounded-xl border text-left transition relative"
                  :class="aiForm.ai_provider === 'OPENAI' ? 'bg-emerald-500/10 border-emerald-500 text-white font-bold' : 'bg-slate-950/60 border-slate-800 text-slate-400 hover:border-slate-700'"
                >
                  <div class="text-sm font-semibold">OpenAI (ChatGPT)</div>
                  <div class="text-[11px] text-slate-400 mt-1">Modelos GPT-4o e GPT-4o mini</div>
                  <span v-if="aiConfigData.has_openai_api_key" class="inline-block mt-2 text-[10px] bg-emerald-500/20 text-emerald-300 px-2 py-0.5 rounded-md font-semibold">
                    ✓ Chave configurada
                  </span>
                </button>
              </div>

              <!-- Input da API Key do Gemini -->
              <div v-if="aiForm.ai_provider === 'GEMINI'" class="space-y-2">
                <label class="block text-xs font-semibold text-slate-300">
                  Google Gemini API Key
                  <span v-if="aiConfigData.has_gemini_api_key" class="text-emerald-400 text-[11px] ml-2 font-normal">
                    (Já configurada com segurança. Preencha apenas se desejar substituir)
                  </span>
                </label>
                <div class="relative">
                  <input
                    type="password"
                    v-model="aiForm.gemini_api_key"
                    :placeholder="aiConfigData.has_gemini_api_key ? '••••••••••••••••••••••••••••••••' : 'Cole sua Gemini API Key aqui (AIzaSy...)'"
                    class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-emerald-500 transition font-mono"
                  />
                </div>
                <p class="text-[11px] text-slate-500">
                  Obtenha gratuitamente em <a href="https://aistudio.google.com/app/apikey" target="_blank" class="text-emerald-400 underline">Google AI Studio</a>.
                </p>
              </div>

              <!-- Input da API Key da OpenAI -->
              <div v-if="aiForm.ai_provider === 'OPENAI'" class="space-y-2">
                <label class="block text-xs font-semibold text-slate-300">
                  OpenAI API Key
                  <span v-if="aiConfigData.has_openai_api_key" class="text-emerald-400 text-[11px] ml-2 font-normal">
                    (Já configurada com segurança. Preencha apenas se desejar substituir)
                  </span>
                </label>
                <div class="relative">
                  <input
                    type="password"
                    v-model="aiForm.openai_api_key"
                    :placeholder="aiConfigData.has_openai_api_key ? '••••••••••••••••••••••••••••••••' : 'Cole sua OpenAI API Key aqui (sk-...)'"
                    class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-emerald-500 transition font-mono"
                  />
                </div>
                <p class="text-[11px] text-slate-500">
                  Obtenha no painel da <a href="https://platform.openai.com/api-keys" target="_blank" class="text-emerald-400 underline">OpenAI Platform</a>.
                </p>
              </div>

              <!-- Seleção do Modelo -->
              <div class="space-y-2">
                <label class="block text-xs font-semibold text-slate-300">Modelo de Inteligência Artificial</label>
                <select
                  v-model="aiForm.ai_model"
                  class="w-full bg-slate-950 border border-slate-800 rounded-xl px-4 py-2.5 text-sm text-white focus:outline-none focus:border-emerald-500 transition"
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
                <label class="block text-xs font-semibold text-slate-300">
                  Instruções Adicionais do Estabelecimento (Prompt)
                </label>
                <textarea
                  v-model="aiForm.system_prompt_custom"
                  rows="4"
                  placeholder="Ex: Não atendemos sem agendamento prévio. Temos estacionamento conveniado gratuito na rua lateral. Em caso de atraso superior a 15 minutos, o horário poderá ser remarcado."
                  class="w-full bg-slate-950 border border-slate-800 rounded-xl p-3 text-sm text-white focus:outline-none focus:border-emerald-500 transition leading-relaxed resize-y"
                ></textarea>
                <p class="text-[11px] text-slate-500">
                  A IA já sabe automaticamente o nome dos seus serviços, preços, profissionais e horários livres. Use este campo para regras comerciais, tom de voz ou recados.
                </p>
              </div>
            </div>
          </div>

          <!-- Coluna 2: Ritmo Humanizado e Salvar -->
          <div class="lg:col-span-5 space-y-6">
            <div class="bg-slate-900/80 border border-slate-800 rounded-2xl p-6 backdrop-blur-md space-y-4">
              <h3 class="text-sm font-bold text-white uppercase tracking-wider text-slate-400 flex items-center gap-2">
                <Sliders class="w-4 h-4 text-emerald-400" />
                Controle de Ritmo Humanizado
              </h3>

              <p class="text-xs text-slate-400 leading-relaxed">
                Evite respostas robóticas instantâneas. O sistema simula o comportamento natural de uma pessoa digitando.
              </p>

              <!-- Tempo de espera inicial -->
              <div class="space-y-2 pt-2 border-t border-slate-800">
                <div class="flex justify-between text-xs">
                  <span class="text-slate-300 font-semibold">Pausa antes de responder</span>
                  <span class="text-emerald-400 font-mono">{{ aiForm.humanized_min_delay_sec }}s a {{ aiForm.humanized_max_delay_sec }}s</span>
                </div>
                <div class="grid grid-cols-2 gap-3">
                  <div>
                    <span class="text-[11px] text-slate-500">Mínimo (seg):</span>
                    <input
                      type="number"
                      min="1"
                      max="10"
                      v-model.number="aiForm.humanized_min_delay_sec"
                      class="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-xs text-white"
                    />
                  </div>
                  <div>
                    <span class="text-[11px] text-slate-500">Máximo (seg):</span>
                    <input
                      type="number"
                      min="1"
                      max="15"
                      v-model.number="aiForm.humanized_max_delay_sec"
                      class="w-full bg-slate-950 border border-slate-800 rounded-lg p-2 text-xs text-white"
                    />
                  </div>
                </div>
              </div>

              <!-- Velocidade de digitação -->
              <div class="space-y-2 pt-2 border-t border-slate-800">
                <div class="flex justify-between text-xs">
                  <span class="text-slate-300 font-semibold">Velocidade do "Digitando..."</span>
                  <span class="text-emerald-400 font-mono">{{ aiForm.typing_speed_chars_sec }} carac/seg</span>
                </div>
                <input
                  type="range"
                  min="15"
                  max="60"
                  v-model.number="aiForm.typing_speed_chars_sec"
                  class="w-full accent-emerald-500"
                />
                <p class="text-[10px] text-slate-500">
                  Respostas mais longas exibirão "digitando..." por mais tempo (delimitado entre 1.5s e 7.5s).
                </p>
              </div>

              <!-- Janela de agrupamento (Debounce) -->
              <div class="space-y-2 pt-2 border-t border-slate-800">
                <div class="flex justify-between text-xs">
                  <span class="text-slate-300 font-semibold">Janela de agrupamento</span>
                  <span class="text-emerald-400 font-mono">{{ aiForm.debounce_window_sec }} segundos</span>
                </div>
                <input
                  type="range"
                  min="2"
                  max="8"
                  v-model.number="aiForm.debounce_window_sec"
                  class="w-full accent-emerald-500"
                />
                <p class="text-[10px] text-slate-500">
                  Se o cliente mandar várias mensagens curtas seguidas, a IA agrupa todas antes de responder.
                </p>
              </div>
            </div>

            <!-- Botão Salvar -->
            <button
              type="submit"
              :disabled="savingAI"
              class="w-full py-3 px-6 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-sm shadow-lg shadow-emerald-500/20 transition flex items-center justify-center gap-2"
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
