<template>
  <Toolbar :class="isSticky ? 'stickyToolbar' : 'customToolbar'">
    <template
      v-for="(_, name) in $slots"
      #[name]="slotProps"
    >
      <slot
        :name="name"
        v-bind="slotProps"
      ></slot>
    </template>
  </Toolbar>
</template>
<script setup lang="ts">
import { onMounted, onUnmounted, ref } from "vue"

const isSticky = ref(false)
const checkIfSticky = () => (isSticky.value = window.scrollY > 100)
onMounted(() => {
  window.addEventListener("scroll", checkIfSticky)
})
onUnmounted(() => {
  window.removeEventListener("scroll", checkIfSticky)
})
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
