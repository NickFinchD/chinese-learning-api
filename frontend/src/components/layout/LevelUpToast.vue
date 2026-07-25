<template>
  <Teleport to="body">
    <Transition name="level-up">
      <div
        v-if="gamification.leveledUpTo"
        class="fixed top-24 left-1/2 z-50 -translate-x-1/2"
        role="status"
      >
        <div
          class="animate-pop-in flex items-center gap-3 rounded-2xl border border-[var(--color-accent)]/50 bg-white/80 px-5 py-3 shadow-lg backdrop-blur-xl dark:border-[var(--color-accent)]/30 dark:bg-gray-900/80"
        >
          <AppIcon
            name="sparkles"
            :size="26"
            class="shrink-0 text-[var(--color-accent)]"
          />

          <div>
            <div class="font-semibold text-gray-900 dark:text-white">
              Новый уровень!
            </div>
            <div class="text-sm text-gray-600 dark:text-gray-400">
              Вы достигли {{ gamification.leveledUpTo }}-го уровня
            </div>
          </div>

          <button
            type="button"
            class="ml-2 shrink-0 rounded-full p-1 text-gray-400 transition hover:bg-white/50 hover:text-gray-600 dark:hover:bg-white/10 dark:hover:text-gray-300"
            aria-label="Закрыть"
            @click="gamification.dismissLevelUp()"
          >
            <AppIcon
              name="x"
              :size="16"
            />
          </button>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { watch } from 'vue'

import { useGamificationStore } from '@/stores/gamification'
import AppIcon from '@/components/base/AppIcon.vue'

const AUTO_DISMISS_MS = 4500

const gamification = useGamificationStore()

let dismissTimeout: ReturnType<typeof setTimeout> | undefined

watch(
  () => gamification.leveledUpTo,
  (level) => {
    clearTimeout(dismissTimeout)

    if (level !== null) {
      dismissTimeout = setTimeout(() => gamification.dismissLevelUp(), AUTO_DISMISS_MS)
    }
  },
)
</script>

<style scoped>
.level-up-enter-active,
.level-up-leave-active {
  transition: opacity 0.25s ease, transform 0.25s ease;
}
.level-up-enter-from,
.level-up-leave-to {
  opacity: 0;
  transform: translate(-50%, -10px);
}
</style>
