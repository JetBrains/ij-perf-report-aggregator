<template>
  <div class="warning-notice">
    <div class="warning-notice-title flex items-center">
      <ExclamationTriangleIcon class="h-5 w-5 mr-2 shrink-0" />
      {{ title }}
    </div>
    <slot />
  </div>
</template>

<script setup lang="ts">
import { ExclamationTriangleIcon } from "@heroicons/vue/20/solid/index"

defineProps<{ title: string }>()
</script>

<style scoped>
/* High-contrast amber notice that follows the app light/dark switch (.dark-mode on <html>). */
.warning-notice {
  border-radius: 0.375rem;
  border: 1px solid #fbbf24; /* amber-400 */
  background-color: #fffbeb; /* amber-50 */
  color: #92400e; /* amber-800 */
  padding: 0.75rem;
}

.warning-notice-title {
  font-weight: 500;
  color: #b45309; /* amber-700 */
  margin-bottom: 0.5rem;
}

.dark-mode .warning-notice {
  border-color: #b45309; /* amber-700 */
  background-color: #451a03; /* amber-950 */
  color: #fde68a; /* amber-200 */
}

.dark-mode .warning-notice-title {
  color: #fcd34d; /* amber-300 */
}
</style>
