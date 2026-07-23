<template>
  <button
    :type="type"
    :disabled="disabled"
    class="w-full rounded-full bg-[#41b3a3] px-4 py-3 font-semibold text-white shadow-lg shadow-[#41b3a3]/30 transition hover:bg-[#41b3a3]/90 disabled:cursor-not-allowed disabled:opacity-50"
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
