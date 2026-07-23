<template>
  <header
    class="flex h-16 items-center justify-between border-b border-white/50 bg-white/30 px-6 backdrop-blur-xl dark:border-white/10 dark:bg-white/5"
  >
    <div class="text-xl font-bold text-[#41b3a3]">
      Wojiao
    </div>

    <div class="flex items-center gap-4">
      <RouterLink
        v-if="gamification.progress"
        to="/app/achievements"
        class="flex items-center gap-2 rounded-full border border-white/50 bg-white/30 px-3 py-1 text-sm backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10"
      >
        <span class="font-semibold text-[#41b3a3]">Ур. {{ gamification.progress.level }}</span>
        <span class="text-gray-400 dark:text-gray-500">·</span>
        <span class="text-gray-700 dark:text-gray-300">{{ gamification.progress.current_level_xp }}/{{ gamification.progress.xp_for_next_level }} XP</span>
      </RouterLink>

      <ThemeToggle />

      <span class="text-sm font-medium text-gray-700 dark:text-gray-300">
        {{ auth.user?.username }}
      </span>

      <button
        class="rounded-full bg-[#e27d60] px-3 py-2 text-white transition hover:bg-[#e27d60]/90"
        @click="onLogout"
      >
        Выйти
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { RouterLink, useRouter } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useGamificationStore } from '@/stores/gamification'
import { logout } from '@/services'

import ThemeToggle from '@/components/base/ThemeToggle.vue'

const auth = useAuthStore()
const gamification = useGamificationStore()
const router = useRouter()

async function onLogout() {
  try {
    await logout()

    auth.logout()

    await router.push('/login')
  } catch (error) {
    console.error(error)
  }
}
</script>
