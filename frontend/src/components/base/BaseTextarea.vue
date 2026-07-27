<template>
  <textarea
    :value="modelValue"
    :placeholder="placeholder"
    :rows="rows"
    class="w-full resize-y rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-gray-900 placeholder-gray-500 outline-none backdrop-blur-md transition focus:border-[var(--color-primary)] focus:bg-white/60 dark:border-white/10 dark:bg-white/5 dark:text-white dark:placeholder-gray-400 dark:focus:bg-white/10"
    @input="onInput"
  />
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    placeholder?: string
    rows?: number
  }>(),
  {
    placeholder: '',
    rows: 3,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

function onInput(event: Event) {
  emit(
    'update:modelValue',
    (event.target as HTMLTextAreaElement).value,
  )
}
</script>
