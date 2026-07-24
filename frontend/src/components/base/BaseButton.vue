<template>
  <button
    :type="type"
    :disabled="disabled"
    class="w-full rounded-full bg-[var(--color-primary)] px-4 py-3 font-semibold text-white shadow-lg shadow-[var(--color-primary)]/30 transition hover:bg-[var(--color-primary)]/90 active:scale-95 disabled:cursor-not-allowed disabled:opacity-50 disabled:active:scale-100"
  >
    <slot />
  </button>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    type?: 'button' | 'submit'
    disabled?: boolean
  }>(),
  {
    type: 'button',
    disabled: false,
  },
)
</script>
