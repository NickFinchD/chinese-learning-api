<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition-opacity duration-150"
      enter-from-class="opacity-0"
      leave-active-class="transition-opacity duration-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black/40 p-4"
        @click.self="$emit('cancel')"
      >
        <div class="w-full max-w-sm rounded-2xl border border-white/50 bg-white/90 p-6 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/95">
          <h2 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
            {{ title }}
          </h2>

          <p class="mb-6 text-sm text-gray-600 dark:text-gray-400">
            {{ message }}
          </p>

          <div class="flex justify-end gap-3">
            <button
              type="button"
              class="rounded-full px-5 py-2.5 text-sm font-semibold text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
              @click="$emit('cancel')"
            >
              Отмена
            </button>

            <button
              type="button"
              :disabled="loading"
              class="rounded-full bg-red-500 px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-red-500/30 transition hover:bg-red-500/90 disabled:cursor-not-allowed disabled:opacity-50"
              @click="$emit('confirm')"
            >
              {{ loading ? 'Удаление...' : 'Удалить' }}
            </button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    open: boolean
    title: string
    message: string
    loading?: boolean
  }>(),
  {
    loading: false,
  },
)

defineEmits<{
  (e: 'cancel'): void
  (e: 'confirm'): void
}>()
</script>
