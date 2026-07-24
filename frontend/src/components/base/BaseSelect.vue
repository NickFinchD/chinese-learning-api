<template>
  <div
    ref="root"
    class="relative"
  >
    <button
      type="button"
      class="flex w-full items-center justify-between gap-3 rounded-xl border border-white/50 bg-white/40 px-4 py-3 text-left text-gray-900 outline-none backdrop-blur-md transition focus:border-[var(--color-primary)] dark:border-white/10 dark:bg-white/5 dark:text-white"
      @click="open = !open"
    >
      <span class="truncate">{{ selectedLabel }}</span>

      <AppIcon
        name="chevron-down"
        :size="16"
        class="text-gray-500 transition dark:text-gray-400"
        :class="{ 'rotate-180': open }"
      />
    </button>

    <Transition
      enter-active-class="transition duration-150 ease-out"
      enter-from-class="opacity-0 -translate-y-1"
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-100 ease-in"
      leave-from-class="opacity-100 translate-y-0"
      leave-to-class="opacity-0 -translate-y-1"
    >
      <ul
        v-if="open"
        class="absolute z-20 mt-2 max-h-64 w-full overflow-auto rounded-xl border border-white/50 bg-white/90 py-1 shadow-xl backdrop-blur-xl dark:border-white/10 dark:bg-slate-900/90"
      >
        <li
          v-for="option in options"
          :key="option.value"
          class="cursor-pointer px-4 py-2.5 text-sm transition"
          :class="option.value === modelValue
            ? 'bg-[var(--color-primary)]/15 font-medium text-[var(--color-primary)] dark:bg-[var(--color-primary)]/20 dark:text-[var(--color-mint)]'
            : 'text-gray-700 hover:bg-white/60 dark:text-gray-300 dark:hover:bg-white/10'"
          @click="select(option.value)"
        >
          {{ option.label }}
        </li>
      </ul>
    </Transition>
  </div>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'

import AppIcon from './AppIcon.vue'

interface Option {
  value: string | number
  label: string
}

const props = defineProps<{
  modelValue: string | number
  options: Option[]
}>()

const emit = defineEmits<{
  (e: 'update:modelValue', value: string | number): void
}>()

const root = ref<HTMLElement | null>(null)
const open = ref(false)

const selectedLabel = computed(() => {
  return props.options.find(option => option.value === props.modelValue)?.label ?? ''
})

function select(value: string | number) {
  emit('update:modelValue', value)

  open.value = false
}

function onClickOutside(event: MouseEvent) {
  if (root.value && !root.value.contains(event.target as Node)) {
    open.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', onClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('click', onClickOutside)
})
</script>
