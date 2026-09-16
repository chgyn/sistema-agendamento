<template>
  <div class="min-h-screen bg-[#0d0e12] text-zinc-100 selection:bg-orange-500 selection:text-zinc-950 flex flex-col justify-between relative overflow-hidden">
    <!-- Efeito de Luz Ambiente -->
    <div class="absolute -top-32 left-1/2 -translate-x-1/2 w-[550px] h-[300px] bg-orange-500/10 rounded-full blur-[130px] pointer-events-none"></div>

    <!-- Header do Estabelecimento -->
    <header class="border-b border-zinc-800/80 bg-[#121318]/80 backdrop-blur-xl sticky top-0 z-40">
      <div class="max-w-4xl mx-auto px-4 py-3.5 flex items-center justify-between">
        <div class="flex items-center space-x-3.5">
          <img
            v-if="bookingStore.tenant?.logo_url"
            :src="bookingStore.tenant.logo_url"
            :alt="bookingStore.tenant.name"
            class="w-11 h-11 rounded-2xl object-cover ring-2 ring-orange-500/40 shadow-glow-sm"
          />
          <div v-else class="w-11 h-11 rounded-2xl bg-orange-500/15 border border-orange-500/30 text-orange-400 flex items-center justify-center font-bold shadow-glow-sm">
            <Scissors class="w-6 h-6" />
          </div>
          <div>
            <h1 class="font-bold text-base text-white tracking-tight leading-none font-display">
              {{ bookingStore.tenant?.name || 'Carregando estabelecimento...' }}
            </h1>
            <p class="text-xs text-zinc-400 mt-1 flex items-center gap-1.5">
              <MapPin class="w-3.5 h-3.5 text-orange-400 shrink-0" />
              <span class="truncate max-w-[240px] sm:max-w-md">
                {{ bookingStore.tenant?.city ? `${bookingStore.tenant.address || ''} — ${bookingStore.tenant.city}, ${bookingStore.tenant.state}` : 'Agendamento Online Oficial' }}
              </span>
            </p>
          </div>
        </div>

        <a
          v-if="bookingStore.tenant?.phone"
          :href="`https://wa.me/${cleanPhone(bookingStore.tenant.phone)}`"
          target="_blank"
          class="hidden sm:flex items-center gap-2 px-3.5 py-2 rounded-xl bg-[#181922] hover:bg-zinc-800 text-zinc-200 hover:text-white border border-zinc-700/60 text-xs font-bold transition shadow-sm"
        >
          <Phone class="w-3.5 h-3.5 text-emerald-400" />
          <span>WhatsApp</span>
        </a>
      </div>
    </header>

    <!-- Main Wizard Container -->
    <main class="flex-1 max-w-4xl w-full mx-auto px-4 py-6 sm:py-10 relative z-10">
      <!-- Loading Inicial do Tenant -->
      <div v-if="isLoadingTenant" class="text-center py-20">
        <Loader2 class="w-10 h-10 animate-spin text-orange-500 mx-auto mb-4" />
        <p class="text-zinc-400 text-sm">Carregando dados da barbearia/salão...</p>
      </div>

      <!-- Erro ao Carregar Estabelecimento -->
      <div v-else-if="bookingStore.bookingError && !bookingStore.tenant" class="max-w-md mx-auto text-center py-16">
        <div class="w-16 h-16 rounded-2xl bg-rose-500/10 text-rose-400 flex items-center justify-center mx-auto mb-4 border border-rose-500/20">
          <AlertTriangle class="w-8 h-8" />
        </div>
        <h2 class="text-xl font-bold text-white mb-2 font-display">Estabelecimento Não Encontrado</h2>
        <p class="text-sm text-zinc-400 mb-6">O link que você acessou pode estar incorreto ou temporariamente indisponível.</p>
        <RouterLink to="/" class="px-5 py-2.5 rounded-xl bg-zinc-800 text-zinc-200 text-sm font-semibold hover:bg-zinc-700 transition">
          Voltar para o início
        </RouterLink>
      </div>

      <!-- Fluxo de Agendamento em 4 Passos + Confirmação -->
      <div v-else class="space-y-6">
        <!-- Barra de Progresso / Stepper EstiloMarca -->
        <div v-if="bookingStore.currentStep < 5" class="glass-panel rounded-2xl p-3.5 sm:p-4.5 border border-zinc-800/80">
          <div class="flex items-center justify-between text-xs font-medium text-zinc-400">
            <button
              @click="bookingStore.currentStep = 1"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 1 ? 'text-orange-400 font-bold' : bookingStore.currentStep > 1 ? 'text-zinc-200 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 1 ? 'bg-orange-500 text-zinc-950 shadow-glow-sm' : 'bg-zinc-800 text-zinc-400']">1</span>
              <span class="hidden sm:inline">Serviço</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-zinc-600" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 2)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 2 ? 'text-orange-400 font-bold' : bookingStore.currentStep > 2 ? 'text-zinc-200 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 2 ? 'bg-orange-500 text-zinc-950 shadow-glow-sm' : 'bg-zinc-800 text-zinc-400']">2</span>
              <span class="hidden sm:inline">Profissional</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-zinc-600" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 3)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 3 ? 'text-orange-400 font-bold' : bookingStore.currentStep > 3 ? 'text-zinc-200 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 3 ? 'bg-orange-500 text-zinc-950 shadow-glow-sm' : 'bg-zinc-800 text-zinc-400']">3</span>
              <span class="hidden sm:inline">Data & Hora</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-zinc-600" />
            <span :class="['flex items-center gap-1.5', bookingStore.currentStep === 4 ? 'text-orange-400 font-bold' : '']">
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 4 ? 'bg-orange-500 text-zinc-950 shadow-glow-sm' : 'bg-zinc-800 text-zinc-400']">4</span>
              <span class="hidden sm:inline">Identificação</span>
            </span>
          </div>
        </div>

        <!-- Alerta de Erro -->
        <div v-if="bookingStore.bookingError" class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/30 text-rose-300 text-sm flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-rose-400 shrink-0 mt-0.5" />
          <div class="flex-1">
            <strong class="font-semibold block mb-0.5">Aviso</strong>
            <span>{{ bookingStore.bookingError }}</span>
          </div>
        </div>

        <!-- PASSO 1: Seleção de Serviço -->
        <div v-if="bookingStore.currentStep === 1" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-white flex items-center gap-2.5 font-display">
              <Sparkles class="w-5 h-5 text-orange-400" />
              Qual serviço você deseja hoje?
            </h2>
            <span class="text-xs text-zinc-400">{{ bookingStore.services.length }} serviços disponíveis</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div
              v-for="svc in bookingStore.services"
              :key="svc.id"
              @click="selectService(svc)"
              :class="[
                'p-5 rounded-2xl border transition-all duration-200 cursor-pointer relative group flex flex-col justify-between',
                bookingStore.selectedService?.id === svc.id
                  ? 'bg-[#1a1b24] border-orange-500 ring-2 ring-orange-500/40 shadow-glow'
                  : 'bg-[#15161d]/85 border-zinc-800/80 hover:border-zinc-700 hover:bg-[#1a1b24]'
              ]"
            >
              <div>
                <div class="flex items-start justify-between gap-3 mb-3">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-orange-500/10 border border-orange-500/20 flex items-center justify-center text-orange-400 shrink-0">
                      <Scissors class="w-5 h-5" />
                    </div>
                    <h3 class="font-bold text-base text-white group-hover:text-orange-400 transition font-display">{{ svc.name }}</h3>
                  </div>
                  <span class="text-sm font-black text-orange-400 bg-orange-950/60 px-3 py-1 rounded-xl border border-orange-800/40 shrink-0">
                    R$ {{ svc.price.toFixed(2).replace('.', ',') }}
                  </span>
                </div>
                <p class="text-xs text-zinc-400 line-clamp-2 mb-4 leading-relaxed">{{ svc.description || 'Atendimento profissional personalizado com excelência.' }}</p>
              </div>

              <div class="flex items-center justify-between pt-3 border-t border-zinc-800/80 text-xs text-zinc-400">
                <span class="flex items-center gap-1.5 font-semibold">
                  <Clock class="w-3.5 h-3.5 text-orange-400" />
                  {{ svc.duration_minutes }} minutos
                </span>
                <span :class="['px-3 py-1 rounded-full text-[11px] font-bold transition', bookingStore.selectedService?.id === svc.id ? 'bg-orange-500 text-zinc-950 shadow-glow-sm' : 'bg-zinc-800 text-zinc-300 group-hover:bg-zinc-700']">
                  {{ bookingStore.selectedService?.id === svc.id ? 'Selecionado ✓' : 'Escolher' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- PASSO 2: Seleção de Profissional -->
        <div v-if="bookingStore.currentStep === 2" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-white flex items-center gap-2.5 font-display">
              <UserCheck class="w-5 h-5 text-orange-400" />
              Escolha seu Barbeiro / Profissional
            </h2>
            <button @click="bookingStore.currentStep = 1" class="text-xs font-semibold text-orange-400 hover:text-orange-300 underline">
              Alterar serviço
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Opção "Qualquer Profissional" -->
            <div
              @click="selectProfessional(null)"
              :class="[
                'p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer flex items-center space-x-3.5',
                bookingStore.selectedProfessional === null
                  ? 'bg-[#1a1b24] border-orange-500 ring-2 ring-orange-500/40 shadow-glow'
                  : 'bg-[#15161d]/85 border-zinc-800/80 hover:border-zinc-700 hover:bg-[#1a1b24]'
              ]"
            >
              <div class="w-12 h-12 rounded-2xl bg-orange-500/15 border border-orange-500/30 flex items-center justify-center text-orange-400 shrink-0">
                <Users class="w-6 h-6" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="font-bold text-sm text-white">Primeiro Disponível</h3>
                  <span class="w-2 h-2 rounded-full bg-emerald-400 animate-pulse"></span>
                </div>
                <p class="text-xs text-zinc-400 mt-0.5">Maior flexibilidade e horários livres</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional === null" class="w-5 h-5 text-orange-400 font-bold" />
            </div>

            <!-- Lista de Profissionais -->
            <div
              v-for="pro in bookingStore.professionals"
              :key="pro.id"
              @click="selectProfessional(pro)"
              :class="[
                'p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer flex items-center space-x-3.5',
                bookingStore.selectedProfessional?.id === pro.id
                  ? 'bg-[#1a1b24] border-orange-500 ring-2 ring-orange-500/40 shadow-glow'
                  : 'bg-[#15161d]/85 border-zinc-800/80 hover:border-zinc-700 hover:bg-[#1a1b24]'
              ]"
            >
              <div class="relative shrink-0">
                <img
                  v-if="pro.avatar_url"
                  :src="pro.avatar_url"
                  :alt="pro.name"
                  class="w-12 h-12 rounded-2xl object-cover ring-2 ring-zinc-700"
                />
                <div v-else class="w-12 h-12 rounded-2xl bg-zinc-800 border border-zinc-700 text-orange-400 flex items-center justify-center font-bold text-lg font-display">
                  {{ pro.name.charAt(0) }}
                </div>
                <!-- Status online / disponível (como na referência) -->
                <span class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 rounded-full bg-emerald-500 border-2 border-[#15161d]"></span>
              </div>

              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-sm text-white truncate font-display">{{ pro.name }}</h3>
                <p class="text-xs text-orange-400 font-semibold truncate">{{ pro.title || pro.specialty || 'Profissional' }}</p>
                <p v-if="pro.specialty" class="text-[11px] text-zinc-400 truncate mt-0.5">{{ pro.specialty }}</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional?.id === pro.id" class="w-5 h-5 text-orange-400 font-bold" />
            </div>
          </div>
        </div>

        <!-- PASSO 3: Seleção de Data e Horário -->
        <div v-if="bookingStore.currentStep === 3" class="space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-white flex items-center gap-2.5 font-display">
              <CalendarIcon class="w-5 h-5 text-orange-400" />
              Escolha a Data & Horário
            </h2>
            <button @click="bookingStore.currentStep = 2" class="text-xs font-semibold text-orange-400 hover:text-orange-300 underline">
              Alterar profissional
            </button>
          </div>

          <!-- Carrossel de Datas Próximas (EstiloMarca) -->
          <div>
            <label class="block text-xs font-bold text-zinc-400 uppercase tracking-wider mb-3">Selecione o Dia:</label>
            <div class="flex gap-2.5 overflow-x-auto pb-3 scrollbar-thin">
              <button
                v-for="d in upcomingDays"
                :key="d.isoString"
                @click="changeDate(d.isoString)"
                :class="[
                  'px-4 py-3.5 rounded-2xl flex flex-col items-center min-w-[76px] border transition-all duration-200 text-center shrink-0 cursor-pointer',
                  bookingStore.selectedDate === d.isoString
                    ? 'bg-orange-500 text-zinc-950 border-orange-400 font-black shadow-glow scale-105'
                    : 'bg-[#16171e]/90 border-zinc-800 text-zinc-300 hover:border-zinc-700 hover:bg-[#1d1f2a]'
                ]"
              >
                <span class="text-[10px] uppercase font-bold tracking-wider" :class="bookingStore.selectedDate === d.isoString ? 'text-zinc-950 font-black' : 'text-zinc-400'">{{ d.dayOfWeekShort }}</span>
                <span class="text-xl font-black my-1 font-display">{{ d.dayOfMonth }}</span>
                <span class="text-[11px] font-semibold" :class="bookingStore.selectedDate === d.isoString ? 'text-zinc-950' : 'text-zinc-400'">{{ d.monthShort }}</span>
              </button>
            </div>
          </div>

          <!-- Grade de Horários Disponíveis -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <label class="text-xs font-bold text-zinc-400 uppercase tracking-wider">
                Horários Livres para {{ formattedSelectedDate }}:
              </label>
              <span v-if="!bookingStore.isLoadingSlots" class="text-xs text-orange-400 font-bold">
                {{ bookingStore.availableSlots.length }} horários disponíveis
              </span>
            </div>

            <!-- Loader de Horários -->
            <div v-if="bookingStore.isLoadingSlots" class="text-center py-12 bg-[#16171e]/60 rounded-2xl border border-zinc-800">
              <Loader2 class="w-6 h-6 animate-spin text-orange-500 mx-auto mb-2" />
              <p class="text-xs text-zinc-400">Consultando agenda em tempo real...</p>
            </div>

            <!-- Sem horários -->
            <div v-else-if="bookingStore.availableSlots.length === 0" class="text-center py-12 bg-[#16171e]/60 rounded-2xl border border-zinc-800 p-6">
              <CalendarX class="w-8 h-8 text-zinc-500 mx-auto mb-2" />
              <p class="text-sm font-bold text-zinc-300">Nenhum horário livre nesta data</p>
              <p class="text-xs text-zinc-500 mt-1 max-w-sm mx-auto">Todos os horários foram preenchidos ou o profissional está de folga. Tente selecionar outro dia acima.</p>
            </div>

            <!-- Slots Disponíveis -->
            <div v-else class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-2.5">
              <button
                v-for="slot in bookingStore.availableSlots"
                :key="slot.start_datetime"
                @click="selectSlot(slot)"
                :class="[
                  'py-3 px-3 rounded-xl border text-center font-bold text-sm transition-all duration-200 cursor-pointer font-display',
                  bookingStore.selectedSlot?.start_datetime === slot.start_datetime
                    ? 'bg-orange-500 text-zinc-950 border-orange-400 shadow-glow font-black ring-2 ring-orange-400'
                    : 'bg-[#16171e]/90 border-zinc-800 text-zinc-200 hover:border-orange-500/50 hover:bg-[#1d1f2a]'
                ]"
              >
                {{ slot.start_time }}
              </button>
            </div>
          </div>
        </div>

        <!-- PASSO 4: Dados do Cliente e Confirmação -->
        <div v-if="bookingStore.currentStep === 4" class="space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-white flex items-center gap-2.5 font-display">
              <CheckCircle2 class="w-5 h-5 text-orange-400" />
              Informações para Contato
            </h2>
            <button @click="bookingStore.currentStep = 3" class="text-xs font-semibold text-orange-400 hover:text-orange-300 underline">
              Alterar horário
            </button>
          </div>

          <!-- Resumo da Escolha (Ticket Style) -->
          <div class="p-5 rounded-2xl bg-[#16171e] border border-zinc-800 text-sm space-y-2.5 relative overflow-hidden">
            <div class="flex items-center justify-between text-zinc-300">
              <span class="text-xs text-zinc-400">Serviço:</span>
              <span class="font-bold text-white font-display">{{ bookingStore.selectedService?.name }} ({{ bookingStore.selectedService?.duration_minutes }} min)</span>
            </div>
            <div class="flex items-center justify-between text-zinc-300">
              <span class="text-xs text-zinc-400">Profissional:</span>
              <span class="font-bold text-orange-400">{{ bookingStore.selectedSlot?.professional_name || 'Primeiro Disponível' }}</span>
            </div>
            <div class="flex items-center justify-between text-zinc-300">
              <span class="text-xs text-zinc-400">Data e Horário:</span>
              <span class="font-bold text-white font-display">{{ formattedSelectedDate }} às {{ bookingStore.selectedSlot?.start_time }}</span>
            </div>
            <div class="flex items-center justify-between pt-3 border-t border-dashed border-zinc-800">
              <span class="text-xs font-bold text-zinc-300 uppercase tracking-wider">Valor Total:</span>
              <span class="text-xl font-black text-orange-400 font-display">R$ {{ bookingStore.selectedService?.price.toFixed(2).replace('.', ',') }}</span>
            </div>
          </div>

          <!-- Formulário -->
          <form class="space-y-4" @submit.prevent="confirmAppointment">
            <div>
              <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Seu Nome Completo *</label>
              <input
                v-model="customerForm.name"
                type="text"
                required
                placeholder="Ex: João da Silva"
                class="w-full px-4 py-3 rounded-xl bg-[#14151c] border border-zinc-800 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">WhatsApp / Celular *</label>
                <input
                  v-model="customerForm.phone"
                  type="tel"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-4 py-3 rounded-xl bg-[#14151c] border border-zinc-800 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">E-mail (opcional)</label>
                <input
                  v-model="customerForm.email"
                  type="email"
                  placeholder="joao@email.com"
                  class="w-full px-4 py-3 rounded-xl bg-[#14151c] border border-zinc-800 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-zinc-300 uppercase tracking-wider mb-1.5">Observações ou Preferências (opcional)</label>
              <textarea
                v-model="customerForm.notes"
                rows="2"
                placeholder="Ex: Preferência por toalha bem quente, corte com tesoura nas pontas..."
                class="w-full px-4 py-3 rounded-xl bg-[#14151c] border border-zinc-800 text-white placeholder-zinc-500 text-sm focus:outline-none focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition resize-none"
              ></textarea>
            </div>

            <button
              type="submit"
              :disabled="bookingStore.isSubmitting"
              class="w-full flex justify-center items-center py-4 px-4 rounded-xl shadow-glow text-base font-bold text-zinc-950 bg-orange-500 hover:bg-orange-400 focus:outline-none focus:ring-2 focus:ring-orange-500 transition transform hover:-translate-y-0.5 active:scale-[0.98] disabled:opacity-50 cursor-pointer"
            >
              <Loader2 v-if="bookingStore.isSubmitting" class="w-5 h-5 animate-spin mr-2" />
              <span>{{ bookingStore.isSubmitting ? 'Confirmando Reserva...' : 'Confirmar Agendamento Agora' }}</span>
            </button>
          </form>
        </div>

        <!-- PASSO 5: Tela de Sucesso / Comprovante Estilo Voucher -->
        <div v-if="bookingStore.currentStep === 5 && bookingStore.confirmedAppointment" class="max-w-lg mx-auto text-center py-6">
          <div class="w-16 h-16 rounded-2xl bg-orange-500/20 border border-orange-500/40 text-orange-400 flex items-center justify-center mx-auto mb-4 shadow-glow">
            <CheckCircle2 class="w-9 h-9" />
          </div>

          <h2 class="text-2xl sm:text-3xl font-black text-white tracking-tight font-display">Agendamento Confirmado!</h2>
          <p class="text-sm text-zinc-400 mt-1 mb-6">
            Sua reserva foi registrada com sucesso. Te esperamos no horário marcado!
          </p>

          <!-- Cartão do Comprovante (Ticket Moderno) -->
          <div class="glass-panel p-6 rounded-2xl text-left border border-zinc-800 space-y-3.5 mb-6">
            <div class="flex items-center justify-between pb-3 border-b border-zinc-800">
              <span class="text-xs uppercase tracking-wider text-zinc-400 font-bold">Código da Reserva</span>
              <span class="text-xs font-mono font-black text-orange-400 bg-orange-950/60 px-2.5 py-1 rounded-md border border-orange-800/40">
                #{{ bookingStore.confirmedAppointment.id?.substring(0, 8) }}
              </span>
            </div>

            <div class="space-y-2.5 text-sm">
              <div class="flex justify-between">
                <span class="text-zinc-400">Cliente:</span>
                <span class="font-bold text-white">{{ bookingStore.confirmedAppointment.customer?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-zinc-400">Serviço:</span>
                <span class="font-bold text-white font-display">{{ bookingStore.confirmedAppointment.service?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-zinc-400">Profissional:</span>
                <span class="font-bold text-orange-400">{{ bookingStore.confirmedAppointment.professional?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-zinc-400">Data e Hora:</span>
                <span class="font-bold text-white font-display">{{ formatConfirmedDate(bookingStore.confirmedAppointment.start_at) }}</span>
              </div>
              <div class="flex justify-between pt-3 border-t border-dashed border-zinc-800">
                <span class="text-zinc-300 font-bold">Valor a Pagar no Local:</span>
                <span class="text-xl font-black text-orange-400 font-display">R$ {{ bookingStore.confirmedAppointment.total_price?.toFixed(2).replace('.', ',') }}</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="bookingStore.resetBooking()"
              class="flex-1 py-3.5 px-4 rounded-xl bg-orange-500 hover:bg-orange-400 text-zinc-950 font-bold text-sm shadow-glow-sm transition cursor-pointer"
            >
              Fazer Novo Agendamento
            </button>
            <RouterLink
              to="/"
              class="py-3.5 px-4 rounded-xl bg-[#16171e] hover:bg-zinc-800 text-zinc-200 font-semibold text-sm transition border border-zinc-700/60 text-center"
            >
              Voltar ao Início
            </RouterLink>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer Simples -->
    <footer class="border-t border-zinc-900 py-4 text-center text-xs text-zinc-500 relative z-10">
      Agendamento seguro powered by <strong class="text-zinc-400 font-semibold">AgendeFácil Go</strong>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  Scissors, MapPin, Phone, Sparkles, UserCheck, Users, Check,
  Calendar as CalendarIcon, Clock, CalendarX, CheckCircle2,
  ChevronRight, AlertCircle, AlertTriangle, Loader2
} from 'lucide-vue-next'
import { format, addDays } from 'date-fns'
import { ptBR } from 'date-fns/locale'
import { useBookingStore, type PublicService, type PublicProfessional, type TimeSlot } from '../../stores/booking'

const route = useRoute()
const bookingStore = useBookingStore()
const isLoadingTenant = ref(true)

const customerForm = reactive({
  name: '',
  phone: '',
  email: '',
  notes: '',
})

const slug = computed(() => (route.params.slug as string) || 'dom-navalha')

onMounted(async () => {
  bookingStore.resetBooking()
  isLoadingTenant.value = true
  await bookingStore.loadTenant(slug.value)
  isLoadingTenant.value = false
})

// Próximos 14 dias para o seletor rápido estilo EstiloMarca
const upcomingDays = computed(() => {
  const list = []
  const today = new Date()
  for (let i = 0; i < 14; i++) {
    const d = addDays(today, i)
    list.push({
      isoString: format(d, 'yyyy-MM-dd'),
      dayOfWeekShort: format(d, 'EEE', { locale: ptBR }),
      dayOfMonth: format(d, 'dd'),
      monthShort: format(d, 'MMM', { locale: ptBR }),
    })
  }
  return list
})

const formattedSelectedDate = computed(() => {
  try {
    const parts = bookingStore.selectedDate.split('-')
    const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]))
    return format(d, "EEEE, dd 'de' MMMM", { locale: ptBR })
  } catch {
    return bookingStore.selectedDate
  }
})

function cleanPhone(phone: string) {
  return phone.replace(/\D/g, '')
}

async function selectService(svc: PublicService) {
  bookingStore.selectedService = svc
  bookingStore.selectedSlot = null
  await bookingStore.loadProfessionalsForService(slug.value, svc.id)
  bookingStore.currentStep = 2
}

async function selectProfessional(pro: PublicProfessional | null) {
  bookingStore.selectedProfessional = pro
  bookingStore.selectedSlot = null
  await bookingStore.fetchAvailability(slug.value)
  bookingStore.currentStep = 3
}

async function changeDate(dateStr: string) {
  bookingStore.selectedDate = dateStr
  bookingStore.selectedSlot = null
  await bookingStore.fetchAvailability(slug.value)
}

function selectSlot(slot: TimeSlot) {
  bookingStore.selectedSlot = slot
  bookingStore.currentStep = 4
}

async function confirmAppointment() {
  await bookingStore.createAppointment(slug.value, {
    name: customerForm.name,
    phone: customerForm.phone,
    email: customerForm.email,
    notes: customerForm.notes,
  })
}

function formatConfirmedDate(dateStr: string) {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return format(d, "dd/MM/yyyy 'às' HH:mm", { locale: ptBR })
  } catch {
    return dateStr
  }
}
</script>
