<template>
  <Toolbar :class="useScrollStore().isScrolled ? 'stickyToolbar' : 'customToolbar'">
    <template #start>
      <slot name="start" />
    </template>
    <template #center>
      <slot name="center" />
    </template>
    <template #end>
      <slot name="end" />
    </template>
  </Toolbar>
</template>
<script setup lang="ts">
import { useScrollListeners, useScrollStore } from "./scrollStore"

useScrollListeners()
</script>
<style scoped>
.customToolbar {
  background-color: transparent;
  border: none;
  padding: 0;
}

.stickyToolbar {
  top: 0rem;
  padding: 0.7rem 0.7rem 0.7rem 0.7rem;
  border-radius: 0;
  position: sticky;
  z-index: 100;
}
</style>
