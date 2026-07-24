<template>
  <button
    v-if="supported"
    type="button"
    :disabled="!hasVoice"
    :title="hasVoice ? `Произнести: ${text}` : 'На этом устройстве не установлен голос для китайского языка'"
    class="inline-flex shrink-0 items-center justify-center rounded-full text-[var(--color-primary)] transition hover:bg-[var(--color-primary)]/10 disabled:cursor-not-allowed disabled:text-gray-300 disabled:hover:bg-transparent dark:text-[var(--color-mint)] dark:hover:bg-[var(--color-mint)]/10 dark:disabled:text-gray-600"
    :class="sizeClass"
    @click.stop.prevent="speak(text)"
  >
    <AppIcon
      name="volume"
      :size="iconSize"
    />
  </button>
</template>

<script setup lang="ts">
import { computed } from 'vue'

import AppIcon from './AppIcon.vue'
import { isSpeechSupported, speak, useChineseVoiceAvailable } from '@/utils/speech'

const props = withDefaults(
  defineProps<{
    text: string
    size?: 'sm' | 'md'
  }>(),
  {
    size: 'md',
  },
)

const supported = isSpeechSupported()
const hasVoice = useChineseVoiceAvailable()

const sizeClass = computed(() => (props.size === 'sm' ? 'h-7 w-7' : 'h-9 w-9'))
const iconSize = computed(() => (props.size === 'sm' ? 16 : 20))
</script>
