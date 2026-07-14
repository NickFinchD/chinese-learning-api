<template>
  <input
    :value="modelValue"
    :type="type"
    :placeholder="placeholder"
    class="w-full rounded-xl border border-gray-300 px-4 py-3 outline-none transition focus:border-blue-500"
    @input="onInput"
  >
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    modelValue: string
    type?: string
    placeholder?: string
  }>(),
  {
    type: 'text',
    placeholder: '',
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void
}>()

function onInput(event: Event) {
  emit(
    'update:modelValue',
    (event.target as HTMLInputElement).value,
  )
}
</script>