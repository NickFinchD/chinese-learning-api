<template>
  <svg
    xmlns="http://www.w3.org/2000/svg"
    viewBox="0 0 24 24"
    :width="size"
    :height="size"
    fill="none"
    stroke="currentColor"
    stroke-width="1.8"
    stroke-linecap="round"
    stroke-linejoin="round"
    class="shrink-0"
    aria-hidden="true"
    v-html="markup"
  />
</template>

<script setup lang="ts">
import { computed } from 'vue'

export type IconName =
  | 'home'
  | 'graduation-cap'
  | 'book-open'
  | 'file-text'
  | 'flask'
  | 'trophy'
  | 'lock'
  | 'settings'
  | 'check-circle'
  | 'x-circle'
  | 'check'
  | 'star'
  | 'sparkles'
  | 'chevron-down'
  | 'sun'
  | 'moon'
  | 'sort'
  | 'refresh'
  | 'menu'
  | 'x'
  | 'volume'
  | 'stop'
  | 'eye'
  | 'eye-off'
  | 'folder'
  | 'plus'
  | 'trash'
  | 'pencil'
  | 'arrow-left'
  | 'clock'

const props = withDefaults(
  defineProps<{
    name: IconName
    size?: number | string
    filled?: boolean
  }>(),
  {
    size: 20,
    filled: false,
  },
)

