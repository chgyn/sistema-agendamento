<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 selection:bg-emerald-500 selection:text-white flex flex-col justify-between">
    <!-- Header do Estabelecimento -->
    <header class="border-b border-slate-800/80 bg-slate-900/80 backdrop-blur-md sticky top-0 z-40">
      <div class="max-w-4xl mx-auto px-4 py-3.5 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <img
            v-if="bookingStore.tenant?.logo_url"
            :src="bookingStore.tenant.logo_url"
            :alt="bookingStore.tenant.name"
            class="w-10 h-10 rounded-xl object-cover ring-2 ring-emerald-500/30"
          />
          <div v-else class="w-10 h-10 rounded-xl bg-emerald-500/20 text-emerald-400 flex items-center justify-center font-bold">
            <Scissors class="w-5 h-5" />
          </div>
          <div>
            <h1 class="font-bold text-base text-white tracking-tight leading-none">
              {{ bookingStore.tenant?.name || 'Carregando estabelecimento...' }}
            </h1>
            <p class="text-xs text-slate-400 mt-1 flex items-center gap-1.5">
              <MapPin class="w-3 h-3 text-emerald-400" />
              <span>{{ bookingStore.tenant?.city ? `${bookingStore.tenant.address || ''} — ${bookingStore.tenant.city}, ${bookingStore.tenant.state}` : 'Agendamento Online' }}</span>
            </p>
          </div>
        </div>

        <a
          v-if="bookingStore.tenant?.phone"
          :href="`https://wa.me/${cleanPhone(bookingStore.tenant.phone)}`"
          target="_blank"
          class="hidden sm:flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/10 text-emerald-400 hover:bg-emerald-500/20 border border-emerald-500/30 text-xs font-medium transition"
        >
          <Phone class="w-3.5 h-3.5" />
          <span>WhatsApp</span>
        </a>
      </div>
    </header>

    <!-- Main Wizard Container -->
    <main class="flex-1 max-w-4xl w-full mx-auto px-4 py-6 sm:py-10">
      <!-- Loading Inicial do Tenant -->
      <div v-if="isLoadingTenant" class="text-center py-20">
        <Loader2 class="w-10 h-10 animate-spin text-emerald-500 mx-auto mb-4" />
        <p class="text-slate-400 text-sm">Carregando dados da barbearia/salão...</p>
      </div>

      <!-- Erro ao Carregar Estabelecimento -->
      <div v-else-if="bookingStore.bookingError && !bookingStore.tenant" class="max-w-md mx-auto text-center py-16">
        <div class="w-16 h-16 rounded-full bg-red-500/10 text-red-400 flex items-center justify-center mx-auto mb-4 border border-red-500/20">
          <AlertTriangle class="w-8 h-8" />
        </div>
        <h2 class="text-xl font-bold text-white mb-2">Estabelecimento Não Encontrado</h2>
        <p class="text-sm text-slate-400 mb-6">O link que você acessou pode estar incorreto ou desativado.</p>
        <RouterLink to="/" class="px-5 py-2.5 rounded-xl bg-slate-800 text-slate-200 text-sm font-medium hover:bg-slate-700 transition">
          Voltar para o início
        </RouterLink>
      </div>

      <!-- Fluxo de Agendamento em 4 Passos + Confirmação -->
      <div v-else class="space-y-6">
        <!-- Barra de Progresso / Stepper -->
        <div v-if="bookingStore.currentStep < 5" class="bg-slate-900/60 border border-slate-800/80 rounded-2xl p-3.5 sm:p-4">
          <div class="flex items-center justify-between text-xs font-medium text-slate-400">
            <button
              @click="bookingStore.currentStep = 1"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 1 ? 'text-emerald-400 font-bold' : bookingStore.currentStep > 1 ? 'text-slate-200 cursor-pointer' : '']"
            >
              <span :class="['w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold', bookingStore.currentStep >= 1 ? 'bg-emerald-500 text-slate-950' : 'bg-slate-800 text-slate-400']">1</span>
              <span>Serviço</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-slate-600" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 2)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 2 ? 'text-emerald-400 font-bold' : bookingStore.currentStep > 2 ? 'text-slate-200 cursor-pointer' : '']"
            >
              <span :class="['w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold', bookingStore.currentStep >= 2 ? 'bg-emerald-500 text-slate-950' : 'bg-slate-800 text-slate-400']">2</span>
              <span>Profissional</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-slate-600" />
            <button
              @click="bookingStore.selectedService && (bookingStore.currentStep = 3)"
              :disabled="!bookingStore.selectedService"
              :class="['flex items-center gap-1.5 transition', bookingStore.currentStep === 3 ? 'text-emerald-400 font-bold' : bookingStore.currentStep > 3 ? 'text-slate-200 cursor-pointer' : '']"
            >
              <span :class="['w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold', bookingStore.currentStep >= 3 ? 'bg-emerald-500 text-slate-950' : 'bg-slate-800 text-slate-400']">3</span>
              <span>Data & Hora</span>
            </button>
            <ChevronRight class="w-3.5 h-3.5 text-slate-600" />
            <span :class="['flex items-center gap-1.5', bookingStore.currentStep === 4 ? 'text-emerald-400 font-bold' : '']">
              <span :class="['w-5 h-5 rounded-full flex items-center justify-center text-[10px] font-bold', bookingStore.currentStep >= 4 ? 'bg-emerald-500 text-slate-950' : 'bg-slate-800 text-slate-400']">4</span>
              <span>Identificação</span>
            </span>
          </div>
        </div>

        <!-- Alerta de Erro -->
        <div v-if="bookingStore.bookingError" class="p-4 rounded-xl bg-red-500/10 border border-red-500/30 text-red-300 text-sm flex items-start gap-3">
          <AlertCircle class="w-5 h-5 text-red-400 shrink-0 mt-0.5" />
          <div class="flex-1">
            <strong class="font-semibold block mb-0.5">Aviso</strong>
            <span>{{ bookingStore.bookingError }}</span>
          </div>
        </div>

        <!-- PASSO 1: Seleção de Serviço -->
        <div v-if="bookingStore.currentStep === 1" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <Sparkles class="w-5 h-5 text-emerald-400" />
              Selecione o Serviço Desejado
            </h2>
            <span class="text-xs text-slate-400">{{ bookingStore.services.length }} serviços disponíveis</span>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <div
              v-for="svc in bookingStore.services"
              :key="svc.id"
              @click="selectService(svc)"
              :class="[
                'p-4 rounded-2xl border transition-all cursor-pointer relative group flex flex-col justify-between',
                bookingStore.selectedService?.id === svc.id
                  ? 'bg-emerald-950/30 border-emerald-500 ring-2 ring-emerald-500/30 shadow-glow'
                  : 'bg-slate-900/70 border-slate-800 hover:border-slate-700 hover:bg-slate-900'
              ]"
            >
              <div>
                <div class="flex items-start justify-between gap-2 mb-2">
                  <h3 class="font-bold text-base text-white group-hover:text-emerald-400 transition">{{ svc.name }}</h3>
                  <span class="text-sm font-extrabold text-emerald-400 bg-emerald-950/60 px-2.5 py-1 rounded-lg border border-emerald-800/40">
                    R$ {{ svc.price.toFixed(2) }}
                  </span>
                </div>
                <p class="text-xs text-slate-400 line-clamp-2 mb-4">{{ svc.description || 'Atendimento profissional personalizado.' }}</p>
              </div>

              <div class="flex items-center justify-between pt-3 border-t border-slate-800/60 text-xs text-slate-400">
                <span class="flex items-center gap-1.5 font-medium">
                  <Clock class="w-3.5 h-3.5 text-emerald-400" />
                  {{ svc.duration_minutes }} minutos
                </span>
                <span :class="['px-2.5 py-1 rounded-full text-[11px] font-semibold transition', bookingStore.selectedService?.id === svc.id ? 'bg-emerald-500 text-slate-950' : 'bg-slate-800 text-slate-300 group-hover:bg-slate-700']">
                  {{ bookingStore.selectedService?.id === svc.id ? 'Selecionado ✓' : 'Escolher' }}
                </span>
              </div>
            </div>
          </div>
        </div>

        <!-- PASSO 2: Seleção de Profissional -->
        <div v-if="bookingStore.currentStep === 2" class="space-y-4">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <UserCheck class="w-5 h-5 text-emerald-400" />
              Com quem você gostaria de ser atendido?
            </h2>
            <button @click="bookingStore.currentStep = 1" class="text-xs text-slate-400 hover:text-white underline">
              Alterar serviço
            </button>
          </div>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3.5">
            <!-- Opção "Qualquer Profissional" -->
            <div
              @click="selectProfessional(null)"
              :class="[
                'p-4 rounded-2xl border transition-all cursor-pointer flex items-center space-x-3.5',
                bookingStore.selectedProfessional === null
                  ? 'bg-emerald-950/30 border-emerald-500 ring-2 ring-emerald-500/30 shadow-glow'
                  : 'bg-slate-900/70 border-slate-800 hover:border-slate-700'
              ]"
            >
              <div class="w-12 h-12 rounded-xl bg-gradient-to-tr from-emerald-500/20 to-teal-400/20 border border-emerald-500/30 flex items-center justify-center text-emerald-400">
                <Users class="w-6 h-6" />
              </div>
              <div class="flex-1">
                <h3 class="font-bold text-sm text-white">Primeiro Profissional Disponível</h3>
                <p class="text-xs text-slate-400">Maior flexibilidade e opções de horários</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional === null" class="w-5 h-5 text-emerald-400" />
            </div>

            <!-- Lista de Profissionais -->
            <div
              v-for="pro in bookingStore.professionals"
              :key="pro.id"
              @click="selectProfessional(pro)"
              :class="[
                'p-4 rounded-2xl border transition-all cursor-pointer flex items-center space-x-3.5',
                bookingStore.selectedProfessional?.id === pro.id
                  ? 'bg-emerald-950/30 border-emerald-500 ring-2 ring-emerald-500/30 shadow-glow'
                  : 'bg-slate-900/70 border-slate-800 hover:border-slate-700'
              ]"
            >
              <img
                v-if="pro.avatar_url"
                :src="pro.avatar_url"
                :alt="pro.name"
                class="w-12 h-12 rounded-xl object-cover ring-1 ring-slate-700"
              />
              <div v-else class="w-12 h-12 rounded-xl bg-slate-800 text-slate-300 flex items-center justify-center font-bold">
                {{ pro.name.charAt(0) }}
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-bold text-sm text-white truncate">{{ pro.name }}</h3>
                <p class="text-xs text-emerald-400 font-medium truncate">{{ pro.title || pro.specialty || 'Profissional' }}</p>
                <p v-if="pro.specialty" class="text-[11px] text-slate-400 truncate mt-0.5">{{ pro.specialty }}</p>
              </div>
              <Check v-if="bookingStore.selectedProfessional?.id === pro.id" class="w-5 h-5 text-emerald-400" />
            </div>
          </div>
        </div>

        <!-- PASSO 3: Seleção de Data e Horário -->
        <div v-if="bookingStore.currentStep === 3" class="space-y-6">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <CalendarIcon class="w-5 h-5 text-emerald-400" />
              Escolha o Dia e Horário
            </h2>
            <button @click="bookingStore.currentStep = 2" class="text-xs text-slate-400 hover:text-white underline">
              Alterar profissional
            </button>
          </div>

          <!-- Carrossel de Datas Próximas (Hoje + 14 dias) -->
          <div>
            <label class="block text-xs font-semibold text-slate-400 uppercase tracking-wider mb-2.5">Selecione o Dia:</label>
            <div class="flex gap-2 overflow-x-auto pb-2 scrollbar-thin">
              <button
                v-for="d in upcomingDays"
                :key="d.isoString"
                @click="changeDate(d.isoString)"
                :class="[
                  'px-3.5 py-3 rounded-2xl flex flex-col items-center min-w-[72px] border transition-all text-center shrink-0',
                  bookingStore.selectedDate === d.isoString
                    ? 'bg-emerald-500 text-slate-950 border-emerald-400 font-bold shadow-lg shadow-emerald-500/20'
                    : 'bg-slate-900/80 border-slate-800 text-slate-300 hover:bg-slate-850 hover:border-slate-700'
                ]"
              >
                <span class="text-[10px] uppercase font-semibold">{{ d.dayOfWeekShort }}</span>
                <span class="text-lg font-extrabold my-0.5">{{ d.dayOfMonth }}</span>
                <span class="text-[10px]">{{ d.monthShort }}</span>
              </button>
            </div>
          </div>

          <!-- Grade de Horários Disponíveis -->
          <div>
            <div class="flex items-center justify-between mb-3">
              <label class="text-xs font-semibold text-slate-400 uppercase tracking-wider">
                Horários Livres para {{ formattedSelectedDate }}:
              </label>
              <span v-if="!bookingStore.isLoadingSlots" class="text-xs text-emerald-400 font-medium">
                {{ bookingStore.availableSlots.length }} horários disponíveis
              </span>
            </div>

            <!-- Loader de Horários -->
            <div v-if="bookingStore.isLoadingSlots" class="text-center py-10 bg-slate-900/40 rounded-2xl border border-slate-800">
              <Loader2 class="w-6 h-6 animate-spin text-emerald-500 mx-auto mb-2" />
              <p class="text-xs text-slate-400">Verificando agenda em tempo real...</p>
            </div>

            <!-- Sem horários -->
            <div v-else-if="bookingStore.availableSlots.length === 0" class="text-center py-10 bg-slate-900/40 rounded-2xl border border-slate-800 p-6">
              <CalendarX class="w-8 h-8 text-slate-500 mx-auto mb-2" />
              <p class="text-sm font-semibold text-slate-300">Nenhum horário livre nesta data</p>
              <p class="text-xs text-slate-500 mt-1">O profissional pode estar de folga ou todos os horários foram preenchidos. Tente selecionar outro dia acima.</p>
            </div>

            <!-- Slots Disponíveis -->
            <div v-else class="grid grid-cols-3 sm:grid-cols-4 md:grid-cols-6 gap-2.5">
              <button
                v-for="slot in bookingStore.availableSlots"
                :key="slot.start_datetime"
                @click="selectSlot(slot)"
                :class="[
                  'py-2.5 px-3 rounded-xl border text-center font-bold text-sm transition-all',
                  bookingStore.selectedSlot?.start_datetime === slot.start_datetime
                    ? 'bg-emerald-500 text-slate-950 border-emerald-400 shadow-glow ring-2 ring-emerald-400'
                    : 'bg-slate-900/90 border-slate-800 text-slate-200 hover:border-emerald-500/50 hover:bg-slate-850'
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
            <h2 class="text-lg font-bold text-white flex items-center gap-2">
              <CheckCircle2 class="w-5 h-5 text-emerald-400" />
              Informações para Contato
            </h2>
            <button @click="bookingStore.currentStep = 3" class="text-xs text-slate-400 hover:text-white underline">
              Alterar horário
            </button>
          </div>

          <!-- Resumo da Escolha -->
          <div class="p-4 rounded-2xl bg-emerald-950/20 border border-emerald-800/40 text-sm space-y-2">
            <div class="flex items-center justify-between text-slate-300">
              <span class="text-xs text-slate-400">Serviço:</span>
              <span class="font-bold text-white">{{ bookingStore.selectedService?.name }} ({{ bookingStore.selectedService?.duration_minutes }} min)</span>
            </div>
            <div class="flex items-center justify-between text-slate-300">
              <span class="text-xs text-slate-400">Profissional:</span>
              <span class="font-semibold text-emerald-400">{{ bookingStore.selectedSlot?.professional_name || 'Profissional Designado' }}</span>
            </div>
            <div class="flex items-center justify-between text-slate-300">
              <span class="text-xs text-slate-400">Data e Horário:</span>
              <span class="font-bold text-white">{{ formattedSelectedDate }} às {{ bookingStore.selectedSlot?.start_time }}</span>
            </div>
            <div class="flex items-center justify-between pt-2 border-t border-emerald-900/50">
              <span class="text-xs font-semibold text-slate-300">Valor Total:</span>
              <span class="text-base font-extrabold text-emerald-400">R$ {{ bookingStore.selectedService?.price.toFixed(2) }}</span>
            </div>
          </div>

          <!-- Formulário -->
          <form class="space-y-4" @submit.prevent="confirmAppointment">
            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Seu Nome Completo *</label>
              <input
                v-model="customerForm.name"
                type="text"
                required
                placeholder="Ex: João da Silva"
                class="w-full px-4 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
              />
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">WhatsApp / Celular *</label>
                <input
                  v-model="customerForm.phone"
                  type="tel"
                  required
                  placeholder="(11) 98765-4321"
                  class="w-full px-4 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                />
              </div>

              <div>
                <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">E-mail (opcional)</label>
                <input
                  v-model="customerForm.email"
                  type="email"
                  placeholder="joao@email.com"
                  class="w-full px-4 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition"
                />
              </div>
            </div>

            <div>
              <label class="block text-xs font-semibold text-slate-300 uppercase tracking-wider mb-1.5">Observações ou Preferências (opcional)</label>
              <textarea
                v-model="customerForm.notes"
                rows="2"
                placeholder="Ex: Preferência por tesoura nas pontas, toalha bem quente..."
                class="w-full px-4 py-2.5 rounded-xl bg-slate-900 border border-slate-700 text-white placeholder-slate-500 text-sm focus:outline-none focus:ring-2 focus:ring-emerald-500 transition resize-none"
              ></textarea>
            </div>

            <button
              type="submit"
              :disabled="bookingStore.isSubmitting"
              class="w-full flex justify-center items-center py-3.5 px-4 rounded-xl shadow-lg shadow-emerald-500/20 text-base font-bold text-slate-950 bg-emerald-500 hover:bg-emerald-400 focus:outline-none focus:ring-2 focus:ring-emerald-500 transition transform hover:-translate-y-0.5 disabled:opacity-50"
            >
              <Loader2 v-if="bookingStore.isSubmitting" class="w-5 h-5 animate-spin mr-2" />
              <span>{{ bookingStore.isSubmitting ? 'Confirmando Reserva...' : 'Confirmar Agendamento Agora' }}</span>
            </button>
          </form>
        </div>

        <!-- PASSO 5: Tela de Sucesso / Comprovante -->
        <div v-if="bookingStore.currentStep === 5 && bookingStore.confirmedAppointment" class="max-w-lg mx-auto text-center py-6">
          <div class="w-16 h-16 rounded-2xl bg-emerald-500/20 border border-emerald-500/40 text-emerald-400 flex items-center justify-center mx-auto mb-4 shadow-glow">
            <CheckCircle2 class="w-8 h-8" />
          </div>

          <h2 class="text-2xl font-black text-white tracking-tight">Agendamento Confirmado!</h2>
          <p class="text-sm text-slate-400 mt-1 mb-6">
            Sua reserva foi registrada no sistema. Te esperamos no horário marcado!
          </p>

          <!-- Cartão do Comprovante -->
          <div class="glass-panel p-6 rounded-2xl text-left border border-slate-800 space-y-3 mb-6">
            <div class="flex items-center justify-between pb-3 border-b border-slate-800">
              <span class="text-xs uppercase tracking-wider text-slate-400 font-semibold">Código do Agendamento</span>
              <span class="text-xs font-mono font-bold text-emerald-400 bg-emerald-950/60 px-2 py-0.5 rounded border border-emerald-800/40">
                #{{ bookingStore.confirmedAppointment.id?.substring(0, 8) }}
              </span>
            </div>

            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-slate-400">Cliente:</span>
                <span class="font-semibold text-white">{{ bookingStore.confirmedAppointment.customer?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-400">Serviço:</span>
                <span class="font-semibold text-white">{{ bookingStore.confirmedAppointment.service?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-400">Profissional:</span>
                <span class="font-semibold text-emerald-400">{{ bookingStore.confirmedAppointment.professional?.name }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-slate-400">Data e Hora:</span>
                <span class="font-bold text-white">{{ formatConfirmedDate(bookingStore.confirmedAppointment.start_at) }}</span>
              </div>
              <div class="flex justify-between pt-2 border-t border-slate-800/80">
                <span class="text-slate-300 font-semibold">Valor a Pagar no Local:</span>
                <span class="text-lg font-black text-emerald-400">R$ {{ bookingStore.confirmedAppointment.total_price?.toFixed(2) }}</span>
              </div>
            </div>
          </div>

          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="bookingStore.resetBooking()"
              class="flex-1 py-3 px-4 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-bold text-sm transition"
            >
              Fazer Novo Agendamento
            </button>
            <RouterLink
              to="/"
              class="py-3 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-medium text-sm transition"
            >
              Voltar ao Início
            </RouterLink>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer Simples -->
    <footer class="border-t border-slate-900 py-4 text-center text-xs text-slate-500">
      Agendamento seguro powered by <strong class="text-slate-400">AgendeFácil Go</strong>
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

// Próximos 14 dias para o seletor rápido
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
