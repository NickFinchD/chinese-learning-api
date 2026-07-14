<template>
  <main class="auth-layout">
    <RouterView />
  </main>
</template>

<script setup lang="ts">
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;

  display: flex;

  justify-content: center;

  align-items: center;

  background: #f5f7fb;
}
</style>