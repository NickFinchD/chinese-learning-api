<template>
  <div>
    <h1 class="mb-8 text-3xl font-bold text-gray-900 dark:text-white">
      Настройки
    </h1>

    <BaseCard class="max-w-md">
      <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
        Аккаунт
      </h2>

      <div class="mb-4">
        <div class="text-sm text-gray-500 dark:text-gray-400">
          Имя пользователя
        </div>
        <div class="text-gray-900 dark:text-white">
          {{ auth.user?.username }}
        </div>
      </div>

      <div class="mb-6">
        <div class="text-sm text-gray-500 dark:text-gray-400">
          Эл. почта
        </div>
        <div class="text-gray-900 dark:text-white">
          {{ auth.user?.email }}
        </div>
      </div>

      <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
        Оформление
      </h2>

      <div class="mb-6 flex items-center justify-between">
        <div class="text-gray-700 dark:text-gray-300">
          Тёмная тема
        </div>
        <ThemeToggle />
      </div>

      <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
        Цветовая гамма
      </h2>

      <div class="mb-6 grid grid-cols-1 gap-3 sm:grid-cols-3">
        <button
          v-for="(paletteOption, id) in PALETTES"
          :key="id"
          type="button"
          class="rounded-xl border p-3 text-left transition"
          :class="palette.paletteId === id
            ? 'border-[var(--color-primary)] bg-[var(--color-primary)]/10'
            : 'border-white/50 bg-white/20 hover:bg-white/40 dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10'"
          @click="palette.setPalette(id as PaletteId)"
        >
          <div class="mb-2 flex overflow-hidden rounded-lg">
            <span
              v-for="hex in Object.values(paletteOption.colors)"
              :key="hex"
              class="h-6 flex-1"
              :style="{ backgroundColor: hex }"
            />
          </div>

          <div class="flex items-center gap-1.5 text-sm font-medium text-gray-800 dark:text-gray-200">
            <AppIcon
              v-if="palette.paletteId === id"
              name="check-circle"
              :size="14"
              filled
              class="text-[var(--color-primary)]"
            />
            {{ paletteOption.label }}
          </div>
        </button>
      </div>

      <h2 class="mb-4 text-xl font-semibold text-gray-900 dark:text-white">
        Китайский шрифт
      </h2>

      <div class="mb-6">
        <BaseSelect
          :model-value="fonts.hanziFont"
          :options="HANZI_FONT_OPTIONS"
          @update:model-value="value => fonts.setHanziFont(value as HanziFont)"
        />

        <div class="font-hanzi mt-4 rounded-xl border border-white/50 bg-white/20 p-4 text-center text-3xl text-gray-900 backdrop-blur-md dark:border-white/10 dark:bg-white/5 dark:text-white">
          你好，世界
        </div>
      </div>

      <BaseButton @click="onLogout">
        Выйти
      </BaseButton>
    </BaseCard>
  </div>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { logout } from '@/services'
import { HANZI_FONT_OPTIONS, useFontsStore } from '@/stores/fonts'
import type { HanziFont } from '@/stores/fonts'
import { PALETTES, usePaletteStore } from '@/stores/palette'
import type { PaletteId } from '@/stores/palette'

import BaseCard from '@/components/base/BaseCard.vue'
import BaseButton from '@/components/base/BaseButton.vue'
import BaseSelect from '@/components/base/BaseSelect.vue'
import ThemeToggle from '@/components/base/ThemeToggle.vue'
import AppIcon from '@/components/base/AppIcon.vue'

const auth = useAuthStore()
const fonts = useFontsStore()
const palette = usePaletteStore()

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
