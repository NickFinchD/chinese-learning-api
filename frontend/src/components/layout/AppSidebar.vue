<template>
  <Transition
    enter-active-class="transition-opacity duration-200"
    enter-from-class="opacity-0"
    leave-active-class="transition-opacity duration-200"
    leave-to-class="opacity-0"
  >
    <div
      v-if="open"
      class="fixed inset-0 z-30 bg-black/40 md:hidden"
      @click="$emit('close')"
    />
  </Transition>

  <aside
    class="fixed inset-y-0 left-0 z-40 w-64 -translate-x-full border-r border-white/50 bg-white/70 backdrop-blur-xl transition-transform duration-200 md:static md:z-auto md:translate-x-0 md:bg-white/30 dark:border-white/10 dark:bg-slate-950/90 dark:md:bg-white/5"
    :class="{ 'translate-x-0': open }"
  >
    <nav class="flex flex-col gap-2 p-4">

      <RouterLink
        to="/app"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="home" />
        Главная
      </RouterLink>

      <RouterLink
        to="/app/courses"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="graduation-cap" />
        Курсы
      </RouterLink>

      <RouterLink
        to="/app/vocabulary"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="book-open" />
        Словарь
      </RouterLink>

      <RouterLink
        to="/app/vocabulary/collections"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="folder" />
        Подборки
      </RouterLink>

      <RouterLink
        to="/app/texts"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="file-text" />
        Тексты
      </RouterLink>

      <RouterLink
        to="/app/tests"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="flask" />
        Тесты
      </RouterLink>

      <RouterLink
        to="/app/achievements"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="trophy" />
        Достижения
      </RouterLink>

      <RouterLink
        to="/app/settings"
        class="flex items-center gap-3 rounded-lg px-4 py-3 text-gray-700 transition hover:bg-white/50 dark:text-gray-300 dark:hover:bg-white/10"
        active-class="bg-[var(--color-primary)]/15 text-[var(--color-primary)] font-semibold dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]"
        @click="$emit('close')"
      >
        <AppIcon name="settings" />
        Настройки
      </RouterLink>

    </nav>
  </aside>
</template>

<script setup lang="ts">
import AppIcon from '@/components/base/AppIcon.vue'

defineProps<{
  open: boolean
}>()

defineEmits<{
  (e: 'close'): void
}>()
</script>
