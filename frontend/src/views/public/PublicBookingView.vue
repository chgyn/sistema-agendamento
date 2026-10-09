<template>
  <div class="min-h-screen bg-[#F8F9FD] text-gray-900 selection:bg-orange-500 selection:text-white flex flex-col justify-between relative overflow-hidden font-sans">
    <!-- Efeito de Luz Ambiente -->
    <div class="absolute -top-32 left-1/2 -translate-x-1/2 w-[550px] h-[300px] bg-orange-500/5 rounded-full blur-[140px] pointer-events-none"></div>

    <!-- Header do Estabelecimento -->
    <header class="border-b border-gray-200/80 bg-white/90 backdrop-blur-xl sticky top-0 z-40 shadow-sm">
      <div class="max-w-4xl mx-auto px-4 py-3.5 flex items-center justify-between">
        <div class="flex items-center space-x-3.5">
          <img
            v-if="bookingStore.tenant?.logo_url"
            :src="bookingStore.tenant.logo_url"
            :alt="bookingStore.tenant.name"
            class="w-11 h-11 rounded-2xl object-cover ring-2 ring-orange-500/30 shadow-sm"
          />
          <div v-else class="w-11 h-11 rounded-2xl bg-orange-50 border border-orange-200 text-orange-600 flex items-center justify-center font-bold shadow-sm">
            <Scissors class="w-6 h-6" />
          </div>
          <div>
            <h1 class="font-bold text-base text-gray-900 tracking-tight leading-none font-display">
              {{ bookingStore.tenant?.name || 'Carregando estabelecimento...' }}
            </h1>
            <p class="text-xs text-gray-500 mt-1 flex items-center gap-1.5">
              <MapPin class="w-3.5 h-3.5 text-orange-500 shrink-0" />
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
          class="hidden sm:flex items-center gap-2 px-3.5 py-2 rounded-xl bg-emerald-50 hover:bg-emerald-100 text-emerald-700 border border-emerald-200 text-xs font-bold transition shadow-sm"
        >
          <Phone class="w-3.5 h-3.5 text-emerald-600" />
          <span>WhatsApp</span>
        </a>
      </div>
    </header>

    <!-- Main Wizard Container -->
    <main class="flex-1 max-w-4xl w-full mx-auto px-4 py-6 sm:py-10 relative z-10">
      <!-- Loading Inicial do Tenant -->
      <div v-if="isLoadingTenant" class="text-center py-20">
        <Loader2 class="w-10 h-10 animate-spin text-orange-500 mx-auto mb-4" />
        <p class="text-gray-500 text-sm">Carregando dados da barbearia/salão...</p>
      </div>

      <!-- Erro ao Carregar Estabelecimento -->
      <div v-else-if="bookingStore.bookingError && !bookingStore.tenant" class="max-w-md mx-auto text-center py-16">
        <div class="w-16 h-16 rounded-2xl bg-rose-50 text-rose-600 flex items-center justify-center mx-auto mb-4 border border-rose-200">
          <AlertTriangle class="w-8 h-8" />
        </div>
        <h2 class="text-xl font-bold text-gray-900 mb-2 font-display">Estabelecimento Não Encontrado</h2>
        <p class="text-sm text-gray-500 mb-6">O link que você acessou pode estar incorreto ou temporariamente indisponível.</p>
        <RouterLink to="/" class="px-5 py-2.5 rounded-xl bg-gray-100 text-gray-700 text-sm font-semibold hover:bg-gray-200 transition border border-gray-200">
          Voltar para o início
        </RouterLink>
      </div>

      <!-- Fluxo de Agendamento em 4 Passos + Confirmação -->
      <div v-else class="space-y-6">
        <!-- Barra de Progresso / Stepper EstiloMarca -->
        <div v-if="bookingStore.currentStep < 5" class="bg-white rounded-2xl p-3.5 sm:p-4.5 border border-gray-200/90 shadow-sm">
          <div class="flex items-center justify-between text-xs font-medium text-gray-500">
            <button
              @click="bookingStore.currentStep = 1"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 1 ? 'text-orange-600 font-bold' : bookingStore.currentStep > 1 ? 'text-gray-800 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 1 ? 'bg-orange-500 text-white shadow-sm' : 'bg-gray-100 text-gray-400']">1</span>
              <span class="hidden sm:inline">Serviço</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-gray-300" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 2)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 2 ? 'text-orange-600 font-bold' : bookingStore.currentStep > 2 ? 'text-gray-800 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 2 ? 'bg-orange-500 text-white shadow-sm' : 'bg-gray-100 text-gray-400']">2</span>
              <span class="hidden sm:inline">Profissional</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-gray-300" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 3)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 3 ? 'text-orange-600 font-bold' : bookingStore.currentStep > 3 ? 'text-gray-800 cursor-pointer' : '']"
            >
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 3 ? 'bg-orange-500 text-white shadow-sm' : 'bg-gray-100 text-gray-400']">3</span>
              <span class="hidden sm:inline">Data & Hora</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-gray-300" />
            <span :class="['flex items-center gap-1.5', bookingStore.currentStep === 4 ? 'text-orange-600 font-bold' : '']">
              <span :class="['w-6 h-6 rounded-full flex items-center justify-center text-[11px] font-black', bookingStore.currentStep >= 4 ? 'bg-orange-500 text-white shadow-sm' : 'bg-gray-100 text-gray-400']">4</span>
              <span class="hidden sm:inline">Identificação</span>
            </span>
          </div>
        </div>

        <!-- Alerta de Erro -->
        <div v-if="bookingStore.bookingError" class="p-4 rounded-xl bg-rose-50 border border-rose-200 text-rose-700 text-sm flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-rose-500 shrink-0 mt-0.5" />
          <div class="flex-1">
            <strong class="font-semibold block mb-0.5">Aviso</strong>
            <span>{{ bookingStore.bookingError }}</span>
          </div>
        </div>

        <!-- PASSO 1: Seleção de Serviço -->
        <div v-if="bookingStore.currentStep === 1" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-gray-900 flex items-center gap-2.5 font-display">
              <Sparkles class="w-5 h-5 text-orange-500" />
              Qual serviço você deseja hoje?
            </h2>
            <span class="text-xs text-gray-500">{{ bookingStore.services.length }} serviços disponíveis</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div
              v-for="svc in bookingStore.services"
              :key="svc.id"
              @click="selectService(svc)"
              :class="[
                'p-5 rounded-2xl border transition-all duration-200 cursor-pointer relative group flex flex-col justify-between bg-white shadow-sm',
                bookingStore.selectedService?.id === svc.id
                  ? 'border-orange-500 ring-2 ring-orange-500/20 bg-orange-50/20 shadow-md'
                  : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
              ]"
            >
              <div>
                <div class="flex items-start justify-between gap-3 mb-3">
                  <div class="flex items-center gap-3">
                    <div class="w-10 h-10 rounded-xl bg-orange-50 border border-orange-200 flex items-center justify-center text-orange-600 shrink-0">
                      <Scissors class="w-5 h-5" />
                    </div>
                    <h3 class="font-bold text-base text-gray-900 group-hover:text-orange-600 transition font-display">{{ svc.name }}</h3>
                  </div>
                  <span class="text-sm font-black text-orange-700 bg-orange-50 px-3 py-1 rounded-xl border border-orange-200 shrink-0">
                    R$ {{ svc.price.toFixed(2).replace('.', ',') }}
                  </span>
                </div>
                <p class="text-xs text-gray-600 line-clamp-2 mb-4 leading-relaxed">{{ svc.description || 'Atendimento profissional personalizado com excelência.' }}</p>
              </div>

              <div class="flex items-center justify-between pt-3 border-t border-gray-100 text-xs text-gray-500">
                <span class="flex items-center gap-1.5 font-semibold text-gray-600">
                  <Clock class="w-3.5 h-3.5 text-orange-500" />
                  {{ svc.duration_minutes }} minutos
                </span>
                <span :class="['px-3 py-1 rounded-full text-[11px] font-bold transition', bookingStore.selectedService?.id === svc.id ? 'bg-orange-500 text-white shadow-sm' : 'bg-gray-100 text-gray-700 group-hover:bg-gray-200']">
                  {{ bookingStore.selectedService?.id === svc.id ? 'Selecionado ✓' : 'Escolher' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- PASSO 2: Seleção de Profissional -->
        <div v-if="bookingStore.currentStep === 2" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-gray-900 flex items-center gap-2.5 font-display">
              <UserCheck class="w-5 h-5 text-orange-500" />
              Escolha seu Barbeiro / Profissional
            </h2>
            <button @click="bookingStore.currentStep = 1" class="text-xs font-semibold text-orange-600 hover:text-orange-700 underline">
              Alterar serviço
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
            <!-- Opção "Qualquer Profissional" -->
            <div
              @click="selectProfessional(null)"
              :class="[
                'p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer flex items-center space-x-3.5 bg-white shadow-sm',
                bookingStore.selectedProfessional === null
                  ? 'border-orange-500 ring-2 ring-orange-500/20 bg-orange-50/20 shadow-md'
                  : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
              ]"
            >
              <div class="w-12 h-12 rounded-2xl bg-orange-50 border border-orange-200 flex items-center justify-center text-orange-600 shrink-0">
                <Users class="w-6 h-6" />
              </div>
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="font-bold text-sm text-gray-900">Primeiro Disponível</h3>
                  <span class="w-2 h-2 rounded-full bg-emerald-500 animate-pulse"></span>
                </div>
                <p class="text-xs text-gray-500 mt-0.5">Maior flexibilidade e horários livres</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional === null" class="w-5 h-5 text-orange-600 font-bold" />
            </div>

            <!-- Lista de Profissionais -->
            <div
              v-for="pro in bookingStore.professionals"
              :key="pro.id"
              @click="selectProfessional(pro)"
              :class="[
                'p-4.5 rounded-2xl border transition-all duration-200 cursor-pointer flex items-center space-x-3.5 bg-white shadow-sm',
                bookingStore.selectedProfessional?.id === pro.id
                  ? 'border-orange-500 ring-2 ring-orange-500/20 bg-orange-50/20 shadow-md'
                  : 'border-gray-200 hover:border-gray-300 hover:shadow-md'
              ]"
            >
              <div class="relative shrink-0">
                <img
                  v-if="pro.avatar_url"
                  :src="pro.avatar_url"
                  :alt="pro.name"
                  class="w-12 h-12 rounded-2xl object-cover ring-2 ring-gray-200"
                />
                <div v-else class="w-12 h-12 rounded-2xl bg-gray-100 border border-gray-200 text-orange-600 flex items-center justify-center font-bold text-lg font-display">
                  {{ pro.name.charAt(0) }}
                </div>
                <!-- Status online / disponível -->
                <span class="absolute -bottom-0.5 -right-0.5 w-3.5 h-3.5 rounded-full bg-emerald-500 border-2 border-white"></span>
              </div>

              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-sm text-gray-900 truncate font-display">{{ pro.name }}</h3>
                <p class="text-xs text-orange-600 font-semibold truncate">{{ pro.title || pro.specialty || 'Profissional' }}</p>
                <p v-if="pro.specialty" class="text-[11px] text-gray-500 truncate mt-0.5">{{ pro.specialty }}</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional?.id === pro.id" class="w-5 h-5 text-orange-600 font-bold" />
            </div>
          </div>
        </div>

        <!-- PASSO 3: Seleção de Data e Horário -->
        <div v-if="bookingStore.currentStep === 3" class="space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg sm:text-xl font-bold text-gray-900 flex items-center gap-2.5 font-display">
              <CalendarIcon class="w-5 h-5 text-orange-500" />
              Escolha a Data & Horário
            </h2>
            <button @click="bookingStore.currentStep = 2" class="text-xs font-semibold text-orange-600 hover:text-orange-700 underline">
              Alterar profissional
            </button>
          </div>

          <!-- Carrossel de Datas Próximas -->
          <div>
            <label class="block text-xs font-bold text-gray-500 uppercase tracking-wider mb-3">Selecione o Dia:</label>
            <div class="flex gap-2.5 overflow-x-auto pb-3 scrollbar-thin">
              <button
                v-for="d in upcomingDays"
                :key="d.isoString"
                @click="changeDate(d.isoString)"
                :class="[
                  'px-4 py-3.5 rounded-2xl flex flex-col items-center min-w-[76px] border transition-all duration-200 text-center shrink-0 cursor-pointer shadow-sm',
                  bookingStore.selectedDate === d.isoString
                    ? 'bg-orange-500 text-white border-orange-500 font-black shadow-md shadow-orange-500/20 scale-105'
                    : 'bg-white border-gray-200 text-gray-700 hover:border-gray-300 hover:bg-gray-50'
                ]"
              >
                <span class="text-[10px] uppercase font-bold tracking-wider" :class="bookingStore.selectedDate === d.isoString ? 'text-white' : 'text-gray-400'">{{ d.dayOfWeekShort }}</span>
                <span class="text-xl font-black my-1 font-display">{{ d.dayOfMonth }}</span>
                <span class="text-[11px] font-semibold" :class="bookingStore.selectedDate === d.isoString ? 'text-white' : 'text-gray-400'">{{ d.monthShort }}</span>
              </button>
            </div>
          </div>

          <!-- Grade de Horários Disponíveis -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <label class="text-xs font-bold text-gray-500 uppercase tracking-wider">
                Horários Livres para {{ formattedSelectedDate }}:
              </label>
              <span v-if="!bookingStore.isLoadingSlots" class="text-xs text-orange-600 font-bold">
                {{ bookingStore.availableSlots.length }} horários disponíveis
              </span>
            </div>

            <!-- Loader de Horários -->
            <div v-if="bookingStore.isLoadingSlots" class="text-center py-12 bg-white rounded-2xl border border-gray-200 shadow-sm">
              <Loader2 class="w-6 h-6 animate-spin text-orange-500 mx-auto mb-2" />
              <p class="text-xs text-gray-500">Consultando agenda em tempo real...</p>
            </div>

            <!-- Sem horários -->
            <div v-else-if="bookingStore.availableSlots.length === 0" class="text-center py-12 bg-white rounded-2xl border border-gray-200 p-6 shadow-sm">
              <CalendarX class="w-8 h-8 text-gray-400 mx-auto mb-2" />
              <p class="text-sm font-bold text-gray-800">Nenhum horário livre nesta data</p>
              <p class="text-xs text-gray-500 mt-1 max-w-sm mx-auto">Todos os horários foram preenchidos ou o profissional está de folga. Tente selecionar outro dia acima.</p>
            </div>

            <!-- Slots Disponíveis -->
            <div v-else class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-2.5">
              <button
                v-for="slot in bookingStore.availableSlots"
                :key="slot.start_datetime"
                @click="selectSlot(slot)"
                :class="[
                  'py-3 px-3 rounded-xl border text-center font-bold text-sm transition-all duration-200 cursor-pointer font-display shadow-sm',
                  bookingStore.selectedSlot?.start_datetime === slot.start_datetime
                    ? 'bg-orange-500 text-white border-orange-500 shadow-md shadow-orange-500/20 font-black ring-2 ring-orange-400'
                    : 'bg-white border-gray-200 text-gray-800 hover:border-orange-400 hover:bg-orange-50/20'
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
            <h2 class="text-lg sm:text-xl font-bold text-gray-900 flex items-center gap-2.5 font-display">
              <CheckCircle2 class="w-5 h-5 text-orange-500" />
              Informações para Contato
            </h2>
            <button @click="bookingStore.currentStep = 3" class="text-xs font-semibold text-orange-600 hover:text-orange-700 underline">
              Alterar horário
            </button>
          </div>

          <!-- Resumo da Escolha (Ticket Style) -->
          <div class="p-5 rounded-2xl bg-white border border-gray-200/90 shadow-sm text-sm space-y-2.5 relative overflow-hidden">
            <div class="flex items-center justify-between text-gray-600">
              <span class="text-xs text-gray-500">Serviço:</span>
              <span class="font-bold text-gray-900 font-display">{{ bookingStore.selectedService?.name }} ({{ bookingStore.selectedService?.duration_minutes }} min)</span>
            </div>
            <div class="flex items-center justify-between text-gray-600">
              <span class="text-xs text-gray-500">Profissional:</span>
              <span class="font-bold text-orange-600">{{ bookingStore.selectedSlot?.professional_name || 'Primeiro Disponível' }}</span>
            </div>
            <div class="flex items-center justify-between text-gray-600">
              <span class="text-xs text-gray-500">Data e Horário:</span>
              <span class="font-bold text-gray-900 font-display">{{ formattedSelectedDate }} às {{ bookingStore.selectedSlot?.start_time }}</span>
            </div>
            <div class="flex items-center justify-between pt-3 border-t border-dashed border-gray-200">
              <span class="text-xs font-bold text-gray-700 uppercase tracking-wider">Valor Total:</span>
              <span class="text-xl font-black text-orange-600 font-display">R$ {{ bookingStore.selectedService?.price.toFixed(2).replace('.', ',') }}</span>
            </div>
          </div>

          <!-- Formulário -->
          <form class="space-y-4" @submit.prevent="confirmAppointment">
            <div>
              <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-1.5">Seu Nome Completo *</label>
              <input
                v-model="customerForm.name"
                type="text"
                required
                placeholder="Ex: João da Silva"
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-1.5">WhatsApp / Celular *</label>
                <input
                  v-model="customerForm.phone"
                  type="tel"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-1.5">E-mail (opcional)</label>
                <input
                  v-model="customerForm.email"
                  type="email"
                  placeholder="joao@email.com"
                  class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-bold text-gray-700 uppercase tracking-wider mb-1.5">Observações ou Preferências (opcional)</label>
              <textarea
                v-model="customerForm.notes"
                rows="2"
                placeholder="Ex: Preferência por toalha bem quente, corte com tesoura nas pontas..."
                class="w-full px-4 py-3 rounded-xl bg-gray-50 border border-gray-300 text-gray-900 placeholder-gray-400 text-sm focus:outline-none focus:bg-white focus:ring-2 focus:ring-orange-500/20 focus:border-orange-500 transition resize-none"
              ></textarea>
            </div>

            <button
              type="submit"
              :disabled="bookingStore.isSubmitting"
              class="w-full flex justify-center items-center py-4 px-4 rounded-xl shadow-md shadow-orange-500/25 text-base font-bold text-white bg-orange-500 hover:bg-orange-600 focus:outline-none focus:ring-2 focus:ring-orange-500 transition transform hover:-translate-y-0.5 active:scale-[0.98] disabled:opacity-50 cursor-pointer"
            >
              <Loader2 v-if="bookingStore.isSubmitting" class="w-5 h-5 animate-spin mr-2" />
              <span>{{ bookingStore.isSubmitting ? 'Confirmando Reserva...' : 'Confirmar Agendamento Agora' }}</span>
            </button>
          </form>
        </div>

        <!-- PASSO 5: Tela de Sucesso / Comprovante Estilo Voucher -->
        <div v-if="bookingStore.currentStep === 5 && bookingStore.confirmedAppointment" class="max-w-lg mx-auto text-center py-6">
          <div class="w-16 h-16 rounded-2xl bg-emerald-50 border border-emerald-200 text-emerald-600 flex items-center justify-center mx-auto mb-4 shadow-sm">
            <CheckCircle2 class="w-9 h-9" />
          </div>

          <h2 class="text-2xl sm:text-3xl font-black text-gray-900 tracking-tight font-display">Agendamento Confirmado!</h2>
          <p class="text-sm text-gray-500 mt-1 mb-6">
            Sua reserva foi registrada com sucesso. Te esperamos no horário marcado!
          </p>

          <!-- Cartão do Comprovante (Ticket Moderno) -->
          <div class="bg-white p-6 rounded-2xl text-left border border-gray-200/90 shadow-md space-y-3.5 mb-6">
            <div class="flex items-center justify-between pb-3 border-b border-gray-100">
              <span class="text-xs uppercase tracking-wider text-gray-500 font-bold">Código da Reserva</span>
              <span class="text-xs font-mono font-black text-orange-700 bg-orange-50 px-2.5 py-1 rounded-md border border-orange-200">
                #{{ bookingStore.confirmedAppointment.id?.substring(0, 8) }}
              </span>
            </div>

            <div class="space-y-2.5 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-500">Cliente:</span>
                <span class="font-bold text-gray-900">{{ bookingStore.confirmedAppointment.customer?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Serviço:</span>
                <span class="font-bold text-gray-900 font-display">{{ bookingStore.confirmedAppointment.service?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Profissional:</span>
                <span class="font-bold text-orange-600">{{ bookingStore.confirmedAppointment.professional?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-500">Data e Hora:</span>
                <span class="font-bold text-gray-900 font-display">{{ formatConfirmedDate(bookingStore.confirmedAppointment.start_at) }}</span>
              </div>
              <div class="flex justify-between pt-3 border-t border-dashed border-gray-200">
                <span class="text-gray-700 font-bold">Valor a Pagar no Local:</span>
                <span class="text-xl font-black text-orange-600 font-display">R$ {{ bookingStore.confirmedAppointment.total_price?.toFixed(2).replace('.', ',') }}</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="bookingStore.resetBooking()"
              class="flex-1 py-3.5 px-4 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-sm shadow-md shadow-orange-500/20 transition cursor-pointer"
            >
              Fazer Novo Agendamento
            </button>
            <RouterLink
              to="/"
              class="py-3.5 px-4 rounded-xl bg-gray-50 hover:bg-gray-100 text-gray-700 font-semibold text-sm transition border border-gray-200 text-center"
            >
              Voltar ao Início
            </RouterLink>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer Simples -->
    <footer class="border-t border-gray-200/80 bg-white/70 py-4 text-center text-xs text-gray-500 relative z-10">
      Agendamento seguro powered by <strong class="text-gray-700 font-semibold">AgendeFácil Go</strong>
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
