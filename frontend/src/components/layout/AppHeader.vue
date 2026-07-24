<template>
  <header
    class="flex h-16 items-center justify-between gap-2 border-b border-white/50 bg-white/30 px-3 backdrop-blur-xl sm:px-6 dark:border-white/10 dark:bg-white/5"
  >
    <div class="flex min-w-0 items-center gap-2">
      <button
        type="button"
        class="flex h-9 w-9 shrink-0 items-center justify-center rounded-lg text-gray-700 transition hover:bg-white/50 md:hidden dark:text-gray-300 dark:hover:bg-white/10"
        @click="$emit('toggle-sidebar')"
      >
        <AppIcon name="menu" />
      </button>

      <RouterLink
        to="/app"
        class="shrink-0 text-xl font-bold text-[var(--color-primary)] transition hover:text-[var(--color-primary)]/80"
      >
        Wojiao
      </RouterLink>
    </div>

    <div class="flex min-w-0 items-center gap-2 sm:gap-4">
      <RouterLink
        v-if="gamification.progress"
        to="/app/achievements"
        class="flex shrink-0 items-center gap-2 rounded-full border border-white/50 bg-white/30 px-3 py-1 text-sm backdrop-blur-md transition hover:bg-white/50 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10"
      >
        <span class="font-semibold text-[var(--color-primary)]">Ур. {{ gamification.progress.level }}</span>
        <span class="hidden text-gray-400 sm:inline dark:text-gray-500">·</span>
        <span class="hidden text-gray-700 sm:inline dark:text-gray-300">{{ gamification.progress.current_level_xp }}/{{ gamification.progress.xp_for_next_level }} XP</span>
      </RouterLink>

      <ThemeToggle />

      <span class="hidden truncate text-sm font-medium text-gray-700 sm:inline dark:text-gray-300">
        {{ auth.user?.username }}
      </span>

      <button
        class="shrink-0 rounded-full bg-[var(--color-secondary)] px-3 py-2 text-[var(--color-secondary-text)] transition hover:bg-[var(--color-secondary)]/90"
        @click="onLogout"
      >
        Выйти
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { RouterLink } from 'vue-router'

import { useAuthStore } from '@/stores/auth'
import { useGamificationStore } from '@/stores/gamification'
import { logout } from '@/services'

import ThemeToggle from '@/components/base/ThemeToggle.vue'
import AppIcon from '@/components/base/AppIcon.vue'

defineEmits<{
  (e: 'toggle-sidebar'): void
}>()

const auth = useAuthStore()
const gamification = useGamificationStore()

async function onLogout() {
  try {
    await logout()

    auth.logout()

    // Full reload (not router.push) so every Pinia store resets — otherwise
    // the next account to log in in this tab would see stale data left over
    // from this session (courses, progress, XP, etc. are in-memory only).
    window.location.href = '/login'
  } catch (error) {
    console.error(error)
  }
}
</script>
