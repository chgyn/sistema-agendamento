<template>
  <div class="min-h-screen bg-[#F8F9FD] text-gray-900 flex flex-col justify-between relative overflow-hidden font-sans">
    <!-- Efeito de Luz Ambiente Superior (sutil em tons quentes) -->
    <div class="absolute -top-40 left-1/2 -translate-x-1/2 w-[600px] h-[350px] bg-orange-500/5 rounded-full blur-[140px] pointer-events-none"></div>

    <!-- Header Público -->
    <header class="border-b border-gray-200/80 bg-white/90 backdrop-blur-xl sticky top-0 z-50 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-orange-500 to-amber-500 flex items-center justify-center shadow-md shadow-orange-500/20">
            <Scissors class="w-5 h-5 text-white font-bold" />
          </div>
          <div>
            <span class="font-bold text-lg tracking-tight text-gray-900 font-display">AgendeFácil</span>
            <span class="text-xs ml-2 px-2 py-0.5 rounded-full bg-orange-50 text-orange-600 border border-orange-200 font-semibold">Multi-Tenant Go</span>
          </div>
        </div>
        <div class="flex items-center space-x-3 sm:space-x-4">
          <RouterLink to="/login" class="text-xs sm:text-sm font-semibold text-gray-600 hover:text-gray-900 transition px-3.5 py-2 rounded-xl hover:bg-gray-100">
            Acessar Painel
          </RouterLink>
          <RouterLink to="/register" class="text-xs sm:text-sm font-bold px-4 py-2 rounded-xl bg-orange-500 hover:bg-orange-600 text-white shadow-md shadow-orange-500/20 transition transform hover:-translate-y-0.5">
            Cadastrar Estabelecimento
          </RouterLink>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <main class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20 flex flex-col items-center text-center relative z-10">
      <div class="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-orange-50 border border-orange-200 text-orange-600 text-xs font-bold mb-6">
        <Sparkles class="w-3.5 h-3.5" />
        Sistema Completo de Agendamento em Go + Vue 3
      </div>

      <h1 class="text-4xl sm:text-6xl font-black tracking-tight max-w-4xl leading-tight text-gray-900 font-display">
        Acelere o agendamento da sua <span class="bg-gradient-to-r from-orange-600 via-amber-500 to-orange-500 bg-clip-text text-transparent">barbearia ou salão</span>
      </h1>
      <p class="mt-6 text-base sm:text-xl text-gray-600 max-w-2xl leading-relaxed">
        Plataforma multi-tenant de alta performance com prevenção nativa de dupla reserva, página pública personalizada e painel administrativo completo.
      </p>

      <!-- Demos em Destaque -->
      <div class="mt-12 w-full max-w-4xl">
        <h2 class="text-xs font-bold uppercase tracking-wider text-gray-500 mb-6">
          Escolha um estabelecimento de demonstração para testar agora
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-left">
          <!-- Tenant 1 Card (Dom Navalha - Barber Theme) -->
          <div class="group relative rounded-2xl p-6 bg-white border border-gray-200/90 shadow-sm hover:shadow-md hover:border-orange-500/50 transition-all duration-300">
            <div class="flex items-center space-x-4 mb-4">
              <img
                src="https://images.unsplash.com/photo-1503951914875-452162b0f3f1?w=100&h=100&fit=crop"
                alt="Dom Navalha"
                class="w-14 h-14 rounded-xl object-cover ring-2 ring-orange-500/30"
              />
              <div>
                <h3 class="font-bold text-lg text-gray-900 group-hover:text-orange-600 transition font-display">Barbearia Dom Navalha</h3>
                <p class="text-xs text-gray-500">Av. Paulista, São Paulo - SP</p>
                <span class="inline-block mt-1 text-[11px] font-mono text-orange-700 bg-orange-50 px-2 py-0.5 rounded-md border border-orange-200">
                  /agendamento/dom-navalha
                </span>
              </div>
            </div>
            <p class="text-sm text-gray-600 mb-6 leading-relaxed">
              Cortes masculinos modernos, barboterapia com toalha quente, pigmentação e visagismo.
            </p>
            <div class="flex items-center gap-3">
              <RouterLink
                to="/agendamento/dom-navalha"
                class="flex-1 text-center py-2.5 px-4 rounded-xl bg-orange-500 hover:bg-orange-600 text-white font-bold text-sm transition shadow-md shadow-orange-500/20"
              >
                Agendar como Cliente
              </RouterLink>
              <RouterLink
                to="/login"
                class="py-2.5 px-4 rounded-xl bg-gray-50 hover:bg-gray-100 text-gray-700 font-semibold text-sm transition border border-gray-200"
              >
                Painel Admin
              </RouterLink>
            </div>
          </div>

          <!-- Tenant 2 Card (Salão Bella Vista - Studio Theme) -->
          <div class="group relative rounded-2xl p-6 bg-white border border-gray-200/90 shadow-sm hover:shadow-md hover:border-pink-500/50 transition-all duration-300">
            <div class="flex items-center space-x-4 mb-4">
              <img
                src="https://images.unsplash.com/photo-1560066984-138dadb4c035?w=100&h=100&fit=crop"
                alt="Salão Bella Vista"
                class="w-14 h-14 rounded-xl object-cover ring-2 ring-pink-500/30"
              />
              <div>
                <h3 class="font-bold text-lg text-gray-900 group-hover:text-pink-600 transition font-display">Salão Bella Vista Studio</h3>
                <p class="text-xs text-gray-500">Ipanema, Rio de Janeiro - RJ</p>
                <span class="inline-block mt-1 text-[11px] font-mono text-pink-700 bg-pink-50 px-2 py-0.5 rounded-md border border-pink-200">
                  /agendamento/bella-vista
                </span>
              </div>
            </div>
            <p class="text-sm text-gray-600 mb-6 leading-relaxed">
              Cortes femininos, escova, mechas, manicure completa e design de sobrancelhas.
            </p>
            <div class="flex items-center gap-3">
              <RouterLink
                to="/agendamento/bella-vista"
                class="flex-1 text-center py-2.5 px-4 rounded-xl bg-pink-600 hover:bg-pink-700 text-white font-bold text-sm transition shadow-md shadow-pink-600/20"
              >
                Agendar como Cliente
              </RouterLink>
              <RouterLink
                to="/login"
                class="py-2.5 px-4 rounded-xl bg-gray-50 hover:bg-gray-100 text-gray-700 font-semibold text-sm transition border border-gray-200"
              >
                Painel Admin
              </RouterLink>
            </div>
          </div>
        </div>
      </div>

      <!-- Feature Grid -->
      <div class="mt-20 grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-5xl text-left">
        <div class="p-5 rounded-2xl bg-white border border-gray-200/90 shadow-sm hover:shadow-md hover:border-gray-300 transition">
          <div class="w-10 h-10 rounded-xl bg-orange-50 border border-orange-200 text-orange-600 flex items-center justify-center mb-3">
            <ShieldCheck class="w-5 h-5" />
          </div>
          <h4 class="font-bold text-gray-900 mb-1 font-display">Zero Dupla Reserva</h4>
          <p class="text-xs text-gray-600 leading-relaxed">
            Transações atômicas com row-level lock no Postgres e validação temporal rígida impedem conflito de horários.
          </p>
        </div>

        <div class="p-5 rounded-2xl bg-white border border-gray-200/90 shadow-sm hover:shadow-md hover:border-gray-300 transition">
          <div class="w-10 h-10 rounded-xl bg-amber-50 border border-amber-200 text-amber-600 flex items-center justify-center mb-3">
            <Layers class="w-5 h-5" />
          </div>
          <h4 class="font-bold text-gray-900 mb-1 font-display">Isolamento Multi-Tenant</h4>
          <p class="text-xs text-gray-600 leading-relaxed">
            Contas 100% isoladas por tenant_id validado em middleware JWT e repositórios com escopo forçado.
          </p>
        </div>

        <div class="p-5 rounded-2xl bg-white border border-gray-200/90 shadow-sm hover:shadow-md hover:border-gray-300 transition">
          <div class="w-10 h-10 rounded-xl bg-yellow-50 border border-yellow-200 text-yellow-600 flex items-center justify-center mb-3">
            <Zap class="w-5 h-5" />
          </div>
          <h4 class="font-bold text-gray-900 mb-1 font-display">Filas Assíncronas Asynq</h4>
          <p class="text-xs text-gray-600 leading-relaxed">
            Worker dedicado em Go integrado com Redis para disparo instantâneo de lembretes e confirmações.
          </p>
        </div>
      </div>
    </main>

    <!-- Footer Público -->
    <footer class="border-t border-gray-200/80 bg-white/80 py-6 text-center text-xs text-gray-500 relative z-10">
      <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span>&copy; 2026 AgendeFácil — Plataforma de Agendamentos Multi-Tenant</span>
        <div class="flex items-center space-x-4">
          <span class="text-orange-600 font-mono">Backend: Go 1.24 + Gin + GORM</span>
          <span class="text-gray-600 font-mono">Frontend: Vue 3 + Vite + Tailwind</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Scissors, Sparkles, ShieldCheck, Layers, Zap } from 'lucide-vue-next'
</script>
