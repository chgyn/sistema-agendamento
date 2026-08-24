<template>
  <div class="min-h-screen bg-slate-950 text-slate-100 flex flex-col justify-between selection:bg-emerald-500 selection:text-white">
    <!-- Header -->
    <header class="border-b border-slate-800/80 bg-slate-900/60 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-emerald-500 to-teal-400 flex items-center justify-center shadow-glow">
            <Scissors class="w-5 h-5 text-slate-950 font-bold" />
          </div>
          <div>
            <span class="font-bold text-lg tracking-tight bg-gradient-to-r from-white via-slate-200 to-slate-400 bg-clip-text text-transparent">AgendeFácil</span>
            <span class="text-xs ml-2 px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20 font-medium">Multi-Tenant Go</span>
          </div>
        </div>
        <div class="flex items-center space-x-4">
          <RouterLink to="/login" class="text-sm font-medium text-slate-300 hover:text-white transition">
            Acessar Painel
          </RouterLink>
          <RouterLink to="/register" class="text-sm font-medium px-4 py-2 rounded-lg bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-semibold shadow-lg shadow-emerald-500/20 transition transform hover:-translate-y-0.5">
            Cadastrar Estabelecimento
          </RouterLink>
        </div>
      </div>
    </header>

    <!-- Hero Section -->
    <main class="flex-1 max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-20 flex flex-col items-center text-center">
      <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-emerald-500/10 border border-emerald-500/30 text-emerald-400 text-xs font-semibold mb-6 animate-pulse">
        <Sparkles class="w-3.5 h-3.5" />
        Sistema Completo de Agendamento em Go + Vue 3
      </div>

      <h1 class="text-4xl sm:text-6xl font-extrabold tracking-tight max-w-4xl leading-tight">
        Acelere o agendamento da sua <span class="bg-gradient-to-r from-emerald-400 via-teal-300 to-cyan-400 bg-clip-text text-transparent">barbearia ou salão</span>
      </h1>
      <p class="mt-6 text-lg sm:text-xl text-slate-400 max-w-2xl">
        Plataforma multi-tenant de alta performance com prevenção nativa de dupla reserva, página pública personalizada e painel administrativo completo.
      </p>

      <!-- Demos em Destaque -->
      <div class="mt-12 w-full max-w-4xl">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-slate-400 mb-6">
          Escolha um estabelecimento de demonstração para testar agora
        </h2>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 text-left">
          <!-- Tenant 1 Card -->
          <div class="group relative rounded-2xl p-6 bg-slate-900/80 border border-slate-800 hover:border-emerald-500/50 transition-all duration-300 hover:shadow-glow">
            <div class="flex items-center space-x-4 mb-4">
              <img
                src="https://images.unsplash.com/photo-1503951914875-452162b0f3f1?w=100&h=100&fit=crop"
                alt="Dom Navalha"
                class="w-14 h-14 rounded-xl object-cover ring-2 ring-emerald-500/30"
              />
              <div>
                <h3 class="font-bold text-lg text-white group-hover:text-emerald-400 transition">Barbearia Dom Navalha</h3>
                <p class="text-xs text-slate-400">Av. Paulista, São Paulo - SP</p>
                <span class="inline-block mt-1 text-[11px] font-mono text-emerald-400 bg-emerald-950/60 px-2 py-0.5 rounded border border-emerald-800/40">
                  /agendamento/dom-navalha
                </span>
              </div>
            </div>
            <p class="text-sm text-slate-400 mb-6">
              Cortes masculinos modernos, barboterapia com toalha quente, pigmentação e visagismo.
            </p>
            <div class="flex items-center gap-3">
              <RouterLink
                to="/agendamento/dom-navalha"
                class="flex-1 text-center py-2.5 px-4 rounded-xl bg-emerald-500 hover:bg-emerald-400 text-slate-950 font-semibold text-sm transition shadow-md"
              >
                Agendar como Cliente
              </RouterLink>
              <RouterLink
                to="/login"
                class="py-2.5 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-medium text-sm transition border border-slate-700"
              >
                Painel Admin
              </RouterLink>
            </div>
          </div>

          <!-- Tenant 2 Card -->
          <div class="group relative rounded-2xl p-6 bg-slate-900/80 border border-slate-800 hover:border-pink-500/50 transition-all duration-300 hover:shadow-glow-pink">
            <div class="flex items-center space-x-4 mb-4">
              <img
                src="https://images.unsplash.com/photo-1560066984-138dadb4c035?w=100&h=100&fit=crop"
                alt="Salão Bella Vista"
                class="w-14 h-14 rounded-xl object-cover ring-2 ring-pink-500/30"
              />
              <div>
                <h3 class="font-bold text-lg text-white group-hover:text-pink-400 transition">Salão Bella Vista Studio</h3>
                <p class="text-xs text-slate-400">Ipanema, Rio de Janeiro - RJ</p>
                <span class="inline-block mt-1 text-[11px] font-mono text-pink-400 bg-pink-950/60 px-2 py-0.5 rounded border border-pink-800/40">
                  /agendamento/bella-vista
                </span>
              </div>
            </div>
            <p class="text-sm text-slate-400 mb-6">
              Cortes femininos, escova, mechas, manicure completa e design de sobrancelhas.
            </p>
            <div class="flex items-center gap-3">
              <RouterLink
                to="/agendamento/bella-vista"
                class="flex-1 text-center py-2.5 px-4 rounded-xl bg-pink-500 hover:bg-pink-400 text-slate-950 font-semibold text-sm transition shadow-md"
              >
                Agendar como Cliente
              </RouterLink>
              <RouterLink
                to="/login"
                class="py-2.5 px-4 rounded-xl bg-slate-800 hover:bg-slate-700 text-slate-200 font-medium text-sm transition border border-slate-700"
              >
                Painel Admin
              </RouterLink>
            </div>
          </div>
        </div>
      </div>

      <!-- Feature Grid -->
      <div class="mt-20 grid grid-cols-1 sm:grid-cols-3 gap-6 max-w-5xl text-left">
        <div class="p-5 rounded-xl bg-slate-900/50 border border-slate-800">
          <div class="w-10 h-10 rounded-lg bg-emerald-500/10 text-emerald-400 flex items-center justify-center mb-3">
            <ShieldCheck class="w-5 h-5" />
          </div>
          <h4 class="font-semibold text-white mb-1">Zero Dupla Reserva</h4>
          <p class="text-xs text-slate-400 leading-relaxed">
            Transações atômicas com row-level lock no Postgres e validação temporal rígida impedem conflito de horários.
          </p>
        </div>

        <div class="p-5 rounded-xl bg-slate-900/50 border border-slate-800">
          <div class="w-10 h-10 rounded-lg bg-teal-500/10 text-teal-400 flex items-center justify-center mb-3">
            <Layers class="w-5 h-5" />
          </div>
          <h4 class="font-semibold text-white mb-1">Isolamento Multi-Tenant</h4>
          <p class="text-xs text-slate-400 leading-relaxed">
            Contas 100% isoladas por tenant_id validado em middleware JWT e repositórios com escopo forçado.
          </p>
        </div>

        <div class="p-5 rounded-xl bg-slate-900/50 border border-slate-800">
          <div class="w-10 h-10 rounded-lg bg-cyan-500/10 text-cyan-400 flex items-center justify-center mb-3">
            <Zap class="w-5 h-5" />
          </div>
          <h4 class="font-semibold text-white mb-1">Filas Assíncronas Asynq</h4>
          <p class="text-xs text-slate-400 leading-relaxed">
            Worker dedicado em Go integrado com Redis para disparo instantâneo de lembretes e confirmações.
          </p>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="border-t border-slate-850 py-6 text-center text-xs text-slate-500">
      <div class="max-w-7xl mx-auto px-4 flex flex-col sm:flex-row items-center justify-between gap-4">
        <span>&copy; 2026 AgendeFácil — Plataforma de Agendamentos Multi-Tenant</span>
        <div class="flex items-center space-x-4">
          <span class="text-emerald-400 font-mono">Backend: Go 1.23 + Gin + GORM</span>
          <span class="text-cyan-400 font-mono">Frontend: Vue 3 + Vite + Tailwind</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { Scissors, Sparkles, ShieldCheck, Layers, Zap } from 'lucide-vue-next'
</script>
