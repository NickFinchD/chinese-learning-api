<template>
  <div class="space-y-2">
    <div
      v-for="(item, index) in modelValue"
      :key="index"
      class="flex items-center gap-2"
    >
      <div class="flex-1">
        <slot
          :item="item"
          :index="index"
          :update="(value: T) => updateAt(index, value)"
        />
      </div>

      <div class="flex shrink-0 gap-1">
        <button
          type="button"
          title="Вверх"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="index === 0"
          @click="move(index, -1)"
        >
          <AppIcon
            name="chevron-down"
            :size="16"
            class="rotate-180"
          />
        </button>

        <button
          type="button"
          title="Вниз"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-gray-600 transition hover:bg-white/60 disabled:opacity-30 dark:text-gray-300 dark:hover:bg-white/10"
          :disabled="index === modelValue.length - 1"
          @click="move(index, 1)"
        >
          <AppIcon
            name="chevron-down"
            :size="16"
          />
        </button>

        <button
          type="button"
          title="Удалить"
          class="flex h-9 w-9 items-center justify-center rounded-lg text-red-500 transition hover:bg-red-500/10 disabled:opacity-30"
          :disabled="modelValue.length <= minItems"
          @click="remove(index)"
        >
          <AppIcon
            name="trash"
            :size="16"
          />
        </button>
      </div>
    </div>

    <button
      type="button"
      class="text-sm font-medium text-[var(--color-primary)] transition hover:text-[var(--color-primary)]/80"
      @click="add"
    >
      + {{ addLabel }}
    </button>
  </div>
</template>

<script setup lang="ts" generic="T">
import AppIcon from '@/components/base/AppIcon.vue'

const props = withDefaults(
  defineProps<{
    modelValue: T[]
    newItem: () => T
    addLabel?: string
    minItems?: number
  }>(),
  {
    addLabel: 'Добавить',
    minItems: 0,
  },
)

const emit = defineEmits<{
  (e: 'update:modelValue', value: T[]): void
}>()

function updateAt(index: number, value: T) {
  const next = [...props.modelValue]
  next[index] = value
  emit('update:modelValue', next)
}

function add() {
  emit('update:modelValue', [...props.modelValue, props.newItem()])
}

function remove(index: number) {
  emit('update:modelValue', props.modelValue.filter((_, i) => i !== index))
}

function move(index: number, direction: number) {
  const target = index + direction

  if (target < 0 || target >= props.modelValue.length) return

  const next = [...props.modelValue]
  const [moved] = next.splice(index, 1)
  next.splice(target, 0, moved)

  emit('update:modelValue', next)
}
</script>