// Hand-authored 24x24 stroke icons (currentColor), kept intentionally simple
// so they read cleanly at small sidebar/inline sizes.
const icons: Record<IconName, string> = {
  'home': '<path d="M4 11.5 12 4l8 7.5" /><path d="M6 10v9a1 1 0 0 0 1 1h3v-6h4v6h3a1 1 0 0 0 1-1v-9" />',
  'graduation-cap': '<path d="M2 8.5 12 4l10 4.5-10 4.5-10-4.5Z" /><path d="M6 10.7v4.6c0 1.5 2.7 3.2 6 3.2s6-1.7 6-3.2v-4.6" /><path d="M21 8.5v6" />',
  'book-open': '<path d="M12 6.5c-1.6-1.2-4-1.8-6.5-1.8-.6 0-1 .4-1 1v11.6c0 .6.4 1 1 1 2.5 0 4.9.6 6.5 1.7 1.6-1.1 4-1.7 6.5-1.7.6 0 1-.4 1-1V5.7c0-.6-.4-1-1-1-2.5 0-4.9.6-6.5 1.8Z" /><path d="M12 6.5v13" />',
  'file-text': '<path d="M7 3.5h7.5L19 8v12a1 1 0 0 1-1 1H7a1 1 0 0 1-1-1v-15.5a1 1 0 0 1 1-1Z" /><path d="M14 3.5V8h5" /><path d="M9 13h6M9 16.5h6" />',
  'flask': '<path d="M10 3.5h4" /><path d="M10.5 3.5v5.7L5.8 18a1.6 1.6 0 0 0 1.4 2.4h9.6a1.6 1.6 0 0 0 1.4-2.4l-4.7-8.8V3.5" /><path d="M7.5 15h9" />',
  'trophy': '<path d="M7 4h10v5.5a5 5 0 0 1-10 0V4Z" /><path d="M7 5.5H4.5a1 1 0 0 0-1 1V8a3 3 0 0 0 3 3" /><path d="M17 5.5h2.5a1 1 0 0 1 1 1V8a3 3 0 0 1-3 3" /><path d="M12 14.5V17" /><path d="M8.5 20.5h7" /><path d="M9.5 17h5l.8 3.5h-6.6l.8-3.5Z" />',
  'lock': '<rect x="5" y="10.5" width="14" height="9.5" rx="1.5" /><path d="M8 10.5V7.5a4 4 0 0 1 8 0v3" />',
  'settings': '<path d="M10.44,2.12 L13.56,2.12 L14.79,5.26 L17.88,3.91 L20.09,6.12 L18.74,9.21 L21.88,10.44 L21.88,13.56 L18.74,14.79 L20.09,17.88 L17.88,20.09 L14.79,18.74 L13.56,21.88 L10.44,21.88 L9.21,18.74 L6.12,20.09 L3.91,17.88 L5.26,14.79 L2.12,13.56 L2.12,10.44 L5.26,9.21 L3.91,6.12 L6.12,3.91 L9.21,5.26 Z" /><circle cx="12" cy="12" r="4.2" />',
  'check-circle': '<circle cx="12" cy="12" r="8.5" /><path d="m8.5 12.3 2.4 2.4 4.6-5.4" />',
  'x-circle': '<circle cx="12" cy="12" r="8.5" /><path d="m9 9 6 6M15 9l-6 6" />',
  'check': '<path d="m5 12.5 4.5 4.5L19 7" />',
  'star': '<path d="m12 3.5 2.7 5.6 6.1.9-4.4 4.3 1 6.1L12 17.3l-5.4 3.1 1-6.1-4.4-4.3 6.1-.9L12 3.5Z" />',
  'sparkles': '<path d="M11.5 3.5 13 8l4.5 1.5L13 11l-1.5 4.5L10 11l-4.5-1.5L10 8l1.5-4.5Z" /><path d="M18.5 14.5 19.3 17l2.2.8-2.2.8-.8 2.4-.8-2.4-2.2-.8 2.2-.8.8-2.5Z" />',
  'chevron-down': '<path d="m6 9 6 6 6-6" />',
  'sun': '<circle cx="12" cy="12" r="4.2" /><path d="M12 3v2M12 19v2M4.2 4.2l1.4 1.4M18.4 18.4l1.4 1.4M3 12h2M19 12h2M4.2 19.8l1.4-1.4M18.4 5.6l1.4-1.4" />',
  'moon': '<path d="M20 14.5A8.5 8.5 0 1 1 9.5 4 6.8 6.8 0 0 0 20 14.5Z" />',
  'sort': '<path d="M7 4v16M7 20l-3.5-3.5M7 20l3.5-3.5" /><path d="M17 20V4M17 4l-3.5 3.5M17 4l3.5 3.5" />',
  'refresh': '<path d="M3 12a9 9 0 0 1 15-6.7L21 8" /><path d="M21 3v5h-5" /><path d="M21 12a9 9 0 0 1-15 6.7L3 16" /><path d="M3 21v-5h5" />',
  'menu': '<path d="M4 6h16M4 12h16M4 18h16" />',
  'x': '<path d="m6 6 12 12M18 6 6 18" />',
  'volume': '<path d="M4 9.5v5h3.5L13 19V5L7.5 9.5H4Z" /><path d="M17 8.5a5 5 0 0 1 0 7" /><path d="M19.5 6a8.5 8.5 0 0 1 0 12" />',
  'stop': '<rect x="5" y="5" width="14" height="14" rx="2" fill="currentColor" stroke="none" />',
  'eye': '<path d="M2.5 12S6 5.5 12 5.5 21.5 12 21.5 12 18 18.5 12 18.5 2.5 12 2.5 12Z" /><circle cx="12" cy="12" r="2.8" />',
  'eye-off': '<path d="M2.5 12S6 5.5 12 5.5c1.6 0 3 .4 4.2 1M21.5 12S18 18.5 12 18.5c-1.6 0-3-.4-4.2-1" /><path d="M9.5 9.6a2.8 2.8 0 0 0 3.9 3.9" /><path d="m3.5 3.5 17 17" />',
  'folder': '<path d="M3.5 6.5a1 1 0 0 1 1-1h4.4l1.6 2h8a1 1 0 0 1 1 1v9a1 1 0 0 1-1 1h-14a1 1 0 0 1-1-1v-11Z" />',
  'plus': '<path d="M12 5v14M5 12h14" />',
  'trash': '<path d="M4.5 7h15" /><path d="M9.5 7V5a1 1 0 0 1 1-1h3a1 1 0 0 1 1 1v2" /><path d="M6.5 7l1 12.5a1 1 0 0 0 1 .9h7a1 1 0 0 0 1-.9l1-12.5" /><path d="M10 11v6M14 11v6" />',
  'pencil': '<path d="M4 20l.9-3.9L15.6 5.4a1.5 1.5 0 0 1 2.1 0l1 1a1.5 1.5 0 0 1 0 2.1L8 19.1 4 20Z" /><path d="m13.8 7 3.2 3.2" />',
  'arrow-left': '<path d="M19 12H5M11 6l-6 6 6 6" />',
  'clock': '<circle cx="12" cy="12" r="8.5" /><path d="M12 7v5l3.5 2" />',
}

const markup = computed(() => {
  const raw = icons[props.name]

  if (props.name === 'star') {
    return props.filled ? raw.replace('<path', '<path fill="currentColor"') : raw
  }

  if (props.name === 'trophy' && props.filled) {
    return raw.replace(
      '<path d="M7 4h10v5.5a5 5 0 0 1-10 0V4Z" />',
      '<path d="M7 4h10v5.5a5 5 0 0 1-10 0V4Z" fill="currentColor" />',
    )
  }

  return raw
})
</script>
