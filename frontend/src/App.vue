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
  // Inicialização de tema: Páginas públicas são SEMPRE Light Mode
  const isInternalAdmin = window.location.pathname.startsWith('/admin')
  if (isInternalAdmin) {
    const adminTheme = localStorage.getItem('admin_theme') || 'dark'
    if (adminTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  } else {
    document.documentElement.classList.remove('dark')
  }

  if (authStore.token) {
    await authStore.fetchMe()
  }
})
</script>
