<template>
  <div class="min-h-screen font-sans antialiased text-[#202224] dark:text-zinc-100">
    <RouterView />
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useAuthStore } from './stores/auth'

const authStore = useAuthStore()

onMounted(async () => {
  // Inicialização do tema (Light por padrão do SaaS Figma ou Dark conforme preferência salva)
  const savedTheme = localStorage.getItem('theme')
  if (savedTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }

  if (authStore.token) {
    await authStore.fetchMe()
  }
})
</script>
