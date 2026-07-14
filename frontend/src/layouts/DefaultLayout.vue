<template>
  <div class="layout">

    <header class="header">
      Chinese Learning
    </header>

    <main class="content">
      <RouterView />
    </main>

  </div>
</template>

<script setup lang="ts">
</script>

<style scoped>
.layout {
  min-height: 100vh;

  display: flex;

  flex-direction: column;
}

.header {
  height: 64px;

  display: flex;

  align-items: center;

  padding: 0 24px;

  background: #2563eb;

  color: white;

  font-weight: 600;
}

.content {
  flex: 1;

  padding: 32px;
}
</style>