<template>
  <Accordion
    class="chart-accordion"
    v-bind="$attrs"
  >
    <slot />
  </Accordion>
</template>

<script setup lang="ts">
import Accordion from "primevue/accordion"

// Forward fallthrough attrs explicitly so the chart-accordion class is not duplicated on the root.
defineOptions({ inheritAttrs: false })
</script>

<!--suppress CssUnusedSymbol -->
<style scoped>
/* Show the expand/collapse chevron to the left of the panel title instead of the right edge.
   PrimeVue's runtime-injected base rule `.p-accordionheader { justify-content: space-between }`
   competes at equal specificity, so `!important` is needed to force left-packing. `order: -1`
   moves the toggle icon ahead of the title; the icon's trailing margin sets the gap (more
   reliable than flex `gap`, which can skip the anonymous text-node title item). */
:deep(.p-accordionheader) {
  justify-content: flex-start !important;
}

:deep(.p-accordionheader-toggle-icon) {
  order: -1;
  margin-inline-end: 0.5rem;
}
</style>
