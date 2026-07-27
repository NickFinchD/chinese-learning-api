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
        @click.self="$emit('close')"
      >
        <div class="max-h-[90vh] w-full max-w-lg overflow-y-auto rounded-2xl border border-white/50 bg-white/90 p-6 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/95">
          <div class="mb-5 flex items-center justify-between">
            <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
              {{ title }}
            </h2>

            <button
              type="button"
              class="flex h-8 w-8 items-center justify-center rounded-lg text-gray-500 transition hover:bg-white/60 dark:text-gray-400 dark:hover:bg-white/10"
              @click="$emit('close')"
            >
              <AppIcon name="x" />
            </button>
          </div>

          <form
            class="space-y-4"
            @submit.prevent="$emit('submit')"
          >
            <slot />

            <p
              v-if="error"
              class="text-sm text-red-600 dark:text-red-400"
            >
              {{ error }}
            </p>

            <div class="flex justify-end gap-3 pt-2">
              <button
                type="button"
                class="rounded-full px-5 py-2.5 text-sm font-semibold text-gray-600 transition hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10"
                @click="$emit('close')"
              >
                Отмена
              </button>

              <button
                type="submit"
                :disabled="saving"
                class="rounded-full bg-[var(--color-primary)] px-5 py-2.5 text-sm font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {{ saving ? 'Сохранение...' : 'Сохранить' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import AppIcon from '@/components/base/AppIcon.vue'

withDefaults(
  defineProps<{
    open: boolean
    title: string
    saving?: boolean
    error?: string
  }>(),
  {
    saving: false,
    error: '',
  },
)

defineEmits<{
  (e: 'close'): void
  (e: 'submit'): void
}>()
</script>
