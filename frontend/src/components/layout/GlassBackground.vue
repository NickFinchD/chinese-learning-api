<template>
  <div class="fixed inset-0 -z-10 overflow-hidden bg-gray-50 transition-colors duration-300 dark:bg-slate-950">
    <div class="animate-blob absolute -left-40 -top-40 h-[28rem] w-[28rem] rounded-full bg-[var(--color-primary)] opacity-40 blur-3xl dark:opacity-30" />
    <div
      class="animate-blob absolute -right-32 top-1/4 h-[26rem] w-[26rem] rounded-full bg-[var(--color-secondary)] opacity-30 blur-3xl dark:opacity-25"
      style="animation-delay: -6s; animation-duration: 24s"
    />
    <div
      class="animate-blob absolute -bottom-40 left-1/4 h-[26rem] w-[26rem] rounded-full bg-[var(--color-mauve)] opacity-30 blur-3xl dark:opacity-20"
      style="animation-delay: -12s; animation-duration: 22s"
    />
    <div
      class="animate-blob absolute bottom-1/4 right-1/4 h-72 w-72 rounded-full bg-[var(--color-accent)] opacity-20 blur-3xl dark:opacity-15"
      style="animation-delay: -3s; animation-duration: 26s"
    />
  </div>
</template>

<script setup lang="ts">
</script>
