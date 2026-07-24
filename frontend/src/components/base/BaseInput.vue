<template>
  <div
    v-if="type === 'password'"
    class="relative"
  >
    <input
      :value="modelValue"
      :type="visible ? 'text' : 'password'"
      :placeholder="placeholder"
      class="w-full rounded-xl border border-white/50 bg-white/40 px-4 py-3 pr-11 text-gray-900 placeholder-gray-500 outline-none backdrop-blur-md transition focus:border-[var(--color-primary)] focus:bg-white/60 dark:border-white/10 dark:bg-white/5 dark:text-white dark:placeholder-gray-400 dark:focus:bg-white/10"
      @input="onInput"
    >

    <button
      type="button"
      tabindex="-1"
      :title="visible ? 'Скрыть пароль' : 'Показать пароль'"
      class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 transition hover:text-gray-600 dark:text-gray-500 dark:hover:text-gray-300"
      @click="visible = !visible"
    >
      <AppIcon :name="visible ? 'eye-off' : 'eye'" />
    </button>
  </div>

  <input
    v-else
    :value="modelValue"
    :type="type"
    :placeholder="placeholder"
    class="w-full rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-gray-900 placeholder-gray-500 outline-none backdrop-blur-md transition focus:border-[var(--color-primary)] focus:bg-white/60 dark:border-white/10 dark:bg-white/5 dark:text-white dark:placeholder-gray-400 dark:focus:bg-white/10"
    @input="onInput"
  >
</template>

<script setup lang="ts">
import { ref } from 'vue'

import AppIcon from './AppIcon.vue'

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

const visible = ref(false)

function onInput(event: Event) {
  emit(
    'update:modelValue',
    (event.target as HTMLInputElement).value,
  )
}
</script>
